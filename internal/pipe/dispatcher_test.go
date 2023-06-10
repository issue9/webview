// SPDX-License-Identifier: MIT

package pipe

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestDispatcher(t *testing.T) {
	a := assert.New(t, false)

	d := NewDispatcher()
	a.Length(d.funcs, 0)

	d.Add(func() {})
	a.Length(d.funcs, 1)

	d.Run()
	a.Length(d.funcs, 0)
}
