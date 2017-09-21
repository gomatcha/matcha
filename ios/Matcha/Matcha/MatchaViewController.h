#import <UIKit/UIKit.h>
#import <MatchaBridge/MatchaBridge.h>
@class MatchaViewPBRoot;
@class MatchaViewNode;
@protocol MatchaChildView;
@protocol MatchaChildViewController;

typedef UIView<MatchaChildView> *(^MatchaViewRegistrationBlock)(MatchaViewNode *);
typedef UIViewController<MatchaChildViewController> *(^MatchaViewControllerRegistrationBlock)(MatchaViewNode *);

@interface MatchaViewController : UIViewController // view.Root
- (id)initWithGoValue:(MatchaGoValue *)value;
- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId viewId:(int64_t)viewId args:(NSArray<MatchaGoValue *> *)args;
+ (void)registerView:(NSString *)viewName block:(MatchaViewRegistrationBlock)block;
+ (void)registerViewController:(NSString *)viewName block:(MatchaViewControllerRegistrationBlock)block;
@end
