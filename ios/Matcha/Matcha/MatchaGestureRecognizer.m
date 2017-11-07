#import "MatchaGestureRecognizer.h"
#import "MatchaProtobuf.h"
#import "MatchaView_Private.h"
#import <UIKit/UIGestureRecognizerSubclass.h>

UIGestureRecognizerState GoGestureStateToIOSState(long long a);
NSString *GestureStateToString(UIGestureRecognizerState s);

@implementation MatchaGestureRecognizer

- (instancetype)init {
    if (self = [super init]) {
        [self addTarget:self action:@selector(action)];
    }
    return self;
}

- (void)action {
    [self.viewNode call:@"gomatcha.io/matcha/pointer Action", nil];
}

- (void)reset {
    [self.viewNode call2:@"gomatcha.io/matcha/pointer Reset", nil];
}

- (void)touchesBegan:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseBegan;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = GoGestureStateToIOSState([self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong);
    
    [self stopTicks];
}

- (void)touchesMoved:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseMoved;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = GoGestureStateToIOSState([self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong);
    
    [self stopTicks];
}

- (void)touchesEnded:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseEnded;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = GoGestureStateToIOSState([self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong);
    
    if (self.state == UIGestureRecognizerStatePossible || self.state == UIGestureRecognizerStateChanged) {
        [self sendTicks];
    } else {
        [self stopTicks];
    }
}

- (void)touchesCancelled:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseCancelled;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = GoGestureStateToIOSState([self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong);
    
    if (self.state == UIGestureRecognizerStatePossible || self.state == UIGestureRecognizerStateChanged) {
        [self sendTicks];
    } else {
        [self stopTicks];
    }
}

// If touch is up but the gesture is still possible, send a event on every screen update with phaseNone so the gesture can cancel.
- (void)sendTicks {
    if (self.timer != NULL) {
        return;
    }
    
    dispatch_source_t timer = dispatch_source_create(DISPATCH_SOURCE_TYPE_TIMER, 0, 0, dispatch_get_main_queue());
    dispatch_source_set_timer(timer, dispatch_time(DISPATCH_TIME_NOW, 0.05*NSEC_PER_SEC), 0.05*NSEC_PER_SEC, 0.05*NSEC_PER_SEC / 5);
    dispatch_source_set_event_handler(timer, ^{
        [self tick];
    });
    dispatch_resume(timer);
    self.timer = timer;
}

- (void)stopTicks {
    if (self.timer == nil) {
        return;
    }
    dispatch_source_cancel(self.timer);
    self.timer = nil;
}

- (void)tick {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:CGPointZero];
    proto.phase = MatchaPointerPBPhase_PhaseNone;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = GoGestureStateToIOSState([self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong);
    
    if (self.state == UIGestureRecognizerStatePossible || self.state == UIGestureRecognizerStateChanged) {
        [self sendTicks];
    } else {
        [self stopTicks];
    }
}

@end

UIGestureRecognizerState GoGestureStateToIOSState(long long a) {
    switch (a) {
        case 0:
            return UIGestureRecognizerStatePossible;
        case 1:
            return UIGestureRecognizerStateChanged;
        case 2:
            return UIGestureRecognizerStateFailed;
        case 3:
            return UIGestureRecognizerStateEnded;
        default:
            return UIGestureRecognizerStateFailed;
    }
}

NSString *GestureStateToString(UIGestureRecognizerState s) {
    switch (s) {
        case UIGestureRecognizerStatePossible:
            return @"Possible";
        case UIGestureRecognizerStateBegan:
            return @"Began";
        case UIGestureRecognizerStateChanged:
            return @"Changed";
        case UIGestureRecognizerStateEnded:
            return @"Recognized/Ended";
        case UIGestureRecognizerStateCancelled:
            return @"Cancelled";
        case UIGestureRecognizerStateFailed:
            return @"Failed";
        default:
            return @"Unknown";
    }
}
