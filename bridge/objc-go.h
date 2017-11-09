// +build matcha,darwin

#ifndef OBJC_GO_H
#define OBJC_GO_H

#import <Foundation/Foundation.h>
#import "go-foreign.h"
#import "go-go.h"

@interface MatchaGoValue : NSObject
- (id)initWithGoRef:(GoRef)ref; // not for external use.
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
@property (nonatomic, readonly) GoRef ref; // not for external use.
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

#endif // OBJC_GO_H
