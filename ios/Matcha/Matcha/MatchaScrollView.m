#import "MatchaScrollView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController_Private.h"
#import "MatchaView_Private.h"

@implementation MatchaScrollView

+ (void)load {
    [MatchaViewController registerView:@"gomatcha.io/matcha/view/scrollview" block:^(MatchaViewNode *node){
        return [[MatchaScrollView alloc] initWithViewNode:node];
    }];
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        self.delegate = self;
        if (@available(iOS 11.0, *)) {
            self.contentInsetAdjustmentBehavior = UIScrollViewContentInsetAdjustmentNever;
        }
    }
    return self;
}

- (void)setNativeState:(NSData *)nativeState {
    MatchaViewPBScrollView *state = [MatchaViewPBScrollView parseFromData:nativeState error:nil];
    if (self.scrollEnabled != state.scrollEnabled) {
        self.scrollEnabled = state.scrollEnabled;
    }
    if (self.showsVerticalScrollIndicator != state.showsVerticalScrollIndicator) {
        self.showsVerticalScrollIndicator = state.showsVerticalScrollIndicator;
    }
    if (self.showsHorizontalScrollIndicator != state.showsHorizontalScrollIndicator) {
        self.showsHorizontalScrollIndicator = state.showsHorizontalScrollIndicator;
    }
    if (self.alwaysBounceVertical != state.vertical) {
        self.alwaysBounceVertical = state.vertical;
    }
    if (self.alwaysBounceHorizontal != state.horizontal) {
        self.alwaysBounceHorizontal = state.horizontal;
    }
}

- (void)scrollViewDidScroll:(UIScrollView *)scrollView {
    // MatchaViewNode changes the scrollOffset. Don't trigger an event back to Go.
    // contentOffset rounds to the nearest 1/screenscale
    if ((fabs(self.contentOffset.x - self.matchaContentOffset.x) < 0.5 && fabs(self.contentOffset.y - self.matchaContentOffset.y) < 0.5)) {
        return;
    }
    self.matchaContentOffset = self.contentOffset;
    
    MatchaViewPBScrollEvent *event = [[MatchaViewPBScrollEvent alloc] init];
    event.contentOffset = [[MatchaLayoutPBPoint alloc] initWithCGPoint:scrollView.contentOffset];
    [self.viewNode call:@"OnScroll", [[MatchaGoValue alloc] initWithData:event.data], nil];
}

@end
