#import "MatchaTabView.h"
#import "MatchaView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController.h"
#import "MatchaView_Private.h"
#import "UIImage+Tint.h"

@implementation MatchaTabView

+ (void)load {
    MatchaRegisterViewController(@"gomatcha.io/matcha/view/tabscreen", ^(MatchaViewNode *node){
        return [[MatchaTabView alloc] initWithViewNode:node];
    });
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super init])) {
        self.tabBar.translucent = NO;
        self.viewNode = viewNode;
        self.delegate = self;
        MatchaConfigureChildViewController(self);
    }
    return self;
}

- (void)setMatchaChildViewControllers:(NSArray<UIViewController *> *)childVCs {
    MatchaiOSPBTabView *pbTabNavigator = (id)[MatchaiOSPBTabView parseFromData:self.nativeState error:nil];
    
    self.tabBar.barTintColor = pbTabNavigator.hasBarColor ? [[UIColor alloc] initWithProtobuf:pbTabNavigator.barColor] : nil;
    self.tabBar.tintColor = pbTabNavigator.hasSelectedColor ? [[UIColor alloc] initWithProtobuf:pbTabNavigator.selectedColor] : nil;
    if ([self.tabBar respondsToSelector:@selector(unselectedItemTintColor)]) {
        self.tabBar.unselectedItemTintColor = pbTabNavigator.hasUnselectedColor ? [[UIColor alloc] initWithProtobuf:pbTabNavigator.unselectedColor] : nil; // TODO(KD): iOS 10.10 only
    }
    
    NSMutableArray *viewControllers = [NSMutableArray array];
    for (NSInteger idx = 0; idx < pbTabNavigator.screensArray.count; idx++) {
        MatchaiOSPBTabChildView *i = pbTabNavigator.screensArray[idx];
        UIViewController *vc = childVCs[idx];
        vc.tabBarItem.title = i.title;
        if ([i.title isEqual:@""]) {
            vc.tabBarItem.imageInsets = UIEdgeInsetsMake(6, 0, -6, 0);
        } else {
            vc.tabBarItem.imageInsets = UIEdgeInsetsZero;
        }
        vc.tabBarItem.badgeValue = i.badge.length == 0 ? nil : i.badge;
        vc.tabBarItem.image = [[UIImage alloc] initWithImageOrResourceProtobuf:i.icon];
        if (pbTabNavigator.hasIconTint) {
            vc.tabBarItem.image = [[vc.tabBarItem.image imageTintedWithColor:[[UIColor alloc] initWithProtobuf:pbTabNavigator.iconTint]] imageWithRenderingMode:UIImageRenderingModeAlwaysOriginal];
        }
        vc.tabBarItem.selectedImage = [[UIImage alloc] initWithImageOrResourceProtobuf:i.selectedIcon];
        if (pbTabNavigator.hasSelectedIconTint) {
            vc.tabBarItem.selectedImage = [[vc.tabBarItem.selectedImage imageTintedWithColor:[[UIColor alloc] initWithProtobuf:pbTabNavigator.selectedIconTint]] imageWithRenderingMode:UIImageRenderingModeAlwaysOriginal];
        }
        
        if (idx == (int)pbTabNavigator.selectedIndex) {
            NSDictionary *attributes = nil;
            if (pbTabNavigator.hasSelectedTextStyle) {
                attributes = [NSAttributedString attributesWithProtobuf:pbTabNavigator.selectedTextStyle];
            }
            [vc.tabBarItem setTitleTextAttributes:attributes forState:UIControlStateNormal];
        } else {
            NSDictionary *attributes = nil;
            if (pbTabNavigator.hasUnselectedTextStyle) {
                attributes = [NSAttributedString attributesWithProtobuf:pbTabNavigator.unselectedTextStyle];
            }
            [vc.tabBarItem setTitleTextAttributes:attributes forState:UIControlStateNormal];
        }
        
        [viewControllers addObject:vc];
    }
    
    self.viewControllers = viewControllers;
    self.selectedIndex = (int)pbTabNavigator.selectedIndex;
}

- (void)tabBarController:(UITabBarController *)tabBarController didSelectViewController:(UIViewController *)viewController {
    MatchaiOSPBTabEvent *event = [[MatchaiOSPBTabEvent alloc] init];
    event.selectedIndex = tabBarController.selectedIndex;
    [self.viewNode call:@"OnSelect", [[MatchaGoValue alloc] initWithData:event.data], nil];
}

- (void)setMatchaChildLayout:(GPBInt64ObjectDictionary *)layoutPaintNodes {
    // no-op
}

@end
