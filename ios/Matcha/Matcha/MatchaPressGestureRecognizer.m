#import "MatchaPressGestureRecognizer.h"
#import "MatchaProtobuf.h"
#import "MatchaNode.h"
#import "MatchaViewController.h"

@interface MatchaPressGestureRecognizer ()
@property (nonatomic, assign) int64_t funcId;
@property (nonatomic, assign) int64_t viewId;
@property (nonatomic, weak) MatchaViewController *viewController;
@property (nonatomic, strong) NSDate *startTime;
@property (nonatomic, assign) BOOL disabled;
@end

@implementation MatchaPressGestureRecognizer

- (id)initWithMatchaVC:(MatchaViewController *)viewController viewId:(int64_t)viewId protobuf:(GPBAny *)pb {
    NSError *error = nil;
    MatchaPBTouchPressRecognizer *pbTapRecognizer = (id)[pb unpackMessageClass:[MatchaPBTouchPressRecognizer class] error:&error];
    if ((self = [super initWithTarget:self action:@selector(action:)])) {
        self.minimumPressDuration = pbTapRecognizer.minDuration.timeInterval;
        self.viewController = viewController;
        self.funcId = pbTapRecognizer.funcId;
        self.viewId = viewId;
    }
    return self;
}

- (void)updateWithProtobuf:(GPBAny *)pb {
    NSError *error = nil;
    MatchaPBTouchPressRecognizer *pbTapRecognizer = (id)[pb unpackMessageClass:[MatchaPBTouchPressRecognizer class] error:&error];
    if (pbTapRecognizer == nil) {
        return;
    }
    self.funcId = pbTapRecognizer.funcId;
}

- (void)disable {
    self.disabled = false;
}

- (void)action:(id)sender {
    if (self.disabled) {
        return;
    }
    
    CGPoint point = [self locationInView:self.view];
    
    MatchaPBTouchPressEvent *event = [[MatchaPBTouchPressEvent alloc] init];
    event.position = [[MatchaLayoutPBPoint alloc] initWithCGPoint:point];
    event.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    if (self.state == UIGestureRecognizerStateBegan) {
        event.kind = MatchaPBTouchEventKind_EventKindChanged;
        self.startTime = [NSDate date];
    } else if (self.state == UIGestureRecognizerStateChanged) {
        event.kind = MatchaPBTouchEventKind_EventKindChanged;
    } else if (self.state == UIGestureRecognizerStateEnded) {
        event.kind = MatchaPBTouchEventKind_EventKindRecognized;
    } else if (self.state == UIGestureRecognizerStateCancelled) {
        event.kind = MatchaPBTouchEventKind_EventKindFailed;
    } else {
        return;
    }
    event.duration = [[GPBDuration alloc] initWithTimeInterval:-self.startTime.timeIntervalSinceNow];
    
    NSData *data = [event data];
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:data];
    
    [self.viewController call:[NSString stringWithFormat:@"%@",@(self.funcId)] viewId:self.viewId args:@[value]];
}

@end
