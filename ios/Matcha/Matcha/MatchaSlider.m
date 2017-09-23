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
    
    self.enabled = view.enabled;
    self.value = view.value;
    self.maximumValue = view.maxValue;
    self.minimumValue = view.minValue;
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
