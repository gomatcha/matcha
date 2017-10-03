#ifndef MOCHIOBJC_H
#define MOCHIOBJC_H

#include <stdbool.h>
#include <stdint.h>

typedef int64_t ObjcRef;
typedef int64_t GoRef;

typedef struct CGoBuffer {
    void *ptr; // UTF8 encoded string
    int64_t len; // length in bytes
} CGoBuffer;

void TestFunc();

ObjcRef MatchaObjcBool(bool v);
bool MatchaObjcToBool(ObjcRef v);
ObjcRef MatchaObjcInt64(int64_t v);
int64_t MatchaObjcToInt64(ObjcRef v);
ObjcRef MatchaObjcFloat64(double v);
double MatchaObjcToFloat64(ObjcRef v);
ObjcRef MatchaObjcGoRef(GoRef v);
GoRef MatchaObjcToGoRef(ObjcRef v);
ObjcRef MatchaObjcString(CGoBuffer str); // Frees the buffer
CGoBuffer MatchaObjcToString(ObjcRef v);
ObjcRef MatchaObjcBytes(CGoBuffer bytes); // Frees the buffer
CGoBuffer MatchaObjcToBytes(ObjcRef v);
ObjcRef MatchaObjcArray(CGoBuffer buf); // Frees the buffer
CGoBuffer MatchaObjcToArray(ObjcRef v);
ObjcRef MatchaForeignBridge(CGoBuffer str); // Frees the buffer

// Call
ObjcRef MatchaObjcCallSentinel();
ObjcRef MatchaObjcCall(ObjcRef v, CGoBuffer str, ObjcRef args);

// Tracker
void MatchaUntrackObjc(ObjcRef key);

// Other
void MatchaForeignPanic();

#endif //MOCHIOBJC_H