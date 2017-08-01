#import <UIKit/UIKit.h>
#import "MatchaView.h"

@interface MatchaSwitchView : UISwitch <MatchaChildView>
@property (nonatomic, weak) MatchaViewNode *viewNode;
@property (nonatomic, strong) MatchaBuildNode *node;
@end
