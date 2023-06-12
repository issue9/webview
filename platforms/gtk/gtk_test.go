// SPDX-License-Identifier: MIT

//go:build linux || openbsd || freebsd || netbsd

package gtk

import (
	"os"
	"testing"

	"github.com/issue9/webview/webviewtest"
)

func TestMain(m *testing.M) {
	d := New(&Options{Debug: true})
	webviewtest.Desktop(d)
	os.Exit(m.Run())
}
