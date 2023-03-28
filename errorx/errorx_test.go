package errorx

import (
	"fmt"
	"testing"
)

func TestNoStackErr(t *testing.T) {
	fmt.Printf("%+v", stack1(false))
}

func TestStackErr(t *testing.T) {
	fmt.Printf("%+v", stack1(true))
}

func stack1(stack bool) error {
	return stack0(stack)
}

func stack0(stack bool) error {
	err := New(4000, "invalid param")
	if stack {
		err = err.WithStack()
	}
	return err
}
