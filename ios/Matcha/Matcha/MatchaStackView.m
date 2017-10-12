#import "MatchaStackView.h"
#import "MatchaView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController.h"
#import <objc/runtime.h>
#import "MatchaView_Private.h"

#define VIEW_ID_KEY @"matchaViewId"

@interface UIViewController (MatchaStackScreen)
- (void)matcha_setViewId:(int64_t)value;
- (int64_t)matcha_viewId;
@end

@implementation UIViewController (MatchaStackScreen)

- (void)matcha_setViewId:(int64_t)value {
    @synchronized (self) {
        objc_setAssociatedObject(self, VIEW_ID_KEY, @(value), OBJC_ASSOCIATION_RETAIN);
    }
}

- (int64_t)matcha_viewId {
    @synchronized (self) {
        return ((NSNumber *)objc_getAssociatedObject(self, VIEW_ID_KEY)).longLongValue;
    }
}

@end

@implementation MatchaStackView

+ (void)load {
    MatchaRegisterViewController(@"gomatcha.io/matcha/view/stacknav", ^(MatchaViewNode *node){
        return [[MatchaStackView alloc] initWithViewNode:node];
    });
    MatchaRegisterViewController(@"gomatcha.io/matcha/view/stacknav Bar", ^(MatchaViewNode *node){
        return [[MatchaStackBar alloc] initWithViewNode:node];
    });
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super init])) {
        self.viewNode = viewNode;
        self.delegate = self;
        MatchaConfigureChildViewController(self);
        self.view.backgroundColor = [UIColor whiteColor];
    }
    return self;
}

- (void)setMatchaChildViewControllers:(NSArray<UIViewController *> *)childVCs {
    MatchaiOSPBStackView *view = (id)[MatchaiOSPBStackView parseFromData:self.nativeState error:nil];
    
    self.navigationBar.barTintColor = view.hasBarColor ? [[UIColor alloc] initWithProtobuf:view.barColor] : nil;
    self.navigationBar.titleTextAttributes = view.hasTitleTextStyle ? [NSAttributedString attributesWithProtobuf:view.titleTextStyle] : nil;
//    self.navigationBar.tintColor = view.hasItemColor ? [[UIColor alloc] initWithProtobuf:view.itemColor]: nil;
    if (view.hasBackTextStyle) {
        [[UIBarButtonItem appearance] setTitleTextAttributes:[NSAttributedString attributesWithProtobuf:view.backTextStyle] forState:UIControlStateNormal];
    } else {
        [[UIBarButtonItem appearance] setTitleTextAttributes:nil forState:UIControlStateNormal];
    }
    
    NSMutableArray *prevIds = [NSMutableArray array];
    for (MatchaiOSPBStackChildView *i in view.childrenArray) {
        [prevIds addObject:@(i.screenId)];
    }
    if ([self.prevIds isEqual:prevIds]) {
        return;
    }
    self.prevIds = prevIds;

    NSMutableArray *viewControllers = [NSMutableArray array];
    for (NSInteger i = 0; i < view.childrenArray.count; i++) {
        MatchaiOSPBStackChildView *childView = view.childrenArray[i];
        MatchaStackBar *bar = (id)childVCs[i * 2];
        UIViewController *vc = childVCs[i * 2 + 1];
        bar.contentViewController = vc;
        [vc matcha_setViewId:childView.screenId];
        [viewControllers addObject:vc];
    }
    
    if (self.viewControllers.count == viewControllers.count) {
        [self setViewControllers:viewControllers animated:NO];
    } else {
        [self setViewControllers:viewControllers animated:YES];
    }
    self.prev = viewControllers;
}

//- (void)navigationController:(UINavigationController *)navigationController willShowViewController:(UIViewController *)viewController animated:(BOOL)animated {
//    NSLog(@"willShow");
//}

- (void)navigationController:(UINavigationController *)navigationController didShowViewController:(UIViewController *)viewController animated:(BOOL)animated {
    [self update];
}

- (void)update {
    NSMutableArray *prevIds = [NSMutableArray array];
    for (UIViewController *i in self.childViewControllers) {
        [prevIds addObject:@(i.matcha_viewId)];
    }
    if ([self.prevIds isEqual:prevIds]) {
        return;
    }
    self.prevIds = prevIds;
    
    GPBInt64Array *array = [[GPBInt64Array alloc] init];
    for (NSNumber *i in prevIds) {
        [array addValue:i.longLongValue];
    }
    MatchaiOSPBStackEvent *event = [[MatchaiOSPBStackEvent alloc] init];
    event.idArray = array;
    [self.viewNode call:@"OnChange", [[MatchaGoValue alloc] initWithData:event.data], nil];
}

- (void)setMatchaChildLayout:(GPBInt64ObjectDictionary *)layoutPaintNodes {
    // no-op
}

@end

@implementation MatchaStackBar

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super init])) {
        self.viewNode = viewNode;
    }
    return self;
}

- (void)setContentViewController:(UIViewController *)contentViewController {
    _contentViewController = contentViewController;
    [self reload];
}

- (void)setNativeState:(NSData *)data {
    _nativeState = data;
    [self reload];
}

