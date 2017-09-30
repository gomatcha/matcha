#import <UIKit/UIKit.h>
#import "MatchaView.h"
@class MatchaViewNode;

@interface MatchaTabView : UITabBarController <MatchaChildViewController, UITabBarControllerDelegate>
@property (nonatomic, weak) MatchaViewNode *viewNode;
@property (nonatomic, strong) NSData *nativeState;
@end
