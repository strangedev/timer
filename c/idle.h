#ifndef _IDLE_H
#define _IDLE_H

#include <X11/Xlib.h>
#include <X11/extensions/scrnsaver.h>

long long getIdleMs() {
  Display *display = XOpenDisplay("");
  if (!display) {
    return -1;
  }

  int error, event;
  if (!XScreenSaverQueryExtension(display, &event, &error)) {
    return -1;
  }

  XScreenSaverInfo info;
  XScreenSaverQueryInfo(display, DefaultRootWindow(display), &info);
  
  return info.idle;
}

#endif