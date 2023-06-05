// SPDX-License-Identifier: MIT

package webview

import "errors"

var (
	errOnlyFuncCanBound      = errors.New("only functions can be bound")
	errBindFuncReturnInvalid = errors.New("bind function may only return a value or value+error")
)

// ErrOnlyFuncCanBound 表示绑定的对象不是方法
func ErrOnlyFuncCanBound() error { return errOnlyFuncCanBound }

// ErrBindFuncReturnInvalid 表示绑定方法的返回值类型不符合要求
func ErrBindFuncReturnInvalid() error { return errBindFuncReturnInvalid }
