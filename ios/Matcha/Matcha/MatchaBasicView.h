#import <UIKit/UIKit.h>
#import "MatchaView.h"
#import "MatchaProtobuf.h"

@interface MatchaBasicView : UIView <MatchaChildView>
@property (nonatomic, weak) MatchaViewNode *viewNode;
@end
