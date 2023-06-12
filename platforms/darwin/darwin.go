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
	"runtime"
	"unsafe"

	"github.com/issue9/webview"
	"github.com/issue9/webview/internal/pipe"
)

func init() {
	runtime.LockOSThread()
}

var (
	dispatcher = pipe.NewDispatcher()
	binder     *pipe.Binder
)

type desktop struct {
	title    string
	position webview.Point
	size     webview.Size
	app      *C.App
}

func New(o *Options) webview.Desktop {
	o = sanitizeOptions(o)

	t := C.CString(o.Title)
	defer C.free(unsafe.Pointer(t))

	wv := C.create_cocoa(C._Bool(o.Debug), C.double(o.Position.X), C.double(o.Position.Y), C.double(o.Size.Width), C.double(o.Size.Height), t)
	d := &desktop{
		title:    o.Title,
		position: o.Position,
		size:     o.Size,
		app:      wv,
	}
	binder = pipe.NewBinder(d, o.Error)

	return d
}

func (d *desktop) SetHTML(html string) {
	t := C.CString(html)
	defer C.free(unsafe.Pointer(t))
	C.set_html(d.app, t)
}

func (d *desktop) Load(url string) {
	t := C.CString(url)
	defer C.free(unsafe.Pointer(t))
	C.load(d.app, t)
}

func (d *desktop) OnLoad(js string) {
	t := C.CString(js)
	defer C.free(unsafe.Pointer(t))
	C.add_user_script(d.app, t)
}

func (d *desktop) Eval(js string) {
	t := C.CString(js)
	defer C.free(unsafe.Pointer(t))
	C.eval(d.app, t)
}

func (d *desktop) Bind(name string, f interface{}) error {
	binder.Bind(name, f)
	return nil
}

func (d *desktop) Dispatch(f func()) {
	dispatcher.Add(f)
	C.dispatch()
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
	C.set_title(d.app, t)
}

func (d *desktop) Position() webview.Point { return d.position }

func (d *desktop) SetPosition(p webview.Point) {
	C.set_position(d.app, C.double(p.X), C.double(p.Y))
	d.position = p
}

func (d *desktop) Size() webview.Size { return d.size }

func (d *desktop) SetSize(s webview.Size, h webview.Hint) {
	switch h {
	case webview.HintFixed:
		p := d.Position()
		C.set_fixed_size(d.app, C.double(p.X), C.double(p.Y), C.double(s.Width), C.double(s.Height))
	case webview.HintMax:
		C.set_max_size(d.app, C.double(s.Width), C.double(s.Height))
	case webview.HintMin:
		C.set_min_size(d.app, C.double(s.Width), C.double(s.Height))
	default: // webview.HintNone
		p := d.Position()
		C.set_frame(d.app, C.double(p.X), C.double(p.Y), C.double(s.Width), C.double(s.Height))
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
