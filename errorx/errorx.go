package errorx

import (
	"fmt"
	"io"
)

type ErrorX struct {
	ErrNo   int
	ErrMsg  string
	Details []string
	*stack
}

func (e *ErrorX) Error() string {
	return fmt.Sprintf("err_no: %d, err_msg: %s, details: %v", e.ErrNo, e.ErrMsg, e.Details)
}

func (e *ErrorX) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.Error())
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func New(errNo int, errMsg string, details ...string) *ErrorX {
	return &ErrorX{ErrNo: errNo, ErrMsg: errMsg, Details: details}
}

func (e *ErrorX) WithDetails(details ...string) *ErrorX {
	x := &ErrorX{
		ErrNo:   e.ErrNo,
		ErrMsg:  e.ErrMsg,
		Details: e.Details,
	}

	x.Details = append(x.Details, details...)
	return x
}

func (e *ErrorX) WithStack() *ErrorX {
	return &ErrorX{
		ErrNo:   e.ErrNo,
		ErrMsg:  e.ErrMsg,
		Details: e.Details,
		stack:   callers(),
	}
}

func WithStack(err error, biz *ErrorX) error {
	if err == nil {
		return nil
	}

	if biz == nil {
		return SystemError.WithDetails(err.Error())
	}

	x := &ErrorX{
		ErrNo:   biz.ErrNo,
		ErrMsg:  biz.ErrMsg,
		Details: biz.Details,
		stack:   biz.stack,
	}

	x.Details = append(x.Details, err.Error())

	if x.stack == nil {
		x.stack = callers()
	}

	return x
}
