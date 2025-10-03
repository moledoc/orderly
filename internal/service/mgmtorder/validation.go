package mgmtorder

import (
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/service/common/validation"
)

func validateTask(task *order.Task) errwrap.Error {
	if task == nil {
		return nil
	}

	if task.ID != nil { // NOTE: ID is required, but when creating we don't allow ID; relevant ID check is done one level up in validation
		err := validation.ValidateID(task.GetID())
		if err != nil {
			return err
		}
	}

	if task.GetState() < 0 || order.Completed < task.GetState() {
		return errwrap.NewError(http.StatusBadRequest, "invalid task.state")
	}

	err := validation.ValidateEmail(task.GetAccountable())
	if err != nil {
		return err
	}

	if len(task.GetObjective()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "invalid task.objective")
	}

	err = validation.ValidateMeta(task.GetMeta())
	if err != nil {
		return err
	}

	return nil
}

func validateSitRep(sitrep *order.SitRep) errwrap.Error {

	if sitrep == nil {
		return nil
	}

	if sitrep.ID != nil { // NOTE: ID is required, but when creating we don't allow ID; relevant ID check is done one level up in validation
		err := validation.ValidateID(sitrep.GetID())
		if err != nil {
			return err
		}
	}

	if sitrep.GetState() < 0 || order.Completed < sitrep.GetState() {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.state")
	}

	if /* sitrep.GetWorkCompleted() < 0 || */ 100 < sitrep.GetWorkCompleted() {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.WorkCompleted")
	}

	if len(sitrep.GetSummary()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "invalid sitrep.Summary")
	}

	err := validation.ValidateMeta(sitrep.GetMeta())
	if err != nil {
		return err
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

	if order.Deadline == nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.deadline")
	}

	for _, subtask := range order.GetSubTasks() {
		err := validateTask(subtask)
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

	return nil
}

func validatePostOrderRequest(req *request.PostOrderRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}
	if req.Order == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil order")
	}

	if req.GetOrder().Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil task")
	}

	if req.GetOrder().GetTask().ID != nil {
		return errwrap.NewError(http.StatusBadRequest, "order.task.id disallowed")
	}

	return validateOrder(req.GetOrder())
}

func validateGetOrderByIDRequest(req *request.GetOrderByIDRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}

	if req.ID == nil {
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

func validateGetOrderVersionsRequest(req *request.GetOrderVersionsRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}

	if req.ID == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}
	return nil
}

func validateGetOrderSubOrdersRequest(req *request.GetOrderSubOrdersRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}

	if req.ID == nil {
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
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}
	if req.Order == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil order")
	}

	if req.GetOrder().Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil task")
	}

	if req.GetOrder().GetTask().ID == nil {
		return errwrap.NewError(http.StatusBadRequest, "order.task.id required")
	}

	return validateOrder(req.GetOrder())
}

func validateDeleteOrderRequest(req *request.DeleteOrderRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}

	if req.ID == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}
	return nil
}

////////

func validatePutSubTaskRequest(req *request.PutSubTaskRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}

	if req.OrderID == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	if req.Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil task")
	}

	if req.GetTask().ID != nil {
		return errwrap.NewError(http.StatusBadRequest, "task.id disallowed")
	}

	return validateTask(req.GetTask())
}

func validatePatchSubTaskRequest(req *request.PatchSubTaskRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}

	if req.OrderID == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	if req.Task == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil task")
	}

	if req.GetTask().ID == nil {
		return errwrap.NewError(http.StatusBadRequest, "task.id required")
	}

	return validateTask(req.GetTask())
}

func validateDeleteSubTaskRequest(req *request.DeleteSubTaskRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}

	if req.OrderID == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	err = validation.ValidateID(req.GetSubTaskID())
	if err != nil {
		return err
	}

	return nil
}

////////

func validatePutSitRepRequest(req *request.PutSitRepRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}

	if req.OrderID == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	if req.SitRep == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil sitrep")
	}

	if req.GetSitRep().ID != nil {
		return errwrap.NewError(http.StatusBadRequest, "sitrep.id disallowed")
	}

	return validateSitRep(req.GetSitRep())
}

func validatePatchSitRepRequest(req *request.PatchSitRepRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}

	if req.OrderID == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil id")
	}

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return err
	}

	if req.SitRep == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil sitrep")
	}

	if req.GetSitRep().ID == nil {
		return errwrap.NewError(http.StatusBadRequest, "sitrep.id required")
	}

	return validateSitRep(req.GetSitRep())
}

func validateDeleteSitRepRequest(req *request.DeleteSitRepRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "nil request")
	}

	if req.OrderID == nil {
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
