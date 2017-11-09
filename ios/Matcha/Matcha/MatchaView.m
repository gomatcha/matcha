#import "MatchaView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController_Private.h"
#import "MatchaSwitchView.h"
#import "MatchaScrollView.h"
#import "MatchaUnknownView.h"
#import "MatchaView_Private.h"
#import "MatchaBuildNode.h"
#import "MatchaGestureRecognizer.h"

UIView<MatchaChildView> *MatchaViewWithNode(MatchaBuildNode *node, MatchaViewNode *viewNode);
static NSLock *sLock = nil;
static NSMutableDictionary *sViewDict = nil;
static NSMutableDictionary *sViewControllerDict = nil;

void MatchaRegisterInit(void);
void MatchaRegisterInit() {
    static dispatch_once_t sOnce = 0;
    dispatch_once(&sOnce, ^{
        sLock = [[NSLock alloc] init];
        sViewDict = [NSMutableDictionary dictionary];
        sViewControllerDict = [NSMutableDictionary dictionary];
    });
}

void MatchaRegisterView(NSString *string, MatchaViewRegistrationBlock block) {
    MatchaRegisterInit();
    [sLock lock];
    sViewDict[string] = block;
    [sLock unlock];
}

void MatchaRegisterViewController(NSString *string, MatchaViewControllerRegistrationBlock block) {
    MatchaRegisterInit();
    [sLock lock];
    sViewControllerDict[string] = block;
    [sLock unlock];
}

UIView<MatchaChildView> *MatchaViewWithNode(MatchaBuildNode *node, MatchaViewNode *viewNode) {
    NSString *name = node.nativeViewName;
    UIView<MatchaChildView> *child = nil;
    
    [sLock lock];
    MatchaViewRegistrationBlock block = sViewDict[name];
    if (block != nil) {
        child = block(viewNode);
    }
    [sLock unlock];
    return child;
}

UIViewController<MatchaChildViewController> *MatchaViewControllerWithNode(MatchaBuildNode *node, MatchaViewNode *viewNode) {
    NSString *name = node.nativeViewName;
    UIViewController<MatchaChildViewController> *child = nil;
    
    [sLock lock];
    MatchaViewControllerRegistrationBlock block = sViewControllerDict[name];
    if (block != nil) {
        child = block(viewNode);
    }
    [sLock unlock];

    return child;
}

@interface MatchaViewNode ()
- (id)initWithParent:(MatchaViewNode *)node rootVC:(MatchaViewController *)rootVC identifier:(NSNumber *)identifier;
@property (nonatomic, strong) UIView<MatchaChildView> *view;
@property (nonatomic, strong) NSDictionary<NSNumber *, UIGestureRecognizer *> *touchRecognizers;

- (void)setRoot:(MatchaViewPBRoot *)root;
@property (nonatomic, strong) UIViewController<MatchaChildViewController> *viewController;
@property (nonatomic, strong) NSMutableDictionary<NSNumber *, MatchaViewNode *> *children;
@property (nonatomic, strong) MatchaViewPBLayoutPaintNode *layoutPaintNode;
@property (nonatomic, strong) MatchaBuildNode *buildNode;
@property (nonatomic, strong) NSNumber *identifier;
@property (nonatomic, weak) MatchaViewNode *parent;
@property (nonatomic, weak) MatchaViewController *rootVC;
@property (nonatomic, strong) MatchaGestureRecognizer *gestureRecognizer;

@property (nonatomic, strong) UIViewController *wrappedViewController;
- (UIViewController *)materializedViewController;
- (UIViewController *)wrappedViewController;
- (UIView *)materializedView;

@property (nonatomic, assign) CGRect frame;
@end

@implementation MatchaViewNode

- (void)call:(NSString *)funcId, ... NS_REQUIRES_NIL_TERMINATION {
    if (self.rootVC.updating) {
        return;
    }
    va_list args;
    va_start(args, funcId);
    NSArray *rlt = [self.rootVC call:funcId viewId:self.identifier.longLongValue args:args];
    va_end(args);
    //return rlt;
    rlt = nil;
}

