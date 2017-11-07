#import <UIKit/UIKit.h>
@class MatchaViewNode;

@interface MatchaGestureRecognizer : UIGestureRecognizer
@property (nonatomic, weak) MatchaViewNode *viewNode;
@property (nonatomic, strong) dispatch_source_t timer;
@end
