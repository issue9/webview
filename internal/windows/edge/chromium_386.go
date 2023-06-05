//go:build windows
// +build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/issue9/webview/internal/windows/w32"
)

func (e *Chromium) Resize() {
	if e.controller == nil {
		return
	}
	var bounds windows.Rect
	w32.GetClientRect(e.hwnd, &bounds)

	e.controller.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(e.controller)),
		uintptr(bounds.Left),
		uintptr(bounds.Top),
		uintptr(bounds.Right),
		uintptr(bounds.Bottom),
	)
}
