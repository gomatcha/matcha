#import <UIKit/UIKit.h>
#import "MatchaView.h"
#import "MatchaProtobuf.h"

@interface MatchaUnknownView : UIView <MatchaChildView>
@property (nonatomic, weak) MatchaViewNode *viewNode;
@property (nonatomic, strong) MatchaBuildNode *node;
@end
