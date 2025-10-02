package local

import (
	"context"
	"net/http"
	"time"

	"github.com/moledoc/orderly/actions"
	"github.com/moledoc/orderly/middleware"
	"github.com/moledoc/orderly/models"
	"github.com/moledoc/orderly/storage"
	"github.com/moledoc/orderly/utils"
)

type StorageOrder map[uint][]*models.Order

func NewStorageOrder() storage.StorageOrderAPI {
	return make(StorageOrder)
}

func (s StorageOrder) Close(ctx context.Context) {

	middleware.SpanStart(ctx, "LocalStorageOrder:Close")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Close")

	s = nil
}

func (s StorageOrder) Read(ctx context.Context, action actions.Action, id uint) ([]*models.Order, models.IError) {
	middleware.SpanStart(ctx, "LocalStorageOrder:Read")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Read")

	if s == nil {
		return nil, models.NewError(http.StatusInternalServerError, "localstorage not initialized for read")
	}

	switch action {
	case actions.READ:
		middleware.SpanStart(ctx, "LocalStorageOrder:Read:READ")
		defer middleware.SpanStop(ctx, "LocalStorageOrder:Read:READ")
		os, ok := s[id]
		if !ok || len(os) == 0 {
			return nil, models.NewError(http.StatusNotFound, "not found during read")
		}
		return []*models.Order{os[len(os)-1]}, nil
	case actions.READVERSIONS:
		middleware.SpanStart(ctx, "LocalStorageOrder:Read:READVERSIONS")
		defer middleware.SpanStop(ctx, "LocalStorageOrder:Read:READVERSIONS")
		os, ok := s[id]
		if !ok || len(os) == 0 {
			return nil, models.NewError(http.StatusNotFound, "not found during read")
		}
		return os, nil
	case actions.READALL:
		middleware.SpanStart(ctx, "LocalStorageOrder:Read:READALL")
		defer middleware.SpanStop(ctx, "LocalStorageOrder:Read:READALL")
		oss := make([]*models.Order, len(s))
		i := 0
		for _, us := range s {
			if len(us) == 0 {
				continue
			}
			oss[i] = us[len(us)-1]
			i += 1
		}
		return oss, nil
	case actions.READSUBTASKS:
		middleware.SpanStart(ctx, "LocalStorageOrder:Read:READSUBTASKS")
		defer middleware.SpanStop(ctx, "LocalStorageOrder:Read:READSUBTASKS")
		ssubtask, ok := s[id]
		if !ok || len(ssubtask) == 0 {
			return nil, models.NewError(http.StatusNotFound, "not found during read")
		}
		subtask := ssubtask[len(ssubtask)-1]

		var oss []*models.Order
		for _, os := range s {
			if len(os) == 0 {
				continue
			}
			if o := os[len(os)-1]; o.ParentOrderID != nil && *o.ParentOrderID == *subtask.ParentOrderID {
				oss = append(oss, o)
			}
		}
		return oss, nil
	default:
		return nil, models.NewError(http.StatusInternalServerError, "undefined read action")
	}
}

