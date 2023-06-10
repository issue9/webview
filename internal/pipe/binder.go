// SPDX-License-Identifier: MIT

package pipe

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strconv"
	"sync"

	"github.com/issue9/webview"
)

type rpcMessage struct {
	ID     int               `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

type Binder struct {
	m        sync.Mutex
	bindings map[string]interface{}
	app      webview.App
	errlog   *log.Logger
}

func NewBinder(app webview.App, errlog *log.Logger) *Binder {
	return &Binder{
		m:        sync.Mutex{},
		bindings: make(map[string]interface{}, 100),
		app:      app,
		errlog:   errlog,
	}
}

// Bind 将 f 以 name 名称绑定在 webview 上
func (b *Binder) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return webview.ErrOnlyFuncCanBound()
	}

	t := v.Type()
	if n := t.NumOut(); n > 2 {
		return webview.ErrBindFuncReturnInvalid()
	} else if n == 2 && !t.Out(1).Implements(errorType) { // 两个参数的第二个必须为 error
		return webview.ErrBindFuncReturnInvalid()
	}

	b.m.Lock()
	b.bindings[name] = f
	b.m.Unlock()

	b.app.OnLoad("(function() { var name = " + jsString(name) + ";" + `
		var RPC = window._rpc = (window._rpc || {nextSeq: 1});
		window[name] = function() {
		  var seq = RPC.nextSeq++;
		  var promise = new Promise(function(resolve, reject) {
			RPC[seq] = {
			  resolve: resolve,
			  reject: reject,
			};
		  });
		  window.external.invoke(JSON.stringify({
			id: seq,
			method: name,
			params: Array.prototype.slice.call(arguments),
		  }));
		  return promise;
		}
	})()`)

	return nil
}

// 调用指定名称的方法
func (b *Binder) call(name string, params ...json.RawMessage) (interface{}, error) {
	b.m.Lock()
	f, ok := b.bindings[name]
	b.m.Unlock()
	if !ok {
		return nil, nil
	}

	v := reflect.ValueOf(f)
	isVariadic := v.Type().IsVariadic()
	numIn := v.Type().NumIn()
	if (isVariadic && len(params) < numIn-1) || (!isVariadic && len(params) != numIn) {
		return nil, errors.New("function arguments mismatch")
	}
	args := []reflect.Value{}
	for i := range params {
		var arg reflect.Value
		if isVariadic && i >= numIn-1 {
			arg = reflect.New(v.Type().In(numIn - 1).Elem())
		} else {
			arg = reflect.New(v.Type().In(i))
		}
		if err := json.Unmarshal(params[i], arg.Interface()); err != nil {
			return nil, err
		}
		args = append(args, arg.Elem())
	}

	res := v.Call(args)
	switch len(res) {
	case 0:
		return nil, nil
	case 1: // One result may be a value, or an error
		if res[0].Type().Implements(errorType) {
			if res[0].Interface() != nil {
				return nil, res[0].Interface().(error)
			}
			return nil, nil
		}
		return res[0].Interface(), nil

	case 2: // Two results: first one is value, second is error
		if !res[1].Type().Implements(errorType) {
			panic("返回的第二个参数只能是 error 类型") // 由 Binds.Bind 确保不会发生此错误
		}
		if res[1].Interface() == nil {
			return res[0].Interface(), nil
		}
		return res[0].Interface(), res[1].Interface().(error)

	default:
		panic("返回参数最多只能有两个") // 由 Binds.Bind 确保不会发生此错误
	}
}

// MessageHandler 处理前端的调用请求
func (b *Binder) MessageHandler(msg string) {
	rpc := rpcMessage{}
	if err := json.Unmarshal([]byte(msg), &rpc); err != nil {
		b.errlog.Printf("invalid RPC message %v", err)
		return
	}

	id := strconv.Itoa(rpc.ID)
	if res, err := b.call(rpc.Method, rpc.Params...); err != nil {
		b.app.Dispatch(func() {
			b.app.Eval("window._rpc[" + id + "].reject(" + jsString(err.Error()) + "); window._rpc[" + id + "] = undefined")
		})
	} else if data, err := json.Marshal(res); err != nil {
		b.app.Dispatch(func() {
			b.app.Eval("window._rpc[" + id + "].reject(" + jsString(err.Error()) + "); window._rpc[" + id + "] = undefined")
		})
	} else {
		b.app.Dispatch(func() {
			b.app.Eval("window._rpc[" + id + "].resolve(" + string(data) + "); window._rpc[" + id + "] = undefined")
		})
	}
}
