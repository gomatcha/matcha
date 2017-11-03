#import "MatchaGestureRecognizer.h"
#import "MatchaProtobuf.h"
#import "MatchaView_Private.h"
#import <UIKit/UIGestureRecognizerSubclass.h>

@implementation MatchaGestureRecognizer

- (void)touchesBegan:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseBegan;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = [self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong;
}

- (void)touchesMoved:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseMoved;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = [self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong;
}

- (void)touchesEnded:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseEnded;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = [self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong;
}

- (void)touchesCancelled:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    MatchaPointerPBEvent *proto = [[MatchaPointerPBEvent alloc] init];
    proto.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    proto.location = [[MatchaLayoutPBPoint alloc] initWithCGPoint:[self locationInView:self.view]];
    proto.phase = MatchaPointerPBPhase_PhaseCancelled;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:proto.data];
    self.state = [self.viewNode call2:@"gomatcha.io/matcha/pointer OnEvent", value, nil][1].toLongLong;
}

@end
