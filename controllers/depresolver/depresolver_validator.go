package depresolver

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// hostNameRegex allows cloud region formats; e.g. af-south-1
	geoTagRegex = "^[a-zA-Z\\-\\d]*$"
	// hostNameRegex is valid as per RFC 1123 that allows hostname segments could start with a digit
	hostNameRegex = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
	// ipAddressRegex matches valid IPv4 addresses
	ipAddressRegex = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"
	// versionNumberRegex matches version in formats 0.1.2, v0.1.2, v0.1.2-alpha
	versionNumberRegex = "^(v){0,1}(0|(?:[1-9]\\d*))(?:\\.(0|(?:[1-9]\\d*))(?:\\.(0|(?:[1-9]\\d*)))?(?:\\-([\\w][\\w\\.\\-_]*))?)?$"
	// k8sNamespaceRegex matches valid kubernetes namespace
	k8sNamespaceRegex = "^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
)

// validator wrapper against field to be verified
type validator struct {
	strValue string
	strArr   []string
	intValue int
	name     string
	err      error
}

// field creates validator
func field(name string, value interface{}) *validator {
	validator := new(validator)
	validator.name = name
	switch v := value.(type) {
	case int:
		validator.intValue = v
	case int8:
		validator.intValue = int(v)
	case int16:
		validator.intValue = int(v)
	case int32:
		validator.intValue = int(v)
	case int64:
		validator.intValue = int(v)
	case string:
		validator.strValue = v
	case []string:
		validator.strArr = v
	default:
		// float32, float64, bool, interface{}, maps, slices, Custom Types
		validator.err = fmt.Errorf("can't parse %v of type %T as int or string", v, v)
	}
	return validator
}

func (v *validator) isNotEmpty() *validator {
	if v.err != nil {
		return v
	}
	if v.strValue == "" {
		v.err = fmt.Errorf("%s is empty", v.name)
	}
	return v
}

func (v *validator) matchRegexp(regex string) *validator {
	if v.err != nil {
		return v
	}
	if v.strValue == "" {
		return v
	}
	r, _ := regexp.Compile(regex)
	if !r.Match([]byte(v.strValue)) {
		v.err = fmt.Errorf(`"%s" does not match given criteria. see: https://www.regextester.com (%s)`, v.strValue, regex)
	}
	return v
}

// matchRegexps returns error if value is not matched by any of regexp
func (v *validator) matchRegexps(regex ...string) *validator {
	if v.err != nil {
		return v
	}
	for _, r := range regex {
		v.err = nil
		if v.matchRegexp(r).err == nil {
			return v
		}
	}
	return v
}

func (v *validator) isHigherThanZero() *validator {
	if v.err != nil {
		return v
	}
	if v.intValue <= 0 {
		v.err = fmt.Errorf(`"%s" is less or equal to zero`, v.name)
	}
	return v
}

func (v *validator) isHigherOrEqualToZero() *validator {
	if v.err != nil {
		return v
	}
	if v.intValue < 0 {
		v.err = fmt.Errorf(`"%s" is less than zero`, v.name)
	}
	return v
}

func (v *validator) isLessOrEqualTo(num int) *validator {
	if v.err != nil {
		return v
	}
	if v.intValue > num {
		v.err = fmt.Errorf(`"%s" is higher than %v`, v.name, num)
	}
	return v
}

func (v *validator) isNotEqualTo(value string) *validator {
	if v.err != nil {
		return v
	}
	if v.strValue == value {
		v.err = fmt.Errorf(`"%s" can't be equal to "%s"`, v.name, value)
	}
	return v
}

func (v *validator) hasItems() *validator {
	if v.err != nil {
		return v
	}
	if len(v.strArr) == 0 {
		v.err = fmt.Errorf(`"%s" can't be empty`, v.name)
	}
	return v
}

func (v *validator) hasUniqueItems() *validator {
	if v.err != nil {
		return v
	}
	m := make(map[string]bool)
	for _, s := range v.strArr {
		m[s] = true
	}
	if len(m) != len(v.strArr) {
		v.err = fmt.Errorf(`"%s" contains redundant values %s`, v.name, v.strArr)
	}
	return v
}

func isNotEmpty(s string) bool {
	return strings.ReplaceAll(s, " ", "") != ""
}
