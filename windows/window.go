// SPDX-License-Identifier: MIT

//go:build windows

package windows

import (
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/issue9/webview/internal/windows/w32"
)

var windowContext = &sync.Map{}

func getWindowContext(wnd uintptr) (*desktop, bool) {
	if x, found := windowContext.Load(wnd); !found {
		return nil, false
	} else {
		d, ok := x.(*desktop)
		return d, ok
	}
}

func setWindowContext(wnd uintptr, data interface{}) {
	windowContext.Store(wnd, data)
}

func (d *desktop) CreateWithOptions(o *Options) error {
	var hinstance windows.Handle
	_ = windows.GetModuleHandleEx(0, nil, &hinstance)

	icow, _, _ := w32.User32GetSystemMetrics.Call(w32.SystemMetricsCxIcon)
	icoh, _, _ := w32.User32GetSystemMetrics.Call(w32.SystemMetricsCyIcon)
	icon, _, _ := w32.User32LoadImageW.Call(uintptr(hinstance), 32512, icow, icoh, 0)

	className, _ := windows.UTF16PtrFromString("webview")
	wc := w32.WndClassExW{
		CbSize:        uint32(unsafe.Sizeof(w32.WndClassExW{})),
		HInstance:     hinstance,
		LpszClassName: className,
		HIcon:         windows.Handle(icon),
		HIconSm:       windows.Handle(icon),
		LpfnWndProc:   windows.NewCallback(wndProc),
	}
	_, _, _ = w32.User32RegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	windowName, _ := windows.UTF16PtrFromString(o.Title)

	d.hwnd, _, _ = w32.User32CreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		0xCF0000, // WS_OVERLAPPEDWINDOW
		uintptr(o.Position.X),
		uintptr(o.Position.Y),
		uintptr(o.Size.Width),
		uintptr(o.Size.Height),
		0,
		0,
		uintptr(hinstance),
		0,
	)
	setWindowContext(d.hwnd, d)

	w32.User32ShowWindow.Call(d.hwnd, w32.SWShow)
	w32.User32UpdateWindow.Call(d.hwnd)
	w32.User32SetFocus.Call(d.hwnd)

	if err := d.chromium.Embed(d.hwnd); err != nil {
		return err
	}

	d.chromium.Resize()
	return nil
}

func wndProc(hwnd, msg, wp, lp uintptr) uintptr {
	if w, ok := getWindowContext(hwnd); ok {
		switch msg {
		case w32.WMMove, w32.WMMoving:
			_ = w.chromium.NotifyParentWindowPositionChanged()
		case w32.WMNCLButtonDown:
			_, _, _ = w32.User32SetFocus.Call(w.hwnd)
			r, _, _ := w32.User32DefWindowProcW.Call(hwnd, msg, wp, lp)
			return r
		case w32.WMSize:
			w.chromium.Resize()
		case w32.WMActivate:
			if wp == w32.WAInactive {
				break
			}
			if w.autofocus {
				w.chromium.Focus()
			}
		case w32.WMClose:
			_, _, _ = w32.User32DestroyWindow.Call(hwnd)
		case w32.WMDestroy:
			w.Close()
		case w32.WMGetMinMaxInfo:
			lpmmi := (*w32.MinMaxInfo)(unsafe.Pointer(lp))
			if w.maxSize.Width > 0 && w.maxSize.Height > 0 {
				lpmmi.PtMaxSize = w32.Point{X: int32(w.maxSize.Width), Y: int32(w.maxSize.Height)}
				lpmmi.PtMaxTrackSize = w32.Point{X: int32(w.maxSize.Width), Y: int32(w.maxSize.Height)}
			}
			if w.minSize.Width > 0 && w.minSize.Height > 0 {
				lpmmi.PtMinTrackSize = w32.Point{X: int32(w.minSize.Width), Y: int32(w.minSize.Height)}
			}
		default:
			r, _, _ := w32.User32DefWindowProcW.Call(hwnd, msg, wp, lp)
			return r
		}
		return 0
	}
	r, _, _ := w32.User32DefWindowProcW.Call(hwnd, msg, wp, lp)
	return r
}
