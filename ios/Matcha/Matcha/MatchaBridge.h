#import <UIKit/UIKit.h>
#import <MatchaBridge/MatchaBridge.h>

@interface NSMapTable (Matcha)
- (id)objectForKeyedSubscript:(id)key;
- (void)setObject:(id)obj forKeyedSubscript:(id)key;
@end
