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
	var inst windows.Handle
	if err := windows.GetModuleHandleEx(0, nil, &inst); err != nil {
		return err
	}

	icon := w32.LoadImage(uintptr(inst))

	className, err := windows.UTF16PtrFromString("webview")
	if err != nil {
		return err
	}

	wc := w32.WndClassExW{
		CbSize:        uint32(unsafe.Sizeof(w32.WndClassExW{})),
		HInstance:     inst,
		LpszClassName: className,
		HIcon:         windows.Handle(icon),
		HIconSm:       windows.Handle(icon),
		LpfnWndProc:   windows.NewCallback(wndProc),
	}
	if err := w32.RegisterClassEx(&wc); err != nil {
		return err
	}

	windowName, err := windows.UTF16PtrFromString(o.Title)
	if err != nil {
		return err
	}

	d.hwnd, err = w32.CreateWindowEx(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		w32.WSOverlappedWindow,
		o.Position,
		o.Size,
		0,
		0,
		uintptr(inst),
		0,
	)
	if err != nil {
		return err
	}

	setWindowContext(d.hwnd, d)

	w32.ShowWindow(d.hwnd, w32.SWShow)
	w32.UpdateWindow(d.hwnd)
	w32.SetFocus(d.hwnd)

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
			w32.SetFocus(w.hwnd)
			return w32.DefWindowProc(hwnd, msg, wp, lp)
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
			w32.DestroyWindow(hwnd)
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
			return w32.DefWindowProc(hwnd, msg, wp, lp)
		}
		return 0
	}
	return w32.DefWindowProc(hwnd, msg, wp, lp)
}
