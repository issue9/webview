// SPDX-License-Identifier: MIT

//go:build windows

// Package w32 将用到的 win32 API 以及相关数据结构封装为 Go 模式
package w32

import (
	"syscall"

	"golang.org/x/sys/windows"
)

const (
	CW_USEDEFAULT = 0x80000000
)

const (
	LR_DEFAULTCOLOR     = 0x0000
	LR_MONOCHROME       = 0x0001
	LR_LOADFROMFILE     = 0x0010
	LR_LOADTRANSPARENT  = 0x0020
	LR_DEFAULTSIZE      = 0x0040
	LR_VGACOLOR         = 0x0080
	LR_LOADMAP3DCOLORS  = 0x1000
	LR_CREATEDIBSECTION = 0x2000
	LR_SHARED           = 0x8000
)

const (
	SystemMetricsCxIcon = 11
	SystemMetricsCyIcon = 12
)

// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-showwindow
const (
	SWHide = 0
	SWShow = 5
)

// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowpos
const (
	SWPNoZOrder     = 0x0004
	SWPNoActivate   = 0x0010
	SWPNoMove       = 0x0002
	SWPFrameChanged = 0x0020
)

const (
	WMDestroy       = 0x0002
	WMMove          = 0x0003
	WMSize          = 0x0005
	WMActivate      = 0x0006
	WMClose         = 0x0010
	WMQuit          = 0x0012
	WMGetMinMaxInfo = 0x0024
	WMNCLButtonDown = 0x00A1
	WMMoving        = 0x0216
	WMApp           = 0x8000
)

const (
	GAParent    = 1
	GARoot      = 2
	GARootOwner = 3
)

const (
	GWLStyle = -16
)

// https://learn.microsoft.com/en-us/windows/win32/winmsg/window-styles
const (
	WSOverlapped       = 0x00000000
	WSMaximizeBox      = 0x00010000
	WSThickFrame       = 0x00040000
	WSCaption          = 0x00C00000
	WSSysMenu          = 0x00080000
	WSMinimizeBox      = 0x00020000
	WSOverlappedWindow = (WSOverlapped | WSCaption | WSSysMenu | WSThickFrame | WSMinimizeBox | WSMaximizeBox)
)

const (
	WAInactive    = 0
	WAActive      = 1
	WAActiveClick = 2
)

type WndClassExW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CnClsExtra    int32
	CbWndExtra    int32
	HInstance     windows.Handle
	HIcon         windows.Handle
	HCursor       windows.Handle
	HbrBackground windows.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       windows.Handle
}

type MinMaxInfo struct {
	PtReserved     Point
	PtMaxSize      Point
	PtMaxPosition  Point
	PtMinTrackSize Point
	PtMaxTrackSize Point
}

type Point struct {
	X, Y int32
}

type Msg struct {
	Hwnd     syscall.Handle
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       Point
	LPrivate uint32
}
