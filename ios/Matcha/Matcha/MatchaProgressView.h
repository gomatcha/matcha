#import <UIKit/UIKit.h>
#import "MatchaView.h"

@interface MatchaProgressView : UIProgressView <MatchaChildView>
@property (nonatomic, weak) MatchaViewNode *viewNode;
@property (nonatomic, strong) MatchaBuildNode *node;
@end
