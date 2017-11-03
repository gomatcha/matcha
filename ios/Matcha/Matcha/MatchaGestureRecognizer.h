#import <UIKit/UIKit.h>
@class MatchaViewNode;

@interface MatchaGestureRecognizer : UIGestureRecognizer
@property (nonatomic, weak) MatchaViewNode *viewNode;
@end
