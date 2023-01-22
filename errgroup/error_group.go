package errgroup

import (
	"fmt"
	"strings"
)

type I struct {
	errs []error
}

func (eg *I) Append(err error) {
	if eg.errs == nil {
		eg.errs = []error{}
	}
	eg.errs = append(eg.errs, err)
}

func (eg *I) AppendGroup(errg *I) {
	if eg.errs == nil {
		eg.errs = []error{}
	}
	eg.errs = append(eg.errs, errg.errs...)
}

func (eg *I) ToErrOrNil() error {
	if len(eg.errs) > 0 {
		return eg
	}
	return nil
}

func (eg *I) Error() string {
	sb := strings.Builder{}
	errCnt := len(eg.errs)
	sb.WriteString(fmt.Sprintf("%d error%s occurred:\n", errCnt, func() string {
		if errCnt > 1 {
			return "s"
		}
		return ""
	}()))

	for _, err := range eg.errs {
		if err != nil && len(err.Error()) > 0 {
			sb.WriteString(fmt.Sprintf("\t* %s\n", err.Error()))
		}
	}
	return sb.String()
}
