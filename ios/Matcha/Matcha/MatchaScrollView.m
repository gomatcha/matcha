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
    }
    return self;
}

- (void)setNativeState:(GPBAny *)nativeState {
    MatchaViewPBScrollView *pbscrollview = (id)[nativeState unpackMessageClass:[MatchaViewPBScrollView class] error:nil];
    self.scrollEnabled = pbscrollview.scrollEnabled;
    self.showsVerticalScrollIndicator = pbscrollview.showsVerticalScrollIndicator;
    self.showsHorizontalScrollIndicator = pbscrollview.showsHorizontalScrollIndicator;
    self.alwaysBounceVertical = pbscrollview.vertical;
    self.alwaysBounceHorizontal = pbscrollview.horizontal;
}

- (void)scrollViewDidScroll:(UIScrollView *)scrollView {
    if (self.viewNode.rootVC.updating || CGPointEqualToPoint(self.contentOffset, self.matchaContentOffset)) {
        return;
    }
    
    MatchaViewPBScrollEvent *event = [[MatchaViewPBScrollEvent alloc] init];
    event.contentOffset = [[MatchaLayoutPBPoint alloc] initWithCGPoint:scrollView.contentOffset];
    [self.viewNode call:@"OnScroll" args:[[MatchaGoValue alloc] initWithData:event.data], nil];
}

@end