- (void)reload {
    if (!self.contentViewController || self.nativeState == nil) {
        return;
    }
    MatchaiOSPBStackBar *bar = [MatchaiOSPBStackBar parseFromData:self.nativeState error:nil];
    
    UINavigationItem *item = self.contentViewController.navigationItem;
    item.title = bar.title;
    item.hidesBackButton = bar.backButtonHidden;
    if (bar.customBackButtonTitle) {
        item.backBarButtonItem = [[UIBarButtonItem alloc] initWithTitle:bar.backButtonTitle style:UIBarButtonItemStylePlain target:nil action:nil];
    } else {
        item.backBarButtonItem = nil;
    }
    NSMutableArray *rightViews = [NSMutableArray array];
    for (MatchaiOSPBStackBarItem *i in bar.rightItemsArray) {
        UIBarButtonItem *item = [[UIBarButtonItem alloc] initWithProtobuf:i];
        [item setTarget:self];
        [item setAction:@selector(onPress:)];
        [rightViews addObject:item];
    }
    item.rightBarButtonItems = rightViews;
    
    NSMutableArray *leftViews = [NSMutableArray array];
    for (MatchaiOSPBStackBarItem *i in bar.leftItemsArray) {
        UIBarButtonItem *item = [[UIBarButtonItem alloc] initWithProtobuf:i];
        [item setTarget:self];
        [item setAction:@selector(onPress:)];
        [leftViews addObject:item];
    }
    item.leftBarButtonItems = leftViews;
}

- (void)setMatchaChildViewControllers:(NSArray<UIViewController *> *)childVCs {
    //    if (bar.hasTitleView) {
    //        self.titleView = childVCs[idx].view;
    //        idx += 1;
    //    } else {
    //        self.titleView = nil;
    //    }
    //    for (NSInteger i = 0; i < bar.rightViewCount; i++) {
    //        UIView *rightView = childVCs[idx].view;
    //        UIBarButtonItem *item = [[UIBarButtonItem alloc] initWithCustomView:rightView];
    //        [rightViews addObject:item];
    //        idx += 1;
    //    }
    //    for (NSInteger i = 0; i < bar.leftViewCount; i++) {
    //        UIView *leftView = childVCs[idx].view;
    //        UIBarButtonItem *item = [[UIBarButtonItem alloc] initWithCustomView:leftView];
    //        [leftViews addObject:item];
    //        idx +=1;
    //    }
}

- (void)setMatchaChildLayout:(NSArray<MatchaViewPBLayoutPaintNode *> *)layoutPaintNodes {
//    NSInteger idx = 0;
//    if (self.titleView) {
//        CGRect f = self.titleView.frame;
//        f.size = ((MatchaViewPBLayoutPaintNode *)layoutPaintNodes[idx]).frame.size;
//        self.titleView.frame = f;
//        idx += 1;
//    }
//    for (NSInteger i = 0; i < self.rightViews.count; i++) {
//        UIBarButtonItem *rightView = self.rightViews[i];
//        CGRect f = rightView.customView.frame;
//        f.size =((MatchaViewPBLayoutPaintNode *)layoutPaintNodes[idx]).frame.size;
//        rightView.customView.frame = f;
//        idx += 1;
//    }
//    for (NSInteger i = 0; i < self.leftViews.count; i++) {
//        UIBarButtonItem *leftView = self.leftViews[i];
//        CGRect f = leftView.customView.frame;
//        f.size =((MatchaViewPBLayoutPaintNode *)layoutPaintNodes[idx]).frame.size;
//        leftView.customView.frame = f;
//        idx += 1;
//    }
}

- (void)onPress:(UIBarButtonItem *)sender {
    [self.viewNode call:sender.onPress, nil];
}

@end

@implementation UIBarButtonItem (Protobuf)

static char defaultHashKey;

- (NSString *)onPress {
    return objc_getAssociatedObject(self, &defaultHashKey) ;
}

- (void)setOnPress:(NSString *)onPress {
    objc_setAssociatedObject(self, &defaultHashKey, onPress, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
}

- (id)initWithProtobuf:(MatchaiOSPBStackBarItem *)proto {
    if (proto.hasImage) {
        UIImage *image = [[UIImage alloc] initWithImageOrResourceProtobuf:proto.image];
        if (!proto.tintsImage) {
            image = [self.image imageWithRenderingMode:UIImageRenderingModeAlwaysOriginal];
        }
        self = [self initWithImage:image style:UIBarButtonItemStylePlain target:nil action:nil];
        self.tintColor = proto.hasTintColor ? [[UIColor alloc] initWithProtobuf:proto.tintColor] : nil;
    } else {
        self = [self initWithTitle:proto.title style:UIBarButtonItemStylePlain target:nil action:nil];
        if (proto.hasTitleStyle) {
            NSDictionary *attributes = [NSAttributedString attributesWithProtobuf:proto.titleStyle];
            [self setTitleTextAttributes:attributes forState:UIControlStateNormal];
            [self setTitleTextAttributes:attributes forState:UIControlStateHighlighted];
        }
    }
    self.enabled = proto.enabled;
    self.onPress = proto.onPress;
    return self;
}
@end
