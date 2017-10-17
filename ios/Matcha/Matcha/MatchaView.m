#import "MatchaView.h"
#import "MatchaProtobuf.h"
#import "MatchaTapGestureRecognizer.h"
#import "MatchaPressGestureRecognizer.h"
#import "MatchaViewController_Private.h"
#import "MatchaSwitchView.h"
#import "MatchaButtonGestureRecognizer.h"
#import "MatchaScrollView.h"
#import "MatchaUnknownView.h"
#import "MatchaView_Private.h"
#import "MatchaBuildNode.h"

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

UIGestureRecognizer *MatchaGestureRecognizerWithPB(int64_t viewId, GPBAny *any, MatchaViewNode *viewNode) {
    if ([any.typeURL isEqual:@"type.googleapis.com/matcha.pointer.TapRecognizer"]) {
        return [[MatchaTapGestureRecognizer alloc] initWithMatchaVC:viewNode.rootVC viewId:viewId protobuf:any];
    } else if ([any.typeURL isEqual:@"type.googleapis.com/matcha.pointer.PressRecognizer"]) {
        return [[MatchaPressGestureRecognizer alloc] initWithMatchaVC:viewNode.rootVC viewId:viewId protobuf:any];
    } else if ([any.typeURL isEqual:@"type.googleapis.com/matcha.pointer.ButtonRecognizer"]) {
        return [[MatchaButtonGestureRecognizer alloc] initWithMatchaVC:viewNode.rootVC viewId:viewId protobuf:any];
    }
    return nil;
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

@property (nonatomic, strong) UIViewController *wrappedViewController;
- (UIViewController *)materializedViewController;
- (UIViewController *)wrappedViewController;
- (UIView *)materializedView;

@property (nonatomic, assign) CGRect frame;
@end

@implementation MatchaViewNode

//- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId, ... NS_REQUIRES_NIL_TERMINATION {
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
        if (self.view) {
            NSMutableArray *addedKeys = [NSMutableArray array];
            NSMutableArray *removedKeys = [NSMutableArray array];
            NSMutableArray *unmodifiedKeys = [NSMutableArray array];
            for (NSNumber *i in self.buildNode.touchRecognizers) {
                GPBAny *child = buildNode.touchRecognizers[i];
                if (child == nil) {
                    [removedKeys addObject:i];
                }
            }
            for (NSNumber *i in buildNode.touchRecognizers) {
                GPBAny *prevChild = self.buildNode.touchRecognizers[i];
                if (prevChild == nil) {
                    [addedKeys addObject:i];
                } else {
                    [unmodifiedKeys addObject:i];
                }
            }
            
            NSMutableDictionary *touchRecognizers = [NSMutableDictionary dictionary];
            for (NSNumber *i in removedKeys) {
                UIGestureRecognizer *recognizer = self.touchRecognizers[i];
                [(id)recognizer disable];
                [self.view removeGestureRecognizer:recognizer];
            }
            for (NSNumber *i in addedKeys) {
                UIGestureRecognizer *recognizer = MatchaGestureRecognizerWithPB(buildNode.identifier.longLongValue, buildNode.touchRecognizers[i], self);
                [self.view addGestureRecognizer:recognizer];
                touchRecognizers[i] = recognizer;
            }
            for (NSNumber *i in unmodifiedKeys) {
                UIGestureRecognizer *recognizer = self.touchRecognizers[i];
                [(id)recognizer updateWithProtobuf:buildNode.touchRecognizers[i]];
                touchRecognizers[i] = recognizer;
            }
            self.touchRecognizers = touchRecognizers;
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
        
        if (self.viewController) {
            // Give view controllers their children's layout objects.
            NSMutableArray<MatchaViewPBLayoutPaintNode *> *layoutPaintNodes = [NSMutableArray array];
            for (MatchaViewNode *i in childrenArray) {
                [layoutPaintNodes addObject:[root.layoutPaintNodes objectForKey:i.identifier.longLongValue]];
            }
            self.viewController.matchaChildLayout = layoutPaintNodes;
        }
    }
    
    // Paint view
    if (pbLayoutPaintNode != nil && pbLayoutPaintNode.paintId != self.layoutPaintNode.paintId) {
        MatchaPaintPBStyle *style = pbLayoutPaintNode.paintStyle;
        
        CGColorRef backgroundColor = MatchaCGColorCreateWithValues(
            style.hasBackgroundColor,
            style.backgroundColorRed,
            style.backgroundColorGreen,
            style.backgroundColorBlue,
            style.backgroundColorAlpha);
        CGColorRef borderColor = MatchaCGColorCreateWithValues(
            style.hasBorderColor,
            style.borderColorRed,
            style.borderColorGreen,
            style.borderColorBlue,
            style.borderColorAlpha);
        CGColorRef shadowColor = MatchaCGColorCreateWithValues(
            style.hasShadowColor,
            style.shadowColorRed,
            style.shadowColorGreen,
            style.shadowColorBlue,
            style.shadowColorAlpha);
        
        if (backgroundColor) {
            self.view.layer.backgroundColor = backgroundColor;
            CFRelease(backgroundColor);
        } else {
            self.view.backgroundColor = [UIColor clearColor];
        }
        self.view.alpha = 1 - style.transparency;
        self.view.layer.borderColor = borderColor;
        self.view.layer.borderWidth = style.borderWidth;
        self.view.layer.cornerRadius = style.cornerRadius;
        self.view.layer.shadowRadius = style.shadowRadius;
        self.view.layer.shadowOffset = CGSizeMake(style.shadowOffsetX, style.shadowOffsetY);
        self.view.layer.shadowColor = shadowColor;
        self.view.layer.shadowOpacity = style.hasShadowColor ? 1 : 0;
        if (style.cornerRadius != 0) {
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
