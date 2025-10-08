package mgmtorder

import (
	"net/http"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/service/common/validation"
)

func ValidateTask(task *order.Task) errwrap.Error {
	if task == nil {
		return nil
	}

	if len(task.ID) > 0 { // NOTE: ID is required, but when creating we don't allow ID; relevant ID check is done one level up in validation
		err := validation.ValidateID(task.GetID())
		if err != nil {
			return err
		}
	}

	if task.GetState() < order.NotStarted || order.Completed < task.GetState() {
		return errwrap.NewError(http.StatusBadRequest, "invalid task.state")
	}

	err := validation.ValidateEmail(task.GetAccountable())
	if err != nil {
		return err
	}

	if len(task.GetObjective()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "invalid task.objective")
	}

	if task.GetDeadline().Equal(time.Time{}) {
		return errwrap.NewError(http.StatusBadRequest, "invalid task.deadline")
	}

	return nil
}

func ValidateSitRep(sitrep *order.SitRep) errwrap.Error {

	if sitrep == nil {
		return nil
	}

	if len(sitrep.ID) > 0 { // NOTE: ID is required, but when creating we don't allow ID; relevant ID check is done one level up in validation
		err := validation.ValidateID(sitrep.GetID())
		if err != nil {
			return err
		}
	}

	if sitrep.GetDateTime().Equal(time.Time{}) {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.datetime")
	}

	if err := validation.ValidateEmail(sitrep.GetBy()); err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.by: %s", err.GetStatusMessage())
	}

	if len(sitrep.GetPing()) > 0 {
		for i, ping := range sitrep.GetPing() {
			if err := validation.ValidateEmail(ping); err != nil {
				return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.%v.ping: %s", i, err.GetStatusMessage())
			}
		}
	}

	if len(sitrep.GetSituation()) == 0 && len(sitrep.GetActions()) == 0 && len(sitrep.GetTBD()) == 0 && len(sitrep.GetIssues()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty sitrep")
	}

	return nil
}

func ValidateOrder(order *order.Order) errwrap.Error {
	if order == nil {
		return nil
	}

	if order.GetTask() == nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.task")
	}

	err := ValidateTask(order.GetTask())
	if err != nil {
		return err
	}

	err = validation.ValidateID(order.GetParentOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.parent_order_id: %s", err.GetStatusMessage())
	}

	for _, delegatedTask := range order.GetDelegatedTasks() {
		err := ValidateTask(delegatedTask)
		if err != nil {
			return err
		}
	}

	for _, sitrep := range order.GetSitReps() {
		err := ValidateSitRep(sitrep)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidatePostOrderRequest(req *request.PostOrderRequest) errwrap.Error {

	if req == nil || req.Order == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if req.GetOrder().Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty task")
	}

	if len(req.GetOrder().GetTask().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "order.task.id disallowed")
	}

	for i, delegatedTask := range req.GetOrder().GetDelegatedTasks() {
		if len(delegatedTask.GetID()) > 0 {
			return errwrap.NewError(http.StatusBadRequest, "order.delegated.%v.id disallowed", i)
		}
	}

	for i, sitrep := range req.GetOrder().GetSitReps() {
		if len(sitrep.GetID()) > 0 {
			return errwrap.NewError(http.StatusBadRequest, "order.sitrep.%v.id disallowed", i)
		}
	}

	return ValidateOrder(req.GetOrder())
}

func ValidateGetOrderByIDRequest(req *request.GetOrderByIDRequest) errwrap.Error {
	if req == nil || len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}
	return nil
}

func ValidateGetOrdersRequest(req *request.GetOrdersRequest) errwrap.Error {
	return nil
}

func ValidateGetOrderSubOrdersRequest(req *request.GetOrderSubOrdersRequest) errwrap.Error {
	if req == nil || len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}
	return nil
}

func ValidatePatchOrderRequest(req *request.PatchOrderRequest) errwrap.Error {
	if req == nil || req.Order == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if req.GetOrder().Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty task")
	}

	if len(req.GetOrder().GetTask().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "order.task.id required")
	}

	return ValidateOrder(req.GetOrder())
}

func ValidateDeleteOrderRequest(req *request.DeleteOrderRequest) errwrap.Error {
	if req == nil || len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}
	return nil
}

////////

func ValidatePutDelegatedTaskRequest(req *request.PutDelegatedTaskRequest) errwrap.Error {
	if req == nil || len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	if req.Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty task")
	}

	if len(req.GetTask().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "task.id disallowed")
	}

	return ValidateTask(req.GetTask())
}

func ValidatePatchDelegatedTaskRequest(req *request.PatchDelegatedTaskRequest) errwrap.Error {
	if req == nil || len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	if req.Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty task")
	}

	if len(req.GetTask().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "task.id required")
	}

	return ValidateTask(req.GetTask())
}

func ValidateDeleteDelegatedTaskRequest(req *request.DeleteDelegatedTaskRequest) errwrap.Error {
	if req == nil || len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	err = validation.ValidateID(req.GetDelegatedTaskID())
	if err != nil {
		return err
	}

	return nil
}

////////

func ValidatePutSitRepRequest(req *request.PutSitRepRequest) errwrap.Error {
	if req == nil || len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	if req.SitRep == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil sitrep")
	}

	if len(req.GetSitRep().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "sitrep.id disallowed")
	}

	return ValidateSitRep(req.GetSitRep())
}

func ValidatePatchSitRepRequest(req *request.PatchSitRepRequest) errwrap.Error {
	if req == nil || len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	if req.SitRep == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil sitrep")
	}

	if len(req.GetSitRep().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "sitrep.id required")
	}

	return ValidateSitRep(req.GetSitRep())
}

func ValidateDeleteSitRepRequest(req *request.DeleteSitRepRequest) errwrap.Error {
	if req == nil || len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	err = validation.ValidateID(req.GetSitRepID())
	if err != nil {
		return err
	}

	return nil
}
