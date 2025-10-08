package mgmtorder

import (
	"net/http"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/service/common/validation"
)

func validateTask(task *order.Task) errwrap.Error {
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

	if task.Deadline.Equal(time.Time{}) {
		return errwrap.NewError(http.StatusBadRequest, "invalid task.deadline")
	}

	return nil
}

func validateSitRep(sitrep *order.SitRep) errwrap.Error {

	if sitrep == nil {
		return nil
	}

	if len(sitrep.ID) > 0 { // NOTE: ID is required, but when creating we don't allow ID; relevant ID check is done one level up in validation
		err := validation.ValidateID(sitrep.GetID())
		if err != nil {
			return err
		}
	}

	if sitrep.GetState() < order.NotStarted || order.Completed < sitrep.GetState() {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.state")
	}

	if 100 < sitrep.GetWorkCompleted() {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.workCompleted")
	}

	if len(sitrep.GetSummary()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.summary")
	}

	return nil
}

func validateOrder(order *order.Order) errwrap.Error {
	if order == nil {
		return nil
	}

	err := validateTask(order.GetTask())
	if err != nil {
		return err
	}

	err = validation.ValidateID(order.GetParentOrderID())
	if err != nil {
		return err
	}

	for _, delegatedTask := range order.GetDelegatedTasks() {
		err := validateTask(delegatedTask)
		if err != nil {
			return err
		}
	}

	for _, sitrep := range order.GetSitReps() {
		err := validateSitRep(sitrep)
		if err != nil {
			return err
		}
	}

	err = validation.ValidateMeta(order.GetMeta())
	if err != nil {
		return err
	}

	return nil
}

func validatePostOrderRequest(req *request.PostOrderRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}
	if req.Order == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty order")
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

	if req.GetOrder().Meta != nil {
		return errwrap.NewError(http.StatusBadRequest, "order.meta disallowed")
	}

	return validateOrder(req.GetOrder())
}

func validateGetOrderByIDRequest(req *request.GetOrderByIDRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}
	return nil
}

func validateGetOrdersRequest(req *request.GetOrdersRequest) errwrap.Error {
	return nil
}

func validateGetOrderSubOrdersRequest(req *request.GetOrderSubOrdersRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}
	return nil
}

func validatePatchOrderRequest(req *request.PatchOrderRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}
	if req.Order == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty order")
	}

	if req.GetOrder().Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty task")
	}

	if len(req.GetOrder().GetTask().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "order.task.id required")
	}

	return validateOrder(req.GetOrder())
}

func validateDeleteOrderRequest(req *request.DeleteOrderRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}
	return nil
}

////////

func validatePutDelegatedTaskRequest(req *request.PutDelegatedTaskRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
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

	return validateTask(req.GetTask())
}

func validatePatchDelegatedTaskRequest(req *request.PatchDelegatedTaskRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
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

	return validateTask(req.GetTask())
}

func validateDeleteDelegatedTaskRequest(req *request.DeleteDelegatedTaskRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
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

func validatePutSitRepRequest(req *request.PutSitRepRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
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

	return validateSitRep(req.GetSitRep())
}

func validatePatchSitRepRequest(req *request.PatchSitRepRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
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

	return validateSitRep(req.GetSitRep())
}

func validateDeleteSitRepRequest(req *request.DeleteSitRepRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetOrderID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
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
