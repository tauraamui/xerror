package xerror

import (
	"testing"

	"github.com/matryer/is"
)

func TestWrappedArgs(t *testing.T) {
	is := is.New(t)

	args := wrappedArgs("wrapped err: %w")
	is.True(len(args) == 1)
	is.Equal(args[0], 1)
}

func TestWrappedArgsAfterOtherArgs(t *testing.T) {
	is := is.New(t)

	args := wrappedArgs("too many buckets: %d/%d, wrapped err: %w")
	is.True(len(args) == 1)
	is.Equal(args[0], 3)
}

func TestWrappedArgsBetweenOtherArgs(t *testing.T) {
	is := is.New(t)

	args := wrappedArgs("service proc: %s failed, wrapped err: %w, lost %d")
	is.True(len(args) == 1)
	is.Equal(args[0], 2)
}

func TestKindPinned(t *testing.T) {
	is := is.New(t)
	err := NewWithKind("NETWORK_ERROR", "something went wrong").Pin()
	x, ok := err.(*x)
	is.True(ok)
	is.True(x.isPinned())
}

func TestWrappedErrorHasKindPinned(t *testing.T) {
	is := is.New(t)
	err := Errorf("wrapped custom err: %w", NewWithKind("WRAPPED_ERROR", "not enough rocks").Pin())
	x, ok := err.(*x)
	is.True(ok)
	is.True(x.isPinned())
}
