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
//    [self.viewNode call:@"gomatcha.io/matcha/pointer Action", nil];
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
}

- (void)touchesMoved:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseMoved;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = GoGestureStateToIOSState([self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong);
}

- (void)touchesEnded:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseEnded;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = GoGestureStateToIOSState([self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong);
}

- (void)touchesCancelled:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseCancelled;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = GoGestureStateToIOSState([self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong);
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
