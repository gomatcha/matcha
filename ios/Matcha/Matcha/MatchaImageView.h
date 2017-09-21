#import <UIKit/UIKit.h>
#import "MatchaView.h"
#import "MatchaProtobuf.h"

@interface MatchaImageView : UIImageView <MatchaChildView>
@property (nonatomic, weak) MatchaViewNode *viewNode;
@end
