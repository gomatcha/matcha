#import "MatchaScrollView.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController.h"

@implementation MatchaScrollView

+ (void)load {
    MatchaRegisterView(@"gomatcha.io/matcha/view/scrollview", ^(MatchaViewNode *node){
        return [[MatchaScrollView alloc] initWithViewNode:node];
    });
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        self.delegate = self;
        self.alwaysBounceVertical = true; // TODO(KD):
    }
    return self;
}

- (void)setNode:(MatchaBuildNode *)value {
    _node = value;
    
    if (self.subviews.count > 0) {
        self.contentSize = ((UIView *)self.subviews[0]).frame.size;
    }
    
    GPBAny *state = value.nativeViewState;
    NSError *error = nil;
    MatchaScrollViewPBView *pbscrollview = (id)[state unpackMessageClass:[MatchaScrollViewPBView class] error:&error];
    if (pbscrollview != nil) {
        self.scrollEnabled = pbscrollview.scrollEnabled;
        self.showsVerticalScrollIndicator = pbscrollview.showsVerticalScrollIndicator;
        self.showsHorizontalScrollIndicator = pbscrollview.showsHorizontalScrollIndicator;
    }
}

- (void)scrollViewDidScroll:(UIScrollView *)scrollView {
    if (self.viewNode.rootVC.updating || CGPointEqualToPoint(self.contentOffset, self.matchaContentOffset)) {
        return;
    }
    
    MatchaScrollViewPBScrollEvent *event = [[MatchaScrollViewPBScrollEvent alloc] init];
    event.contentOffset = [[MatchaLayoutPBPoint alloc] initWithCGPoint:scrollView.contentOffset];
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:event.data];
    [self.viewNode.rootVC call:@"OnScroll" viewId:self.node.identifier.longLongValue args:@[value]];
}

@end
