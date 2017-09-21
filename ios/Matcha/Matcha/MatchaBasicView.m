#import "MatchaBasicView.h"

@implementation MatchaBasicView

+ (void)load {
    [MatchaViewController registerView:@"" block:^(MatchaViewNode *node){
        return [[MatchaBasicView alloc] initWithViewNode:node];
    }];
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
    }
    return self;
}

- (void)setNativeState:(GPBAny *)nativeState {
    // no-op
}

@end

