webview
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/issue9/webview/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/webview)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/webview)](https://pkg.go.dev/github.com/issue9/webview)
![Go version](https://img.shields.io/github/go-mod/go-version/issue9/webview)
======

基于 webview 技术的应用开发框架，目前支持以下平台：

- darwin：支持 macOS 10.13
- windows：采用 webview2，支持 windows 10、windows 11。
- GTK：所有支持 GTK 的平台，需要安装 GTK 3 和 WebKit2GTK 2.22 以上版本

windows 相关代码主要来自 [go-webview2](https://github.com/jchv/go-webview2)，
GTK 和 windows 则参考了 [webview](https://github.com/webview/webview)

安装
----

```shell
go get github.com/issue9/webview
```

版权
----

本项目采用 [MIT](http://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
