package xerror_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/pkg/errors"
	"github.com/tauraamui/xerror"
)

type tableTest interface {
	preTest(*testing.T)
	run(*testing.T)
}

type xerrorWrappedErrorIsTest struct {
	skip           bool
	title          string
	parentError    error
	wrappedError   error
	shouldNotMatch bool
}

func (x xerrorWrappedErrorIsTest) preTest(t *testing.T) {
	if len(x.title) == 0 {
		t.Error("table tests must all have titles")
	}

	if x.skip {
		t.Skip()
	}
}

func (x xerrorWrappedErrorIsTest) run(t *testing.T) {
	t.Run(x.title, func(t *testing.T) {
		is := is.NewRelaxed(t)

		r := errors.Is(x.parentError, x.wrappedError)
		is.True((x.shouldNotMatch && !r) || (!x.shouldNotMatch && r))
	})
}

var nativeErrorType = errors.New("native error type")
var customErrorType = xerror.New("custom error type")
var deepChildErrorType = xerror.New("deep error type")

func TestWrappedErrorsFoundWithIsError(t *testing.T) {
	tests := []xerrorWrappedErrorIsTest{
		{
			title:        "error is resolves native error wrapped by native error type",
			parentError:  fmt.Errorf("parent error: %w", nativeErrorType),
			wrappedError: nativeErrorType,
		},
		{
			title:        "error is resolves custom error wrapped by native error type",
			parentError:  fmt.Errorf("parent error: %w", customErrorType),
			wrappedError: customErrorType,
		},
		{
			title:        "error is resolves custom error with params wrapped by native error type",
			parentError:  fmt.Errorf("parent error: %w", customErrorType.WithParam("music", "heavy-metal")),
			wrappedError: customErrorType,
		},
		{
			title:        "error is resolves custom error with stack trace wrapped by native error type",
			parentError:  fmt.Errorf("parent error: %w", customErrorType.WithStackTrace()),
			wrappedError: customErrorType,
		},
		{
			title: "error is resolves deeply buried error in custom error tree wrapped by native error type",
			parentError: fmt.Errorf(
				"parent error: %w", xerror.Errorf(
					"sub-1 parent error: %w", xerror.Errorf(
						"sub-2 parent error: %w", deepChildErrorType,
					)),
			),
			wrappedError: deepChildErrorType,
		},
		{
			title:        "error is resolves custom error with params wrapped in custom error",
			parentError:  xerror.Errorf("parent error: %w", customErrorType.WithParam("music", "heavy-metal")),
			wrappedError: customErrorType,
		},
		{
			title:        "error is resolves custom error wrapped in custom error with params",
			parentError:  xerror.Errorf("parent error: %w", customErrorType).WithParam("holiday", "spain"),
			wrappedError: customErrorType,
		},
		{
			title:        "error is resolves custom error wrapped in custom error with stack trace",
			parentError:  xerror.Errorf("parent error: %w", customErrorType).WithStackTrace(),
			wrappedError: customErrorType,
		},
		{
			title:        "error is resolves custom error with stack trace wrapped in custom error",
			parentError:  xerror.Errorf("parent error: %w", customErrorType.WithStackTrace()),
			wrappedError: customErrorType,
		},
		{
			title: "error is resolves deeply buried error wrapped in custom errors tree",
			parentError: xerror.Errorf(
				"parent error: %w", xerror.Errorf(
					"sub-1 parent error: %w", xerror.Errorf(
						"sub-2 parent error: %w", deepChildErrorType,
					)),
			),
			wrappedError: deepChildErrorType,
		},
		{
			title:          "custom error with no wrapped error",
			parentError:    xerror.Errorf("no wrapped errors here"),
			wrappedError:   customErrorType,
			shouldNotMatch: true,
		},
		{
			title:          "custom error with different wrapped error",
			parentError:    xerror.Errorf("native error right?: %w", nativeErrorType),
			wrappedError:   customErrorType,
			shouldNotMatch: true,
		},
	}

	for _, tt := range tests {
		runTest(t, tt)
	}
}

func TestNewErrorGivesErrInstance(t *testing.T) {
	is := is.New(t)

	err := xerror.New("")
	is.True(err != nil)
}

const TestError = xerror.Kind("test_error")
const TestParamsError = xerror.Kind("test_params_error")

func TestErrorMsgMatchesGivenErrorMsg(t *testing.T) {
	t.Run("error msg matches given initial msg but doesn't include context", func(t *testing.T) {
		is := is.New(t)

		err := xerror.NewWithKind("MUTABLE_ERR", "test error message")
		is.Equal(err.ErrorMsg(), "test error message")
	})

	t.Run("error msg matches msg which replaces initial msg", func(t *testing.T) {
		is := is.New(t)

		err := xerror.NewWithKind("MUTABLE_ERR", "test error message").Msg("replaced message!")
		is.Equal(err.ErrorMsg(), "replaced message!")
	})
}

func TestErrorToError(t *testing.T) {
	is := is.New(t)

	nativeErr := xerror.NewWithKind("NATIVE_ERR", "native error").ToError()
	is.True(nativeErr != nil)
	is.Equal(nativeErr.Error(), "Kind: NATIVE_ERR | native error")
}

type xerrorExpectedStringTest struct {
	skip       bool
	title      string
	err        error
	expected   string
	customEval func(string) error
}

func (x xerrorExpectedStringTest) preTest(t *testing.T) {
	if len(x.title) == 0 {
		t.Error("table tests must all have titles")
	}

	if x.skip {
		t.Skip()
	}
}

func (x xerrorExpectedStringTest) run(t *testing.T) {
	t.Run(x.title, func(t *testing.T) {
		is := is.NewRelaxed(t)

		if x.customEval != nil {
			is.NoErr(x.customEval(x.err.Error()))
			return
		}

		is.Equal(x.err.Error(), x.expected)
	})
}

