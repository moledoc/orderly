package mgmtorder

import (
	"net/http"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/service/common/validation"
)

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

	if !validation.IsIgnoreEmpty(sitrep.GetSituation()+
		sitrep.GetActions()+
		sitrep.GetTODO()+
		sitrep.GetIssues(), ignore) &&
		//
		len(sitrep.GetSituation()) == 0 &&
		len(sitrep.GetActions()) == 0 &&
		len(sitrep.GetTODO()) == 0 &&
		len(sitrep.GetIssues()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty sitrep")
	}

	return nil
}

func ValidateOrder(o *order.Order, ignore validation.IgnoreField) errwrap.Error {
	if o == nil {
		return nil
	}

	err := validation.ValidateID(o.GetID())
	if !validation.IsFieldIgnored(validation.IgnoreID, ignore) && err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.id: %s", err.GetStatusMessage())
	}

	if !validation.IsIgnoreEmpty(o.GetState(), ignore) && o.GetState() < 0 || order.Done < o.GetState() {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.state")
	}

	err = validation.ValidateID(o.GetAccountableID())
	if !validation.IsIgnoreEmpty(o.GetAccountableID(), ignore) && err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.accountable: %s", err.GetStatusMessage())
	}

	if !validation.IsIgnoreEmpty(o.GetObjective(), ignore) && len(o.GetObjective()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.objective")
	}

	if !validation.IsIgnoreEmpty(o.GetDeadline(), ignore) && o.GetDeadline().Equal(time.Time{}) {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.deadline")
	}

	err = validation.ValidateID(o.GetParentOrderID())
	if !validation.IsFieldIgnored(validation.IgnoreID, ignore) && err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.parent_order_id: %s", err.GetStatusMessage())
	}

	for i, delegatedOrder := range o.GetDelegatedOrders() {
		err := ValidateOrder(delegatedOrder, ignore)
		if err != nil {
			return errwrap.NewError(http.StatusBadRequest, "invalid order.delegated_orders.%v: %s", i, err.GetStatusMessage())
		}
	}

	for i, sitrep := range o.GetSitReps() {
		err := ValidateSitRep(sitrep, ignore)
		if err != nil {
			return errwrap.NewError(http.StatusBadRequest, "invalid order.sitreps.%v: %s", i, err.GetStatusMessage())
		}
	}

	return nil
}

func ValidatePostOrderRequest(req *request.PostOrderRequest) errwrap.Error {

	if req == nil || req.Order == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetOrder().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "order.id disallowed")
	}

	for i, delegatedOrder := range req.GetOrder().GetDelegatedOrders() {
		if len(delegatedOrder.GetID()) > 0 {
			return errwrap.NewError(http.StatusBadRequest, "order.delegated_orders.%v.id disallowed", i)
		}
	}

	for i, sitrep := range req.GetOrder().GetSitReps() {
		if len(sitrep.GetID()) > 0 {
			return errwrap.NewError(http.StatusBadRequest, "order.sitreps.%v.id disallowed", i)
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

	if len(req.GetParentOrderID()) > 0 {
		err := validation.ValidateID(req.GetParentOrderID())
		if err != nil {
			return errwrap.NewError(http.StatusBadRequest, "invalid parent_order_id: %s", err.GetStatusMessage())
		}
	}

	if len(req.GetAccountableID()) > 0 {
		err := validation.ValidateID(req.GetAccountableID())
		if err != nil {
			return errwrap.NewError(http.StatusBadRequest, "invalid accountable: %s", err.GetStatusMessage())
		}
	}

	return nil
}

func ValidatePatchOrderRequest(req *request.PatchOrderRequest) errwrap.Error {
	if req == nil || req.Order == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.Order.GetID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order.task.id: %s", err.GetStatusMessage())
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

func ValidatePutDelegatedOrderRequest(req *request.PutDelegatedOrdersRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	if len(req.GetOrders()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty tasks")
	}

	for i, tsk := range req.GetOrders() {
		if len(tsk.GetID()) > 0 {
			return errwrap.NewError(http.StatusBadRequest, "tasks.%v.id disallowed", i)
		}

		err := ValidateOrder(tsk, validation.IgnoreID)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidatePatchDelegatedOrderRequest(req *request.PatchDelegatedOrdersRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	if len(req.GetOrders()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty tasks")
	}

	for i, tsk := range req.GetOrders() {
		err := validation.ValidateID(tsk.GetID())
		if err != nil {
			return errwrap.NewError(http.StatusBadRequest, "invalid tasks.%v.id: %s", i, err.GetStatusMessage())
		}

		err = ValidateOrder(tsk, validation.IgnoreEmpty)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateDeleteDelegatedOrderRequest(req *request.DeleteDelegatedOrdersRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	if len(req.GetDelegatedOrderIDs()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty delegated_task_ids")
	}

	for i, id := range req.GetDelegatedOrderIDs() {
		err = validation.ValidateID(id)
		if err != nil {
			return errwrap.NewError(http.StatusBadRequest, "invalid delegated_task_ids.%v: %s", i, err.GetStatusMessage())
		}
	}

	return nil
}

////////

func ValidatePutSitRepRequest(req *request.PutSitRepsRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	if len(req.GetSitReps()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty sitreps")
	}

	for i, sitrep := range req.GetSitReps() {
		if len(sitrep.GetID()) > 0 {
			return errwrap.NewError(http.StatusBadRequest, "sitreps.%v.id disallowed", i)
		}

		err := ValidateSitRep(sitrep, validation.IgnoreID)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidatePatchSitRepRequest(req *request.PatchSitRepsRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	if len(req.GetSitReps()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty sitreps")
	}

	for i, sitrep := range req.GetSitReps() {
		err := validation.ValidateID(sitrep.GetID())
		if err != nil {
			return errwrap.NewError(http.StatusBadRequest, "invalid sitreps.%v.id: %s", i, err.GetStatusMessage())
		}

		err = ValidateSitRep(sitrep, validation.IgnoreEmpty)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateDeleteSitRepRequest(req *request.DeleteSitRepsRequest) errwrap.Error {

	err := validation.ValidateID(req.GetOrderID())
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid order_id: %s", err.GetStatusMessage())
	}

	if len(req.GetSitRepIDs()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty sitrep_ids")
	}

	for i, sitrep := range req.GetSitRepIDs() {
		err = validation.ValidateID(sitrep)
		if err != nil {
			return errwrap.NewError(http.StatusBadRequest, "invalid sitrep_ids.%v: %s", i, err.GetStatusMessage())
		}
	}

	return nil
}
