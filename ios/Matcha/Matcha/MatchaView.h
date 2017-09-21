#import <UIKit/UIKit.h>
#import <Matcha/MatchaBuildNode.h>
#import <Matcha/MatchaViewController.h>
@class MatchaViewConfig;
@class MatchaViewController;
@class MatchaViewNode;

@protocol MatchaChildView <NSObject>
- (id)initWithViewNode:(MatchaViewNode *)viewNode;
- (void)setNode:(MatchaBuildNode *)node;
@end

@protocol MatchaChildViewController <NSObject>
- (id)initWithViewNode:(MatchaViewNode *)viewNode;
- (void)setNode:(MatchaBuildNode *)node;
- (void)setMatchaChildViewControllers:(NSArray<UIViewController *> *)childVCs;
- (void)setMatchaChildLayout:(NSMutableArray<MatchaViewPBLayoutPaintNode *> *)layoutPaintNodes;
@end

@interface MatchaViewNode : NSObject
- (MatchaViewController *)rootVC;
- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId args:(MatchaGoValue *)args, ... NS_REQUIRES_NIL_TERMINATION;
@end
