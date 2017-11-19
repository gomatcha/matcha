#import <UIKit/UIKit.h>
#import <Matcha/Matcha.h>

//! Project version number for CustomView.
FOUNDATION_EXPORT double CustomViewVersionNumber;

//! Project version string for CustomView.
FOUNDATION_EXPORT const unsigned char CustomViewVersionString[];

@interface CameraView : UIView <MatchaChildView>
@property (nonatomic, strong) MatchaViewNode *viewNode;
@end
