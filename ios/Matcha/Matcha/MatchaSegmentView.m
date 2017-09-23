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

- (void)setNativeState:(NSData *)nativeState {
    MatchaiOSPBSegmentView *view = [MatchaiOSPBSegmentView parseFromData:nativeState error:nil];
    
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
    [self.viewNode call:@"OnChange", [[MatchaGoValue alloc] initWithData:event.data], nil];
}

@end
