#ifndef MOCHIGO_H
#define MOCHIGO_H

#import <Foundation/Foundation.h>
#include "matchaobjc.h"
@class MatchaGoValue;

GoRef matchaGoBool(bool);
bool matchaGoToBool(GoRef);
GoRef matchaGoInt(int);
GoRef matchaGoInt64(int64_t);
int64_t matchaGoToInt64(GoRef);
GoRef matchaGoUint64(uint64_t);
uint64_t matchaGoToUint64(GoRef);
GoRef matchaGoFloat64(double);
double matchaGoToFloat64(GoRef);
GoRef matchaGoString(CGoBuffer); // Frees the buffer
CGoBuffer matchaGoToString(GoRef);
GoRef matchaGoBytes(CGoBuffer); // Frees the buffer
CGoBuffer matchaGoToBytes(GoRef);

GoRef matchaGoArray();
int64_t matchaGoArrayLen(GoRef);
GoRef matchaGoArrayAppend(GoRef, GoRef);
GoRef matchaGoArrayAt(GoRef, int64_t);

GoRef matchaGoMap();
GoRef matchaGoMapKeys(GoRef);
GoRef matchaGoMapGet(GoRef map, GoRef key);
void matchaGoMapSet(GoRef map, GoRef key, GoRef value);

GoRef matchaGoType(CGoBuffer); // Frees the buffer
GoRef matchaGoFunc(CGoBuffer); // Frees the buffer

bool matchaGoIsNil(GoRef);
bool matchaGoEqual(GoRef, GoRef);
GoRef matchaGoElem(GoRef);
GoRef matchaGoCall(GoRef, CGoBuffer, GoRef);
GoRef matchaGoField(GoRef, CGoBuffer);
void matchaGoFieldSet(GoRef, CGoBuffer, GoRef);

void matchaGoUntrack(GoRef);

@interface MatchaGoBridge : NSObject
+ (MatchaGoBridge *)sharedBridge;
@end

@interface MatchaGoValue : NSObject
- (id)initWithGoRef:(GoRef)ref;
- (id)initWithBool:(BOOL)v;
- (id)initWithInt:(int)v;
- (id)initWithLongLong:(long long)v;
- (id)initWithUnsignedLongLong:(unsigned long long)v;
- (id)initWithDouble:(double)v;
- (id)initWithString:(NSString *)v;
- (id)initWithData:(NSData *)v;
- (id)initWithArray:(NSArray<MatchaGoValue *> *)v;
- (id)initWithType:(NSString *)typeName;
- (id)initWithFunc:(NSString *)funcName;
@property (nonatomic, readonly) GoRef ref;
- (BOOL)toBool;
- (long long)toLongLong;
- (unsigned long long)toUnsignedLongLong;
- (double)toDouble;
- (NSString *)toString;
- (NSData *)toData;
- (NSArray *)toArray;
- (NSMapTable *)toMapTable;
// - (NSDictionary *)toDictionary;
- (BOOL)isNil;
- (BOOL)isEqual:(MatchaGoValue *)value;
- (MatchaGoValue *)elem;
- (NSArray<MatchaGoValue *> *)call:(NSString *)method args:(NSArray<MatchaGoValue *> *)args; // pass in nil for the method to call a closure.
- (MatchaGoValue *)field:(NSString *)name;
- (void)setField:(NSString *)name value:(MatchaGoValue *)value;
- (MatchaGoValue *)objectForKeyedSubscript:(NSString *)key;
- (void)setObject:(MatchaGoValue *)object forKeyedSubscript:(NSString *)key;
@end

#endif // MOCHIGO_H