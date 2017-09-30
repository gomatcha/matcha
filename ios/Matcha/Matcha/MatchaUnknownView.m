#import "MatchaUnknownView.h"

@implementation MatchaUnknownView

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        
        self.label = [[UILabel alloc] init];
        self.label.text = @"Unknown View";
        self.label.font = [UIFont boldSystemFontOfSize:13];
        self.label.textColor = [UIColor whiteColor];
        [self addSubview:self.label];
        
    }
    return self;
}

- (void)setNativeState:(NSData *)nativeState {
    // no-op
}

- (void)layoutSubviews {
    [super layoutSubviews];
    CGRect b = self.bounds;
    
    [self.label sizeToFit];
    CGRect f = self.label.frame;
    f.origin.y = 0;
    f.origin.x = b.origin.x;
    f.size.width = b.size.width;
    self.label.frame = f;
    
    // LayoutSubview happens after painting
    self.layer.backgroundColor = [UIColor redColor].CGColor;
}

@end

