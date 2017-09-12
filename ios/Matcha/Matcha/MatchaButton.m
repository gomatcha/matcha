#import "MatchaButton.h"
#import "MatchaViewController.h"
#import "MatchaProtobuf.h"

@implementation MatchaButton

+ (void)load {
    MatchaRegisterView(@"gomatcha.io/matcha/view/button", ^(MatchaViewNode *node){
        return [[MatchaButton alloc] initWithViewNode:node];
    });
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        self.button = [UIButton buttonWithType:UIButtonTypeSystem];
        [self.button addTarget:self action:@selector(onPress) forControlEvents:UIControlEventTouchUpInside];
        [self addSubview:self.button];
    }
    return self;
}

- (void)setNode:(MatchaBuildNode *)value {
    _node = value;
    MatchaViewPbButton *pbbutton = (id)[value.nativeViewState unpackMessageClass:[MatchaViewPbButton class] error:NULL];
    
    if (pbbutton.hasColor) {
        self.button.tintColor = [[UIColor alloc] initWithProtobuf:pbbutton.color];
    } else {
        self.button.tintColor = nil;
    }
    self.button.titleLabel.font = [UIFont systemFontOfSize:20];
    [self.button setTitle:pbbutton.str forState:UIControlStateNormal];
    self.button.enabled = pbbutton.enabled;
}

- (void)layoutSubviews {
    self.button.frame = self.bounds;
}

- (void)onPress {
    [self.viewNode.rootVC call:@"OnPress" viewId:self.node.identifier.longLongValue args:@[]];
}

@end
