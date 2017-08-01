#import "MatchaNode.h"
#import "MatchaProtobuf.h"

@interface MatchaNodeRoot ()
@property (nonatomic, strong) GPBInt64ObjectDictionary<MatchaViewPBLayoutPaintNode*> *layoutPaintNodes;
@property (nonatomic, strong) GPBInt64ObjectDictionary<MatchaViewPBBuildNode*> *buildNodes;
@property (nonatomic, strong) NSMutableDictionary<NSString*, GPBAny*> *middleware;
@end

@implementation MatchaNodeRoot
- (id)initWithProtobuf:(MatchaViewPBRoot *)pbroot {
    if ((self = [super init])) {
        self.layoutPaintNodes = pbroot.layoutPaintNodes;
        self.buildNodes = pbroot.buildNodes;
        self.middleware = pbroot.middleware;
    }
    return self;
}
@end

@interface MatchaBuildNode ()
@property (nonatomic, strong) GPBInt64Array *childIds;
@property (nonatomic, strong) NSMutableDictionary<NSString*, GPBAny*> *nativeValues;
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
        
        GPBAny *any = self.nativeValues[@"gomatcha.io/matcha/touch"];
        NSError *error = nil;
        MatchaPBTouchRecognizerList *recognizerList = (id)[any unpackMessageClass:[MatchaPBTouchRecognizerList class] error:&error];
        if (error == nil) {
            NSMutableDictionary *touchRecognizers = [NSMutableDictionary dictionary];
            for (MatchaPBTouchRecognizer *i in recognizerList.recognizersArray) {
                touchRecognizers[@(i.id_p)] = i.recognizer;
            }
            self.touchRecognizers = touchRecognizers;
        }
    }
    return self;
}

@end
