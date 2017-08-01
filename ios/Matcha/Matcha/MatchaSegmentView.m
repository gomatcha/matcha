#import "MatchaSegmentView.h"

@implementation MatchaSegmentView

+ (void)load {
    MatchaRegisterView(@"gomatcha.io/matcha/view/segmentview", ^(MatchaViewNode *node){
        return [[MatchaSegmentView alloc] initWithViewNode:node];
    });
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        [self addTarget:self action: @selector(onChange:) forControlEvents:UIControlEventValueChanged];
    }
    return self;
}

- (void)setNode:(MatchaBuildNode *)value {
    _node = value;
    MatchaSegmentViewPbView *view = (id)[value.nativeViewState unpackMessageClass:[MatchaSegmentViewPbView class] error:nil];
    
    [self removeAllSegments];
    for (NSInteger i = 0; i < view.titlesArray.count; i++) {
        [self insertSegmentWithTitle:view.titlesArray[i] atIndex:i animated:NO];
    }
    self.selectedSegmentIndex = (int)view.value;
    self.enabled = view.enabled;
    self.momentary = view.momentary;
}

- (void)onChange:(id)sender {
    MatchaSegmentViewPbEvent *event = [[MatchaSegmentViewPbEvent alloc] init];
    event.value = self.selectedSegmentIndex;
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:event.data];
    
    [self.viewNode.rootVC call:@"OnChange" viewId:self.node.identifier.longLongValue args:@[value]];
}

@end
