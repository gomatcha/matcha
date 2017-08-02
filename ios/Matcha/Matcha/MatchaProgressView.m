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
    MatchaProgressViewPBView *view = (id)[value.nativeViewState unpackMessageClass:[MatchaProgressViewPBView class] error:nil];
    self.progress = view.progress;
    self.tintColor = [[UIColor alloc] initWithProtobuf:view.progressColor];
}

@end
