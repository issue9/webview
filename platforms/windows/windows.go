// SPDX-License-Identifier: MIT

//go:build windows

// Package windows windows 平台实现
//
// https://learn.microsoft.com/zh-cn/microsoft-edge/webview2/concepts/overview-features-apis?tabs=win32cpp
package windows

import (
	"log"
	"sync"

	"golang.org/x/sys/windows"

	"github.com/issue9/webview"
	"github.com/issue9/webview/internal/binds"
	"github.com/issue9/webview/internal/windows/edge"
	"github.com/issue9/webview/internal/windows/w32"
)

type desktop struct {
	hwnd       uintptr
	mainThread uintptr
	chromium   *edge.Chromium
	title      string
	position   webview.Point
	size       webview.Size
	maxSize    webview.Size
	minSize    webview.Size
	autofocus  bool
	errlog     *log.Logger

	binds     *binds.Binds
	m         sync.Mutex
	dispatchq []func()
}

func New(o *Options) (webview.Desktop, error) {
	o = sanitizeOptions(o)

	d := &desktop{
		mainThread: uintptr(windows.GetCurrentThreadId()),
		position:   o.Position,
		size:       o.Size,
		autofocus:  o.AutoFocus,
		errlog:     o.Error,
	}

	d.binds = binds.New(d)

	chromium := edge.NewChromium()
	chromium.MessageCallback = d.binds.MessageHandler
	chromium.DataPath = o.DataPath
	chromium.SetPermission(edge.CoreWebView2PermissionKindClipboardRead, edge.CoreWebView2PermissionStateAllow)
	d.chromium = chromium

	if err := d.createWindow(o); err != nil {
		return nil, err
	}

	settings, err := chromium.GetSettings()
	if err != nil {
		return nil, err
	}

	if err = settings.PutAreDefaultContextMenusEnabled(o.Debug); err != nil {
		return nil, err
	}

	if err = settings.PutAreDevToolsEnabled(o.Debug); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *desktop) Load(url string) { d.chromium.Navigate(url) }

func (d *desktop) SetHTML(html string) { d.chromium.NavigateToString(html) }

func (d *desktop) Run() {
	var msg w32.Msg
	for {
		w32.GetMessage(&msg, 0, 0, 0)
		if msg.Message == w32.WMApp {
			d.m.Lock()
			q := append([]func(){}, d.dispatchq...)
			d.dispatchq = []func(){}
			d.m.Unlock()
			for _, v := range q {
				v()
			}
		} else if msg.Message == w32.WMQuit {
			return
		}
		r := w32.GetAncestor(uintptr(msg.Hwnd), w32.GARoot)
		if w32.IsDialogMessage(r, &msg) {
			continue
		}
		w32.TranslateMessage(&msg)
		w32.DispatchMessage(&msg)
	}
}

func (d *desktop) Close() { w32.PostQuitMessage(0) }

func (d *desktop) OnLoad(js string) { d.chromium.Init(js) }

func (d *desktop) Eval(js string) { d.chromium.Eval(js) }

func (d *desktop) Dispatch(f func()) {
	d.m.Lock()
	d.dispatchq = append(d.dispatchq, f)
	d.m.Unlock()
	w32.PostThreadMessage(d.mainThread, w32.WMApp, 0, 0)
}

func (d *desktop) Bind(name string, f interface{}) error {
	return d.binds.Bind(name, f)
}

func (d *desktop) Title() string { return d.title }

func (d *desktop) SetTitle(title string) {
	if err := w32.SetWindowText(d.hwnd, title); err != nil {
		d.errlog.Println(err)
	} else {
		d.title = title
	}
}

func (d *desktop) Size() webview.Size { return d.size }

func (d *desktop) SetSize(s webview.Size, hint webview.Hint) {
	index := w32.GWLStyle
	style := w32.GetWindowLongPtr(d.hwnd, index)
	if hint == webview.HintFixed {
		style &^= (w32.WSThickFrame | w32.WSMaximizeBox)
	} else {
		style |= (w32.WSThickFrame | w32.WSMaximizeBox)
	}
	w32.SetWindowLongPtr(d.hwnd, index, style)

	if hint == webview.HintMax {
		d.maxSize = s
	} else if hint == webview.HintMin {
		d.minSize = s
	} else {
		p := d.Position()
		r := windows.Rect{
			Left:   int32(p.X),
			Top:    int32(p.Y),
			Right:  int32(s.Width + p.X),
			Bottom: int32(s.Height + p.Y),
		}
		w32.AdjustWindowRect(&r, w32.WSOverlappedWindow, false)
		w32.SetWindowPos(d.hwnd, 0, p, s, w32.SWPNoZOrder|w32.SWPNoActivate|w32.SWPNoMove|w32.SWPFrameChanged)
		d.chromium.Resize()
		d.size = s // 保存 size
	}
}

func (d *desktop) Position() webview.Point { return d.position }

func (d *desktop) SetPosition(p webview.Point) {
	w32.SetWindowPos(d.hwnd, 0, p, d.Size(), w32.SWPNoZOrder|w32.SWPNoActivate|w32.SWPNoMove|w32.SWPFrameChanged)
	d.position = p
}