func (s StorageOrder) Write(ctx context.Context, action actions.Action, order *models.Order) (*models.Order, models.IError) {

	middleware.SpanStart(ctx, "LocalStorageOrder:Write")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Write")

	if s == nil {
		return nil, models.NewError(http.StatusInternalServerError, "localstorage not initialized for write")
	}
	if order == nil {
		return nil, models.NewError(http.StatusInternalServerError, "invalid order object in write")
	}

	var os []*models.Order
	var ok bool
	if order.Task.ID != nil {
		os, ok = s[*order.Task.ID]
	}

	switch action {

	case actions.CREATE:
		middleware.SpanStart(ctx, "LocalStorageOrder:Write:CREATE")
		defer middleware.SpanStop(ctx, "LocalStorageOrder:Write:CREATE")
		if ok || len(os) > 0 {
			return nil, models.NewError(http.StatusConflict, "already exists during write")
		}
		id := uint(len(s) + 1)
		order.Task.ID = &id
		now := time.Now().UTC()
		order.Task.Meta = &models.Meta{
			Version: uint(1),
			Created: now,
			Updated: now,
		}
		for i, subtask := range order.SubTasks {
			subtask.ID = utils.Ptr(uint(i + 1))
		}
		s[id] = append(s[id], order)
		return order, nil

	case actions.UPDATE:
		middleware.SpanStart(ctx, "LocalStorageOrder:Write:UPDATE")
		defer middleware.SpanStop(ctx, "LocalStorageOrder:Write:UPDATE")
		if !ok || len(os) == 0 {
			return nil, models.NewError(http.StatusNotFound, "not found during write")
		}
		var updOrder models.Order = utils.Deref(os[len(os)-1].Clone())
		updated := false

		updateTaskFunc := func(updated bool, task *models.Task, updTask *models.Task) bool {
			upd := false
			if task.State != nil && utils.Deref(updTask.State) != utils.Deref(task.State) {
				updTask.State = task.State
				upd = true
			}
			if task.Accountable != nil && utils.Deref(updTask.Accountable) != utils.Deref(task.Accountable) {
				updTask.Accountable = task.Accountable
				upd = true
			}
			if task.Objective != nil && utils.Deref(updTask.Objective) != utils.Deref(task.Objective) {
				updTask.Objective = task.Objective
				upd = true
			}
			return updated || upd
		}

		updated = updateTaskFunc(updated, order.Task, updOrder.Task)

		if order.ParentOrderID != nil && utils.Deref(updOrder.ParentOrderID) != utils.Deref(order.ParentOrderID) {
			updOrder.ParentOrderID = order.ParentOrderID
			updated = true
		}
		if order.Deadline != nil && utils.Deref(updOrder.Deadline) != utils.Deref(order.Deadline) {
			updOrder.Deadline = order.Deadline
			updated = true
		}

		for _, subtask := range order.SubTasks {
			if subtask.ID != nil { // NOTE: existing subtask
				upd := false
				for j, updSubTask := range updOrder.SubTasks {
					if utils.Deref(subtask.ID) == utils.Deref(updSubTask.ID) {
						upd = updateTaskFunc(upd, subtask, updSubTask)
						if upd {
							updatedSubtask := utils.Deref(updSubTask)
							updatedSubtask.Meta = &models.Meta{
								Version: updSubTask.Meta.Version + 1,
								Created: updSubTask.Meta.Created,
								Updated: time.Now().UTC(),
							}
							updOrder.SubTasks[j] = &updatedSubtask
						}
						break
					}
				}

				updated = updated || upd
			} else { // NOTE: new subtask
				updated = true
				now := time.Now().UTC()
				updOrder.SubTasks = append(updOrder.SubTasks, &models.Task{
					ID:          utils.Ptr(uint(1 + len(updOrder.SubTasks))),
					State:       subtask.State,
					Accountable: subtask.Accountable,
					Objective:   subtask.Objective,
					Meta: &models.Meta{
						Version: uint(1),
						Created: now,
						Updated: now,
					},
				})
			}
		}

		if updated {
			now := time.Now().UTC()
			updOrder.Task.Meta = &models.Meta{
				Version: updOrder.Task.Meta.Version + 1,
				Created: updOrder.Task.Meta.Created,
				Updated: now,
			}
			s[*order.Task.ID] = append(s[*order.Task.ID], &updOrder)
		}
		os = s[*order.Task.ID]
		return os[len(os)-1], nil

	case actions.DELETESOFT:
		middleware.SpanStart(ctx, "LocalStorageOrder:Write:SOFTDELETE")
		defer middleware.SpanStop(ctx, "LocalStorageOrder:Write:SOFTDELETE")
		if ok {
			for _, u := range os {
				u.Task.Meta.Deleted = true
			}
		}
		return nil, nil

	case actions.DELETEHARD:
		middleware.SpanStart(ctx, "LocalStorageOrder:Write:HARDDELETE")
		defer middleware.SpanStop(ctx, "LocalStorageOrder:Write:HARDDELETE")
		if ok {
			delete(s, *order.Task.ID)
		}
		return nil, nil

	default:
		return nil, models.NewError(http.StatusInternalServerError, "undefined write action")
	}
}
