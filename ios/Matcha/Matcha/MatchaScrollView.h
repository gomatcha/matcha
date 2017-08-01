#import <UIKit/UIKit.h>
#import "MatchaView.h"

@interface MatchaScrollView : UIScrollView <MatchaChildView, UIScrollViewDelegate>
@property (nonatomic, weak) MatchaViewNode *viewNode;
@property (nonatomic, strong) MatchaBuildNode *node;
@property (nonatomic, assign) CGPoint matchaContentOffset;
@property (nonatomic, assign) BOOL scrollEvents;
@end
