package response

import (
	"fmt"
	"github.com/penglin1995/webkit/errorx"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Mode string

const (
	ModeDev        Mode = "dev"
	ModeProduction Mode = "production"
)

type Responder interface {
	Fail(ctx *gin.Context, err error)
	Succeed(ctx *gin.Context, data interface{})
}

func New(mode Mode) Responder {
	return &responder{runMode: mode}
}

type responder struct {
	runMode Mode
}

func (r *responder) Fail(ctx *gin.Context, err error) {
	stack := strings.Split(fmt.Sprintf("%+v", err), "\n")

	ex, ok := errors.Cause(err).(*errorx.ErrorX)
	if !ok {
		ex = errorx.SystemError.WithDetails("backend should use errorX!!!")
	}

	resp := gin.H{
		"err_no":  ex.BizNo,
		"err_msg": ex.BizMsg,
	}

	if r.runMode == ModeDev {
		resp["details"] = ex.Details
		resp["stack"] = stack
	}

	logrus.StandardLogger()
	logrus.Errorf("%+v", err)

	ctx.JSON(http.StatusOK, resp)
}

func (r *responder) Succeed(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"err_no":  0,
		"err_msg": "",
		"data":    data,
	})
}
