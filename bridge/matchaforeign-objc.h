// +build matcha,darwin

#ifndef MOCHIFOREIGN_OBJC_H
#define MOCHIFOREIGN_OBJC_H

#import <Foundation/Foundation.h>
@class MatchaGoValue;

@interface MatchaObjcBridge : NSObject
+ (MatchaObjcBridge *)sharedBridge;
- (void)setObject:(id<NSObject>)obj forKey:(NSString *)string;
- (id<NSObject>)objectForKey:(NSString *)string;
@end

// Tracker
FgnRef MatchaForeignTrack(id value);
id MatchaForeignGet(FgnRef key);

// Utilities
NSString *MatchaCGoBufferToNSString(CGoBuffer buf); // Frees the buffer.
CGoBuffer MatchaNSStringToCGoBuffer(NSString *str); // Allocates a buffer.
NSData *MatchaCGoBufferToNSData(CGoBuffer buf); // Frees the buffer.
CGoBuffer MatchaNSDataToCGoBuffer(NSData *data); // Allocates a buffer.
NSArray<MatchaGoValue *> *MatchaCGoBufferToNSArray(CGoBuffer buf); // Frees the buffer.
NSArray<id> *MatchaCGoBufferToNSArray2(CGoBuffer buf); // Frees the buffer.
CGoBuffer MatchaNSArrayToCGoBuffer(NSArray<MatchaGoValue *> *array); // Allocates a buffer.
CGoBuffer MatchaNSArrayToCGoBuffer2(NSArray *array); // Allocates a buffer.

#endif //MOCHIFOREIGN_OBJC_H