#import "MatchaBridge.h"

@implementation MatchaGoValue (Matcha)

- (id)initWithCGPoint:(CGPoint)point {
    if ((self = [self initWithType:@"layout.Point"].elem)) {
        self[@"X"] = [[MatchaGoValue alloc] initWithDouble:point.x];
        self[@"Y"] = [[MatchaGoValue alloc] initWithDouble:point.y];
    }
    return self;
}

- (id)initWithCGSize:(CGSize)size {
     if ((self = [self initWithType:@"layout.Point"].elem)) {
         self[@"X"] = [[MatchaGoValue alloc] initWithDouble:size.width];
         self[@"Y"] = [[MatchaGoValue alloc] initWithDouble:size.height];
     }
     return self;
}

- (id)initWithCGRect:(CGRect)rect {
    if ((self = [self initWithType:@"layout.Rect"].elem)) {
        self[@"Min"] = [[MatchaGoValue alloc] initWithCGPoint:rect.origin];
        self[@"Max"] = [[MatchaGoValue alloc] initWithCGPoint:CGPointMake(rect.origin.x + rect.size.width, rect.origin.y + rect.size.height)];
    }
    return self;
}

- (CGPoint)toCGPoint {
    CGPoint point;
    point.x = self[@"X"].toDouble;
    point.y = self[@"Y"].toDouble;
    return point;
}

- (CGSize)toCGSize {
    CGSize size;
    size.width = self[@"X"].toDouble;
    size.height = self[@"Y"].toDouble;
    return size;
}

- (CGRect)toCGRect {
    MatchaGoValue *min = self[@"Min"];
    MatchaGoValue *max = self[@"Max"];
    CGRect rect;
    rect.origin.x = min[@"X"].toDouble;
    rect.origin.y = min[@"Y"].toDouble;
    rect.size.width = max[@"X"].toDouble - rect.origin.x;
    rect.size.height = max[@"Y"].toDouble - rect.origin.y;
    return rect;
}

- (UIEdgeInsets)toUIEdgeInsets {
    UIEdgeInsets insets;
    insets.top = self[@"Top"].toDouble;
    insets.bottom = self[@"Bottom"].toDouble;
    insets.right = self[@"Right"].toDouble;
    insets.left = self[@"Left"].toDouble;
    return insets;
}
@end

@implementation NSMapTable (Matcha) 

- (id)objectForKeyedSubscript:(id)key {
    return [self objectForKey:key];
}

- (void)setObject:(id)obj forKeyedSubscript:(id)key {
    if (obj != nil) {
        [self setObject:obj forKey:key];
    } else {
        [self removeObjectForKey:key];
    }
}

@end
