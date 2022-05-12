#ifndef _IDLE_H
#define _IDLE_H

#include <stdbool.h>
#include <stdint.h>
#include <X11/Xlib.h>

int init(Display* display);

Bool getIdleMs(Display* display, int64_t* idleMs);

#endif
