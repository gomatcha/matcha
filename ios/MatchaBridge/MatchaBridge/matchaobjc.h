#ifndef MOCHIOBJC_H
#define MOCHIOBJC_H

#import <Foundation/Foundation.h>

int MatchaTest();

typedef int64_t ObjcRef;
typedef int64_t GoRef;

typedef struct CGoBuffer {
    void *ptr; // UTF8 encoded string
    int64_t len; // length in bytes
} CGoBuffer;

@interface MatchaObjcBridge : NSObject
+ (MatchaObjcBridge *)sharedBridge;
@end

ObjcRef MatchaObjcBridge_();

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

ObjcRef MatchaObjcArray();
int64_t MatchaObjcArrayLen(ObjcRef v);
void MatchaObjcArrayAppend(ObjcRef v, ObjcRef a);
ObjcRef MatchaObjcArrayAt(ObjcRef v, int64_t index);

// ObjcRef MatchaObjcDict();
// ObjcRef MatchaObjcDictKeys(ObjcRef v);
// ObjcRef MatchaObjcDictGet(ObjcRef v, ObjcRef key);
// ObjcRef MatchaObjcDictSet(ObjcRef v, ObjcRef key, ObjCRef value);

// Call
ObjcRef MatchaObjcCallSentinel();
ObjcRef MatchaObjcCall(ObjcRef v, CGoBuffer str, ObjcRef args);

// Tracker
ObjcRef MatchaTrackObjc(id value);
id MatchaGetObjc(ObjcRef key);
void MatchaUntrackObjc(ObjcRef key);

// Utilities
NSString *MatchaCGoBufferToNSString(CGoBuffer buf); // Frees the buffer.
CGoBuffer MatchaNSStringToCGoBuffer(NSString *str); // Allocates a buffer.
NSData *MatchaCGoBufferToNSData(CGoBuffer buf); // Frees the buffer.
CGoBuffer MatchaNSDataToCGoBuffer(NSData *data); // Allocates a buffer.


// ObjcRef MatchaObjcWithGo(GoRef v);
// GoRef MatchaObjcToGo(ObjcRef v);

#endif //MOCHIOBJC_H