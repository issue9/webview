// SPDX-License-Identifier: MIT

package webview

// Renderer 基于 webview 渲染的基本接口
type Renderer interface {
	// 直接将内容设置为 HTML
	SetHTML(html string)

	// 加载指定地址的页面
	//
	// url 可以是本地或是网络地址
	Load(url string)
}

// App 基于 webview 应用的基本接口
type App interface {
	Renderer

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

// Desktop 基于 webview 桌面应用的接口
type Desktop interface {
	App

	// 标题
	Title() string
	SetTitle(string)

	// 窗口大小
	Size() Size
	SetSize(Size, Hint)

	// 窗口位置
	Position() Point
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
	HintFixed
	HintMin
	HintMax
)
