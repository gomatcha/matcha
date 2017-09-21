#import "MatchaView_Private.h"

void MatchaConfigureChildViewController(UIViewController *vc) {
    vc.edgesForExtendedLayout=UIRectEdgeNone;
    vc.extendedLayoutIncludesOpaqueBars=NO;
    vc.automaticallyAdjustsScrollViewInsets=NO;
}
