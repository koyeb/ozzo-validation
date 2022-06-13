// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"regexp"
)

// ErrMatchInvalid is the error that returns in case of invalid format.
var ErrMatchInvalid = NewError("validation_match_invalid", "must be in a valid format")

// Match returns a validation rule that checks if a value matches the specified regular expression.
// This rule should only be used for validating strings and byte slices, or a validation error will be reported.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func Match(re *regexp.Regexp) MatchRule {
	return MatchRule{
		re:     re,
		err:    ErrMatchInvalid,
		invert: false,
	}
}

// NotMatch returns a validation rule that checks if a value doesn't matches the specified regular expression.
// This rule should only be used for validating strings and byte slices, or a validation error will be reported.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func NotMatch(re *regexp.Regexp) MatchRule {
	return MatchRule{
		re:     re,
		err:    ErrMatchInvalid,
		invert: true,
	}
}

// MatchRule is a validation rule that checks if a value matches the specified regular expression.
type MatchRule struct {
	re     *regexp.Regexp
	err    Error
	invert bool
}

// Validate checks if the given value is valid or not.
func (r MatchRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil {
		return nil
	}

	isString, str, isBytes, bs := StringOrBytes(value)

	switch {
	case isString:
		if str == "" {
			return nil
		}
		x := r.invert
		y := r.re.MatchString(str)
		if (x || y) && !(x && y) {
			return nil
		}

	case isBytes:
		if len(bs) == 0 {
			return nil
		}
		x := r.invert
		y := r.re.Match(bs)
		if (x || y) && !(x && y) {
			return nil
		}
	}

	return r.err
}

// Error sets the error message for the rule.
func (r MatchRule) Error(message string) MatchRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r MatchRule) ErrorObject(err Error) MatchRule {
	r.err = err
	return r
}
