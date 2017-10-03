#ifndef MOCHIOBJC_H
#define MOCHIOBJC_H

#include <stdbool.h>
#include <stdint.h>

typedef int64_t FgnRef;
typedef int64_t GoRef;

typedef struct CGoBuffer {
    void *ptr; // UTF8 encoded string
    int64_t len; // length in bytes
} CGoBuffer;

void TestFunc();

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
FgnRef MatchaForeignCallSentinel();
FgnRef MatchaForeignCall(FgnRef v, CGoBuffer str, FgnRef args);

// Tracker
void MatchaForeignObjc(FgnRef key);

// Other
void MatchaForeignPanic();

@interface MatchaObjcBridge : NSObject
+ (MatchaObjcBridge *)sharedBridge;
- (void)setObject:(id<NSObject>)obj forKey:(NSString *)string;
- (id<NSObject>)objectForKey:(NSString *)string;
@end

#endif //MOCHIOBJC_H
