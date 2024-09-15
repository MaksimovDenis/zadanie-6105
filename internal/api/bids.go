package api

import (
	"fmt"
	"net/http"
	"time"

	"git.codenrock.com/zadanie-6105/internal/storage/queries"
	"git.codenrock.com/zadanie-6105/pkg/helper"
	"git.codenrock.com/zadanie-6105/pkg/protocol/oapi"

	"github.com/gin-gonic/gin"
	"github.com/kak-tus/nan"
)

func (api *API) CreateBid(ctx *gin.Context) {
	var req oapi.CreateBidJSONBody

	if err := ctx.BindJSON(&req); err != nil {
		api.l.Error().Err(err).Msg("failed to unmarshall create bids body")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	params := queries.CreateBidParams{
		Name:           req.Name,
		Description:    nan.String(req.Description),
		Status:         queries.NullStatusType{StatusType: queries.StatusType(req.Status), Valid: true},
		TenderID:       nan.Int32(req.TenderId),
		OrganizationID: nan.Int32(req.OrganizationId),
		Username:       req.CreatorUsername,
	}

	tx, err := api.storage.BeginTx(ctx.Request.Context())
	if err != nil {
		api.l.Error().Err(err).Msg("failed to begin transaction")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Ошибка транзакции.")

		return
	}

	queriesWithTx := api.storage.Queries.WithTx(tx)

	bid, err := queriesWithTx.CreateBid(ctx.Request.Context(), params)
	if err != nil {
		_ = tx.Rollback(ctx)

		api.l.Error().Err(err).Msg("failed to add new bid to storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	err = queriesWithTx.CreateBidsHistory(ctx.Request.Context(), bid.ID)
	if err != nil {
		_ = tx.Rollback(ctx)

		api.l.Error().Err(err).Msg("failed to add new bid to storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Ошибка транзакции.")

		return
	}

	if err := tx.Commit(ctx); err != nil {
		api.l.Error().Err(err).Msg("failed to commit transaction")
		helper.CustomErrorResponse(ctx, http.StatusInternalServerError, "Ошибка транзакции.")

		return
	}

	res := oapi.Bid{
		Id:          bid.ID,
		Name:        bid.Name,
		Description: bid.Description.String,
		Status:      oapi.BidStatus(bid.Status.StatusType),
		AuthorType:  oapi.BidAuthorType(bid.AuthorType.String),
		AuthorId:    bid.UserID.Int32,
		TenderId:    bid.TenderID.Int32,
		Version:     bid.Version.Int32,
		CreatedAt:   bid.CreatedAt.Time.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) GetUserBids(ctx *gin.Context, params oapi.GetUserBidsParams) {
	limit := int32(5)
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}

	offset := int32(0)
	if params.Offset != nil && *params.Offset > 0 {
		offset = *params.Offset
	}

	username := ""
	if params.Username != nil {
		username = *params.Username
	}

	args := queries.GetUserBidsParams{
		Username: username,
		Limit:    limit,
		Offset:   offset,
	}

	bids, err := api.storage.Queries.GetUserBids(ctx.Request.Context(), args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to get bids from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	res := make([]oapi.Bid, len(bids))

	for idx, value := range bids {
		res[idx] = oapi.Bid{
			Id:          value.ID,
			Name:        value.Name,
			Description: value.Description.String,
			Status:      oapi.BidStatus(value.Status.StatusType),
			AuthorType:  oapi.BidAuthorType(value.AuthorType.String),
			AuthorId:    value.UserID.Int32,
			TenderId:    value.TenderID.Int32,
			Version:     value.Version.Int32,
			CreatedAt:   value.CreatedAt.Time.Format(time.RFC3339),
		}
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) GetBidsForTender(ctx *gin.Context, tenderId int32, params oapi.GetBidsForTenderParams) {
	if tenderId <= 0 {
		api.l.Error().Msg("invalid tender id arguments")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	limit := int32(5)
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}

	offset := int32(0)
	if params.Offset != nil && *params.Offset > 0 {
		offset = *params.Offset
	}

	args := queries.GetBidsForTenderParams{
		TenderID: nan.Int32(tenderId),
		Limit:    limit,
		Offset:   offset,
		Username: params.Username,
	}

	tenders, err := api.storage.Queries.GetBidsForTender(ctx.Request.Context(), args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to get tenders from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	res := make([]oapi.Bid, len(tenders))

	for idx, value := range tenders {
		res[idx] = oapi.Bid{
			Id:          value.ID,
			Name:        value.Name,
			Description: value.Description.String,
			Status:      oapi.BidStatus(value.Status.StatusType),
			AuthorType:  oapi.BidAuthorType(value.AuthorType.String),
			AuthorId:    value.UserID.Int32,
			TenderId:    value.TenderID.Int32,
			Version:     value.Version.Int32,
			CreatedAt:   value.CreatedAt.Time.Format(time.RFC3339),
		}
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) GetBidStatus(ctx *gin.Context, bidId int32, params oapi.GetBidStatusParams) {
	if bidId <= 0 {
		api.l.Error().Msg("invalid bid id arguments")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	args := queries.GetBidStatusParams{
		ID:       bidId,
		Username: params.Username,
	}

	status, err := api.storage.Queries.GetBidStatus(ctx.Request.Context(), args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to get bid status from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	res := oapi.BidStatus(status.StatusType)

	ctx.JSON(http.StatusOK, res)
}

func (api *API) UpdateBidStatus(ctx *gin.Context, bidId int32, params oapi.UpdateBidStatusParams) {
	if bidId <= 0 {
		api.l.Error().Msg("invalid bid id argument")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	args := queries.UpdateBidStatusParams{
		ID:       bidId,
		Username: params.Username,
		Status:   queries.NullStatusType{StatusType: queries.StatusType(params.Status), Valid: true},
	}

	updatedStatus, err := api.storage.Queries.UpdateBidStatus(ctx.Request.Context(), args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to update status in storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	res := oapi.Bid{
		Id:          updatedStatus.ID,
		Name:        updatedStatus.Name,
		Description: updatedStatus.Description.String,
		Status:      oapi.BidStatus(updatedStatus.Status.StatusType),
		AuthorType:  oapi.BidAuthorType(updatedStatus.AuthorType.String),
		AuthorId:    updatedStatus.UserID.Int32,
		TenderId:    updatedStatus.TenderID.Int32,
		Version:     updatedStatus.Version.Int32,
		CreatedAt:   updatedStatus.CreatedAt.Time.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) EditBid(ctx *gin.Context, bidId int32, params oapi.EditBidParams) {
	if bidId <= 0 {
		api.l.Error().Msg("invalid bid id argument")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest,
			"Данные неправильно сформированы или не соответствуют требованиям.")

		return
	}

	var req oapi.EditBidJSONBody

	if err := ctx.BindJSON(&req); err != nil {
		api.l.Error().Err(err).Msg("failed to unmarshall edit bid body")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest,
			"Данные неправильно сформированы или не соответствуют требованиям.")

		return
	}

	var description nan.NullString
	if req.Description != nil {
		description = nan.NullString{String: *req.Description, Valid: false}
	}

	args := queries.EditBidParams{
		ID:          bidId,
		Name:        *req.Name,
		Description: description,
		Username:    params.Username,
	}

	tx, err := api.storage.BeginTx(ctx.Request.Context())
	if err != nil {
		api.l.Error().Err(err).Msg("failed to begin transaction")
		helper.CustomErrorResponse(ctx, http.StatusInternalServerError, "Ошибка транзакции.")

		return
	}

	queriesWithTx := api.storage.Queries.WithTx(tx)

	updatedBid, err := queriesWithTx.EditBid(ctx.Request.Context(), args)
	if err != nil {
		_ = tx.Rollback(ctx)

		api.l.Error().Err(err).Msg("failed to add new bid to storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest,
			"Данные неправильно сформированы или не соответствуют требованиям.")

		return
	}

	err = queriesWithTx.CreateTenderHistory(ctx.Request.Context(), int32(updatedBid.ID))
	if err != nil {
		_ = tx.Rollback(ctx)

		api.l.Error().Err(err).Msg("failed to add new bid in history table to storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest,
			"Данные неправильно сформированы или не соответствуют требованиям.")

		return
	}

	if err := tx.Commit(ctx); err != nil {
		api.l.Error().Err(err).Msg("failed to commit transaction")
		helper.CustomErrorResponse(ctx, http.StatusInternalServerError, "Ошибка транзакции.")

		return
	}

	res := oapi.Bid{
		Id:          updatedBid.ID,
		Name:        updatedBid.Name,
		Description: updatedBid.Description.String,
		Status:      oapi.BidStatus(updatedBid.Status.StatusType),
		AuthorType:  oapi.BidAuthorType(updatedBid.AuthorType.String),
		AuthorId:    updatedBid.UserID.Int32,
		TenderId:    updatedBid.TenderID.Int32,
		Version:     updatedBid.Version.Int32,
		CreatedAt:   updatedBid.CreatedAt.Time.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) RollbackBid(ctx *gin.Context, bidId int32, version int32, params oapi.RollbackBidParams) {
	if bidId <= 0 {
		api.l.Error().Msg("invalid bid id argument")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	args := queries.RollbackBidParams{
		ID:       bidId,
		Version:  nan.Int32(version),
		Username: params.Username,
	}

	bid, err := api.storage.Queries.RollbackBid(ctx.Request.Context(), args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to rollback bid from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	res := oapi.Bid{
		Id:          bid.ID,
		Name:        bid.Name,
		Description: bid.Description.String,
		Status:      oapi.BidStatus(bid.Status.StatusType),
		AuthorType:  oapi.BidAuthorType(bid.AuthorType.String),
		AuthorId:    bid.UserID.Int32,
		TenderId:    bid.TenderID.Int32,
		Version:     bid.Version.Int32,
		CreatedAt:   bid.CreatedAt.Time.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) SubmitBidFeedback(ctx *gin.Context, bidId int32, params oapi.SubmitBidFeedbackParams) {
	if bidId <= 0 {
		api.l.Error().Msg("invalid bid id argument")
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid bid id argument"})

		return
	}

	args := queries.SubmitBidFeedbackParams{
		ID:       bidId,
		Username: params.Username,
		Feedback: params.BidFeedback,
	}

	bid, err := api.storage.Queries.SubmitBidFeedback(ctx, args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to rollback bid from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Отзыв не может быть отправлен.")

		return
	}

	res := oapi.Bid{
		Id:          bid.ID,
		Name:        bid.Name,
		Description: bid.Description.String,
		Status:      oapi.BidStatus(bid.Status.StatusType),
		AuthorType:  oapi.BidAuthorType(bid.AuthorType.String),
		AuthorId:    bid.UserID.Int32,
		TenderId:    bid.TenderID.Int32,
		Version:     bid.Version.Int32,
		CreatedAt:   bid.CreatedAt.Time.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) GetBidReviews(ctx *gin.Context, tenderId int32, params oapi.GetBidReviewsParams) {
	if tenderId <= 0 {
		api.l.Error().Msg("invalid tender id arguments")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	limit := int32(5)
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}

	offset := int32(0)
	if params.Offset != nil && *params.Offset > 0 {
		offset = *params.Offset
	}

	args := queries.GetBidReviewsParams{
		TenderID:   nan.Int32(tenderId),
		Username:   params.AuthorUsername,
		Username_2: params.RequesterUsername,
		Limit:      limit,
		Offset:     offset,
	}

	reviews, err := api.storage.Queries.GetBidReviews(ctx, args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to rollback bid from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	fmt.Println(reviews)

	res := make([]oapi.BidReview, len(reviews))

	for idx, value := range reviews {
		res[idx] = oapi.BidReview{
			Id:          value.ID,
			Description: value.Feedback,
			CreatedAt:   value.CreatedAt.Time.Format(time.RFC3339),
		}
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) SubmitBidDecision(ctx *gin.Context, bidId int32, params oapi.SubmitBidDecisionParams) {
	if bidId <= 0 {
		api.l.Error().Msg("invalid tender id arguments")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Решение не может быть отправлено.")

		return
	}

	args := queries.SubmitBidDecisionParams{
		ID:       bidId,
		Column2:  string(params.Decision),
		Username: params.Username,
	}

	bid, err := api.storage.Queries.SubmitBidDecision(ctx, args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to rollback bid from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Решение не может быть отправлено.")

		return
	}

	res := oapi.Bid{
		Id:          bid.ID,
		Name:        bid.Name,
		Description: bid.Description.String,
		Status:      oapi.BidStatus(bid.Status.StatusType),
		AuthorType:  oapi.BidAuthorType(bid.AuthorType.String),
		AuthorId:    bid.UserID.Int32,
		TenderId:    bid.TenderID.Int32,
		Version:     bid.Version.Int32,
		CreatedAt:   bid.CreatedAt.Time.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, res)
}
