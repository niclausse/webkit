package response

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/niclausse/errorx"
	"github.com/niclausse/webkit/v2/mode"
	"github.com/niclausse/zlog"
	"github.com/pkg/errors"
	"net/http"
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

type Result struct {
	ErrNo   int      `json:"errNo"`
	ErrMsg  string   `json:"errMsg"`
	Data    NullData `json:"data"`
	Details []string `json:"details,omitempty"`
}

type NullData struct {
	Valid bool
	Data  interface{}
}

func (nd *NullData) MarshalJSON() ([]byte, error) {
	if !nd.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(nd.Data)
}

func (r *responder) Fail(ctx *gin.Context, err error) {
	ex, ok := errors.Cause(err).(*errorx.ErrorX)
	if !ok {
		ex = errorx.SystemErr.WithDetails("backend should use errorX!!!")
	}

	resp := &Result{
		ErrNo:  ex.ErrNo,
		ErrMsg: ex.ErrMsg,
	}

	if r.runMode == mode.DevelopMode {
		resp.Details = ex.Details
	}

	zlog.WithContext(ctx.Request.Context()).Errorf("%+v", err)

	requestId, _ := ctx.Request.Context().Value(zlog.ContextKeyRequestID).(string)
	ctx.Header("X-Request-Id", requestId)
	ctx.JSON(http.StatusOK, resp)
}

func (r *responder) Succeed(ctx *gin.Context, data interface{}) {
	requestId, _ := ctx.Request.Context().Value(zlog.ContextKeyRequestID).(string)
	ctx.Header("X-Request-Id", requestId)
	ctx.JSON(http.StatusOK, &Result{
		ErrNo:  0,
		ErrMsg: "",
		Data: NullData{
			Valid: true,
			Data:  data,
		},
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
