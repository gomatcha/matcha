#import "MatchaUnknownView.h"

@implementation MatchaUnknownView

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        self.backgroundColor = [UIColor redColor];
    }
    return self;
}

- (void)setNode:(MatchaBuildNode *)value {
    _node = value;
}

@end

