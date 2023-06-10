// SPDX-License-Identifier: MIT

//go:build darwin

package darwin

import (
	"log"

	"github.com/issue9/webview"
	"github.com/issue9/webview/internal/presets"
)

// Options 初始 webview 的选项
type Options struct {
	// Debug 调试模式
	Debug bool

	// Title 标题
	Title string

	// Position 初始位置
	Position webview.Point

	// Size 初始大小
	Size webview.Size

	// Error 错误日志输出
	//
	// 部分非致命的错误经由此输出，如果为空，则采用 log.Default() 。
	Error *log.Logger
}

func sanitizeOptions(o *Options) *Options {
	if o == nil {
		o = &Options{}
	}

	if o.Title == "" {
		o.Title = presets.Title
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
