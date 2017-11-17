#import <Foundation/Foundation.h>
#import <MatchaBridge/MatchaBridge.h>
#import <Matcha/Matcha.h>

// Bridging with Go example. See examples/bridge/bridgeexample.go
@interface ObjcBridge : NSObject
- (MatchaGoValue *)callWithGoValues:(MatchaGoValue *)param;
- (NSString *)callWithForeignValues:(NSString *)param;
@end
