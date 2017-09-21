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

- (void)setNode:(MatchaBuildNode *)value {
    _node = value;
    GPBAny *state = value.nativeViewState;
    NSError *error = nil;
    MatchaViewPbSwitchView *view = (id)[state unpackMessageClass:[MatchaViewPbSwitchView class] error:&error];
    if (view != nil) {
        [self setOn:view.value animated:true];
        self.enabled = view.enabled;
    }
}

- (void)onChange:(id)sender {
    MatchaViewPbSwitchEvent *event = [[MatchaViewPbSwitchEvent alloc] init];
    event.value = self.on;
    [self.viewNode call:@"OnChange" args:[[MatchaGoValue alloc] initWithData:event.data], nil];
}

@end
