#import "MatchaTabView.h"
#import "MatchaView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController.h"
#import "MatchaView_Private.h"

@implementation MatchaTabView

+ (void)load {
    MatchaRegisterViewController(@"gomatcha.io/matcha/view/tabscreen", ^(MatchaViewNode *node){
        return [[MatchaTabView alloc] initWithViewNode:node];
    });
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super init])) {
        self.viewNode = viewNode;
        self.delegate = self;
        MatchaConfigureChildViewController(self);
    }
    return self;
}

- (void)setMatchaChildViewControllers:(NSArray<UIViewController *> *)childVCs {
    MatchaiOSPBTabView *pbTabNavigator = (id)[self.nativeState unpackMessageClass:[MatchaiOSPBTabView class] error:nil];
    
    self.tabBar.barTintColor = pbTabNavigator.hasBarColor ? [[UIColor alloc] initWithProtobuf:pbTabNavigator.barColor] : nil;
    self.tabBar.tintColor = pbTabNavigator.hasBarColor ? [[UIColor alloc] initWithProtobuf:pbTabNavigator.selectedColor] : nil;
    if ([self.tabBar respondsToSelector:@selector(unselectedItemTintColor)]) {
        self.tabBar.unselectedItemTintColor = pbTabNavigator.hasUnselectedColor ? [[UIColor alloc] initWithProtobuf:pbTabNavigator.unselectedColor] : nil; // TODO(KD): iOS 10.10 only
    }
    if (pbTabNavigator.hasUnselectedTextStyle) {
        [[UITabBarItem appearance] setTitleTextAttributes:[NSAttributedString attributesWithProtobuf:pbTabNavigator.unselectedTextStyle] forState:UIControlStateNormal];
    }
    if (pbTabNavigator.hasSelectedTextStyle) {
        [[UITabBarItem appearance] setTitleTextAttributes:[NSAttributedString attributesWithProtobuf:pbTabNavigator.selectedTextStyle] forState:UIControlStateSelected];
    }

    
    NSMutableArray *viewControllers = [NSMutableArray array];
    for (NSInteger idx = 0; idx < pbTabNavigator.screensArray.count; idx++) {
        MatchaiOSPBTabChildView *i = pbTabNavigator.screensArray[idx];
        UIViewController *vc = childVCs[idx];
        vc.tabBarItem.title = i.title;
        vc.tabBarItem.badgeValue = i.badge.length == 0 ? nil : i.badge;
        vc.tabBarItem.image = [[UIImage alloc] initWithImageOrResourceProtobuf:i.icon];
        vc.tabBarItem.selectedImage = [[UIImage alloc] initWithImageOrResourceProtobuf:i.selectedIcon];
        [viewControllers addObject:vc];
    }
    
    self.viewControllers = viewControllers;
    self.selectedIndex = (int)pbTabNavigator.selectedIndex;
}

- (void)tabBarController:(UITabBarController *)tabBarController didSelectViewController:(UIViewController *)viewController {
    MatchaiOSPBTabEvent *event = [[MatchaiOSPBTabEvent alloc] init];
    event.selectedIndex = tabBarController.selectedIndex;
    [self.viewNode call:@"OnSelect" args:[[MatchaGoValue alloc] initWithData:event.data], nil];
}

- (void)setMatchaChildLayout:(GPBInt64ObjectDictionary *)layoutPaintNodes {
    // no-op
}

@end
