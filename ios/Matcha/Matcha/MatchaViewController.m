#import "MatchaViewController.h"
#import "MatchaViewController_Private.h"
#import "MatchaView.h"
#import "MatchaBridge.h"
#import "MatchaBuildNode.h"
#import "MatchaObjcBridge.h"
#import "MatchaProtobuf.h"
#import "MatchaView_Private.h"

@interface MatchaViewController ()
@property (nonatomic, assign) NSInteger identifier;
@property (nonatomic, strong) MatchaViewNode *viewNode;
@property (nonatomic, strong) MatchaGoValue *goValue;
@property (nonatomic, assign) CGRect lastFrame;
@property (nonatomic, assign) BOOL loaded;
@property (nonatomic, assign) BOOL statusbarhidden;
@property (nonatomic, assign) UIStatusBarStyle statusbarstyle;
@property (nonatomic, assign) BOOL updating;
@end

@implementation MatchaViewController

- (id)initWithGoValue:(MatchaGoValue *)value2 {
    if ((self = [super initWithNibName:nil bundle:nil])) {
        MatchaGoValue *value = [[[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/view NewRoot"] call:nil, value2, nil][0];
        self.goValue = value;
        self.identifier = (int)[value call:@"Id", nil][0].toLongLong;
        self.viewNode = [[MatchaViewNode alloc] initWithParent:nil rootVC:self identifier:@([value call:@"ViewId", nil][0].toLongLong)];
        self.edgesForExtendedLayout = UIRectEdgeNone;
        self.extendedLayoutIncludesOpaqueBars=NO;
        self.automaticallyAdjustsScrollViewInsets=NO;
        
        [MatchaObjcBridge_X configure];
        [[MatchaObjcBridge_X viewControllers] setObject:self forKey:@(self.identifier)];
    }
    return self;
}

- (void)dealloc {
    [self.goValue call:@"Stop", nil];
}

- (void)viewDidLayoutSubviews {
    if (!CGRectEqualToRect(self.lastFrame, self.view.frame)) {
        self.lastFrame = self.view.frame;
        
        MatchaGoValue *width = [[MatchaGoValue alloc] initWithDouble:self.view.frame.size.width];
        MatchaGoValue *height = [[MatchaGoValue alloc] initWithDouble:self.view.frame.size.height];
        [self.goValue call:@"SetSize", width, height, nil];
    }
}

- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId viewId:(int64_t)viewId args2:(NSArray *)args {
    MatchaGoValue *goValue = [[MatchaGoValue alloc] initWithString:funcId];
    MatchaGoValue *goViewId = [[MatchaGoValue alloc] initWithLongLong:viewId];
    MatchaGoValue *goArgs = [[MatchaGoValue alloc] initWithArray:args];
    return [self.goValue call:@"Call", goValue, goViewId, goArgs, nil];
}

- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId viewId:(int64_t)viewId args:(va_list)args {
    MatchaGoValue *goValue = [[MatchaGoValue alloc] initWithString:funcId];
    MatchaGoValue *goViewId = [[MatchaGoValue alloc] initWithLongLong:viewId];
    NSMutableArray *array = [NSMutableArray array];
    id arg = nil;
    while ((arg = va_arg(args, id))) {
        [array addObject:arg];
    }
    return [self.goValue call:@"Call", goValue, goViewId, [[MatchaGoValue alloc] initWithArray:array], nil];
}

- (void)update:(MatchaViewPBRoot *)root {
    self.updating = true;
    [self.viewNode setRoot:root];
    
    GPBAny *any = root.middleware[@"gomatcha.io/matcha/app activity"];
    if (any) {
        MatchaAppPBActivityIndicator *indicator = (id)[any unpackMessageClass:[MatchaAppPBActivityIndicator class] error:NULL];
        [UIApplication.sharedApplication setNetworkActivityIndicatorVisible:indicator.visible];
    }
    
    any = root.middleware[@"gomatcha.io/matcha/app statusbar"];
    if (any) {
        MatchaAppPBStatusBar *statusBar = (id)[any unpackMessageClass:[MatchaAppPBStatusBar class] error:NULL];
        UIStatusBarStyle style = 0;
        if (statusBar.style == MatchaAppPBStatusBarStyle_StatusBarStyleDefault) {
            style = UIStatusBarStyleDefault;
        } else if (statusBar.style == MatchaAppPBStatusBarStyle_StatusBarStyleLight) {
            style = UIStatusBarStyleLightContent;
        } else if (statusBar.style == MatchaAppPBStatusBarStyle_StatusBarStyleDark) {
            style = UIStatusBarStyleDefault;
        }
        self.statusbarstyle = style;
        self.statusbarhidden = statusBar.hidden;
        [self setNeedsStatusBarAppearanceUpdate];
    }
    
    if (!self.loaded) {
        self.loaded = TRUE;
        UIView *view = self.viewNode.view ?: self.viewNode.viewController.view;
        [self.view addSubview:view];
        self.view.autoresizingMask = UIViewAutoresizingFlexibleWidth|UIViewAutoresizingFlexibleHeight;
        view.frame = self.view.bounds;
    }
    self.updating = false;
}

- (UIStatusBarStyle)preferredStatusBarStyle {
    return self.statusbarstyle;
}

- (BOOL)prefersStatusBarHidden {
    return self.statusbarhidden;
}

- (void)printViewHierarchy {
    [self.goValue call:@"PrintDebug", nil];
}

- (void)setPrintViewHierarchyOnUpdate:(BOOL)val {
    _printViewHierarchyOnUpdate = val;
    [self.goValue call:@"SetPrintDebug", [[MatchaGoValue alloc] initWithBool:val], nil];
}

- (void)motionEnded:(UIEventSubtype)motion withEvent:(UIEvent *)event {
    if (motion == UIEventSubtypeMotionShake) {
        [[[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/application OnShake"] call:nil, nil];
    }
}

- (BOOL)canBecomeFirstResponder {
    return YES;
}

+ (void)registerView:(NSString *)viewName block:(MatchaViewRegistrationBlock)block {
    MatchaRegisterView(viewName, block);
}

+ (void)registerViewController:(NSString *)viewName block:(MatchaViewControllerRegistrationBlock)block {
    MatchaRegisterViewController(viewName, block);
}

@end
