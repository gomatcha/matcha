#import "MatchaView.h"
#import "MatchaProtobuf.h"
#import "MatchaTapGestureRecognizer.h"
#import "MatchaPressGestureRecognizer.h"
#import "MatchaViewController.h"
#import "MatchaSwitchView.h"
#import "MatchaButtonGestureRecognizer.h"
#import "MatchaScrollView.h"

static NSLock *sLock = nil;
static NSMutableDictionary *sViewDict = nil;
static NSMutableDictionary *sViewControllerDict = nil;

void MatchaRegisterInit();
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
    if ([any.typeURL isEqual:@"type.googleapis.com/matcha.touch.TapRecognizer"]) {
        return [[MatchaTapGestureRecognizer alloc] initWithMatchaVC:viewNode.rootVC viewId:viewId protobuf:any];
    } else if ([any.typeURL isEqual:@"type.googleapis.com/matcha.touch.PressRecognizer"]) {
        return [[MatchaPressGestureRecognizer alloc] initWithMatchaVC:viewNode.rootVC viewId:viewId protobuf:any];
    } else if ([any.typeURL isEqual:@"type.googleapis.com/matcha.touch.ButtonRecognizer"]) {
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

@implementation MatchaViewNode

- (id)initWithParent:(MatchaViewNode *)node rootVC:(MatchaViewController *)rootVC identifier:(NSNumber *)identifier {
    if ((self = [super init])) {
        self.parent = node;
        self.identifier = identifier;
        self.rootVC = rootVC; 
    }
    return self;
}

- (void)setRoot:(MatchaNodeRoot *)root {
    MatchaViewPBLayoutPaintNode *pbLayoutPaintNode = [root.layoutPaintNodes objectForKey:self.identifier.longLongValue];
    
    MatchaViewPBBuildNode *pbBuildNode = [root.buildNodes objectForKey:self.identifier.longLongValue];
    MatchaBuildNode *buildNode = nil;
    if (pbBuildNode != nil) {
        buildNode = [[MatchaBuildNode alloc] initWithProtobuf:pbBuildNode];
    }
//    NSAssert(self.buildNode == nil || [self.buildNode.nativeViewName isEqual:buildNode.nativeViewName], @"Node with different name");
    
    if (self.view == nil && self.viewController == nil) {
        self.view = MatchaViewWithNode(buildNode, self);
        self.viewController = MatchaViewControllerWithNode(buildNode, self);
        if (self.view == nil && self.viewController == nil) {
            NSLog(@"Cannot find corresponding view or view controller for node: %@", buildNode.nativeViewName);
        }
        self.view.autoresizingMask = UIViewAutoresizingNone;
        self.view.backgroundColor = [UIColor clearColor];
    }
    
    // Build children
    NSDictionary<NSNumber *, MatchaViewNode *> *children = nil;
    NSMutableArray *addedKeys = [NSMutableArray array];
    NSMutableArray *removedKeys = [NSMutableArray array];
    NSMutableArray *unmodifiedKeys = [NSMutableArray array];
    if (buildNode != nil && ![buildNode.buildId isEqual:self.buildNode.buildId]) {
        for (NSNumber *i in self.children) {
            MatchaBuildNode *child = [root.buildNodes objectForKey:i.longLongValue];
            if (child == nil) {
                [removedKeys addObject:i];
            }
        }
        for (NSInteger i = 0; i < buildNode.childIds.count; i++) {
            NSNumber *childId = @([buildNode.childIds valueAtIndex:i]);
            MatchaViewNode *prevChild = self.children[childId];
            if (prevChild == nil) {
                [addedKeys addObject:childId];
            } else {
                [unmodifiedKeys addObject:childId];
            }
        }
        
        // Add/remove child nodes
        NSMutableDictionary<NSNumber *, MatchaViewNode *> *mutChildren = [NSMutableDictionary dictionary];
        for (NSNumber *i in addedKeys) {
            mutChildren[i] = [[MatchaViewNode alloc] initWithParent:self rootVC:self.rootVC identifier:i];
        }
        for (NSNumber *i in unmodifiedKeys) {
            mutChildren[i] = self.children[i];
        }
        children = mutChildren;
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
            self.view.node = buildNode;
        } else if (self.viewController) {
            self.viewController.node = buildNode;
            
            NSMutableDictionary<NSNumber *, UIViewController *> *childVCs = [NSMutableDictionary dictionary];
            for (NSNumber *i in children) {
                MatchaViewNode *child = children[i];
                childVCs[i] = child.wrappedViewController;
            }
            self.viewController.matchaChildViewControllers = childVCs;
        }
        
        // Add/remove subviews
        for (NSNumber *i in addedKeys) {
            MatchaViewNode *child = children[i];
            child.view.node = [[MatchaBuildNode alloc] initWithProtobuf:[root.buildNodes objectForKey:i.longLongValue]];
            
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
                [child.view removeFromSuperview];
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
        if ([self.parent.view isKindOfClass:[MatchaScrollView class]]) {
            MatchaScrollView *scrollView = (MatchaScrollView *)self.parent.view;
            
            CGPoint origin = f.origin;
            origin.x *= -1;
            origin.y *= -1;
            f.origin = CGPointZero;
            self.materializedView.frame = f;
            scrollView.matchaContentOffset = origin;
            scrollView.contentOffset = origin;
            
        } else if (self.parent.viewController == nil) {
            // let view controllers do their own layout
            if (!CGRectEqualToRect(f, self.frame)) {
                self.materializedView.frame = f;
                self.frame = f;
            }
        } else if (self.viewController) {
            self.viewController.matchaChildLayout = root.layoutPaintNodes;
        }
    }
    
    // Paint view
    if (pbLayoutPaintNode != nil && pbLayoutPaintNode.paintId != self.layoutPaintNode.paintId) {
        if (pbLayoutPaintNode.paintStyle.hasBackgroundColor) {
            self.view.layer.backgroundColor = MatchaCGColorWithProtobuf(pbLayoutPaintNode.paintStyle.backgroundColor);
        } else {
            self.view.backgroundColor = [UIColor clearColor];
        }
        
        self.view.alpha = 1 - pbLayoutPaintNode.paintStyle.transparency;
        self.view.layer.borderColor = MatchaCGColorWithProtobuf(pbLayoutPaintNode.paintStyle.borderColor);
        self.view.layer.borderWidth = pbLayoutPaintNode.paintStyle.borderWidth;
        self.view.layer.cornerRadius = pbLayoutPaintNode.paintStyle.cornerRadius;
        self.view.layer.shadowRadius = pbLayoutPaintNode.paintStyle.shadowRadius;
        self.view.layer.shadowOffset = pbLayoutPaintNode.paintStyle.shadowOffset.toCGSize;
        self.view.layer.shadowColor = MatchaCGColorWithProtobuf(pbLayoutPaintNode.paintStyle.shadowColor);
        self.view.layer.shadowOpacity = pbLayoutPaintNode.paintStyle.hasShadowColor ? 0 : 1;
        if (pbLayoutPaintNode.paintStyle.cornerRadius != 0) {
            self.view.clipsToBounds = YES; // TODO(KD): Be better about this...
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
