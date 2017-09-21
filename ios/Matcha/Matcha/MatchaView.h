#import <UIKit/UIKit.h>
#import <Matcha/MatchaViewController.h>
@class GPBAny;
@class MatchaViewPBLayoutPaintNode;
@class MatchaViewConfig;
@class MatchaViewController;
@class MatchaViewNode;

@protocol MatchaChildView <NSObject>
- (id)initWithViewNode:(MatchaViewNode *)viewNode; // viewNode should be weakly retained
- (void)setNativeState:(GPBAny *)nativeState;
@end

@protocol MatchaChildViewController <NSObject>
- (id)initWithViewNode:(MatchaViewNode *)viewNode; // viewNode should be weakly retained
- (void)setNativeState:(GPBAny *)nativeState;
- (void)setMatchaChildViewControllers:(NSArray<UIViewController *> *)childVCs;
- (void)setMatchaChildLayout:(NSArray<MatchaViewPBLayoutPaintNode *> *)layoutPaintNodes;
@end

@interface MatchaViewNode : NSObject
- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId args:(MatchaGoValue *)args, ... NS_REQUIRES_NIL_TERMINATION;
@end
