#import <UIKit/UIKit.h>
#import <MatchaBridge/MatchaBridge.h>
@class MatchaViewPBRoot;
@class MatchaViewNode;
@protocol MatchaChildView;
@protocol MatchaChildViewController;

typedef UIView<MatchaChildView> *(^MatchaViewRegistrationBlock)(MatchaViewNode *);
typedef UIViewController<MatchaChildViewController> *(^MatchaViewControllerRegistrationBlock)(MatchaViewNode *);

@interface MatchaViewController : UIViewController
- (id)initWithGoValue:(MatchaGoValue *)value;
- (void)printViewHierarchy;
@property (nonatomic, assign) BOOL printViewHierarchyOnUpdate;
+ (void)registerView:(NSString *)viewName block:(MatchaViewRegistrationBlock)block;
+ (void)registerViewController:(NSString *)viewName block:(MatchaViewControllerRegistrationBlock)block;
@end
