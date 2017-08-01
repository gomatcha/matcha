#import <UIKit/UIKit.h>
#import "MatchaView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController.h"

@interface MatchaSegmentView : UISegmentedControl <MatchaChildView>
@property (nonatomic, weak) MatchaViewNode *viewNode;
@property (nonatomic, strong) MatchaBuildNode *node;
@end
