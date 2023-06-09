//go:build go1.18 && windows
// +build go1.18,windows

package edge

import "syscall"

// Call calls a COM procedure.
//
//go:uintptrescapes
func (p ComProc) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	// The magic uintptrescapes comment is needed to prevent moving uintptr(unsafe.Pointer(p)) so calls to .Call() also
	// satisfy the unsafe.Pointer rule "(4) Conversion of a Pointer to a uintptr when calling syscall.Syscall."
	// Otherwise it might be that pointers get moved, especially pointer onto the Go stack which might grow dynamically.
	// See https://pkg.go.dev/unsafe#Pointer and https://github.com/golang/go/issues/34474
	return syscall.SyscallN(uintptr(p), a...)
}
