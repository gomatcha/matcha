#import "MatchaMiddleware.h"

//NSMutableArray *_MatchaMiddlewares() {
//    static NSMutableArray *sMiddlewares = nil;
//    static dispatch_once_t sOnce = 0;
//    dispatch_once(&sOnce, ^{
//        sMiddlewares = [NSMutableArray array];
//    });
//    return sMiddlewares;
//}
//
//void MatchaRegisterMiddleware(MatchaMiddlewareBlock block) {
//    NSMutableArray *arr = _MatchaMiddlewares();
//    [arr addObject:block];
//}
//
//NSArray<MatchaMiddlewareBlock> *MatchaMiddlewares() {
//    return _MatchaMiddlewares();
//}
