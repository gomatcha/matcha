// +build matcha,darwin

#include <stdio.h>
#include <stdint.h>
#include <string.h>
#import <Foundation/Foundation.h>
#include "matchago-objc.h"
#include "matchaforeign-objc.h"

@interface MatchaGoBridge ()
@property (nonatomic, strong) MatchaGoValue *rootObject;
@end

@implementation MatchaGoBridge

+ (MatchaGoBridge *)sharedBridge {
    static MatchaGoBridge *sBridge = nil;
    static dispatch_once_t sOnce;
    dispatch_once (&sOnce, ^{
        sBridge = [[MatchaGoBridge alloc] init];
    });
    return sBridge;
}

@end

@implementation MatchaGoValue {
    GoRef _ref;
}

@synthesize ref = _ref;

- (id)initWithGoRef:(GoRef)ref {
    if ((self = [super init])) {
        _ref = ref;
    }
    return self;
}

- (id)initWithObject:(id)v {
    return [self initWithGoRef:matchaGoForeign(MatchaTrackObjc(v))];
}

- (id)initWithBool:(BOOL)v {
    return [self initWithGoRef:matchaGoBool(v)];
}

- (id)initWithInt:(int)v {
    return [self initWithGoRef:matchaGoInt(v)];
}

- (id)initWithLongLong:(long long)v {
    return [self initWithGoRef:matchaGoInt64(v)];
}

- (id)initWithUnsignedLongLong:(unsigned long long)v {
    return [self initWithGoRef:matchaGoUint64(v)];
}

- (id)initWithDouble:(double)v {
    return [self initWithGoRef:matchaGoFloat64(v)];
}

- (id)initWithString:(NSString *)v {
    CGoBuffer buf = MatchaNSStringToCGoBuffer(v);
    return [self initWithGoRef:matchaGoString(buf)];
}

- (id)initWithData:(NSData *)v {
    CGoBuffer buf = MatchaNSDataToCGoBuffer(v);
    return [self initWithGoRef:matchaGoBytes(buf)];
}

- (id)initWithArray:(NSArray<MatchaGoValue *> *)v {
    GoRef ref = matchaGoArray();
    for (MatchaGoValue *i in v) {
        GoRef prev = ref;
        ref = matchaGoArrayAppend(ref, i.ref);
        matchaGoUntrack(prev); // Must manually untrack
    }
    return [self initWithGoRef:ref];
}

- (id)initWithType:(NSString *)v {
    CGoBuffer buf = MatchaNSStringToCGoBuffer(v);
    return [self initWithGoRef:matchaGoType(buf)];
}

- (id)initWithFunc:(NSString *)v {
    CGoBuffer buf = MatchaNSStringToCGoBuffer(v);
    return [self initWithGoRef:matchaGoFunc(buf)];
}

- (id)toObject {
    return MatchaGetObjc(matchaGoToForeign(_ref));
}

- (BOOL)toBool {
    return matchaGoToBool(_ref);
}

- (long long)toLongLong {
    return matchaGoToInt64(_ref);
}

- (unsigned long long)toUnsignedLongLong {
    return matchaGoToUint64(_ref);
}

- (double)toDouble {
    return matchaGoToFloat64(_ref);
}

- (NSString *)toString {
    return MatchaCGoBufferToNSString(matchaGoToString(_ref));
}

- (NSData *)toData {
    return MatchaCGoBufferToNSData(matchaGoToBytes(_ref));
}

- (NSArray *)toArray {
    return MatchaCGoBufferToNSArray(matchaGoArrayBuffer(_ref));
}

- (NSMapTable *)toMapTable {
    NSMapTable *mapTable = [NSMapTable strongToStrongObjectsMapTable];
    MatchaGoValue *array = [[MatchaGoValue alloc] initWithGoRef:matchaGoMapKeys(_ref)];
    for (MatchaGoValue *key in array.toArray) {
        MatchaGoValue *value = [[MatchaGoValue alloc] initWithGoRef:matchaGoMapGet(_ref, key.ref)];
        [mapTable setObject:value forKey:key];
    }
    return mapTable;
}

- (BOOL)isEqual:(MatchaGoValue *)value {
    return matchaGoEqual(_ref, value.ref);
}

- (BOOL)isNil {
    return matchaGoIsNil(_ref);
}

- (MatchaGoValue *)elem {
    return [[MatchaGoValue alloc] initWithGoRef:matchaGoElem(_ref)];
}

- (NSArray<MatchaGoValue *> *)call:(NSString *)method, ... {
    NSMutableArray *array = [[NSMutableArray alloc] init];
    va_list args;
    va_start(args, method);
    NSArray *rlt = [self call:method args:args];
    va_end(args);
    
    return rlt;
}

- (NSArray<MatchaGoValue *> *)call:(NSString *)method args:(va_list)args {
    NSMutableArray *array = [NSMutableArray array];
    id arg = nil;
    while ((arg = va_arg(args, id))) {
        [array addObject:arg];
    }

    CGoBuffer argsBuffer = MatchaNSArrayToCGoBuffer(array);
    CGoBuffer rlt = matchaGoCall(_ref, MatchaNSStringToCGoBuffer(method), argsBuffer);
    return MatchaCGoBufferToNSArray(rlt);
}

- (MatchaGoValue *)field:(NSString *)name {
    GoRef rlt = matchaGoField(_ref, MatchaNSStringToCGoBuffer(name));
    return [[MatchaGoValue alloc] initWithGoRef:rlt];
}

- (void)setField:(NSString *)name value:(MatchaGoValue *)value {
    matchaGoFieldSet(_ref, MatchaNSStringToCGoBuffer(name), value.ref);
}

- (MatchaGoValue *)objectForKeyedSubscript:(NSString *)key {
    return [self field:key];
}

- (void)setObject:(MatchaGoValue *)object forKeyedSubscript:(NSString *)key {
    [self setField:key value:object];
}

- (void)dealloc {
    matchaGoUntrack(_ref);
}

@end