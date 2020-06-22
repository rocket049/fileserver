// +build minimal

#pragma once

#ifndef GO_QTWIDGETS_H
#define GO_QTWIDGETS_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct QtWidgets_PackedString { char* data; long long len; void* ptr; };
struct QtWidgets_PackedList { void* data; long long len; };

#ifdef __cplusplus
}
#endif

#endif