// +build matcha,darwin

#include <stdio.h>
#include <stdint.h>
#include <string.h>
#import <Foundation/Foundation.h>
#include "go-foreign.h"
#include "objc-foreign.h"
#include "go-go.h"
#include "objc-go.h"

@interface MatchaObjcBridge ()
@property (nonatomic, strong) NSMutableDictionary<NSString *, id<NSObject>> *dictionary;
@end

@implementation MatchaObjcBridge

+ (id)sharedBridge {
    static MatchaObjcBridge *sBridge = nil;
    static dispatch_once_t sOnce;
    dispatch_once (&sOnce, ^{
        sBridge = [[MatchaObjcBridge alloc] init];
    });
    return sBridge;
}

- (id)init {
    if ((self = [super init])) {
        self.dictionary = [NSMutableDictionary dictionary];
    }
    return self;
}

- (void)setObject:(id<NSObject>)obj forKey:(NSString *)string {
    self.dictionary[string] = obj;
}

- (id<NSObject>)objectForKey:(NSString *)string {
    return self.dictionary[string];
}

@end

@interface MatchaTracker : NSObject {
    NSMapTable *_mapTable;
    int64_t _maxKey;
}
@end

@implementation MatchaTracker

+ (MatchaTracker *)sharedTracker {
    static MatchaTracker *sTracker = nil;
    static dispatch_once_t sOnce;
    dispatch_once (&sOnce, ^{
        sTracker = [[MatchaTracker alloc] init];
    });
    return sTracker;
}

- (id)init {
    if ((self = [super init])) {
        _mapTable = [[NSMapTable alloc] initWithKeyOptions:NSPointerFunctionsObjectPersonality|NSPointerFunctionsStrongMemory 
            valueOptions:NSPointerFunctionsObjectPersonality|NSPointerFunctionsStrongMemory capacity:0];
        _maxKey = 0;
    }
    return self;
}

- (FgnRef)track:(id)object {
    if (object == nil) {
        return 0;
    }
    @synchronized (self) {
        _maxKey += 1;
        [_mapTable setObject:object forKey:@(_maxKey)];
        return _maxKey;
    }
}

- (void)untrack:(FgnRef)key {
    if (key == 0) {
        return;
    }
    @synchronized (self) {
        id keyObj = @(key);
        id object = [_mapTable objectForKey:keyObj];
        if (object == nil) {
            NSLog(@"UntrackError");
            @throw @"Untrack error. No corresponding object for key.";
        }
        [_mapTable removeObjectForKey:keyObj];
    }
}

- (id)get:(FgnRef)key {
    if (key == 0) {
        return nil;
    }
    @synchronized (self) {
        id object = [_mapTable objectForKey:(id)@(key)];
        if (object == nil) {
            @throw @"Get error. No corresponding object for key";
        }
        return object;
    }
}

@end

FgnRef MatchaForeignBool(bool v) {
    return MatchaForeignTrack(@(v));
}

bool MatchaForeignToBool(FgnRef v) {
    NSNumber *val = MatchaForeignGet(v);
    return val.boolValue;
}

FgnRef MatchaForeignInt64(int64_t v) {
    return MatchaForeignTrack(@(v));
}

int64_t MatchaForeignToInt64(FgnRef v) {
    NSNumber *val = MatchaForeignGet(v);
    return val.longLongValue;
}

FgnRef MatchaForeignFloat64(double v) {
    return MatchaForeignTrack(@(v));
}

double MatchaForeignToFloat64(FgnRef v) {
    NSNumber *val = MatchaForeignGet(v);
    return val.doubleValue;
}

FgnRef MatchaForeignGoRef(GoRef v) {
    return MatchaForeignTrack([[MatchaGoValue alloc] initWithGoRef:v]);
}

GoRef MatchaForeignToGoRef(FgnRef v) {
    MatchaGoValue *val = MatchaForeignGet(v);
    return val.ref;
}

FgnRef MatchaForeignString(CGoBuffer cstr) {
    return MatchaForeignTrack(MatchaCGoBufferToNSString(cstr));
}

CGoBuffer MatchaForeignToString(FgnRef v) {
    return MatchaNSStringToCGoBuffer(MatchaForeignGet(v));
}

