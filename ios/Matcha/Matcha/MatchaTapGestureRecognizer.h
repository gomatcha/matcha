#import <UIKit/UIKit.h>
#import "MatchaProtobuf.h"
@class MatchaViewRoot;
@class MatchaViewController;

@interface MatchaTapGestureRecognizer : UITapGestureRecognizer
- (id)initWithMatchaVC:(MatchaViewController *)viewRoot viewId:(int64_t)viewId protobuf:(GPBAny *)pb;
- (void)disable;
- (void)updateWithProtobuf:(GPBAny *)pb;
@end
