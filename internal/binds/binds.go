// SPDX-License-Identifier: MIT

// Package binds 后端与前端的绑定功能
package binds

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strconv"
	"sync"

	"github.com/issue9/webview"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

type Binds struct {
	m        sync.Mutex
	bindings map[string]interface{}
	app      webview.App
}

func New(app webview.App) *Binds {
	return &Binds{
		m:        sync.Mutex{},
		bindings: make(map[string]interface{}, 100),
		app:      app,
	}
}

// Bind 将 f 以 name 名称绑定在 webview 上
func (b *Binds) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return webview.ErrOnlyFuncCanBound()
	}
	if n := v.Type().NumOut(); n > 2 {
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
func (b *Binds) call(name string, params ...json.RawMessage) (interface{}, error) {
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
			return nil, errors.New("second return value must be an error")
		}
		if res[1].Interface() == nil {
			return res[0].Interface(), nil
		}
		return res[0].Interface(), res[1].Interface().(error)

	default:
		return nil, errors.New("unexpected number of return values")
	}
}

// MessageHandler 处理前端的调用请求
func (b *Binds) MessageHandler(msg string) {
	rpc := rpcMessage{}
	if err := json.Unmarshal([]byte(msg), &rpc); err != nil {
		log.Printf("invalid RPC message: %v", err)
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
