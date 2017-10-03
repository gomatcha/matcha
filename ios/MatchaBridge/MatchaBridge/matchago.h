#import <Foundation/Foundation.h>
@class MatchaGoValue;

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
- (id)initWithType:(NSString *)typeName;
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
- (BOOL)isEqual:(MatchaGoValue *)value;
- (MatchaGoValue *)elem;
- (NSArray<MatchaGoValue *> *)call:(NSString *)method, ... NS_REQUIRES_NIL_TERMINATION; // pass in nil for the method to call a closure. varargs should be of MatchaGoValue *.
- (NSArray<MatchaGoValue *> *)call:(NSString *)method args:(va_list)args;
- (MatchaGoValue *)field:(NSString *)name;
- (void)setField:(NSString *)name value:(MatchaGoValue *)value;
- (MatchaGoValue *)objectForKeyedSubscript:(NSString *)key;
- (void)setObject:(MatchaGoValue *)object forKeyedSubscript:(NSString *)key;
@end
