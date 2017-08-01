#import "MatchaButtonGestureRecognizer.h"
#import <UIKit/UIGestureRecognizerSubclass.h>
#import "MatchaViewController.h"
#import "MatchaProtobuf.h"

@interface MatchaButtonGestureRecognizer () <UIGestureRecognizerDelegate>
@property (nonatomic, assign) int64_t funcId;
@property (nonatomic, assign) int64_t viewId;
@property (nonatomic, weak) MatchaViewController *viewController;
@property (nonatomic, assign) BOOL disabled;
@property (nonatomic, assign) BOOL inside;
@property (nonatomic, assign) BOOL ignoresScroll;
@property (nonatomic, assign) UIGestureRecognizerState lastState;
@end


@implementation MatchaButtonGestureRecognizer

- (id)initWithMatchaVC:(MatchaViewController *)viewController viewId:(int64_t)viewId protobuf:(GPBAny *)pbany {
    NSError *error = nil;
    MatchaPBTouchButtonRecognizer *pb = (id)[pbany unpackMessageClass:[MatchaPBTouchButtonRecognizer class] error:&error];
    if (pb == nil) {
        return nil;
    }
    if ((self = [super initWithTarget:self action:@selector(action:)])) {
        self.delegate = self;
        self.viewController = viewController;
        self.funcId = pb.onEvent;
        self.viewId = viewId;
    }
    return self;
}

- (void)updateWithProtobuf:(GPBAny *)pbany {
    NSError *error = nil;
    MatchaPBTouchButtonRecognizer *pb = (id)[pbany unpackMessageClass:[MatchaPBTouchButtonRecognizer class] error:&error];
    if (pb == nil) {
        return;
    }
    self.funcId = pb.onEvent;
}

- (void)disable {
    self.disabled = false;
}

- (void)action:(id)sender {
//    CGPoint point = [self locationInView:self.view];
//    [self action:sender inside:CGRectContainsPoint(self.view.bounds, point)];
}

- (void)action:(id)sender inside:(bool)inside {
    if (self.disabled) {
        return;
    }
    
    BOOL prevInside = self.inside;
    self.inside = inside;
    self.lastState = self.state;
    
    MatchaPBTouchButtonEvent *event = [[MatchaPBTouchButtonEvent alloc] init];
    event.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    event.inside = self.inside;
    if (self.state == UIGestureRecognizerStatePossible && self.inside != prevInside) { // Only update if inside has changed
        event.kind = MatchaPBTouchEventKind_EventKindPossible;
    } else if (self.state == UIGestureRecognizerStateRecognized) {
        event.kind = MatchaPBTouchEventKind_EventKindRecognized;
    } else if (self.state == UIGestureRecognizerStateFailed) {
        event.kind = MatchaPBTouchEventKind_EventKindFailed;
    } else {
        return;
    }
    
    NSData *data = [event data];
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:data];
    
    [self.viewController call:[NSString stringWithFormat:@"%@", @(self.funcId)] viewId:self.viewId args:@[value]];
}

- (BOOL)gestureRecognizer:(UIGestureRecognizer *)gestureRecognizer shouldReceiveTouch:(UITouch *)touch {
    UIView *view = touch.view;
    while (true) {
        if (view == self.view) {
            break;
        }
        for (UIGestureRecognizer *i in view.gestureRecognizers) {
            if ([i isKindOfClass:[MatchaButtonGestureRecognizer class]]) {
                return false;
            }
        }
    }
    return YES;
}

#pragma mark - STuff

- (void)touchesBegan:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    [super touchesBegan:touches withEvent:event];
    if (self.state != UIGestureRecognizerStatePossible) {
        return;
    }
    self.inside = false; // reset the inside prop
    if (touches.count != 1) {
        self.state = UIGestureRecognizerStateFailed;
        [self action:self inside:[self touchesInside:touches]];
    } else {
        [self action:self inside:[self touchesInside:touches]];
    }
}

- (void)touchesMoved:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    [super touchesMoved:touches withEvent:event];
    if (self.state != UIGestureRecognizerStatePossible) {
        return;
    }
    if (touches.count != 1) {
        self.state = UIGestureRecognizerStateFailed;
        [self action:self inside:[self touchesInside:touches]];
    } else {
        [self action:self inside:[self touchesInside:touches]];
    }
}

- (void)touchesCancelled:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    [super touchesCancelled:touches withEvent:event];
    if (self.state != UIGestureRecognizerStatePossible) {
        return;
    }
    
    self.state = UIGestureRecognizerStateFailed;
    [self action:self inside:false];
}

- (void)touchesEnded:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    [super touchesEnded:touches withEvent:event];
    if (self.state != UIGestureRecognizerStatePossible) {
        return;
    }
    
    if (touches.count != 1) {
        self.state = UIGestureRecognizerStateFailed;
        [self action:self inside:[self touchesInside:touches]];
        return;
    }
    if ([self touchesInside:touches]) {
        self.state = UIGestureRecognizerStateRecognized;
        [self action:self inside:true];
    } else {
        self.state = UIGestureRecognizerStateFailed;
        [self action:self inside:false];
    }
}

- (BOOL)touchesInside:(NSSet *)touches {
    if (touches.count != 1) {
        return false;
    }
    CGPoint point = [touches.anyObject locationInView:self.view];
    return CGRectContainsPoint(self.view.bounds, point);
}

- (void)reset {
    // Make sure we acknowledge the failed state.
    if (self.state == UIGestureRecognizerStateFailed && self.state != self.lastState) {
        [self action:self inside:false];
    }
    [super reset];
}

@end
