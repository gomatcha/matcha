#import <UIKit/UIKit.h>
#import "MatchaView.h"
@class MatchaBuildNode;

void MatchaConfigureChildViewController(UIViewController *vc);
UIView<MatchaChildView> *MatchaViewWithNode(MatchaBuildNode *node, MatchaViewNode *viewNode);
UIGestureRecognizer *MatchaGestureRecognizerWithPB(int64_t viewId, GPBAny *any, MatchaViewNode *viewNode);
UIViewController<MatchaChildViewController> *MatchaViewControllerWithNode(MatchaBuildNode *node, MatchaViewNode *viewNode);
void MatchaRegisterView(NSString *string, MatchaViewRegistrationBlock block);
void MatchaRegisterViewController(NSString *string, MatchaViewControllerRegistrationBlock block);

@interface MatchaViewNode (Private)
- (id)initWithParent:(MatchaViewNode *)node rootVC:(MatchaViewController *)rootVC identifier:(NSNumber *)identifier;
- (void)setRoot:(MatchaViewPBRoot *)root;
- (UIViewController<MatchaChildViewController> *)viewController;
- (UIView<MatchaChildView> *)view;
- (MatchaViewController *)rootVC;
- (NSArray<MatchaGoValue *> *)call2:(NSString *)funcId, ... NS_REQUIRES_NIL_TERMINATION;
@end
