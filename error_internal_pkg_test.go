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
