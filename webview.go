// SPDX-License-Identifier: MIT

package webview

// WebView 定义了基于 webview 应用的必要接口
type WebView interface {
	// 直接将内容设置为 HTML
	SetHTML(html string)

	// 加载指定地址的页面
	//
	// url 可以是本地或是网络地址
	Load(url string)

	// 新页面加载时执行的 JS
	OnLoad(js string)

	// 计算 JS 结果并返回
	Eval(js string)

	// Bind 绑定方法至前端
	//
	// f 必须是一个函数，反加值可以是单个值，或是两值，如果是两个值，那么其第二个必须得是 error。
	Bind(name string, f interface{}) error

	// 切换界面语言
	//
	// 比如右键菜单等
	//SetLocale(string)

	// Run 运行程序
	Run()

	// 关闭服务
	Close()
}
