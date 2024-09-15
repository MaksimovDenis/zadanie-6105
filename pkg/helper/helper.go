package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CustomErrorResponse(ctx *gin.Context, statusCode int, reason string) {
	switch statusCode {
	case http.StatusBadRequest:
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": reason})
	case http.StatusUnauthorized:
		ctx.JSON(http.StatusUnauthorized, gin.H{"reason": "Пользователь не существует или некорректен."})
	case http.StatusForbidden:
		ctx.JSON(http.StatusForbidden, gin.H{"reason": "Недостаточно прав для выполнения действия."})
	case http.StatusNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"reason": "Нет информации об объекте"})
	case http.StatusInternalServerError:
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": "Сервер не готов обрабатывать запросы, если ответ статусом 500 или любой другой, кроме 200."}) // nolint
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": "Сервер не готов обрабатывать запросы, если ответ статусом 500 или любой другой, кроме 200."}) // nolint
	}
}
