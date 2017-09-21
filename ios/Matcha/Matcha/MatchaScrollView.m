#import "MatchaScrollView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController_Private.h"

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

- (void)setNode:(MatchaBuildNode *)value {
    _node = value;
    
    GPBAny *state = value.nativeViewState;
    NSError *error = nil;
    MatchaViewPBScrollView *pbscrollview = (id)[state unpackMessageClass:[MatchaViewPBScrollView class] error:&error];
    if (pbscrollview != nil) {
        self.scrollEnabled = pbscrollview.scrollEnabled;
        self.showsVerticalScrollIndicator = pbscrollview.showsVerticalScrollIndicator;
        self.showsHorizontalScrollIndicator = pbscrollview.showsHorizontalScrollIndicator;
        self.alwaysBounceVertical = pbscrollview.vertical;
        self.alwaysBounceHorizontal = pbscrollview.horizontal;
    }
}

- (void)scrollViewDidScroll:(UIScrollView *)scrollView {
    if (self.viewNode.rootVC.updating || CGPointEqualToPoint(self.contentOffset, self.matchaContentOffset)) {
        return;
    }
    
    MatchaViewPBScrollEvent *event = [[MatchaViewPBScrollEvent alloc] init];
    event.contentOffset = [[MatchaLayoutPBPoint alloc] initWithCGPoint:scrollView.contentOffset];
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:event.data];
    [self.viewNode.rootVC call:@"OnScroll" viewId:self.node.identifier.longLongValue args:@[value]];
}

@end
