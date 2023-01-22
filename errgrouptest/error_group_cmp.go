package errgrouptest

import (
	"errors"
	"fmt"
	"strings"
)

func MatchGroupErrors(expErrs []string, actualErr error) error {
	return resolveErrors(expErrs, errorListFromErr(actualErr))
}

func errorListFromErr(err error) []string {
	var errs []string
	if err != nil {
		errs = strings.Split(strings.TrimSpace(strings.ReplaceAll(err.Error(), "\t", "")), "\n")
	}
	return errs
}

func resolveErrors(exp, act []string) error {
	sb := strings.Builder{}

	if len(act) > 0 {
		act = act[1:]
	}

	renderNonExistentErrors(&sb, exp, act)
	renderUnexpectedErrors(&sb, exp, act)

	failureMsg := sb.String()
	if len(failureMsg) > 0 {
		return errors.New(failureMsg)
	}
	return nil
}

func renderNonExistentErrors(sb *strings.Builder, exp, act []string) {
	if str := renderErrors("non-existent expected error(s)", exp, act); len(str) > 0 {
		sb.WriteString(str)
	}
}

func renderUnexpectedErrors(sb *strings.Builder, exp, act []string) {
	if str := renderErrors("unexpected error(s)", act, exp); len(str) > 0 {
		sb.WriteString(str)
	}
}

func renderErrors(header string, exp, act []string) string {
	sb := strings.Builder{}
	errs := missingErrors(exp, act)
	if c := len(errs); c > 0 {
		sb.WriteString(fmt.Sprintf("%d %s:\n", c, header))

		for _, ee := range errs {
			sb.WriteString(fmt.Sprintf("\t%s\n", ee))
		}
	}
	return sb.String()
}

func missingErrors(exp, act []string) []string {
	missing := []string{}
	for _, ee := range exp {
		foundMatch := false
		for _, ae := range act {
			if ee == ae {
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			missing = append(missing, ee)
		}
	}
	return missing
}
