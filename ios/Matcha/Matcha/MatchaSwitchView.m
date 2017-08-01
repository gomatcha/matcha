#import "MatchaSwitchView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController.h"

@implementation MatchaSwitchView

+ (void)load {
    MatchaRegisterView(@"gomatcha.io/matcha/view/switch", ^(MatchaViewNode *node){
        return [[MatchaSwitchView alloc] initWithViewNode:node];
    });
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
    MatchaPBSwitchViewView *view = (id)[state unpackMessageClass:[MatchaPBSwitchViewView class] error:&error];
    if (view != nil) {
        self.on = view.value;
    }
}

- (void)onChange:(id)sender {
    MatchaPBSwitchViewEvent *event = [[MatchaPBSwitchViewEvent alloc] init];
    event.value = self.on;
    
    NSData *data = [event data];
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:data];
    
    [self.viewNode.rootVC call:@"OnChange" viewId:self.node.identifier.longLongValue args:@[value]];
}

@end
