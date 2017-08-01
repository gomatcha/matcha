#import "MatchaProgressView.h"
#import "MatchaProtobuf.h"

@implementation MatchaProgressView

+ (void)load {
    MatchaRegisterView(@"gomatcha.io/matcha/view/progressview", ^(MatchaViewNode *node){
        return [[MatchaProgressView alloc] initWithViewNode:node];
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
    MatchaProgressViewPBView *view = (id)[state unpackMessageClass:[MatchaProgressViewPBView class] error:&error];
    if (view != nil) {
        self.progress = view.progress;
    }
}

@end
