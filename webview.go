// SPDX-License-Identifier: MIT

package webview

// App 基于 webview 应用的基本接口
type App interface {
	// SetHTML 直接将内容设置为 HTML
	SetHTML(html string)

	// Load 加载指定地址的页面
	//
	// url 可以是本地或是网络地址
	Load(url string)

	// OnLoad 新页面加载时执行的 JS
	OnLoad(js string)

	// Bind 绑定方法至前端
	//
	// f 必须是一个函数，反加值可以是单个值，或是两值，如果是两个值，那么其第二个必须得是 error。
	Bind(name string, f interface{}) error

	// Run 运行程序
	Run()

	// Close 关闭服务
	Close()
}

// Desktop 基于 webview 桌面应用的接口
type Desktop interface {
	App

	//Title 获取标题
	Title() string

	// SetTitle 设置窗口标题
	SetTitle(string)

	// Size 获取窗口大小
	Size() Size

	// SetSize 调整窗口的大小
	SetSize(Size, Hint)

	// Position 获取窗口位置
	Position() Point

	// SetPosition 移动窗口的位置
	SetPosition(Point)
}

type Point struct {
	X, Y int
}

type Size struct {
	Width, Height int
}

type Hint int8

const (
	HintNone Hint = iota
	HintMin
	HintMax
)
