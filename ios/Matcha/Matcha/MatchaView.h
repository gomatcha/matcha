#import <UIKit/UIKit.h>
#import <Matcha/MatchaNode.h>
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
- (void)setMatchaChildViewControllers:(NSDictionary<NSNumber *, UIViewController *> *)childVCs;
- (void)setMatchaChildLayout:(GPBInt64ObjectDictionary *)layoutPaintNodes;
@end

typedef UIView<MatchaChildView> *(^MatchaViewRegistrationBlock)(MatchaViewNode *);
typedef UIViewController<MatchaChildViewController> *(^MatchaViewControllerRegistrationBlock)(MatchaViewNode *);

UIGestureRecognizer *MatchaGestureRecognizerWithPB(int64_t viewId, GPBAny *any, MatchaViewNode *viewNode);
UIView<MatchaChildView> *MatchaViewWithNode(MatchaBuildNode *node, MatchaViewNode *viewNode);
UIViewController<MatchaChildViewController> *MatchaViewControllerWithNode(MatchaBuildNode *node, MatchaViewNode *viewNode);
void MatchaRegisterView(NSString *string, MatchaViewRegistrationBlock block);
void MatchaRegisterViewController(NSString *string, MatchaViewControllerRegistrationBlock block);

@interface MatchaViewNode : NSObject
- (id)initWithParent:(MatchaViewNode *)node rootVC:(MatchaViewController *)rootVC identifier:(NSNumber *)identifier;
@property (nonatomic, strong) UIView<MatchaChildView> *view;
@property (nonatomic, strong) NSDictionary<NSNumber *, UIGestureRecognizer *> *touchRecognizers;

@property (nonatomic, strong) UIViewController<MatchaChildViewController> *viewController;
@property (nonatomic, strong) NSDictionary<NSNumber *, MatchaViewNode *> *children;
- (void)setRoot:(MatchaNodeRoot *)root;
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
