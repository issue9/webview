// SPDX-License-Identifier: MIT

//go:build windows

package windows

import (
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/issue9/webview/webviewtest"
)

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	d, err := New(nil)
	a.NotError(err).NotNil(d)
	webviewtest.Desktop(d)
}
