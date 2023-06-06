// SPDX-License-Identifier: MIT

package binds

import (
	"encoding/json"
	"strings"
)

type rpcMessage struct {
	ID     int               `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

func jsString(v string) string {
	return `"` + strings.ReplaceAll(v, "\"", "\\\"") + `"`
}
