#include "x11.h"

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <X11/Xlib.h>
#include <X11/extensions/scrnsaver.h>

#define NoDisplay -1
#define NoXScreenSaver -2

int init(Display* display) {
  if (!display) {
    printf("Failed to open the display.\n");

    return NoDisplay;
  }

  int error, event;
  if (!XScreenSaverQueryExtension(display, &event, &error)) {
    printf("The XScreenSaver extension is not present.");

    if (error) {
      printf(" (code %d)\n", error);
      
      return error;
    }

    printf("\n");

    return NoXScreenSaver;
  }

  return Success;
}

Bool getIdleMs(Display* display, int64_t* idleMs) {
  int ok;
  
  XScreenSaverInfo *info = XScreenSaverAllocInfo();
  Window rootWindow = DefaultRootWindow(display);

  ok = XScreenSaverQueryInfo(display, rootWindow, info);
  if (!ok) {
    printf("Failed to perform the XScreenSaver query.\n");
    
    XFree(info);
    return false;
  }
  
  *idleMs = info->idle;
  
  XFree(info);
  return true;
}
