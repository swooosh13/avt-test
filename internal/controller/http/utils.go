package rest

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swooosh13/avito-test/pkg/pagination"
)

func parseID(c *gin.Context, param string) (int64, error) {
	idParam := c.Param(param)
	if idParam == "" {
		return 0, errors.New("empty id param")
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id param")
	}

	return id, nil
}

func parsePaginationParams(c *gin.Context) pagination.Params {
	var (
		params pagination.Params
		err    error
	)

	params.Limit, err = strconv.ParseUint(c.Query("limit"), 10, 64)
	if err != nil {
		params.Limit = pagination.DefaultLimit
	}

	params.Offset, err = strconv.ParseUint(c.Query("offset"), 10, 64)
	if err != nil {
		params.Offset = pagination.DefaultOffset
	}

	return params
}
