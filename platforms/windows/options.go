// SPDX-License-Identifier: MIT

//go:build windows

package windows

import (
	"log"

	"github.com/issue9/webview"
	"github.com/issue9/webview/internal/presets"
	"github.com/issue9/webview/internal/windows/w32"
)

// Options 初始 webview 的选项
type Options struct {
	// Debug 调试模式
	Debug bool

	// DataPath 指定 webview 的数据路径
	DataPath string

	// AutoFocus 当窗口获得焦点时组件自动获取焦点
	AutoFocus bool

	// Title 标题
	Title string

	// Position 初始位置
	Position webview.Point

	// Size 初始大小
	Size webview.Size

	// Style 窗口样式
	Style Style

	// Error 错误日志输出
	//
	// 部分非致命的错误经由此输出，如果为空，则采用 log.Default() 。
	Error *log.Logger
}

type Style = int

// 窗口样式定义，可参考[官方文档]
//
// [官方文档] https://learn.microsoft.com/en-us/windows/win32/winmsg/window-styles
const (
	WSOverlapped       Style = 0x00000000
	WSMaximizeBox            = 0x00010000
	WSThickFrame             = 0x00040000
	WSCaption                = 0x00C00000
	WSSysMenu                = 0x00080000
	WSMinimizeBox            = 0x00020000
	WSOverlappedWindow       = (WSOverlapped | WSCaption | WSSysMenu | WSThickFrame | WSMinimizeBox | WSMaximizeBox)
)

func sanitizeOptions(o *Options) *Options {
	if o == nil {
		o = &Options{}
	}

	if o.Title == "" {
		o.Title = presets.Title
	}

	if o.Position.X == 0 {
		o.Position.X = w32.CW_USEDEFAULT
	}
	if o.Position.Y == 0 {
		o.Position.Y = w32.CW_USEDEFAULT
	}

	if o.Size.Width == 0 {
		o.Size.Width = presets.Width
	}

	if o.Size.Height == 0 {
		o.Size.Height = presets.Height
	}

	if o.Error == nil {
		o.Error = log.Default()
	}

	return o
}
