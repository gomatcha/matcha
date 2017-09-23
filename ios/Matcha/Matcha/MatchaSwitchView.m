#import "MatchaSwitchView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController.h"

@implementation MatchaSwitchView

+ (void)load {
    [MatchaViewController registerView:@"gomatcha.io/matcha/view/switch" block:^(MatchaViewNode *node){
        return [[MatchaSwitchView alloc] initWithViewNode:node];
    }];
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        [self addTarget:self action: @selector(onChange:) forControlEvents: UIControlEventValueChanged];
    }
    return self;
}

- (void)setNativeState:(NSData *)nativeState {
    MatchaViewPbSwitchView *view = [MatchaViewPbSwitchView parseFromData:nativeState error:nil];
    [self setOn:view.value animated:true];
    self.enabled = view.enabled;
}

- (void)onChange:(id)sender {
    MatchaViewPbSwitchEvent *event = [[MatchaViewPbSwitchEvent alloc] init];
    event.value = self.on;
    [self.viewNode call:@"OnChange", [[MatchaGoValue alloc] initWithData:event.data], nil];
}

@end
