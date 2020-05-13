package pgot

import (
	"reflect"
	"strconv"
	"text/template"
)

var (
	// This should include all publicly accessible custom got functions
	funcMap = template.FuncMap{
		"lnp":   lnp,
		"toInt": toInt,
	}
)

// The lnp (link new page) function converts the supplied url into a link which
// opens a new page.  If label is blank (aka "") then the url is displayed,
// otherwise the label text is used when displaying the link.
func lnp(label, url string) string {
	if label == "" {
		return "<a href=\"" + url + "\" target=\"_blank\">" + url + "</a>"
	} else {
		return "<a href=\"" + url + "\" target=\"_blank\">" + label + "</a>"
	}
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// Leveraged from https://github.com/spf13/cast
// caste.go
// Copyright (c) 2014 Steve Francia
// toInt casts an interface to an int type.
func toInt(i interface{}) int {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		return s
	case int64:
		return int(s)
	case int32:
		return int(s)
	case int16:
		return int(s)
	case int8:
		return int(s)
	case uint:
		return int(s)
	case uint64:
		return int(s)
	case uint32:
		return int(s)
	case uint16:
		return int(s)
	case uint8:
		return int(s)
	case float64:
		return int(s)
	case float32:
		return int(s)
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return int(v)
		}
		panic("unable to cast variable to int")
	case bool:
		if s {
			return 1
		}
		return 0
	case nil:
		return 0
	default:
		panic("unable to cast variable to int")
	}
}
