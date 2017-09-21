#import "CustomView.h"
#import "Customview.pbobjc.h"

@implementation CustomView

+ (void)load {
    [MatchaViewController registerView:@"github.com/overcyn/customview" block:^(MatchaViewNode *node){
        return [[CustomView alloc] initWithViewNode:node];
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
    CustomViewProtoView *view = (id)[state unpackMessageClass:[CustomViewProtoView class] error:&error];
    if (view != nil) {
        [self setOn:view.value animated:true];
        self.enabled = view.enabled;
    }
}

- (void)onChange:(id)sender {
    CustomViewProtoEvent *event = [[CustomViewProtoEvent alloc] init];
    event.value = self.on;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:event.data];
    [self.viewNode.rootVC call:@"OnChange" viewId:self.node.identifier.longLongValue args:@[value]];
}

@end
