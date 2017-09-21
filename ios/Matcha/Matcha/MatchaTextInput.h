#import <UIKit/UIKit.h>
#import "MatchaView.h"

@interface MatchaTextInput : UITextView <MatchaChildView, UITextViewDelegate>
@property (nonatomic, weak) MatchaViewNode *viewNode;
@property (nonatomic, assign) bool hasFocus;
@property (nonatomic, strong) NSAttributedString *attrStr2;
@property (nonatomic, assign) bool multiline;
@end
