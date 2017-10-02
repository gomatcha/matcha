#import "MatchaSlider.h"
#import "MatchaViewController.h"

@implementation MatchaSlider

+ (void)load {
    [MatchaViewController registerView:@"gomatcha.io/matcha/view/slider" block:^(MatchaViewNode *node){
        return [[MatchaSlider alloc] initWithViewNode:node];
    }];
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        [self addTarget:self action:@selector(onChange:forEvent:) forControlEvents:UIControlEventValueChanged];
    }
    return self;
}

- (void)setNativeState:(NSData *)nativeState {
    MatchaViewPbSlider *view = [MatchaViewPbSlider parseFromData:nativeState error:nil];
    if (self.enabled != view.enabled) {
        self.enabled = view.enabled;
    }
    if (self.value != view.value) {
        self.value = view.value;
    }
    if (self.maximumValue != view.maxValue) {
        self.maximumValue = view.maxValue;
    }
    if (self.minimumValue != view.minValue) {
        self.minimumValue = view.minValue;
    }
}

- (void)setAlpha:(CGFloat)alpha {
    // UISlider.enabled sets the alpha don't allow MatchaViewNode to reset it back to 1.
    if (self.enabled == false && alpha > 0.99) {
        return;
    }
    [super setAlpha:alpha];
}

- (void)onChange:(id)sender forEvent:(UIEvent *)e {
    MatchaViewPbSliderEvent *event = [[MatchaViewPbSliderEvent alloc] init];
    event.value = self.value;
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:event.data];
    [self.viewNode call:@"OnValueChange", value, nil];
    if (e.allTouches.anyObject.phase == UITouchPhaseEnded) {
        [self.viewNode call:@"OnSubmit", value, nil];
    }
}

@end
