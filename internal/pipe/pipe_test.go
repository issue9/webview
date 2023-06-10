// SPDX-License-Identifier: MIT

package pipe

import (
	"encoding/json"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestJSString(t *testing.T) {
	a := assert.New(t, false)

	val := "abc"
	b, err := json.Marshal(val)
	a.NotError(err).Equal(string(b), jsString(val))

	val = "abc\""
	b, err = json.Marshal(val)
	a.NotError(err).Equal(string(b), jsString(val))

	val = "abc'"
	b, err = json.Marshal(val)
	a.NotError(err).Equal(string(b), jsString(val))
}
