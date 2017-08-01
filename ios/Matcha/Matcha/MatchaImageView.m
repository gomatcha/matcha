#import "MatchaImageView.h"

@implementation MatchaImageView

+ (void)load {
    MatchaRegisterView(@"gomatcha.io/matcha/view/imageview", ^(MatchaViewNode *node){
        return [[MatchaImageView alloc] initWithViewNode:node];
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
    MatchaImageViewPBView *pbimageview = (id)[state unpackMessageClass:[MatchaImageViewPBView class] error:&error];
    
    UIImage *image = [[UIImage alloc] initWithImageOrResourceProtobuf:pbimageview.image];
    
    switch (pbimageview.resizeMode) {
        case MatchaImageViewPBResizeMode_GPBUnrecognizedEnumeratorValue:
        case MatchaImageViewPBResizeMode_Fit:
            self.contentMode = UIViewContentModeScaleAspectFit;
            break;
        case MatchaImageViewPBResizeMode_Fill:
            self.contentMode = UIViewContentModeScaleAspectFill;
            break;
        case MatchaImageViewPBResizeMode_Stretch:
            self.contentMode = UIViewContentModeScaleToFill;
            break;
        case MatchaImageViewPBResizeMode_Center:
            self.contentMode = UIViewContentModeCenter;
            break;
    }
    if (pbimageview.hasTint) {
        self.tintColor = [[UIColor alloc] initWithProtobuf:pbimageview.tint];
        image = [image imageWithRenderingMode:UIImageRenderingModeAlwaysTemplate];
    }
    
    if (![self.image isEqual:image]) {
        self.image = image;
    }
}

@end
