#import "MatchaTextView.h"
#import "MatchaViewController.h"

@implementation MatchaTextView

+ (void)load {
    [MatchaViewController registerView:@"gomatcha.io/matcha/view/textview" block:^(MatchaViewNode *node){
        return [[MatchaTextView alloc] initWithViewNode:node];
    }];
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
    }
    return self;
}

- (void)setNativeState:(GPBAny *)nativeState {
    MatchaPBStyledText *text = (id)[nativeState unpackMessageClass:[MatchaPBStyledText class] error:nil];
    NSAttributedString *attrString = [[NSAttributedString alloc] initWithProtobuf:text];
    self.attributedText = attrString;
    self.numberOfLines = 0;
}

@end
