#import <UIKit/UIKit.h>
#import "MatchaView.h"
#import "MatchaProtobuf.h"

@interface MatchaSlider : UISlider <MatchaChildView>
@property (nonatomic, weak) MatchaViewNode *viewNode;
@end
