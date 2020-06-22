// +build minimal

#pragma once

#ifndef GO_QTCORE_H
#define GO_QTCORE_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct QtCore_PackedString { char* data; long long len; void* ptr; };
struct QtCore_PackedList { void* data; long long len; };

#ifdef __cplusplus
}
#endif

#endif