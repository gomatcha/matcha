#import <UIKit/UIKit.h>
#import <Matcha/MatchaBuildNode.h>
#import <Matcha/MatchaViewController.h>
@class MatchaViewConfig;
@class MatchaViewController;
@class MatchaViewNode;

@protocol MatchaChildView <NSObject>
- (id)initWithViewNode:(MatchaViewNode *)viewNode;
- (void)setNode:(MatchaBuildNode *)node;
@end

@protocol MatchaChildViewController <NSObject>
- (id)initWithViewNode:(MatchaViewNode *)viewNode;
- (void)setNode:(MatchaBuildNode *)node;
- (void)setMatchaChildViewControllers:(NSArray<UIViewController *> *)childVCs;
- (void)setMatchaChildLayout:(NSMutableArray<MatchaViewPBLayoutPaintNode *> *)layoutPaintNodes;
@end

@interface MatchaViewNode : NSObject
- (id)initWithParent:(MatchaViewNode *)node rootVC:(MatchaViewController *)rootVC identifier:(NSNumber *)identifier;
@property (nonatomic, strong) UIView<MatchaChildView> *view;
@property (nonatomic, strong) NSDictionary<NSNumber *, UIGestureRecognizer *> *touchRecognizers;

@property (nonatomic, strong) UIViewController<MatchaChildViewController> *viewController;
@property (nonatomic, strong) NSMutableDictionary<NSNumber *, MatchaViewNode *> *children;
- (void)setRoot:(MatchaViewPBRoot *)root;
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
