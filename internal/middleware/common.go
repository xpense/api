package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const ErrorMesageBadID = "missing/not-a-number ID in request"

const idKey = "id"

type CommonMiddleware interface {
	SetIDParamToContext(ctx *gin.Context)
}

type commonMiddleware struct{}

func NewCommonMiddleware() CommonMiddleware {
	return &commonMiddleware{}
}

func (c *commonMiddleware) SetIDParamToContext(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": ErrorMesageBadID,
		})
		return
	}

	ctx.Set(idKey, uint(id))
	ctx.Next()
}

func GetIDParamFromContext(ctx *gin.Context) uint {
	id, _ := ctx.Get(idKey)
	return id.(uint)
}
