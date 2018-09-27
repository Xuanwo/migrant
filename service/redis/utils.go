package redis

import "strings"

func convert(s string) []interface{} {
	r := make([]interface{}, 0)
	for _, v := range strings.Split(s, " ") {
		if v == "" {
			continue
		}
		r = append(r, v)
	}

	return r
}
