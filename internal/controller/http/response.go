package rest

import "github.com/gin-gonic/gin"

type response struct {
	Data interface{} `json:"data"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func (h *Handler) errorResponse(ctx *gin.Context, status int, message string, err error) {
	h.logger.Error().Err(err).Msg(message)

	ctx.AbortWithStatusJSON(status, errorResponse{
		Message: message,
	})
}

type listRange struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