- (NSArray<MatchaGoValue *> *)call2:(NSString *)funcId, ... NS_REQUIRES_NIL_TERMINATION {
    va_list args;
    va_start(args, funcId);
    NSArray<MatchaGoValue *> *rlt = [self.rootVC call:funcId viewId:self.identifier.longLongValue args:args];
    va_end(args);
    return rlt[0].toArray;
}

- (id)initWithParent:(MatchaViewNode *)node rootVC:(MatchaViewController *)rootVC identifier:(NSNumber *)identifier {
    if ((self = [super init])) {
        self.parent = node;
        self.identifier = identifier;
        self.rootVC = rootVC; 
    }
    return self;
}

- (void)setRoot:(MatchaViewPBRoot *)root {
    MatchaViewPBLayoutPaintNode *pbLayoutPaintNode = [root.layoutPaintNodes objectForKey:self.identifier.longLongValue];
    
    MatchaViewPBBuildNode *pbBuildNode = [root.buildNodes objectForKey:self.identifier.longLongValue];
    MatchaBuildNode *buildNode = nil;
    if (pbBuildNode != nil) {
        buildNode = [[MatchaBuildNode alloc] initWithProtobuf:pbBuildNode];
    }
    
    // Create view
    if (self.view == nil && self.viewController == nil) {
        self.view = MatchaViewWithNode(buildNode, self);
        self.viewController = MatchaViewControllerWithNode(buildNode, self);
        if (self.view == nil && self.viewController == nil) {
            self.view = [[MatchaUnknownView alloc] initWithViewNode:self];
            NSLog(@"Cannot find corresponding view or view controller for node: %@", buildNode.nativeViewName);
        }
        self.materializedView.autoresizingMask = UIViewAutoresizingNone;
        self.view.backgroundColor = [UIColor clearColor];
    }
    
    // Build children
    NSMutableArray<MatchaViewNode *> *childrenArray = [NSMutableArray array];
    NSMutableDictionary<NSNumber *, MatchaViewNode *> *children = [NSMutableDictionary dictionary];
    NSMutableArray *addedKeys = [NSMutableArray array];
    NSMutableArray *removedKeys = [NSMutableArray array];
    NSMutableArray *unmodifiedKeys = [NSMutableArray array];
    if (buildNode != nil && ![buildNode.buildId isEqual:self.buildNode.buildId]) {        
        for (NSNumber *i in self.children) {
            MatchaViewPBBuildNode *child = [root.buildNodes objectForKey:i.longLongValue];
            if (child == nil) {
                [removedKeys addObject:i];
            }
        }
        for (NSInteger i = 0; i < buildNode.childIds.count; i++) {
            NSNumber *childId = @([buildNode.childIds valueAtIndex:i]);
            MatchaViewNode *prevChild = self.children[childId];
            if (prevChild == nil) {
                [addedKeys addObject:childId];
                MatchaViewNode *n = [[MatchaViewNode alloc] initWithParent:self rootVC:self.rootVC identifier:childId];
                [childrenArray addObject:n];
                children[childId] = n;
            } else {
                [unmodifiedKeys addObject:childId];
                [childrenArray addObject:prevChild];
                children[childId] = prevChild;
            }
        }
    } else {
        children = self.children;
    }
    
    // Update children
    for (NSNumber *i in children) {
        MatchaViewNode *child = children[i];
        [child setRoot:root];
    }
    
    if (buildNode != nil && ![buildNode.buildId isEqual:self.buildNode.buildId]) {
        // Update the views with native values
        if (self.view) {
            self.view.nativeState = buildNode.nativeViewState;
        } else if (self.viewController) {
            self.viewController.nativeState = buildNode.nativeViewState;
            
            NSMutableArray<UIViewController *> *childVCs = [NSMutableArray array];
            for (MatchaViewNode *i in childrenArray) {
                [childVCs addObject:i.wrappedViewController];
            }
            self.viewController.matchaChildViewControllers = childVCs;
        }
        
        // Add/remove subviews
        for (NSNumber *i in addedKeys) {
            MatchaViewNode *child = children[i];
            // child.view.node = [[MatchaBuildNode alloc] initWithProtobuf:[root.buildNodes objectForKey:i.longLongValue]];
            
            if (self.viewController) {
                // no-op. The view controller will handle this itself.
            } else if (child.view) {
                [self.materializedView addSubview:child.view];
            } else if (child.viewController) {
//                [self.materializedViewController addChildViewController:child.viewController]; // TODO(KD): Why can't I add as a child view controller?
                [self.materializedView addSubview:child.viewController.view];
            }
        }
        for (NSNumber *i in removedKeys) {
            MatchaViewNode *child = self.children[i];
            if (self.viewController) {
                // no-op
            } else if (child.view) {
                [child.view removeFromSuperview];
            } else if (child.viewController) {
                [child.materializedView removeFromSuperview];
                [child.viewController removeFromParentViewController];
            }
        }
        
        // Update gesture recognizers
        if (buildNode.nativeValues[@"gomatcha.io/matcha/pointer Gestures"] != nil) {
            if (self.gestureRecognizer == nil) {
                self.gestureRecognizer = [[MatchaGestureRecognizer alloc] init];
                self.gestureRecognizer.viewNode = self;
                [self.materializedView addGestureRecognizer:self.gestureRecognizer];
            }
        } else {
            if (self.gestureRecognizer != nil) {
                [self.materializedView removeGestureRecognizer:self.gestureRecognizer];
                self.gestureRecognizer = nil;
            }
        }
    }

    // Layout subviews
    if (pbLayoutPaintNode != nil && pbLayoutPaintNode.layoutId != self.layoutPaintNode.layoutId) {
        if (self.view) {
            for (NSInteger i = 0; i < pbLayoutPaintNode.childOrderArray.count; i++) {
                NSNumber *key = @([pbLayoutPaintNode.childOrderArray valueAtIndex:i]);
                UIView *subview = children[key].view;
                if ([self.view.subviews indexOfObject:subview] != i) {
                    [self.view insertSubview:subview atIndex:i];
                }
            }
        }
        
        CGRect f = pbLayoutPaintNode.frame;
        if (self.parent == nil) {
        } else if (self.parent.viewController) {
            // let view controllers do their own layout.
        } else if ([self.parent.view isKindOfClass:[MatchaScrollView class]]) {
            
            MatchaScrollView *scrollView = (MatchaScrollView *)self.parent.view;
            CGPoint origin = f.origin;
            origin.x *= -1;
            origin.y *= -1;
            f.origin = CGPointZero;
            self.materializedView.frame = f;
            
            if ((fabs(origin.x - scrollView.matchaContentOffset.x) < 0.5 && fabs(origin.y - scrollView.matchaContentOffset.y) < 0.5)) {
                scrollView.matchaContentOffset = origin;
                scrollView.contentOffset = origin;
                scrollView.contentSize = f.size;
            } else {
                // If Go has independently changed the content offset, cancel any acceleration. Otherwise the view will continue scrolling.
                scrollView.matchaContentOffset = origin;
                [scrollView setContentOffset:origin animated:NO];
                scrollView.contentSize = f.size;
            }
        } else {
            if (!CGRectEqualToRect(f, self.frame)) {
                self.materializedView.frame = f;
                self.frame = f;
            }
        }
        
//        if (self.viewController) {
//            // Give view controllers their children's layout objects.
//            NSMutableArray<MatchaViewPBLayoutPaintNode *> *layoutPaintNodes = [NSMutableArray array];
//            for (MatchaViewNode *i in childrenArray) {
//                [layoutPaintNodes addObject:[root.layoutPaintNodes objectForKey:i.identifier.longLongValue]];
//            }
//            self.viewController.matchaChildLayout = layoutPaintNodes;
//        }
    }
    
    // Paint view
    if (pbLayoutPaintNode != nil && pbLayoutPaintNode.paintId != self.layoutPaintNode.paintId) {
        
        CGColorRef backgroundColor = MatchaCGColorCreateWithValues(
            pbLayoutPaintNode.hasBackgroundColor,
            pbLayoutPaintNode.backgroundColorRed,
            pbLayoutPaintNode.backgroundColorGreen,
            pbLayoutPaintNode.backgroundColorBlue,
            pbLayoutPaintNode.backgroundColorAlpha);
        CGColorRef borderColor = MatchaCGColorCreateWithValues(
            pbLayoutPaintNode.hasBorderColor,
            pbLayoutPaintNode.borderColorRed,
            pbLayoutPaintNode.borderColorGreen,
            pbLayoutPaintNode.borderColorBlue,
            pbLayoutPaintNode.borderColorAlpha);
        CGColorRef shadowColor = MatchaCGColorCreateWithValues(
            pbLayoutPaintNode.hasShadowColor,
            pbLayoutPaintNode.shadowColorRed,
            pbLayoutPaintNode.shadowColorGreen,
            pbLayoutPaintNode.shadowColorBlue,
            pbLayoutPaintNode.shadowColorAlpha);
        
        if (backgroundColor) {
            self.view.layer.backgroundColor = backgroundColor;
            CFRelease(backgroundColor);
        } else {
            self.view.backgroundColor = [UIColor clearColor];
        }
        self.view.alpha = 1 - pbLayoutPaintNode.transparency;
        self.view.layer.borderColor = borderColor;
        self.view.layer.borderWidth = pbLayoutPaintNode.borderWidth;
        self.view.layer.cornerRadius = pbLayoutPaintNode.cornerRadius;
        self.view.layer.shadowRadius = pbLayoutPaintNode.shadowRadius;
        self.view.layer.shadowOffset = CGSizeMake(pbLayoutPaintNode.shadowOffsetX, pbLayoutPaintNode.shadowOffsetY);
        self.view.layer.shadowColor = shadowColor;
        self.view.layer.shadowOpacity = pbLayoutPaintNode.hasShadowColor ? 1 : 0;
        if (pbLayoutPaintNode.cornerRadius != 0) {
            self.view.clipsToBounds = YES; // TODO(KD): Be better about this...
        }
        if (borderColor) {
            CFRelease(borderColor);
        }
        if (shadowColor) {
            CFRelease(shadowColor);
        }
    }
    
    if (pbLayoutPaintNode != nil) {
        _layoutPaintNode = pbLayoutPaintNode;
    }
    if (buildNode != nil) {
        _buildNode = buildNode;
    }
    self.children = children;
}

- (UIViewController *)materializedViewController {
    UIViewController *vc = nil;
    MatchaViewNode *viewNode = self;
    while (vc == nil && viewNode != nil) {
        viewNode = self.parent;
        vc = viewNode.viewController;
    }
    if (vc == nil) {
        vc = self.rootVC;
    }
    return vc;
}

- (UIViewController *)wrappedViewController {
    if (_wrappedViewController) {
        return _wrappedViewController;
    }
    
    if (self.viewController) {
        _wrappedViewController = self.viewController;
        return _wrappedViewController;
    }
    _wrappedViewController = [[UIViewController alloc] initWithNibName:nil bundle:nil];
    _wrappedViewController.view = self.view;
    MatchaConfigureChildViewController(_wrappedViewController);
    return _wrappedViewController;
}

- (UIView *)materializedView {
    return self.viewController.view ?: self.view;
}

@end

void MatchaConfigureChildViewController(UIViewController *vc) {
    vc.edgesForExtendedLayout=UIRectEdgeNone;
    vc.extendedLayoutIncludesOpaqueBars=NO;
    vc.automaticallyAdjustsScrollViewInsets=NO;
}
