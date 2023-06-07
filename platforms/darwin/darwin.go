// SPDX-License-Identifier: MIT

//go:build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit
#include "darwin.h"
*/
import "C"
import (
	"unsafe"

	"github.com/issue9/webview"
)

type desktop struct {
	title    string
	position webview.Point
	size     webview.Size

	wv *C.CocoaWebView
}

func New(o *Options) webview.Desktop {
	t := C.CString(o.Title)
	defer C.free(unsafe.Pointer(t))

	return &desktop{
		title:    o.Title,
		position: o.Position,
		size:     o.Size,
		wv:       C.create_cocoa(C.double(o.Position.X), C.double(o.Position.Y), C.double(o.Size.Width), C.double(o.Size.Height), t),
	}
}

func (d *desktop) SetHTML(html string) {
	t := C.CString(html)
	defer C.free(unsafe.Pointer(t))
	C.set_html(d.wv, t)
}

func (d *desktop) Load(url string) {
	t := C.CString(url)
	defer C.free(unsafe.Pointer(t))
	C.load(d.wv, t)
}

func (d *desktop) OnLoad(js string) {
	// TODO
}

func (d *desktop) Eval(js string) {
	t := C.CString(js)
	defer C.free(unsafe.Pointer(t))
	C.eval(d.wv, t)
}

func (d *desktop) Bind(name string, f interface{}) error {
	// TODO
	return nil
}

func (d *desktop) Run() {
	C.run()
}

func (d *desktop) Close() {
	C.terminate()
}

func (d *desktop) Title() string { return d.title }

func (d *desktop) SetTitle(title string) {
	t := C.CString(title)
	defer C.free(unsafe.Pointer(t))
	C.set_title(d.wv, t)
}

func (d *desktop) Position() webview.Point { return d.position }

func (d *desktop) SetPosition(p webview.Point) {
	C.set_position(d.wv, C.double(p.X), C.double(p.Y))
	d.position = p
}

func (d *desktop) Size() webview.Size { return d.size }

func (d *desktop) SetSize(s webview.Size, h webview.Hint) {
	switch h {
	case webview.HintFixed:
	case webview.HintMax:
		C.set_max_size(d.wv, C.double(s.Width), C.double(s.Height))
	case webview.HintMin:
		C.set_min_size(d.wv, C.double(s.Width), C.double(s.Height))
	default: // webview.HintNone
		p := d.Position()
		C.set_frame(d.wv, true, C.double(p.X), C.double(p.Y), C.double(s.Width), C.double(s.Height))
		d.size = s
	}
}
