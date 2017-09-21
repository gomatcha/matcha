#import "MatchaViewController.h"

@interface MatchaViewController (Private)
+ (NSPointerArray *)viewControllers;
+ (MatchaViewController *)viewControllerWithIdentifier:(NSInteger)identifier;
- (void)update:(MatchaViewPBRoot *)node;
@property (nonatomic, readonly) NSInteger identifier;
@property (nonatomic, readonly) BOOL updating;
@end
