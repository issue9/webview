// SPDX-License-Identifier: MIT

package loader

import _ "embed"

//go:embed arm64/WebView2Loader.dll
var WebView2Loader []byte
