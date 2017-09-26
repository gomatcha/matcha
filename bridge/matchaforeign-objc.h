// +build matcha,darwin

#ifndef MOCHIFOREIGN_OBJC_H
#define MOCHIFOREIGN_OBJC_H

#import <Foundation/Foundation.h>

@interface MatchaObjcBridge : NSObject
+ (MatchaObjcBridge *)sharedBridge;
- (void)setObject:(id<NSObject>)obj forKey:(NSString *)string;
- (id<NSObject>)objectForKey:(NSString *)string;
@end

// Tracker
ObjcRef MatchaTrackObjc(id value);
id MatchaGetObjc(ObjcRef key);

// Utilities
NSString *MatchaCGoBufferToNSString(CGoBuffer buf); // Frees the buffer.
CGoBuffer MatchaNSStringToCGoBuffer(NSString *str); // Allocates a buffer.
NSData *MatchaCGoBufferToNSData(CGoBuffer buf); // Frees the buffer.
CGoBuffer MatchaNSDataToCGoBuffer(NSData *data); // Allocates a buffer.

#endif //MOCHIFOREIGN_OBJC_H