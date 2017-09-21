#import <UIKit/UIKit.h>
#import "MatchaView.h"
#import "MatchaProtobuf.h"

@interface MatchaButton : UIView <MatchaChildView>
@property (nonatomic, strong) UIButton *button;
@property (nonatomic, weak) MatchaViewNode *viewNode;
@end

