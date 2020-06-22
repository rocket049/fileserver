// +build minimal

#pragma once

#ifndef GO_QTGUI_H
#define GO_QTGUI_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct QtGui_PackedString { char* data; long long len; void* ptr; };
struct QtGui_PackedList { void* data; long long len; };

#ifdef __cplusplus
}
#endif

#endif