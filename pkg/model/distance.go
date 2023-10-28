package model

import (
	"errors"
)

type DistanceNorm string

const (
	DistanceNormSubstring   DistanceNorm = "substring"
	DistanceNormLevenshtein DistanceNorm = "levenshtein"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *DistanceNorm) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *DistanceNorm) Set(v string) error {
	switch v {
	case "substring", "levenshtein":
		*e = DistanceNorm(v)
		return nil
	default:
		return errors.New("must be one of: `substring` or `levenshtein`")
	}
}

// Type is only used in help text
func (e *DistanceNorm) Type() string {
	return "DistanceNorm"
}
