#import <UIKit/UIKit.h>
#import "MatchaView.h"

void MatchaConfigureChildViewController(UIViewController *vc);
UIView<MatchaChildView> *MatchaViewWithNode(MatchaBuildNode *node, MatchaViewNode *viewNode);
UIGestureRecognizer *MatchaGestureRecognizerWithPB(int64_t viewId, GPBAny *any, MatchaViewNode *viewNode);
UIViewController<MatchaChildViewController> *MatchaViewControllerWithNode(MatchaBuildNode *node, MatchaViewNode *viewNode);
void MatchaRegisterView(NSString *string, MatchaViewRegistrationBlock block);
void MatchaRegisterViewController(NSString *string, MatchaViewControllerRegistrationBlock block);
