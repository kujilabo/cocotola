package factory_test

import (
	"testing"

	"github.com/kujilabo/cocotola/lib/factory"
)

type A interface {
	RunA()
}
type a struct {
}

func NewA() (A, error) {
	return &a{}, nil
}

func (a *a) RunA() {}

func Test_Factory1(t *testing.T) {
	// factory.Factory1(NewA)
	factory.PrintSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	factory.PrintSlice([]string{"a", "b", "c", "d"})
}
