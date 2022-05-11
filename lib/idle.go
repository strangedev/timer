package lib

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lX11 -lXext -lXss
// #include "idle.h"
import "C"
import "time"

func GetIdleTime() (time.Duration, error) {
	idleMs := int64(C.getIdleMs())

	if idleMs < 0 {
		return 0, &X11Error{}
	}

	return time.Duration(idleMs * 1e6), nil
}
