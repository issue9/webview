// SPDX-License-Identifier: MIT

//go:build windows

package w32

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/issue9/webview"
)

var (
	shlwapi                  = windows.NewLazySystemDLL("shlwapi")
	shlwapiSHCreateMemStream = shlwapi.NewProc("SHCreateMemStream")

	user32                   = windows.NewLazySystemDLL("user32")
	user32GetSystemMetrics   = user32.NewProc("GetSystemMetrics")
	user32LoadImageW         = user32.NewProc("LoadImageW")
	user32RegisterClassExW   = user32.NewProc("RegisterClassExW")
	user32CreateWindowExW    = user32.NewProc("CreateWindowExW")
	user32DestroyWindow      = user32.NewProc("DestroyWindow")
	user32ShowWindow         = user32.NewProc("ShowWindow")
	user32UpdateWindow       = user32.NewProc("UpdateWindow")
	user32SetWindowPos       = user32.NewProc("SetWindowPos")
	user32SetWindowTextW     = user32.NewProc("SetWindowTextW")
	user32GetWindowLongPtrW  = user32.NewProc("GetWindowLongPtrW")
	user32SetWindowLongPtrW  = user32.NewProc("SetWindowLongPtrW")
	user32DefWindowProcW     = user32.NewProc("DefWindowProcW")
	user32AdjustWindowRect   = user32.NewProc("AdjustWindowRect")
	user32GetMessageW        = user32.NewProc("GetMessageW")
	user32TranslateMessage   = user32.NewProc("TranslateMessage")
	user32DispatchMessageW   = user32.NewProc("DispatchMessageW")
	user32PostQuitMessage    = user32.NewProc("PostQuitMessage")
	user32PostThreadMessageW = user32.NewProc("PostThreadMessageW")
	user32IsDialogMessage    = user32.NewProc("IsDialogMessage")
	user32GetClientRect      = user32.NewProc("GetClientRect")
	user32GetAncestor        = user32.NewProc("GetAncestor")
	user32SetFocus           = user32.NewProc("SetFocus")
)

func GetAncestor(h uintptr, flag uint) uintptr {
	r, _, _ := user32GetAncestor.Call(h, uintptr(flag))
	return r
}

func GetClientRect(h uintptr, rect *windows.Rect) {
	user32GetClientRect.Call(h, uintptr(unsafe.Pointer(rect)))
}

func IsDialogMessage(h uintptr, msg *Msg) bool {
	r, _, _ := user32IsDialogMessage.Call(h, uintptr(unsafe.Pointer(msg)))
	return r > 0
}

func PostThreadMessage(thread uintptr, msg uint, w, l uintptr) {
	user32PostThreadMessageW.Call(thread, uintptr(msg), w, l)
}

func PostQuitMessage(code int) { user32PostQuitMessage.Call(uintptr(code)) }

func DispatchMessage(m *Msg) { user32DispatchMessageW.Call(uintptr(unsafe.Pointer(m))) }

func TranslateMessage(m *Msg) { user32TranslateMessage.Call(uintptr(unsafe.Pointer(m))) }

func GetMessage(m *Msg, h uintptr, min, max uint) int {
	r, _, _ := user32GetMessageW.Call(uintptr(unsafe.Pointer(m)), h, uintptr(min), uintptr(max))
	return int(r)
}

func AdjustWindowRect(rect *windows.Rect, style int, menu bool) {
	m := 0
	if menu {
		m = 1
	}
	user32AdjustWindowRect.Call(uintptr(unsafe.Pointer(rect)), uintptr(style), uintptr(m))
}

func DefWindowProc(hwnd, msg, wp, lp uintptr) uintptr {
	r, _, _ := user32DefWindowProcW.Call(hwnd, msg, wp, lp)
	return r
}

func SetWindowLongPtr(h uintptr, index, style int) {
	user32SetWindowLongPtrW.Call(h, uintptr(index), uintptr(style))
}

func GetWindowLongPtr(h uintptr, index int) int {
	style, _, _ := user32GetWindowLongPtrW.Call(h, uintptr(index))
	return int(style)
}

func SetWindowText(h uintptr, title string) error {
	t, err := windows.UTF16FromString(title)
	if err != nil {
		return err
	}

	r, _, err := user32SetWindowTextW.Call(h, uintptr(unsafe.Pointer(&t[0])))
	if r == 0 {
		return err
	}
	return nil
}

func SetWindowPos(h, z uintptr, p webview.Point, size webview.Size, flag uint) {
	user32SetWindowPos.Call(h, z, uintptr(p.X), uintptr(p.Y), uintptr(size.Width), uintptr(size.Height), uintptr(flag))
}

func SetFocus(h uintptr) { user32SetFocus.Call(h) }

func UpdateWindow(h uintptr) { user32UpdateWindow.Call(h) }

func ShowWindow(h uintptr, style int) { user32ShowWindow.Call(h, uintptr(style)) }

func DestroyWindow(h uintptr) { user32DestroyWindow.Call(h) }

func CreateWindowEx(
	exStyle, className, windowName, style uintptr, point webview.Point, size webview.Size, parent, menu, inst, param uintptr) (uintptr, error) {
	ret, _, err := user32CreateWindowExW.Call(
		exStyle, className, windowName, style,
		uintptr(point.X), uintptr(point.Y), uintptr(size.Width), uintptr(size.Height),
		parent, menu, inst, param)

	if ret == 0 {
		return 0, err
	}
	return ret, nil
}

func RegisterClassEx(wc *WndClassExW) error {
	ret, _, err := user32RegisterClassExW.Call(uintptr(unsafe.Pointer(wc)))
	if ret == 0 {
		return err
	}
	return nil
}

func LoadImage(instance uintptr) uintptr {
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-loadimagew

	w := GetSystemMetrics(SystemMetricsCxIcon)
	h := GetSystemMetrics(SystemMetricsCyIcon)
	icon, _, _ := user32LoadImageW.Call(instance, 32512, w, h, 0)
	return icon
}

func GetSystemMetrics(v uintptr) uintptr {
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getsystemmetrics
	// If the function fails, the return value is 0.
	// GetLastError does not provide extended error information.

	r, _, _ := user32GetSystemMetrics.Call(v)
	return r
}

func SHCreateMemStream(data []byte) (uintptr, error) {
	ret, _, err := shlwapiSHCreateMemStream.Call(
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
	)
	if ret == 0 {
		return 0, err
	}

	return ret, nil
}
