#import "MatchaTapGestureRecognizer.h"
#import <MatchaBridge/MatchaBridge.h>
#import "MatchaBuildNode.h"
#import "MatchaProtobuf.h"
#import "MatchaBridge.h"
#import "MatchaViewController_Private.h"

@interface MatchaTapGestureRecognizer ()
@property (nonatomic, assign) int64_t funcId;
@property (nonatomic, assign) int64_t viewId;
@property (nonatomic, weak) MatchaViewController *viewController;
@property (nonatomic, assign) bool disabled;
@end

@implementation MatchaTapGestureRecognizer

- (id)initWithMatchaVC:(MatchaViewController *)viewController viewId:(int64_t)viewId protobuf:(GPBAny *)pb {
    NSError *error = nil;
    MatchaPointerPBTapRecognizer *pbTapRecognizer = (id)[pb unpackMessageClass:[MatchaPointerPBTapRecognizer class] error:&error];
    if (pbTapRecognizer == nil) {
        return nil;
    }
    if ((self = [super initWithTarget:self action:@selector(action:)])) {
        self.numberOfTapsRequired = (int)pbTapRecognizer.count;
        self.viewController = viewController;
        self.funcId = pbTapRecognizer.onEvent;
        self.viewId = viewId;
    }
    return self;
}

- (void)disable {
    self.disabled = true;
}

- (void)updateWithProtobuf:(GPBAny *)pb {
    NSError *error = nil;
    MatchaPointerPBTapRecognizer *pbTapRecognizer = (id)[pb unpackMessageClass:[MatchaPointerPBTapRecognizer class] error:&error];
    if (pbTapRecognizer == nil) {
        return;
    }
    self.funcId = pbTapRecognizer.onEvent;
}

- (void)action:(id)sender {
    if (self.disabled) {
        return;
    }
    
    CGPoint point = [self locationInView:self.view];
    
    MatchaPointerPBTapEvent *event = [[MatchaPointerPBTapEvent alloc] init];
    event.position = [[MatchaLayoutPBPoint alloc] initWithCGPoint:point];
    event.timestamp = [[GPBTimestamp alloc] initWithDate:[NSDate date]];
    event.kind = MatchaPointerPBEventKind_EventKindRecognized;
    
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:event.data];
    [self.viewController call:[NSString stringWithFormat:@"gomatcha.io/matcha/touch %@",@(self.funcId)] viewId:self.viewId args2:@[value]];
}

@end
