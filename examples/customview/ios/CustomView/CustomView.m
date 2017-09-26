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

- (void)setNativeState:(NSData *)nativeState {
    CustomViewProtoView *view = [CustomViewProtoView parseFromData:nativeState error:nil];
    [self setOn:view.value animated:true];
    self.enabled = view.enabled;
}

- (void)onChange:(id)sender {
    CustomViewProtoEvent *event = [[CustomViewProtoEvent alloc] init];
    event.value = self.on;
    [self.viewNode call:@"OnChange", [[MatchaGoValue alloc] initWithData:event.data], nil];
}

@end
