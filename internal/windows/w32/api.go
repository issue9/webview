// SPDX-License-Identifier: MIT

//go:build windows

package w32

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	shlwapi                  = windows.NewLazySystemDLL("shlwapi")
	shlwapiSHCreateMemStream = shlwapi.NewProc("SHCreateMemStream")

	user32                 = windows.NewLazySystemDLL("user32")
	user32GetSystemMetrics = user32.NewProc("GetSystemMetrics")
	User32LoadImageW       = user32.NewProc("LoadImageW")
)

func LoadImage(instance uintptr) uintptr {
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-loadimagew

	w := GetSystemMetrics(SystemMetricsCxIcon)
	h := GetSystemMetrics(SystemMetricsCyIcon)
	icon, _, _ := User32LoadImageW.Call(instance, 32512, w, h, 0)
	return icon
}

func GetSystemMetrics(v uintptr) uintptr {
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getsystemmetrics
	// If the function fails, the return value is 0.
	// GetLastError does not provide extended error information.

	r, _, _ := user32GetSystemMetrics.Call(v)
	return r
}

func SHCreateMemStream(data []byte) (uintptr, error) {
	ret, _, err := shlwapiSHCreateMemStream.Call(
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
	)
	if ret == 0 {
		return 0, err
	}

	return ret, nil
}
