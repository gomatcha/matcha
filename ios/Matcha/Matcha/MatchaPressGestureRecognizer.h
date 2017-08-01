#import <UIKit/UIKit.h>
@class GPBAny;
@class MatchaViewController;

@interface MatchaPressGestureRecognizer : UILongPressGestureRecognizer
- (id)initWithMatchaVC:(MatchaViewController *)viewRoot viewId:(int64_t)viewId protobuf:(GPBAny *)pb;
- (void)updateWithProtobuf:(GPBAny *)pb;
- (void)disable;
@end
