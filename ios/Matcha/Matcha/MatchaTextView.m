#import "MatchaTextView.h"

@implementation MatchaTextView

+ (void)load {
    MatchaRegisterView(@"gomatcha.io/matcha/view/textview", ^(MatchaViewNode *node){
        return [[MatchaTextView alloc] initWithViewNode:node];
    });
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
    }
    return self;
}

- (void)setNode:(MatchaBuildNode *)value {
    _node = value;
    GPBAny *state = value.nativeViewState;
    NSError *error = nil;
    MatchaPBStyledText *text = (id)[state unpackMessageClass:[MatchaPBStyledText class] error:&error];
    if (text != nil) {
        NSAttributedString *attrString = [[NSAttributedString alloc] initWithProtobuf:text];
        self.attributedText = attrString;
        self.numberOfLines = 0;
    }
}

@end
