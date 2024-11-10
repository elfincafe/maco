package mmaco

import (
	"reflect"
	"strings"
)

type (
	meta struct {
		kind     reflect.Kind
		short    string
		long     string
		required bool
		desc     string
		value    string
		handler  string
	}
)

func newMeta() *meta {
	m := new(meta)
	m.kind = reflect.Invalid
	m.short = ""
	m.long = ""
	m.required = false
	m.desc = ""
	m.value = ""
	m.handler = ""
	return m
}

func getMetas(t reflect.Type) map[string]*meta {
	metas := map[string]*meta{}
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		kind := t.Field(i).Type.Kind()
		content := t.Field(i).Tag.Get(tagName)
		if content == "" {
			continue
		}
		meta := getMeta(content)
		meta.kind = kind
		metas[name] = meta
	}
	return metas
}

func getMeta(content string) *meta {
	meta := newMeta()
	contents := strings.Split(content, ",")
	key := ""
	for _, v := range contents {
		t := strings.TrimLeft(v, " \t\v\r\n\f")
		if strings.HasPrefix(strings.ToLower(t), "short=") {
			meta.short = strings.TrimSpace(t[6:])
			key = "short"
		} else if strings.HasPrefix(strings.ToLower(t), "long=") {
			meta.long = strings.TrimSpace(t[5:])
			key = "long"
		} else if strings.HasPrefix(strings.ToLower(t), "desc=") {
			meta.desc = t[5:]
			key = "desc"
		} else if strings.HasPrefix(strings.ToLower(t), "value=") {
			meta.value = strings.TrimSpace(t[6:])
			key = "value"
		} else if strings.HasPrefix(strings.ToLower(t), "handler=") {
			meta.handler = strings.TrimSpace(t[8:])
			key = "handler"
		} else if strings.ToLower(strings.TrimSpace(t)) == "required" {
			meta.required = true
		} else if key == "desc" {
			meta.desc += "," + v // concatinate variable "v" not "t"
		}
	}
	return meta
}