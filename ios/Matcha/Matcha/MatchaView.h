#import <UIKit/UIKit.h>
#import <Matcha/MatchaViewController.h>
@class GPBAny;
@class MatchaViewPBLayoutPaintNode;
@class MatchaViewNode;

@protocol MatchaChildView <NSObject>
- (id)initWithViewNode:(MatchaViewNode *)viewNode; // viewNode should be weakly retained
- (void)setNativeState:(NSData *)nativeState;
@end

@protocol MatchaChildViewController <NSObject>
- (id)initWithViewNode:(MatchaViewNode *)viewNode; // viewNode should be weakly retained
- (void)setNativeState:(NSData *)nativeState;
- (void)setMatchaChildViewControllers:(NSArray<UIViewController *> *)childVCs;
//- (void)setMatchaChildLayout:(NSArray<MatchaViewPBLayoutPaintNode *> *)layoutPaintNodes;
@end

@interface MatchaViewNode : NSObject
//- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId, ... NS_REQUIRES_NIL_TERMINATION; // varargs should be of MatchaGoValue *
- (void)call:(NSString *)funcId, ... NS_REQUIRES_NIL_TERMINATION; // varargs should be of MatchaGoValue *
@end
