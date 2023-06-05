// SPDX-License-Identifier: MIT

//go:build windows

package windows

import (
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
}

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

	return o
}
