// SPDX-License-Identifier: MIT

// Package webviewtest 用于 webview 测试
package webviewtest

import (
	"log"

	"github.com/issue9/webview"
)

// Desktop 测试 webview.Desktop 接口的实现
//
// NOTE: 部分系统需要在 TestMain 方法执行才行。
func Desktop(w webview.Desktop) {
	// BUG(caixw): 以下代码会让程序暂停
	//time.Sleep(time.Second)

	w.SetTitle("Hello")

	w.SetSize(webview.Size{Width: 200, Height: 500}, webview.HintNone)

	app(w)
}

// App 测试 webview.App 接口的实现
func App(w webview.App) {
	app(w)
}

// 代码主要来源 https://github.com/webview/webview/blob/899018ad0e5cc22a18cd734393ccae4d55e3b2b4/webview_test.go#L10
func app(w webview.App) {
	w.Bind("noop", func() string {
		log.Println("hello")
		return "hello"
	})
	w.Bind("add", func(a, b int) int {
		return a + b
	})
	w.Bind("quit", func() {
		w.Close()
	})
	w.SetHTML(`<!doctype html>
		<html>
			<body>hello</body>
			<script>
				window.onload = function() {
					document.body.innerText = ` + "`hello, ${navigator.userAgent}`" + `;
					noop().then(function(res) {
						console.log('noop res', res);
						add(1, 2).then(function(res) {
							console.log('add res', res);
							quit();
						});
					});
				};
			</script>
		</html>
	)`)
	w.Run()
}
