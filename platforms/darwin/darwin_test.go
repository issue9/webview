// SPDX-License-Identifier: MIT

//go:build darwin

package darwin

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
