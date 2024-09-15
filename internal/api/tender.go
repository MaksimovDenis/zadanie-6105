package api

import (
	"net/http"
	"time"

	"git.codenrock.com/zadanie-6105/internal/storage/queries"
	"git.codenrock.com/zadanie-6105/pkg/helper"
	"git.codenrock.com/zadanie-6105/pkg/protocol/oapi"

	"github.com/gin-gonic/gin"
	"github.com/kak-tus/nan"
)

func (api *API) CreateTender(ctx *gin.Context) {
	var req oapi.CreateTenderJSONBody

	if err := ctx.BindJSON(&req); err != nil {
		api.l.Error().Err(err).Msg("failed to unmarshall create tender body")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	params := queries.CreateTenderParams{
		Name:           req.Name,
		Description:    nan.String(req.Description),
		ServiceType:    nan.String(string(req.ServiceType)),
		Status:         queries.NullStatusType{StatusType: queries.StatusType(req.Status), Valid: true},
		OrganizationID: nan.Int32(req.OrganizationId),
		Username:       req.CreatorUsername,
	}

	tx, err := api.storage.BeginTx(ctx.Request.Context())
	if err != nil {
		api.l.Error().Err(err).Msg("failed to begin transaction")
		helper.CustomErrorResponse(ctx, http.StatusInternalServerError, "Ошибка транзакции.")

		return
	}

	queriesWithTx := api.storage.Queries.WithTx(tx)

	tender, err := queriesWithTx.CreateTender(ctx.Request.Context(), params)
	if err != nil {
		_ = tx.Rollback(ctx)

		api.l.Error().Err(err).Msg("failed to add new tender to storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	err = queriesWithTx.CreateTenderHistory(ctx.Request.Context(), int32(tender.ID))
	if err != nil {
		_ = tx.Rollback(ctx)

		api.l.Error().Err(err).Msg("failed to add new tender in history table to storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	if err := tx.Commit(ctx); err != nil {
		api.l.Error().Err(err).Msg("failed to commit transaction")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Ошибка транзакции.")

		return
	}

	res := oapi.Tender{
		Id:          tender.ID,
		Name:        tender.Name,
		Description: tender.Description.String,
		Status:      oapi.TenderStatus(tender.Status.StatusType),
		ServiceType: oapi.TenderServiceType(tender.ServiceType.String),
		Verstion:    tender.Version.Int32,
		CreatedAt:   tender.CreatedAt.Time.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) GetTenders(ctx *gin.Context, params oapi.GetTendersParams) {
	limit := int32(5)
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}

	offset := int32(0)
	if params.Offset != nil && *params.Offset > 0 {
		offset = *params.Offset
	}

	var paramsServiceType []string

	if params.ServiceType != nil && len(*params.ServiceType) > 0 {
		for _, value := range *params.ServiceType {
			if value != "" {
				paramsServiceType = append(paramsServiceType, string(value))
			}
		}
	}

	if len(paramsServiceType) == 0 {
		paramsServiceType = nil
	}

	args := queries.GetTendersParams{
		Limit:   limit,
		Offset:  offset,
		Column1: paramsServiceType,
	}

	tenders, err := api.storage.Queries.GetTenders(ctx.Request.Context(), args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to get tenders from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	res := make([]oapi.Tender, len(tenders))

	for idx, value := range tenders {
		res[idx] = oapi.Tender{
			Id:          value.ID,
			Name:        value.Name,
			Description: value.Description.String,
			Status:      oapi.TenderStatus(value.Status.StatusType),
			ServiceType: oapi.TenderServiceType(value.ServiceType.String),
			Verstion:    value.Version.Int32,
			CreatedAt:   value.CreatedAt.Time.Format(time.RFC3339),
		}
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) GetUserTenders(ctx *gin.Context, params oapi.GetUserTendersParams) {
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

	args := queries.GetUserTendersParams{
		Username: username,
		Limit:    limit,
		Offset:   offset,
	}

	tenders, err := api.storage.Queries.GetUserTenders(ctx.Request.Context(), args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to get tenders from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	res := make([]oapi.Tender, len(tenders))

	for idx, value := range tenders {
		res[idx] = oapi.Tender{
			Id:          value.ID,
			Name:        value.Name,
			Description: value.Description.String,
			Status:      oapi.TenderStatus(value.Status.StatusType),
			ServiceType: oapi.TenderServiceType(value.ServiceType.String),
			Verstion:    value.Version.Int32,
			CreatedAt:   value.CreatedAt.Time.Format(time.RFC3339),
		}
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) GetTenderStatus(ctx *gin.Context, tenderId int32, params oapi.GetTenderStatusParams) {
	if tenderId <= 0 {
		api.l.Error().Msg("invalid tender id argument")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	username := ""
	if params.Username != nil {
		username = *params.Username
	}

	args := queries.GetTenderStatusParams{
		ID:       tenderId,
		Username: username,
	}

	status, err := api.storage.Queries.GetTenderStatus(ctx.Request.Context(), args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to get status from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	res := oapi.TenderStatus(status.StatusType)

	ctx.JSON(http.StatusOK, res)
}

func (api *API) UpdateTenderStatus(ctx *gin.Context, tenderId int32, params oapi.UpdateTenderStatusParams) {
	if tenderId <= 0 {
		api.l.Error().Msg("invalid tender id argument")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	args := queries.UpdateTenderStatusParams{
		ID:       tenderId,
		Username: params.Username,
		Status:   queries.NullStatusType{StatusType: queries.StatusType(params.Status), Valid: true},
	}

	updatedTender, err := api.storage.Queries.UpdateTenderStatus(ctx.Request.Context(), args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to update status in storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	res := oapi.Tender{
		Id:          updatedTender.ID,
		Name:        updatedTender.Name,
		Description: updatedTender.Description.String,
		Status:      oapi.TenderStatus(updatedTender.Status.StatusType),
		ServiceType: oapi.TenderServiceType(updatedTender.ServiceType.String),
		Verstion:    updatedTender.Version.Int32,
		CreatedAt:   updatedTender.CreatedAt.Time.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) EditTender(ctx *gin.Context, tenderId int32, params oapi.EditTenderParams) {
	if tenderId <= 0 {
		api.l.Error().Msg("invalid tender id argument")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest,
			"Данные неправильно сформированы или не соответствуют требованиям.")

		return
	}

	var req oapi.EditTenderJSONBody

	if err := ctx.BindJSON(&req); err != nil {
		api.l.Error().Err(err).Msg("failed to unmarshall create tender body")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest,
			"Данные неправильно сформированы или не соответствуют требованиям.")

		return
	}

	var description nan.NullString
	if req.Description != nil {
		description = nan.NullString{String: *req.Description, Valid: true}
	}

	var serviceType nan.NullString
	if req.ServiceType != nil {
		serviceType = nan.NullString{String: string(*req.ServiceType), Valid: true}
	}

	args := queries.EditTenderParams{
		Name:        *req.Name,
		Description: description,
		ServiceType: serviceType,
		ID:          tenderId,
		Username:    params.Username,
	}

	tx, err := api.storage.BeginTx(ctx.Request.Context())
	if err != nil {
		api.l.Error().Err(err).Msg("failed to begin transaction")
		helper.CustomErrorResponse(ctx, http.StatusInternalServerError, "Ошибка транзакции.")

		return
	}

	queriesWithTx := api.storage.Queries.WithTx(tx)

	tender, err := queriesWithTx.EditTender(ctx.Request.Context(), args)
	if err != nil {
		_ = tx.Rollback(ctx)
		helper.CustomErrorResponse(ctx, http.StatusBadRequest,
			"Данные неправильно сформированы или не соответствуют требованиям.")
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": "Failed to edit tender"})

		return
	}

	err = queriesWithTx.CreateTenderHistory(ctx.Request.Context(), int32(tender.ID))
	if err != nil {
		_ = tx.Rollback(ctx)
		helper.CustomErrorResponse(ctx, http.StatusBadRequest,
			"Данные неправильно сформированы или не соответствуют требованиям.")
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": "Failed to edit tender"})

		return
	}

	if err := tx.Commit(ctx); err != nil {
		api.l.Error().Err(err).Msg("failed to commit transaction")
		helper.CustomErrorResponse(ctx, http.StatusInternalServerError, "Ошибка транзакции.")

		return
	}

	res := oapi.Tender{
		Id:          tender.ID,
		Name:        tender.Name,
		Description: tender.Description.String,
		Status:      oapi.TenderStatus(tender.Status.StatusType),
		ServiceType: oapi.TenderServiceType(tender.ServiceType.String),
		Verstion:    tender.Version.Int32,
		CreatedAt:   tender.CreatedAt.Time.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) RollbackTender(ctx *gin.Context, tenderId int32, version int32, params oapi.RollbackTenderParams) {
	if tenderId <= 0 || version <= 0 {
		api.l.Error().Msg("invalid arguments")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	args := queries.RollbackTenderParams{
		ID:       tenderId,
		Version:  nan.Int32(version),
		Username: params.Username,
	}

	tender, err := api.storage.Queries.RollbackTender(ctx.Request.Context(), args)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to rollback tender from storage")
		helper.CustomErrorResponse(ctx, http.StatusBadRequest, "Неверный формат запроса или его параметры.")

		return
	}

	res := oapi.Tender{
		Id:          tender.ID,
		Name:        tender.Name,
		Description: tender.Description.String,
		Status:      oapi.TenderStatus(tender.Status.StatusType),
		ServiceType: oapi.TenderServiceType(tender.ServiceType.String),
		Verstion:    tender.Version.Int32,
		CreatedAt:   tender.CreatedAt.Time.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, res)
}
