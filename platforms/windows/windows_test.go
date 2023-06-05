// SPDX-License-Identifier: MIT

//go:build windows

package windows

import (
	"encoding/json"
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/issue9/webview/webviewtest"
)

func TestJSString(t *testing.T) {
	a := assert.New(t, false)

	val := "abc"
	b, err := json.Marshal(val)
	a.NotError(err).Equal(string(b), jsString(val))

	val = "abc\""
	b, err = json.Marshal(val)
	a.NotError(err).Equal(string(b), jsString(val))

	val = "abc'"
	b, err = json.Marshal(val)
	a.NotError(err).Equal(string(b), jsString(val))
}

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	d, err := New(nil)
	a.NotError(err).NotNil(d)
	webviewtest.Desktop(d)
}
