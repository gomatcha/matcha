#import "MatchaBasicView.h"

@implementation MatchaBasicView

+ (void)load {
    MatchaRegisterView(@"", ^(MatchaViewNode *node){
        return [[MatchaBasicView alloc] initWithViewNode:node];
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
}

@end

