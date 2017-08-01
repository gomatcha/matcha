#import "MatchaDeadlockLogger.h"
#import "MatchaBridge.h"

@interface MatchaDeadlockLogger ()
@property (nonatomic, strong) dispatch_queue_t queue;
@property (nonatomic, strong) dispatch_source_t timer;
@property (nonatomic, strong) MatchaGoValue *printStackFunc;
@end

@implementation MatchaDeadlockLogger

+ (instancetype)sharedLogger {
    static MatchaDeadlockLogger *sLogger = nil;
    static dispatch_once_t sOnce = 0;
    dispatch_once(&sOnce, ^{
        sLogger = [[MatchaDeadlockLogger alloc] init];
    });
    return sLogger;
}

- (id)init {
    if ((self = [super init])) {
        self.printStackFunc = [[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/internal printStack"];
        __weak typeof(self) weakSelf = self;
        self.queue = dispatch_queue_create(NULL, DISPATCH_QUEUE_SERIAL);
        self.timer = dispatch_source_create(DISPATCH_SOURCE_TYPE_TIMER, 0, 0, self.queue);
        dispatch_source_set_timer(self.timer, dispatch_time(DISPATCH_TIME_NOW, 5*NSEC_PER_SEC), 5*NSEC_PER_SEC, 5*NSEC_PER_SEC / 5);
        dispatch_source_set_event_handler(self.timer, ^{
            [weakSelf timerDidTick];
        });
        dispatch_resume(self.timer);
    }
    return self;
}

- (void)timerDidTick {
    __block bool flag = false;
    dispatch_async(dispatch_get_main_queue(), ^{
        flag = true;
    });
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 0.5 * NSEC_PER_SEC), self.queue, ^{
        if (!flag) {
            NSLog(@"Deadlock detected!");
            [self.printStackFunc call:nil args:nil];
        }
    });
}

@end
