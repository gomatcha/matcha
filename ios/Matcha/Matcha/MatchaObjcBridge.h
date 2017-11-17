#import <UIKit/UIKit.h>
#import <MatchaBridge/MatchaBridge.h>

@interface MatchaObjcBridge_X : NSObject
+ (NSMapTable *)viewControllers;
- (MatchaGoValue *)sizeForAttributedString:(NSData *)data maxLines:(int)maxLines;
- (bool)updateId:(NSInteger)identifier withProtobuf:(NSData *)protobuf;
- (NSString *)assetsDir;
- (MatchaGoValue *)imageForResource:(NSString *)path;
- (MatchaGoValue *)propertiesForResource:(NSString *)path;
- (void)displayAlert:(NSData *)protobuf;
- (BOOL)openURL:(NSString *)url;
- (int)orientation;
@end
