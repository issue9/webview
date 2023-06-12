// SPDX-License-Identifier: MIT

//go:build linux || openbsd || freebsd || netbsd

package gtk

import (
	"log"
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/issue9/webview"
	"github.com/issue9/webview/internal/presets"
)

func TestSanitizeOptions(t *testing.T) {
	a := assert.New(t, false)

	o := sanitizeOptions(nil)
	a.Equal(o.Title, presets.Title).
		Equal(o.Error, log.Default())

	o = sanitizeOptions(&Options{
		Size: webview.Size{Height: 1000, Width: 1000},
	})
	a.Equal(o.Title, presets.Title).
		Equal(o.Size.Height, 1000).
		Equal(o.Size.Width, 1000)
}