FgnRef MatchaForeignBytes(CGoBuffer buf) {
    return MatchaForeignTrack(MatchaCGoBufferToNSData(buf));
}

CGoBuffer MatchaForeignToBytes(FgnRef v) {
    NSData *data = MatchaForeignGet(v);
    return MatchaNSDataToCGoBuffer(data);
}

FgnRef MatchaForeignArray(CGoBuffer buf) {
    NSArray *array = MatchaCGoBufferToNSArray2(buf);
    return MatchaForeignTrack(array);
}

CGoBuffer MatchaForeignToArray(FgnRef v) {
    NSArray *val = MatchaForeignGet(v);
    return MatchaNSArrayToCGoBuffer2(val);
}

FgnRef MatchaForeignBridge(CGoBuffer str) {
    NSString *string = MatchaCGoBufferToNSString(str);
    MatchaObjcBridge *root = [[MatchaObjcBridge sharedBridge] objectForKey:string];;
    return MatchaForeignTrack(root);
}

// Call

@interface MatchaNilSentinel : NSObject
@end
@implementation MatchaNilSentinel
@end

FgnRef MatchaForeignCall(FgnRef v, CGoBuffer cstr, CGoBuffer arguments) {
    id obj = MatchaForeignGet(v);
    NSArray *args = MatchaCGoBufferToNSArray2(arguments);
    NSString *str = MatchaCGoBufferToNSString(cstr);
    SEL sel = NSSelectorFromString(str);
    NSMethodSignature *sig = [[obj class] instanceMethodSignatureForSelector:sel];
    if (sig == nil) {
        NSLog(@"MatchaForeignCall with nil signature: %@, %@, %@", obj, str, args);
    }
    
    // Build invocation.
    NSInvocation *inv = [NSInvocation invocationWithMethodSignature:sig];
    inv.selector = sel;
    inv.target = obj;
    for (int i=0; i < args.count; i++) {
        id argObj = args[i];
        NSNumber *num = (NSNumber *)argObj;
        const char *type = [sig getArgumentTypeAtIndex:i+2];
        
        switch (type[0]) {
        case 'c': {
            char arg = num.charValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'i': {
            int arg = num.intValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 's': {
            short arg = num.shortValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'l': {
            long arg = num.longValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'q': {
            long long arg = num.longLongValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'C': {
            unsigned char arg = num.unsignedCharValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'I': {
            unsigned int arg = num.unsignedIntValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'S': {
            unsigned short arg = num.unsignedShortValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'L': {
            unsigned long arg = num.unsignedLongValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'Q': {
            unsigned long long arg = num.unsignedLongLongValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'f': {
            float arg = num.floatValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'd': {
            double arg = num.doubleValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case 'B': {
            bool arg = num.boolValue;
            [inv setArgument:&arg atIndex:i+2];
            break;
        }
        case '@': {
            if ([argObj isKindOfClass:[MatchaNilSentinel class]]) {
                id nilObject = nil;
                [inv setArgument:&nilObject atIndex:i+2];
            } else {
                [inv setArgument:&argObj atIndex:i+2];
            }
            break;
        }
        default: {
            @throw @"MatchaForeignCall: Unsupported argument type";
        }
        }
    }
    
    // Invoke.
    [inv invoke];

    // Get return value.
    const char *type = [sig methodReturnType];
    id ret = nil;
    switch (type[0]) {
    case 'c': {
        char v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 'C': {
        unsigned char v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 'i': {
        int v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 'I': {
        unsigned int v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 's': {
        short v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 'S': {
        unsigned short v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 'l': {
        long v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 'L': {
        unsigned long v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 'q': {
        long long v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 'Q': {
        unsigned long long v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 'f': {
        float v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
    case 'd': {
        double v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }
     case 'B': {
        bool v;
        [inv getReturnValue:&v];
        ret = @(v);
        break;
    }   
    case 'v': {
        ret = nil;
        break;
    }
    case '@': {
        void *v = nil;
        [inv getReturnValue:&v];
        ret = (__bridge id)v;
        break;
    }
    default: {
        @throw @"MatchaForeignCall: Unsupported return type";
    }
    }
    return MatchaForeignTrack(ret);
}

// Tracker

FgnRef MatchaForeignTrack(id value) {
    return [[MatchaTracker sharedTracker] track:value];
}

id MatchaForeignGet(FgnRef key) {
    return [[MatchaTracker sharedTracker] get:key];
}

void MatchaForeignUntrack(FgnRef key) {
    [[MatchaTracker sharedTracker] untrack:key];
}

// Other

void MatchaForeignPanic() {
    @throw [NSException exceptionWithName:@"Golang Panic" reason:@"" userInfo:nil];
}

// Utilities

NSString *MatchaCGoBufferToNSString(CGoBuffer cstr) {
    if (cstr.len == 0) {
        return @"";
    }
    return [[NSString alloc] initWithBytesNoCopy:cstr.ptr length:cstr.len encoding:NSUTF8StringEncoding freeWhenDone:YES];
}

CGoBuffer MatchaNSStringToCGoBuffer(NSString *str) {
    int len = [str lengthOfBytesUsingEncoding:NSUTF8StringEncoding];
    if (len == 0) {
        return (CGoBuffer){0};
    }

    char *buf = (char *)malloc(len);
    assert(buf != NULL);
    [str getBytes:buf maxLength:len usedLength:NULL encoding:NSUTF8StringEncoding options:0 range:NSMakeRange(0, str.length) remainingRange:NULL];
  
    CGoBuffer cstr;
    cstr.ptr = buf;
    cstr.len = len;
    return cstr;
}

NSData *MatchaCGoBufferToNSData(CGoBuffer buf) {
    if (buf.len == 0) {
        return [NSData data];
    }
    return [[NSData alloc] initWithBytesNoCopy:buf.ptr length:buf.len freeWhenDone:YES];
}

CGoBuffer MatchaNSDataToCGoBuffer(NSData *data) {
    int len = [data length];
    if (len == 0) {
        return (CGoBuffer){0};
    }

    char *buf = (char *)malloc(len);
    assert(buf != NULL);
    [data getBytes:buf length:len];
  
    CGoBuffer cstr;
    cstr.ptr = buf;
    cstr.len = len;
    return cstr;
}

NSArray<MatchaGoValue *> *MatchaCGoBufferToNSArray(CGoBuffer buf) {
    NSMutableArray *array = [NSMutableArray array];
    char *data = buf.ptr;
    for (NSInteger i = 0; i < buf.len/8; i++) {
        GoRef ref = 0;
        memcpy(&ref, data, 8);
        [array addObject:[[MatchaGoValue alloc] initWithGoRef:ref]];
        data += 8;
    }
    free(buf.ptr);
    return array;
}

NSArray<id> *MatchaCGoBufferToNSArray2(CGoBuffer buf) {
    NSMutableArray *array = [NSMutableArray array];
    char *data = buf.ptr;
    for (NSInteger i = 0; i < buf.len/8; i++) {
        FgnRef ref = 0;
        memcpy(&ref, data, 8);
        [array addObject:MatchaForeignGet(ref)];
        data += 8;
    }
    free(buf.ptr);
    return array;
}

CGoBuffer MatchaNSArrayToCGoBuffer(NSArray<MatchaGoValue *> *array) {
    if (array.count == 0) {
        return (CGoBuffer){0};
    }
    
    char *buf = (char *)malloc(array.count * 8);
    char *data = buf;
    assert(buf != NULL);
    for (int i = 0; i < array.count; i++) {
        int64_t ref = array[i].ref;
        memcpy(data, &ref, 8);
        data += 8;
    }
    
    CGoBuffer cstr;
    cstr.ptr = buf;
    cstr.len = array.count * 8;
    return cstr;
}

CGoBuffer MatchaNSArrayToCGoBuffer2(NSArray *array) {
    if (array.count == 0) {
        return (CGoBuffer){0};
    }
    
    char *buf = (char *)malloc(array.count * 8);
    char *data = buf;
    assert(buf != NULL);
    for (int i = 0; i < array.count; i++) {
        FgnRef ref = MatchaForeignTrack(array[i]);
        memcpy(data, &ref, 8);
        data += 8;
    }
    
    CGoBuffer cstr;
    cstr.ptr = buf;
    cstr.len = array.count * 8;
    return cstr;
}