func TestNewErrorAndErrorfsOutputsExpectedString(t *testing.T) {
	tests := []xerrorExpectedStringTest{
		{
			title:    "simple new error just prints out msg string",
			err:      xerror.New("fake db update failed"),
			expected: "fake db update failed",
		},
		{
			title:    "new error with param prints out msg string with not assigned kind and with param",
			err:      xerror.New("fake db update failed").WithParam("trace-request-id", "efw4fv32f"),
			expected: "Kind: N/A | fake db update failed, Params: [trace-request-id: {efw4fv32f}]",
		},
		{
			title:    "new error with kind prints out msg string with assigned kind",
			err:      xerror.NewWithKind(TestError, "fake db update failed"),
			expected: "Kind: TEST_ERROR | fake db update failed",
		},
		{
			title:    "new error with kind and param prints out msg string with assigned kind and with param",
			err:      xerror.NewWithKind(TestParamsError, "fake db update failed").WithParam("trace-request-id", "wdgrte4fg"),
			expected: "Kind: TEST_PARAMS_ERROR | fake db update failed, Params: [trace-request-id: {wdgrte4fg}]",
		},
		{
			title:    "new error with kind but then has kind changed prints out msg string with new assigned kind",
			err:      xerror.NewWithKind(TestParamsError, "fake db update failed").AsKind("CHANGED_KIND"),
			expected: "Kind: CHANGED_KIND | fake db update failed",
		},
		{
			title:    "new error with not assigned kind but then has kind changed prints out msg string with new assigned kind",
			err:      xerror.New("fake db update failed").AsKind("FROM_NONE_TO_CHANGED_KIND"),
			expected: "Kind: FROM_NONE_TO_CHANGED_KIND | fake db update failed",
		},
		{
			title:    "new error with msg but has msg updated prints out msg string with new content",
			err:      xerror.New("fake db update failed").Msg("err msg content replaced with this!"),
			expected: "err msg content replaced with this!",
		},
		{
			title: "new error with not assigned kind and with params prints out msg string with not assigned kind and with params",
			err: xerror.New("fake db update failed").WithParams(
				map[string]interface{}{
					"trace-request-id": "fr495fre",
					"request-ip":       "39.49.13.45",
				},
			),
			expected: "Kind: N/A | fake db update failed, Params: [trace-request-id: {fr495fre} | request-ip: {39.49.13.45}]",
			customEval: func(s string) error {
				if !strings.Contains(s, "Kind: N/A | fake db update failed, Params: [") {
					return errors.New("error msg does not include header section")
				}

				if !strings.Contains(s, "trace-request-id: {fr495fre}") {
					return errors.New("error msg params do not contain trace request id entry")
				}

				if !strings.Contains(s, "request-ip: {39.49.13.45}") {
					return errors.New("error msg params do not contain request ip entry")
				}

				return nil
			},
		},
		{
			title: " new error with not assigned kind has param and then with params and prints out expected string",
			err: xerror.New("fake db update failed").WithParam("fruit-type", "peach").WithParams(
				map[string]interface{}{
					"trace-request-id": "fr495fre",
					"request-ip":       "39.49.13.45",
				},
			),
			// keeping unused param here to be clear what we expect
			// the msg to look like, if we pretend that maps are always
			// in key insertion order, which they are not.
			expected: "Kind: N/A | fake db update failed, Params: [fruit-type: {peach} | trace-request-id: {fr495fre} | request-ip: {39.49.13.45}]",
			customEval: func(s string) error {
				if !strings.Contains(s, "Kind: N/A | fake db update failed, Params: [") {
					return errors.New("error msg does not include header section")
				}

				if !strings.Contains(s, "fruit-type: {peach}") {
					return errors.New("error msg params do not contain peach entry")
				}

				if !strings.Contains(s, "trace-request-id: {fr495fre}") {
					return errors.New("error msg params do not contain trace request id entry")
				}

				if !strings.Contains(s, "request-ip: {39.49.13.45}") {
					return errors.New("error msg params do not contain request ip entry")
				}

				return nil
			},
		},
		{
			title:    "errorf prints out msg string",
			err:      xerror.Errorf("too many seconds %d/60 elapsed", 112),
			expected: "too many seconds 112/60 elapsed",
		},
		{
			title:    "errorf with stack trace prints out msg string",
			err:      xerror.Errorf("too many seconds %d/60 elapsed", 112).WithStackTrace(),
			expected: "too many seconds 112/60 elapsed",
			customEval: func(s string) error {
				if !strings.Contains(s, "too many seconds 112/60 elapsed") {
					return errors.New("error msg does not include header section")
				}

				if !strings.Contains(s, "/xerror.(*x).format") {
					return errors.New("error msg does probably not contain stack trace")
				}

				return nil
			},
		},
		{
			title:    "errorf wrapped in native error includes xerr's context information",
			err:      fmt.Errorf("wrapped err: %w", xerror.Errorf("not enough rocks %d/30 in bucket", 19).AsKind("CUSTOM_ERROR")),
			expected: "wrapped err: Kind: CUSTOM_ERROR | not enough rocks 19/30 in bucket",
		},
		{
			title:    "custom error wrapped in custom error includes each one's context information",
			err:      xerror.Errorf("wrapped custom err: %w", xerror.NewWithKind("WRAPPED_ERROR", "not enough rocks")),
			expected: "wrapped custom err: Kind: WRAPPED_ERROR | not enough rocks",
		},
	}

	for _, tt := range tests {
		runTest(t, tt)
	}
}

func runTest(t *testing.T, tt tableTest) {
	tt.preTest(t)
	tt.run(t)
}
