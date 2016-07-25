// +build xlib
// generated by stringer -type=xEventType; DO NOT EDIT

package xlib

import "fmt"

const _xEventType_name = "xKeyPressxKeyReleasexButtonPressxButtonReleasexMotionNotifyxEnterNotifyxLeaveNotifyxFocusInxFocusOutxKeymapNotifyxExposexGraphicsExposexNoExposexVisibilityNotifyxCreateNotifyxDestroyNotifyxUnmapNotifyxMapNotifyxMapRequestxReparentNotifyxConfigureNotifyxConfigureRequestxGravityNotifyxResizeRequestxCirculateNotifyxCirculateRequestxPropertyNotifyxSelectionClearxSelectionRequestxSelectionNotifyxColormapNotifyxClientMessagexMappingNotifyxGenericEvent"

var _xEventType_index = [...]uint16{0, 9, 20, 32, 46, 59, 71, 83, 91, 100, 113, 120, 135, 144, 161, 174, 188, 200, 210, 221, 236, 252, 269, 283, 297, 313, 330, 345, 360, 377, 393, 408, 422, 436, 449}

func (i xEventType) String() string {
	i -= 2
	if i < 0 || i+1 >= xEventType(len(_xEventType_index)) {
		return fmt.Sprintf("xEventType(%d)", i+2)
	}
	return _xEventType_name[_xEventType_index[i]:_xEventType_index[i+1]]
}