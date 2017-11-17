#import "ExampleObjcBridge.h"
#import <MatchaBridge/MatchaBridge.h>

@implementation ObjcBridge

+ (void)load {
    static dispatch_once_t sOnce = 0;
    dispatch_once(&sOnce, ^{
        // Register the ObjcBridge object with Matcha.
        [[MatchaObjcBridge sharedBridge] setObject:[[ObjcBridge alloc] init] forKey:@"gomatcha.io/matcha/example/bridge"];
    });
}

- (MatchaGoValue *)callWithGoValues:(MatchaGoValue *)param {
    NSString *string = [NSString stringWithFormat:@"Hello %@!", param.toString];
    return [[MatchaGoValue alloc] initWithString:string];
}

- (NSString *)callWithForeignValues:(NSString *)param {
    return [NSString stringWithFormat:@"Hello %@!", param];
}

- (NSString *)callGoFunctionWithForeignValues {
    MatchaGoValue *func = [[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/examples/bridge callWithGoValues"];
    return [func call:@"", [[MatchaGoValue alloc] initWithString:@"Yuki"], nil][0].toString;
}

- (NSString *)callGoFunctionWithGoValues {
    MatchaGoValue *func = [[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/examples/bridge callWithForeignValues"];
    return (NSString *)[func call:@"", [[MatchaGoValue alloc] initWithObject:@"Ame"], nil][0].toObject;
}

@end
