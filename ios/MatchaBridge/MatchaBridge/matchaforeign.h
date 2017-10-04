#ifndef MOCHIOBJC_H
#define MOCHIOBJC_H

@interface MatchaObjcBridge : NSObject
+ (MatchaObjcBridge *)sharedBridge;
- (void)setObject:(id<NSObject>)obj forKey:(NSString *)string;
- (id<NSObject>)objectForKey:(NSString *)string;
@end

#endif //MOCHIOBJC_H
