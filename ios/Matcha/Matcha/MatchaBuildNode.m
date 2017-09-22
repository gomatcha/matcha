#import "MatchaBuildNode.h"
#import "MatchaProtobuf.h"

@interface MatchaBuildNode ()
@property (nonatomic, strong) GPBInt64Array *childIds;
@property (nonatomic, strong) NSMutableDictionary<NSString*, NSData *> *nativeValues;
@property (nonatomic, strong) NSString *nativeViewName;
@property (nonatomic, strong) GPBAny *nativeViewState;
@property (nonatomic, strong) NSNumber *identifier;
@property (nonatomic, strong) NSNumber *buildId;
@property (nonatomic, strong) NSDictionary<NSNumber *, GPBAny *> *touchRecognizers;
@end

@implementation MatchaBuildNode

- (id)initWithProtobuf:(MatchaViewPBBuildNode *)node {
    if ((self = [super init])) {
        self.identifier = @(node.id_p);
        self.buildId = @(node.buildId);
        self.nativeViewName = node.bridgeName;
        self.nativeViewState = node.bridgeValue;
        self.nativeValues = node.values;
        self.childIds = node.childrenArray;
        
        NSData *data = self.nativeValues[@"gomatcha.io/matcha/touch"];
        NSError *error = nil;
        MatchaPointerPBRecognizerList *recognizerList = [MatchaPointerPBRecognizerList parseFromData:data error:&error];
        if (error == nil) {
            NSMutableDictionary *touchRecognizers = [NSMutableDictionary dictionary];
            for (MatchaPointerPBRecognizer *i in recognizerList.recognizersArray) {
                touchRecognizers[@(i.id_p)] = i.recognizer;
            }
            self.touchRecognizers = touchRecognizers;
        }
    }
    return self;
}

@end
