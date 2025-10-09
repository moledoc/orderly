package mgmtorder

import (
	"net/http"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/service/common/validation"
)

func ValidateTask(task *order.Task, ignore validation.IgnoreField) errwrap.Error {
	if task == nil {
		return nil
	}

	err := validation.ValidateID(task.GetID())
	if !validation.IsFieldIgnored(validation.IgnoreID, ignore) && err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid task.id: %s", err.GetStatusMessage())
	}

	if !validation.IsIgnoreEmpty(task.GetState(), ignore) && task.GetState() < order.NotStarted || order.Completed < task.GetState() {
		return errwrap.NewError(http.StatusBadRequest, "invalid task.state")
	}

	err = validation.ValidateEmail(task.GetAccountable())
	if !validation.IsIgnoreEmpty(task.GetAccountable(), ignore) && err != nil {
		return err
	}

	if !validation.IsIgnoreEmpty(task.GetObjective(), ignore) && len(task.GetObjective()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "invalid task.objective")
	}

	if !validation.IsIgnoreEmpty(task.GetDeadline(), ignore) && task.GetDeadline().Equal(time.Time{}) {
		return errwrap.NewError(http.StatusBadRequest, "invalid task.deadline")
	}

	return nil
}

func ValidateSitRep(sitrep *order.SitRep, ignore validation.IgnoreField) errwrap.Error {

	if sitrep == nil {
		return nil
	}

	err := validation.ValidateID(sitrep.GetID())
	if !validation.IsFieldIgnored(validation.IgnoreID, ignore) && err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.id: %s", err.GetStatusMessage())
	}

	if !validation.IsIgnoreEmpty(sitrep.GetDateTime(), ignore) && sitrep.GetDateTime().Equal(time.Time{}) {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.datetime")
	}

	err = validation.ValidateEmail(sitrep.GetBy())
	if !validation.IsIgnoreEmpty(sitrep.GetBy(), ignore) && err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.by: %s", err.GetStatusMessage())
	}

	if len(sitrep.GetPing()) > 0 {
		for i, ping := range sitrep.GetPing() {
			if err := validation.ValidateEmail(ping); err != nil {
				return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.%v.ping: %s", i, err.GetStatusMessage())
			}
		}
	}

	if !validation.IsIgnoreEmpty(sitrep.GetSituation()+
		sitrep.GetActions()+
		sitrep.GetTBD()+
		sitrep.GetIssues(), ignore) &&
		//
		len(sitrep.GetSituation()) == 0 &&
		len(sitrep.GetActions()) == 0 &&
		len(sitrep.GetTBD()) == 0 &&
		len(sitrep.GetIssues()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty sitrep")
	}

	return nil
}

func ValidateOrder(order *order.Order, ignore validation.IgnoreField) errwrap.Error {
	if order == nil {
		return nil
	}

	if !validation.IsIgnoreEmpty(order.GetTask(), ignore) && order.GetTask() == nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.task")
	}

	err := ValidateTask(order.GetTask(), ignore)
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.task: %s", err.GetStatusMessage())
	}

	err = validation.ValidateID(order.GetParentOrderID())
	if !validation.IsIgnoreEmpty(order.GetParentOrderID(), ignore) && err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.parent_order_id: %s", err.GetStatusMessage())
	}

	for i, delegatedTask := range order.GetDelegatedTasks() {
		err := ValidateTask(delegatedTask, ignore)
		if err != nil {
			return errwrap.NewError(http.StatusBadRequest, "invalid order.delegated_task.%v: %s", i, err.GetStatusMessage())
		}
	}

	for i, sitrep := range order.GetSitReps() {
		err := ValidateSitRep(sitrep, ignore)
		if err != nil {
			return errwrap.NewError(http.StatusBadRequest, "invalid order.sitrep.%v: %s", i, err.GetStatusMessage())
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

	return ValidateOrder(req.GetOrder(), validation.IgnoreID)
}

func ValidateGetOrderByIDRequest(req *request.GetOrderByIDRequest) errwrap.Error {

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}
	return nil
}

func ValidateGetOrdersRequest(req *request.GetOrdersRequest) errwrap.Error {
	return nil
}

func ValidateGetOrderSubOrdersRequest(req *request.GetOrderSubOrdersRequest) errwrap.Error {

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}
	return nil
}

func ValidatePatchOrderRequest(req *request.PatchOrderRequest) errwrap.Error {
	if req == nil || req.Order == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetOrder().GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "order.task.id missing")
	}

	return ValidateOrder(req.GetOrder(), validation.IgnoreEmpty)
}

func ValidateDeleteOrderRequest(req *request.DeleteOrderRequest) errwrap.Error {

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}
	return nil
}

////////

func ValidatePutDelegatedTaskRequest(req *request.PutDelegatedTaskRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	if req.Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty task")
	}

	if len(req.GetTask().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "task.id disallowed")
	}

	return ValidateTask(req.GetTask(), validation.IgnoreID)
}

func ValidatePatchDelegatedTaskRequest(req *request.PatchDelegatedTaskRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	if req.Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty task")
	}

	if len(req.GetTask().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "task.id required")
	}

	return ValidateTask(req.GetTask(), validation.IgnoreNothing)
}

func ValidateDeleteDelegatedTaskRequest(req *request.DeleteDelegatedTaskRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	err = validation.ValidateID(req.GetDelegatedTaskID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid delegated_task_id: %s", err.GetStatusMessage())
	}

	return nil
}

////////

func ValidatePutSitRepRequest(req *request.PutSitRepRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	if req.SitRep == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty sitrep")
	}

	if len(req.GetSitRep().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "sitrep.id disallowed")
	}

	return ValidateSitRep(req.GetSitRep(), validation.IgnoreID)
}

func ValidatePatchSitRepRequest(req *request.PatchSitRepRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	if req.SitRep == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty sitrep")
	}

	if len(req.GetSitRep().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "sitrep.id required")
	}

	return ValidateSitRep(req.GetSitRep(), validation.IgnoreNothing)
}

func ValidateDeleteSitRepRequest(req *request.DeleteSitRepRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	err = validation.ValidateID(req.GetSitRepID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep_id: %s", err.GetStatusMessage())
	}

	return nil
}
