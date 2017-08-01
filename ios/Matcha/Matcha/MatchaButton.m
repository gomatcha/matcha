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
    MatchaButtonPBView *pbbutton = (id)[value.nativeViewState unpackMessageClass:[MatchaButtonPBView class] error:NULL];
    
    NSAttributedString *string = [[NSAttributedString alloc] initWithProtobuf:pbbutton.styledText];
    [self.button setAttributedTitle:string forState:UIControlStateNormal];
}

- (void)layoutSubviews {
    self.button.frame = self.bounds;
}

- (void)onPress {
    [self.viewNode.rootVC call:@"OnPress" viewId:self.node.identifier.longLongValue args:@[]];
}

@end
