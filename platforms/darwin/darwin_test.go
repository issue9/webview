// SPDX-License-Identifier: MIT

//go:build darwin

package darwin

import (
	"testing"

	"github.com/issue9/assert/v3"
	"github.com/issue9/webview/webviewtest"
)

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	d := New(&Options{})
	a.NotNil(d)
	webviewtest.Desktop(d)
}
