//
//  CustomView.m
//  CustomView
//
//  Created by Kevin Dang on 7/18/17.
//  Copyright © 2017 Matcha. All rights reserved.
//

#import "CustomView.h"

@implementation CustomView

+ (void)load {
    MatchaRegisterView(@"github.com/overcyn/customview", ^(MatchaViewNode *node){
        return [[CustomView alloc] initWithViewNode:node];
    });
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        self.switchView = [[UISwitch alloc] init];
        [self addSubview:self.switchView];
    }
    return self;
}

- (void)setNode:(MatchaBuildNode *)value {
    _node = value;
}

- (void)layoutSubviews {
    [super layoutSubviews];
    self.switchView.frame = self.bounds;
}

@end
