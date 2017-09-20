#import "MatchaUnknownView.h"

@implementation MatchaUnknownView

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        super.backgroundColor = [UIColor redColor];
        
        self.label = [[UILabel alloc] init];
        self.label.text = @"Unknown";
        self.label.textColor = [UIColor whiteColor];
        self.label.textAlignment = NSTextAlignmentCenter;
        [self addSubview:self.label];
        
    }
    return self;
}

- (void)setNode:(MatchaBuildNode *)value {
    _node = value;
}

- (void)setBackgroundColor:(UIColor *)backgroundColor {
    // no-op
}

- (void)layoutSubviews {
    [super layoutSubviews];
    self.label.frame = self.bounds;
}

@end

