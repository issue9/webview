// SPDX-License-Identifier: MIT

// Package pipe 前后端的通信通道
package pipe

import (
	"reflect"
	"strings"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func jsString(v string) string {
	return `"` + strings.ReplaceAll(v, "\"", "\\\"") + `"`
}
