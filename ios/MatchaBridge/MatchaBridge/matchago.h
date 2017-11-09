#import <Foundation/Foundation.h>
@class MatchaGoValue;

void matchaTestFunc(void);

@interface MatchaGoValue : NSObject
- (id)initWithObject:(id)v;
- (id)initWithBool:(BOOL)v;
- (id)initWithInt:(int)v;
- (id)initWithLongLong:(long long)v;
- (id)initWithUnsignedLongLong:(unsigned long long)v;
- (id)initWithDouble:(double)v;
- (id)initWithString:(NSString *)v;
- (id)initWithData:(NSData *)v;
- (id)initWithArray:(NSArray<MatchaGoValue *> *)v;
- (id)initWithFunc:(NSString *)funcName;
- (id)toObject;
- (BOOL)toBool;
- (long long)toLongLong;
- (unsigned long long)toUnsignedLongLong;
- (double)toDouble;
- (NSString *)toString;
- (NSData *)toData;
- (NSArray *)toArray;
- (BOOL)isNil;
- (NSArray<MatchaGoValue *> *)call:(NSString *)method, ... NS_REQUIRES_NIL_TERMINATION; // pass in nil for the method to call a closure. varargs should be of MatchaGoValue *.
- (NSArray<MatchaGoValue *> *)call:(NSString *)method args:(va_list)args;
@end
