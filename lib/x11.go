package lib

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lX11 -lXext -lXss
// #include "x11.h"
import "C"
import (
	"time"
)

type X11 struct {
	Display *C.Display
}

func NewX11() X11 {
	return X11{
		Display: C.XOpenDisplay(C.CString("")),
	}
}

func (x *X11) Init() error {
	if errno := C.init(x.Display); errno != C.Success {
		return &X11Error{}
	}

	return nil
}

func (x *X11) GetIdleTime() (time.Duration, error) {
	var idleMs C.int64_t
	if ok := C.getIdleMs(x.Display, &idleMs); ok == 0 {
		return 0, &X11Error{}
	}

	return time.Duration(idleMs * 1e6), nil
}
