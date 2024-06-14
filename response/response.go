package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/niclausse/webkit/errorx"
	"github.com/niclausse/webkit/mode"
	"github.com/niclausse/webkit/zlog"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

var std Responder = &responder{runMode: mode.DevelopMode}

func CustomizeResponder(_std Responder) {
	if _std != nil {
		std = _std
	}
}

type Responder interface {
	SetMode(mode mode.Mode)
	Fail(ctx *gin.Context, err error)
	Succeed(ctx *gin.Context, data interface{})
}

type responder struct {
	runMode mode.Mode
}

func (r *responder) SetMode(mode mode.Mode) {
	r.runMode = mode
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

	if r.runMode == mode.DevelopMode {
		resp["details"] = ex.Details
		resp["stack"] = stack
	}

	zlog.Errorf("%+v, sid: %s", err, ctx.GetString("sid"))

	ctx.JSON(http.StatusOK, resp)
}

func (r *responder) Succeed(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"errNo":  0,
		"errMsg": "",
		"data":   data,
	})
}

func SetMode(mode mode.Mode) {
	std.SetMode(mode)
}

func Fail(ctx *gin.Context, err error) {
	std.Fail(ctx, err)
}

func Succeed(ctx *gin.Context, data interface{}) {
	std.Succeed(ctx, data)
}
