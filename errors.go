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
//
// 绑字的函数其返回类型可以是以下几种类型：
//   - 无返回；
//   - 1 个任意值；
//   - 2 个值，其中第二个返回值必须得是 error 类型；
//
// 其它情况会返回此错误。
func ErrBindFuncReturnInvalid() error { return errBindFuncReturnInvalid }
