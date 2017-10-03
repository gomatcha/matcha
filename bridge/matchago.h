// +build matcha

#ifndef MOCHIGO_H
#define MOCHIGO_H

#include "matchaforeign.h"
#include <stdbool.h>

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
GoRef matchaGoType(CGoBuffer);
GoRef matchaGoFunc(CGoBuffer);

bool matchaGoIsNil(GoRef);
bool matchaGoEqual(GoRef, GoRef);
GoRef matchaGoElem(GoRef);
CGoBuffer matchaGoCall(GoRef, CGoBuffer, CGoBuffer);
GoRef matchaGoField(GoRef, CGoBuffer);
void matchaGoFieldSet(GoRef, CGoBuffer, GoRef);

void matchaGoUntrack(GoRef);

#endif // MOCHIGO_H