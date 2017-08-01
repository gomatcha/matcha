#import <UIKit/UIKit.h>
#import <MatchaBridge/MatchaBridge.h>
@class MatchaPaintOptions;
@class MatchaLayoutGuide;
@class MatchaViewPBNode;
@class MatchaViewPBRoot;
@class MatchaLayoutPBGuide;
@class MatchaPaintPBStyle;
@class MatchaPBRecognizer;
@class MatchaViewPBLayoutPaintNode;
@class GPBInt64ObjectDictionary;
@class MatchaViewPBLayoutPaintNode;
@class MatchaViewPBBuildNode;
@class GPBInt64Array;
@class GPBAny;

@interface MatchaNodeRoot : NSObject // view.root
- (id)initWithProtobuf:(MatchaViewPBRoot *)data;
@property (nonatomic, readonly) GPBInt64ObjectDictionary *layoutPaintNodes;
@property (nonatomic, readonly) GPBInt64ObjectDictionary *buildNodes;
@property (nonatomic, readonly) NSMutableDictionary<NSString*, GPBAny*> *middleware;
@end

@interface MatchaBuildNode : NSObject
- (id)initWithProtobuf:(MatchaViewPBBuildNode *)node;
@property (nonatomic, readonly) GPBInt64Array *childIds;
@property (nonatomic, readonly) NSMutableDictionary<NSString*, GPBAny*> *nativeValues;
@property (nonatomic, readonly) NSString *nativeViewName;
@property (nonatomic, readonly) GPBAny *nativeViewState;
@property (nonatomic, readonly) NSNumber *identifier;
@property (nonatomic, readonly) NSNumber *buildId;
@property (nonatomic, readonly) NSDictionary<NSNumber *, GPBAny *> *touchRecognizers;
@end
