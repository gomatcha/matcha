#import "MatchaViewController.h"
#import "MatchaView.h"
#import "MatchaBridge.h"
#import "MatchaNode.h"
#import "MatchaObjcBridge.h"
#import "MatchaProtobuf.h"

@interface MatchaViewController ()
@property (nonatomic, assign) NSInteger identifier;
@property (nonatomic, strong) MatchaViewNode *viewNode;
@property (nonatomic, strong) MatchaGoValue *goValue;
@property (nonatomic, assign) CGRect lastFrame;
@property (nonatomic, assign) BOOL loaded;
@property (nonatomic, assign) BOOL statusbarhidden;
@property (nonatomic, assign) UIStatusBarStyle statusbarstyle;
@end

@implementation MatchaViewController

+ (NSPointerArray *)viewControllers {
    static NSPointerArray *sPointerArray;
    static dispatch_once_t sOnce;
    dispatch_once(&sOnce, ^{
        sPointerArray = [NSPointerArray weakObjectsPointerArray];
    });
    return sPointerArray;
}

+ (MatchaViewController *)viewControllerWithIdentifier:(NSInteger)identifier {
    for (MatchaViewController *i in [self viewControllers]) {
        if (i.identifier == identifier) {
            return i;
        }
    }
    return nil;
}

- (id)initWithGoValue:(MatchaGoValue *)value2 {
    if ((self = [super initWithNibName:nil bundle:nil])) {
        [MatchaObjcBridge_X configure];
        [[MatchaObjcBridge sharedBridge] setObject:[MatchaObjcBridge_X new] forKey:@""];
        
        MatchaGoValue *value = [[[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/view NewRoot"] call:nil args:@[value2]][0];
        self.goValue = value;
        self.identifier = (int)[value call:@"Id" args:nil][0].toLongLong;
        [[MatchaViewController viewControllers] addPointer:(__bridge void *)self];
        self.viewNode = [[MatchaViewNode alloc] initWithParent:nil rootVC:self identifier:@([value call:@"ViewId" args:nil][0].toLongLong)];
        self.edgesForExtendedLayout = UIRectEdgeNone;
        self.extendedLayoutIncludesOpaqueBars=NO;
        self.automaticallyAdjustsScrollViewInsets=NO;
    }
    return self;
}

- (void)dealloc {
    [self.goValue call:@"Stop" args:nil];
}

- (void)viewDidLayoutSubviews {
    if (!CGRectEqualToRect(self.lastFrame, self.view.frame)) {
        self.lastFrame = self.view.frame;
        
        MatchaGoValue *width = [[MatchaGoValue alloc] initWithDouble:self.view.frame.size.width];
        MatchaGoValue *height = [[MatchaGoValue alloc] initWithDouble:self.view.frame.size.height];
        [self.goValue call:@"SetSize" args:@[width, height]];
    }
}

- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId viewId:(int64_t)viewId args:(NSArray<MatchaGoValue *> *)args {
    MatchaGoValue *goValue = [[MatchaGoValue alloc] initWithString:funcId];
    MatchaGoValue *goViewId = [[MatchaGoValue alloc] initWithLongLong:viewId];
    MatchaGoValue *goArgs = [[MatchaGoValue alloc] initWithArray:args];
    return [self.goValue call:@"Call" args:@[goValue, goViewId, goArgs]];
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
@end

void MatchaConfigureChildViewController(UIViewController *vc) {
    vc.edgesForExtendedLayout=UIRectEdgeNone;
    vc.extendedLayoutIncludesOpaqueBars=NO;
    vc.automaticallyAdjustsScrollViewInsets=NO;
}

bool MatchaColorEqualToColor(MatchaColor a, MatchaColor b) {
    return a.red == b.red && a.blue == b.blue && a.green == b.green && a.alpha == b.alpha;
}
