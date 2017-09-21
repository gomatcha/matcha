#import "MatchaSegmentView.h"

@implementation MatchaSegmentView

+ (void)load {
    [MatchaViewController registerView:@"gomatcha.io/matcha/view/segmentview" block:^(MatchaViewNode *node){
        return [[MatchaSegmentView alloc] initWithViewNode:node];
    }];
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
    MatchaiOSPBSegmentView *view = (id)[value.nativeViewState unpackMessageClass:[MatchaiOSPBSegmentView class] error:nil];
    
    [self removeAllSegments];
    for (NSInteger i = 0; i < view.titlesArray.count; i++) {
        [self insertSegmentWithTitle:view.titlesArray[i] atIndex:i animated:NO];
    }
    self.selectedSegmentIndex = (int)view.value;
    self.enabled = view.enabled;
    self.momentary = view.momentary;
}

- (void)onChange:(id)sender {
    MatchaiOSPBSegmentViewEvent *event = [[MatchaiOSPBSegmentViewEvent alloc] init];
    event.value = self.selectedSegmentIndex;
    [self.viewNode call:@"OnChange" args:[[MatchaGoValue alloc] initWithData:event.data], nil];
}

@end
