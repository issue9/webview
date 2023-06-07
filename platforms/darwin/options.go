// SPDX-License-Identifier: MIT

package darwin

import (
	"log"

	"github.com/issue9/webview"
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
