package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/niclausse/webkit/consts"
	"github.com/niclausse/webkit/errorx"
	"github.com/niclausse/webkit/zlog"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

var std = &responder{runMode: consts.DevelopMode}

func SetMode(mode consts.Mode) {
	std.runMode = mode
}

type Responder interface {
	Fail(ctx *gin.Context, err error)
	Succeed(ctx *gin.Context, data interface{})
}

func New(mode consts.Mode) Responder {
	return &responder{runMode: mode}
}

type responder struct {
	runMode consts.Mode
}

func (r *responder) Fail(ctx *gin.Context, err error) {
	stack := strings.Split(fmt.Sprintf("%+v", err), "\n")

	ex, ok := errors.Cause(err).(*errorx.ErrorX)
	if !ok {
		ex = errorx.SystemErr.WithDetails("backend should use errorX!!!")
	}

	resp := gin.H{
		"errNo":  ex.ErrNo,
		"errMsg": ex.ErrMsg,
	}

	if r.runMode == consts.DevelopMode {
		resp["details"] = ex.Details
		resp["stack"] = stack
	}

	zlog.Errorf("%+v", err)

	ctx.JSON(http.StatusOK, resp)
}

func (r *responder) Succeed(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"errNo":  0,
		"errMsg": "",
		"data":   data,
	})
}

func Fail(ctx *gin.Context, err error) {
	std.Fail(ctx, err)
}

func Succeed(ctx *gin.Context, data interface{}) {
	std.Succeed(ctx, data)
}
