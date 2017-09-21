#import "MatchaImageView.h"
#import "MatchaViewController.h"

@implementation MatchaImageView

+ (void)load {
    [MatchaViewController registerView:@"gomatcha.io/matcha/view/imageview" block:^(MatchaViewNode *node){
        return [[MatchaImageView alloc] initWithViewNode:node];
    }];
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
    }
    return self;
}

- (void)setNativeState:(GPBAny *)nativeState {
    MatchaViewPBImageView *view = (id)[nativeState unpackMessageClass:[MatchaViewPBImageView class] error:nil];
    
    UIImage *image = [[UIImage alloc] initWithImageOrResourceProtobuf:view.image];
    switch (view.resizeMode) {
        case MatchaViewPBImageResizeMode_GPBUnrecognizedEnumeratorValue:
        case MatchaViewPBImageResizeMode_Fit:
            self.contentMode = UIViewContentModeScaleAspectFit;
            break;
        case MatchaViewPBImageResizeMode_Fill:
            self.contentMode = UIViewContentModeScaleAspectFill;
            break;
        case MatchaViewPBImageResizeMode_Stretch:
            self.contentMode = UIViewContentModeScaleToFill;
            break;
        case MatchaViewPBImageResizeMode_Center:
            self.contentMode = UIViewContentModeCenter;
            break;
    }
    if (view.hasTint) {
        self.tintColor = [[UIColor alloc] initWithProtobuf:view.tint];
        image = [image imageWithRenderingMode:UIImageRenderingModeAlwaysTemplate];
    }
    
    if (![self.image isEqual:image]) {
        self.image = image;
    }
}

@end
