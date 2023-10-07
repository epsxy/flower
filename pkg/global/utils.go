package global

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func RemoveConstraints(s []string) []string {
	var r []string
	for _, e := range s {
		reConstraint := regexp.MustCompile(ConstraintRemovalRegexp)
		reIndex := regexp.MustCompile(IndexRemovalRegexp)
		if !reConstraint.MatchString(e) && !reIndex.MatchString(e) {
			r = append(r, e)
		}
	}
	return r
}

func RemoveComments(s []string) []string {
	var r []string
	for _, i := range s {
		if !strings.Contains(i, "--") {
			r = append(r, i)
		} else {
			r = append(r, strings.SplitN(i, "--", 2)[0])
		}
	}
	return r
}

func SplitAndGetFields(raw string) []string {
	var res []string
	rawArr := strings.Split(raw, ",")
	for _, e := range rawArr {
		res = append(res, Unquote(e))
	}
	return res
}

func Unquote(raw string) string {
	raw = strings.ReplaceAll(raw, "'", "")
	raw = strings.ReplaceAll(raw, "`", "")
	raw = strings.ReplaceAll(raw, " ", "")
	return raw
}

func CleanUpString(s string) string {
	re := regexp.MustCompile(`^\s+`)
	return re.ReplaceAllString(s, "")
}

func CleanUpLine(s []string) []string {
	var r []string
	for _, i := range s {
		if i != "" {
			re := regexp.MustCompile(`^\s+`)
			r = append(r, re.ReplaceAllString(i, ""))
		}

	}
	return r
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func PrintStruct(object interface{}) {
	aJSON, _ := json.Marshal(object)
	fmt.Printf("%s\n", string(aJSON))
}
