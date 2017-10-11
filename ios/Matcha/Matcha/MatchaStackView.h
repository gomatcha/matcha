#import <UIKit/UIKit.h>
#import "MatchaView.h"
@class MatchaViewNode;
@class GPBInt64Array;
@class MatchaiOSPBStackBarItem;

@interface MatchaStackView : UINavigationController <MatchaChildViewController, UINavigationControllerDelegate>
- (id)initWithViewNode:(MatchaViewNode *)viewNode;
@property (nonatomic, weak) MatchaViewNode *viewNode;
@property (nonatomic, strong) NSData *nativeState;

//Internal
@property (nonatomic, strong) NSArray<NSNumber *> *prevIds;
@property (nonatomic, strong) NSArray *prev;
@end

@interface MatchaStackBar : UIViewController <MatchaChildViewController>
- (id)initWithViewNode:(MatchaViewNode *)viewNode;
@property (nonatomic, weak) MatchaViewNode *viewNode;
@property (nonatomic, strong) NSData *nativeState;
@property (nonatomic, weak) UIViewController *contentViewController;
@end

@interface UIBarButtonItem (Protobuf)
@property (nonatomic, strong) NSString *onPress;
- (id)initWithProtobuf:(MatchaiOSPBStackBarItem *)proto;
@end
