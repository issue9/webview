// SPDX-License-Identifier: MIT

//go:build linux || openbsd || freebsd || netbsd

// Package gtk GTK 平台实现
package gtk

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk.h"
*/
import "C"
import (
	"unsafe"

	"github.com/issue9/webview"
	"github.com/issue9/webview/internal/pipe"
)

var (
	dispatcher = pipe.NewDispatcher()
	binder     *pipe.Binder
)

type desktop struct {
	title    string
	position webview.Point
	size     webview.Size

	app *C.App
}

func New(o *Options) webview.Desktop {
	o = sanitizeOptions(o)
	p := o.Position
	s := o.Size

	title := C.CString(o.Title)
	defer C.free(unsafe.Pointer(title))

	d := &desktop{
		title:    o.Title,
		position: o.Position,
		size:     o.Size,

		app: C.create_gtk(C._Bool(o.Debug), C.int(p.X), C.int(p.Y), C.int(s.Width), C.int(s.Height), title),
	}

	binder = pipe.NewBinder(d, o.Error)

	return d
}

func (d *desktop) SetHTML(html string) {
	t := C.CString(html)
	defer C.free(unsafe.Pointer(t))
	C.load_html(d.app, t)
}

func (d *desktop) Load(url string) {
	t := C.CString(url)
	defer C.free(unsafe.Pointer(t))
	C.load(d.app, t)
}

func (d *desktop) OnLoad(js string) {
	t := C.CString(js)
	defer C.free(unsafe.Pointer(t))
	C.add_script(d.app, t)
}

func (d *desktop) Eval(js string) {
	t := C.CString(js)
	defer C.free(unsafe.Pointer(t))
	C.eval(d.app, t)
}

func (d *desktop) Bind(name string, f interface{}) error {
	return binder.Bind(name, f)
}

func (d *desktop) Dispatch(f func()) {
	dispatcher.Add(f)
	C.dispatch()
}

func (d *desktop) Run() {
	C.run(d.app)
}

func (d *desktop) Close() {
	C.quit(d.app)
}

func (d *desktop) Title() string { return d.title }

func (d *desktop) SetTitle(t string) {
	s := C.CString(t)
	defer C.free(unsafe.Pointer(s))
	C.set_title(d.app, s)
}

func (d *desktop) Position() webview.Point { return d.position }

func (d *desktop) SetPosition(p webview.Point) {
	C.move(d.app, C.int(p.X), C.int(p.Y))
	d.position = p
}

func (d *desktop) Size() webview.Size { return d.size }

func (d *desktop) SetSize(s webview.Size, h webview.Hint) {
	switch h {
	case webview.HintFixed:
		C.set_fixed_size(d.app, C.int(s.Width), C.int(s.Height))
	case webview.HintMax:
		C.set_max_size(d.app, C.int(s.Width), C.int(s.Height))
	case webview.HintMin:
		C.set_min_size(d.app, C.int(s.Width), C.int(s.Height))
	default: // webview.HintNone
		C.set_size(d.app, C.int(s.Width), C.int(s.Height))
		d.size = s
	}
}

//export dispatchCallback
func dispatchCallback() {
	dispatcher.Run()
}

//export messageCallback
func messageCallback(msg *C.char) {
	binder.MessageHandler(C.GoString(msg))
}
