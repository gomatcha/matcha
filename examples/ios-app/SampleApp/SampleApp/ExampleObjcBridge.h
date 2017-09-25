#import <Foundation/Foundation.h>
#import <MatchaBridge/MatchaBridge.h>
#import <Matcha/Matcha.h>

@interface ObjcBridge : NSObject
- (MatchaGoValue *)callWithGoValues:(MatchaGoValue *)param;
- (NSString *)callWithForeignValues:(long long)param;
@end
