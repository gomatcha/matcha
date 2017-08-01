#import <UIKit/UIKit.h>
@class MatchaViewController;
@class GPBAny;

@interface MatchaButtonGestureRecognizer : UIGestureRecognizer
- (id)initWithMatchaVC:(MatchaViewController *)viewRoot viewId:(int64_t)viewId protobuf:(GPBAny *)pb;
- (void)disable;
- (void)updateWithProtobuf:(GPBAny *)pb;
@end
