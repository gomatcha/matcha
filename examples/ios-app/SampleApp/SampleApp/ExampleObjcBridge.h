#import <Foundation/Foundation.h>
#import <MatchaBridge/MatchaBridge.h>
#import <Matcha/Matcha.h>

// Bridging with Go example. See examples/bridge/bridgeexample.go
@interface ExampleObjcBridge : NSObject
- (MatchaGoValue *)callWithGoValues:(MatchaGoValue *)param;
- (NSString *)callWithForeignValues:(NSString *)param;
- (NSString *)callGoFunctionWithForeignValues;
- (NSString *)callGoFunctionWithGoValues;
@end
