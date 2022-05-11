#include "idle.h"

#include <stdio.h>
#include <X11/Xlib.h>
#include <X11/extensions/scrnsaver.h>

long long getIdleMs() {
  Display *display = XOpenDisplay("");
  if (!display) {
    printf("Failed to open the display.\n");

    return -1;
  }

  int error, event;
  if (!XScreenSaverQueryExtension(display, &event, &error)) {
    printf("The XScreenSaver extension is not present.\n");

    return -1;
  }

  XScreenSaverInfo info;
  XScreenSaverQueryInfo(display, DefaultRootWindow(display), &info);
  
  return info.idle;
}