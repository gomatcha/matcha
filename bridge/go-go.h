// +build matcha

#ifndef GO_GO_H
#define GO_GO_H

#include "go-foreign.h"
#include <stdbool.h>

bool matchaGoIsNil(GoRef);
GoRef matchaGoBool(bool);
bool matchaGoToBool(GoRef);
GoRef matchaGoInt(int);
GoRef matchaGoInt64(int64_t);
int64_t matchaGoToInt64(GoRef);
GoRef matchaGoUint64(uint64_t);
uint64_t matchaGoToUint64(GoRef);
GoRef matchaGoFloat64(double);
double matchaGoToFloat64(GoRef);
GoRef matchaGoString(CGoBuffer);
CGoBuffer matchaGoToString(GoRef);
GoRef matchaGoBytes(CGoBuffer);
CGoBuffer matchaGoToBytes(GoRef);
GoRef matchaGoArray(CGoBuffer);
CGoBuffer matchaGoToArray(GoRef);
GoRef matchaGoForeign(FgnRef);
FgnRef matchaGoToForeign(GoRef);
GoRef matchaGoFunc(CGoBuffer);

CGoBuffer matchaGoCall(GoRef, CGoBuffer, CGoBuffer);
void matchaGoUntrack(GoRef);
void matchaTestFunc(void);

#endif // GO_GO_H