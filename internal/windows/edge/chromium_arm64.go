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

	words := (*[2]uintptr)(unsafe.Pointer(&bounds))
	e.controller.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(e.controller)),
		words[0],
		words[1],
	)
}
