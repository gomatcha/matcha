#import "ExampleObjcBridge.h"
#import <MatchaBridge/MatchaBridge.h>

@implementation ObjcBridge

+ (void)load {
    static dispatch_once_t sOnce = 0;
    dispatch_once(&sOnce, ^{
        ObjcBridge *b = [[ObjcBridge alloc] init];
        [[MatchaObjcBridge sharedBridge] setObject:b forKey:@"gomatcha.io/matcha/example"];

        MatchaGoValue *func1 = [[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/examples/bridge callWithGoValues"];
        NSString *str1 = [func1 call:@"", [[MatchaGoValue alloc] initWithLongLong:123], nil][0].toString;

        MatchaGoValue *func2 = [[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/examples/bridge callWithForeignValues"];
        NSString *str2 = (NSString *)[func2 call:@"", [[MatchaGoValue alloc] initWithObject:@456], nil][0].toObject;

        str1 = nil;
        str2 = nil;
    });
}

- (MatchaGoValue *)callWithGoValues:(MatchaGoValue *)param {
    NSString *string = [NSString stringWithFormat:@"%lld", param.toLongLong];
    return [[MatchaGoValue alloc] initWithString:string];
}

- (NSString *)callWithForeignValues:(long long)param {
    return [NSString stringWithFormat:@"%lld", param];
}

@end
