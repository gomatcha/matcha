#import <UIKit/UIKit.h>
#import <MatchaBridge/MatchaBridge.h>

@interface MatchaGoValue (Matcha)
- (id)initWithCGPoint:(CGPoint)point;
- (id)initWithCGSize:(CGSize)size;
- (id)initWithCGRect:(CGRect)rect;
- (CGPoint)toCGPoint;
- (CGSize)toCGSize;
- (CGRect)toCGRect;
- (UIEdgeInsets)toUIEdgeInsets;
@end

@interface NSMapTable (Matcha)
- (id)objectForKeyedSubscript:(id)key;
- (void)setObject:(id)obj forKeyedSubscript:(id)key;
@end
