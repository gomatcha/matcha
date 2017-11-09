#ifndef GO_FOREIGN_H
#define GO_FOREIGN_H

#include <stdbool.h>
#include <stdint.h>

typedef int64_t FgnRef;
typedef int64_t GoRef;

typedef struct CGoBuffer {
    void *ptr; // UTF8 encoded string
    int64_t len; // length in bytes
} CGoBuffer;

FgnRef MatchaForeignNil();
bool MatchaForeignIsNil(FgnRef v);
FgnRef MatchaForeignBool(bool v);
bool MatchaForeignToBool(FgnRef v);
FgnRef MatchaForeignInt64(int64_t v);
int64_t MatchaForeignToInt64(FgnRef v);
FgnRef MatchaForeignFloat64(double v);
double MatchaForeignToFloat64(FgnRef v);
FgnRef MatchaForeignGoRef(GoRef v);
GoRef MatchaForeignToGoRef(FgnRef v);
FgnRef MatchaForeignString(CGoBuffer str); // Frees the buffer
CGoBuffer MatchaForeignToString(FgnRef v);
FgnRef MatchaForeignBytes(CGoBuffer bytes); // Frees the buffer
CGoBuffer MatchaForeignToBytes(FgnRef v);
FgnRef MatchaForeignArray(CGoBuffer buf); // Frees the buffer
CGoBuffer MatchaForeignToArray(FgnRef v);
FgnRef MatchaForeignBridge(CGoBuffer str); // Frees the buffer

// Call
FgnRef MatchaForeignCall(FgnRef v, CGoBuffer str, CGoBuffer args);

// Tracker
void MatchaForeignUntrack(FgnRef key);
int64_t MatchaForeignTrackerCount();

// Other
void MatchaForeignPanic();

#endif //GO_FOREIGN_H