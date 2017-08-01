#import "MatchaTabScreen.h"
#import "MatchaView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController.h"

@implementation MatchaTabScreen

+ (void)load {
    MatchaRegisterViewController(@"gomatcha.io/matcha/view/tabscreen", ^(MatchaViewNode *node){
        return [[MatchaTabScreen alloc] initWithViewNode:node];
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

- (void)setMatchaChildViewControllers:(NSDictionary<NSNumber *, UIViewController *> *)childVCs {
    GPBAny *state = self.node.nativeViewState;
    MatchaTabScreenPBView *pbTabNavigator = (id)[state unpackMessageClass:[MatchaTabScreenPBView class] error:nil];
    
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
    for (MatchaTabScreenPBChildView *i in pbTabNavigator.screensArray) {
        UIViewController *vc = childVCs[@(i.id_p)];
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
    MatchaTabScreenPBEvent *event = [[MatchaTabScreenPBEvent alloc] init];
    event.selectedIndex = tabBarController.selectedIndex;
    
    NSData *data = [event data];
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:data];
    
    [self.viewNode.rootVC call:@"OnSelect" viewId:self.node.identifier.longLongValue args:@[value]];
}

- (void)setMatchaChildLayout:(GPBInt64ObjectDictionary *)layoutPaintNodes {
    // no-op
}

@end
