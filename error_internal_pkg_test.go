package xerror

import "testing"

func TestWrappedArgs(t *testing.T) {
	wrappedArgs("wrapped err: %w")
}
