package xerror

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type XError interface {
	Error() string
	ErrorMsg() string
}

type I interface {
	AsKind(Kind) I
	IsKind(Kind) bool
	Pin() XError
	isPinned() bool
	ToError() error
	XError
	Msg(string) I
	WithStackTrace() I
	// WithParam will append the param key/value pair passed in
	// to the internal map.
	WithParam(string, interface{}) I
	// WithParams will set or merge given params map to use
	// as the param values. Passing nil instead of a map will
	// clear the params completely.
	WithParams(map[string]interface{}) I
}

type Kind string

const NA = Kind("N/A")

type x struct {
	wrapped      *x
	kind         Kind
	pinned       bool
	causeErr     error
	errMsg       string
	error        error
	stackTrace   bool
	stackPrinter func(error) string
	params       map[string]interface{}
}

func Errorf(format string, values ...any) I {
	x := newFromError(fmt.Errorf(format, values...), extractWrapped(format, values))
	return x
}

func extractWrapped(format string, values []any) *x {
	args := wrappedArgs(format)
	if len(args) > 0 && len(values) >= args[0]-1 {
		if ww, ok := values[args[0]-1].(*x); ok {
			return ww
		}
	}

	return nil
}

func wrappedArgs(format string) []int {
	args := []int{}
	end := len(format)

	verbCount := 0
	for i := 0; i < end; {
		b := format[i]
		foundVerb := false
		if b == '%' {
			if i+1 >= end {
				break
			}
			foundVerb = true
			verbCount++
			i++
		}

		if foundVerb && format[i] == 'w' {
			args = append(args, verbCount)
		}

		i++
	}

	return args
}

func newFromError(e error, w *x) I {
	var cause error
	if c := errors.Unwrap(e); c != nil {
		cause = c
	}
	return &x{
		wrapped:  w,
		error:    e,
		causeErr: cause,
		kind:     NA, errMsg: e.Error(), stackTrace: false, pinned: false,
		stackPrinter: func(e error) string {
			return fmt.Sprintf("%+v", e)
		},
	}
}

func New(es string) I {
	return NewWithKind(NA, es)
}

func NewWithKind(k Kind, es string) I {
	i := x{
		kind: k, errMsg: es, stackTrace: false,
		stackPrinter: func(e error) string {
			return fmt.Sprintf("%+v", e)
		},
	}
	i.format()
	return &i
}

func (x *x) format() {
	err := errors.New(x.toString())
	x.error = err
}

func (x *x) Is(e error) bool {
	if errors.Is(x.error, e) || errors.Is(x.causeErr, e) {
		return true
	}

	if err := errors.Unwrap(x.error); err != nil {
		return errors.Is(err, e)
	}

	return false
}

func (x *x) isPinned() bool {
	if yes := x.pinned; yes {
		return yes
	}

	if x.wrapped != nil {
		return x.wrapped.isPinned()
	}

	return false
}

func (x *x) IsKind(k Kind) bool {
	if yes := x.kind == k; yes {
		return yes
	}

	if x.wrapped != nil {
		return x.wrapped.IsKind(k)
	}

	return false
}

func (x *x) AsKind(k Kind) I {
	defer x.format()
	x.kind = k
	return x
}

func (x *x) Pin() XError {
	defer x.format()
	x.pinned = true
	return x
}

func (x *x) ToError() error {
	return x.error
}

func (x *x) Error() string {
	x.format()

	if x.stackTrace {
		return x.stackPrinter(x.error)
	}
	return x.error.Error()
}

func (x *x) ErrorMsg() string {
	return x.errMsg
}

func (x *x) Msg(m string) I {
	defer x.format()
	x.errMsg = m
	return x
}

func (x *x) WithStackTrace() I {
	x.stackTrace = true
	x.format()
	return x
}

func (x *x) WithParams(p map[string]interface{}) I {
	defer x.format()

	if x.params != nil {
		mergeParams(x.params, p)
		return x
	}

	x.params = p
	return x
}

func mergeParams(p1, p2 map[string]interface{}) {
	for k, v := range p2 {
		p1[k] = v
	}
}

// WithParam will append the param key/value pair passed in
// to the internal map.
func (x *x) WithParam(key string, v interface{}) I {
	defer x.format()
	if x.params == nil {
		x.params = map[string]interface{}{}
	}
	x.params[key] = v
	return x
}

func (x *x) resolveKind() Kind {
	k := x.kind

	pk, wasPinned := x.pinnedChildKind()
	if wasPinned {
		return pk
	}

	return k
}

func (x *x) pinnedChildKind() (Kind, bool) {
	if x.wrapped == nil {
		return x.kind, x.pinned
	}

	return x.wrapped.pinnedChildKind()
}

func (x *x) toString() string {
	logMsg := x.errMsg
	kind := x.resolveKind()
	if !x.pinned && kind != NA || len(x.params) > 0 {
		logMsg = fmt.Sprintf("Kind: %s | %s", strings.ToUpper(string(kind)), x.errMsg)
	}

	params := []string{}
	for k, v := range x.params {
		params = append(params, fmt.Sprintf("%s: {%+v}", k, v))
	}

	return fmt.Sprintf("%s%s", logMsg, func() string {
		if len(params) != 0 {
			return fmt.Sprintf(", Params: [%+v]", strings.Join(params, " | "))
		}
		return ""
	}())
}
