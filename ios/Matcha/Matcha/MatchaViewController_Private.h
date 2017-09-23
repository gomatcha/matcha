#import "MatchaViewController.h"

@interface MatchaViewController (Private)
+ (NSPointerArray *)viewControllers;
+ (MatchaViewController *)viewControllerWithIdentifier:(NSInteger)identifier;
- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId viewId:(int64_t)viewId args:(va_list)args;
- (NSArray<MatchaGoValue *> *)call:(NSString *)funcId viewId:(int64_t)viewId args2:(NSArray *)args;
- (void)update:(MatchaViewPBRoot *)node;
@property (nonatomic, readonly) NSInteger identifier;
@property (nonatomic, readonly) BOOL updating;
@end
