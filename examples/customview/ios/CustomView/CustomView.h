//
//  CustomView.h
//  CustomView
//
//  Created by Kevin Dang on 7/18/17.
//  Copyright © 2017 Matcha. All rights reserved.
//

#import <UIKit/UIKit.h>
#import <Matcha/Matcha.h>

//! Project version number for CustomView.
FOUNDATION_EXPORT double CustomViewVersionNumber;

//! Project version string for CustomView.
FOUNDATION_EXPORT const unsigned char CustomViewVersionString[];

// In this header, you should import all the public headers of your framework using statements like #import <CustomView/PublicHeader.h>

@interface CustomView : UIView <MatchaChildView>
@property (nonatomic, strong) MatchaViewNode *viewNode;
@property (nonatomic, strong) MatchaBuildNode *node;
@property (nonatomic, strong) UISwitch *switchView;
@end
