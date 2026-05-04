// Package subscription implements subscription management domain logic.
//
// It provides HTTP handlers, service layer, repository layer, and DTO/domain
// mapping for creating, updating, deleting, listing, and aggregating
// subscription data.
//
// The package follows a layered architecture:
//   - transport (HTTP handlers)
//   - service (business logic)
//   - repository (PostgreSQL persistence)
//   - mapping (DTO ↔ domain conversion)
package subscription

//	@Summary	Create subscription
//	@Tags		subscription
//	@ID			create-subscription
//	@Accept		json
//	@Produce	json
//	@Param		sub	body		SubReq			true	"User ID must be uuid\nDate format: MM-YYYY"
//	@Success	201	{object}	SubResp			"Created subscription"
//	@Failure	400	{object}	validation.Resp	"Bad request"
//	@Failure	404	{string}	string			"Not found"
//	@Failure	500	{string}	string			"Internal server error"
//	@Router		/subscriptions [post].

//	@Summary		Delete subscription
//	@Description	Delete by subscription ID.
//	@Tags			subscription
//	@ID				delete-subscription
//	@Param			id	path	string	true	"Subscription ID"
//	@Success		204	"no content"
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		404	{string}	string	"Not found"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/subscriptions/{id} [delete].

//	@Summary		Get subscription
//	@Description	Get by subscription ID.
//	@Tags			subscription
//	@ID				get-subscription
//	@Produce		json
//	@Param			id	path		string	true	"Subscription ID"
//	@Success		200	{object}	SubResp	"Subscription"
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		404	{string}	string	"Not found"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/subscriptions/{id} [get].

//	@Summary		Update subscription
//	@Description	Update by subscription ID.
//	@Tags			subscription
//	@ID				update-subscription
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Subscription ID"
//	@Param			sub	body		UpdateSubReq	true	"Subscription"
//	@Success		200	{object}	SubResp			"User ID must be uuid"
//	@Failure		400	{object}	validation.Resp	"Bad request"
//	@Failure		404	{string}	string			"Not found"
//	@Failure		500	{string}	string			"Internal server error"
//	@Router			/subscriptions/{id} [patch].

//	@Summary		List subscriptions
//	@Description	Get paginated list of subscriptions with optional filters
//	@Tags			subscription
//	@ID				subscription-list
//	@Produce		json
//	@Param			page			query		int			true	"Page number (1-based)"
//	@Param			limit			query		int			true	"Items per page (max: 100)"
//	@Param			service_name	query		string		false	"filter by service name"
//	@Param			user_id			query		string		false	"filter by user ID"
//	@Success		200				{object}	SubListResp	"List subscriptions"
//	@Failure		400				{string}	string		"Bad request"
//	@Failure		500				{string}	string		"Internal server error"
//	@Router			/subscriptions [get].

//	@Summary		Get subscription summa
//	@Description	Get total_price of all subscriptions.
//	@Tags			subscription
//	@ID				get-subscription-sum
//	@Produce		json
//	@Param			start_date		query		string			true	"Date format: MM-YYYY"
//	@Param			end_date		query		string			true	"Date format: MM-YYYY"
//	@Param			service_name	query		string			false	"filter by service name"
//	@Param			user_id			query		string			false	"filter by user ID"
//	@Success		200				{object}	SubSumResp		"Subscription sum"
//	@Failure		400				{object}	validation.Resp	"Bad request"
//	@Failure		500				{string}	string			"Internal server error"
//	@Router			/subscriptions/sum [get].
