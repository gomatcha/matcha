#import <UIKit/UIKit.h>
#import <MatchaBridge/MatchaBridge.h>
@class MatchaViewPBRoot;
@class MatchaViewNode;
@protocol MatchaChildView;
@protocol MatchaChildViewController;

@interface MatchaViewController : UIViewController // view.Root
+ (NSPointerArray *)viewControllers;
+ (MatchaViewController *)viewControllerWithIdentifier:(NSInteger)identifier;

- (id)initWithGoValue:(MatchaGoValue *)value;
- (void)update:(MatchaViewPBRoot *)node;
- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId viewId:(int64_t)viewId args:(NSArray<MatchaGoValue *> *)args;
@property (nonatomic, readonly) NSInteger identifier;
@property (nonatomic, assign) BOOL updating;
@end

typedef UIView<MatchaChildView> *(^MatchaViewRegistrationBlock)(MatchaViewNode *);
typedef UIViewController<MatchaChildViewController> *(^MatchaViewControllerRegistrationBlock)(MatchaViewNode *);
