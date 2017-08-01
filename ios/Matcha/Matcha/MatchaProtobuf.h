#import <UIKit/UIKit.h>
#import <MatchaBridge/MatchaBridge.h>
#import <Matcha/MatchaViewController.h>
#import <Protobuf/Protobuf.h>
#import "View.pbobjc.h"
#import "Layout.pbobjc.h"
#import "Text.pbobjc.h"
#import "Scrollview.pbobjc.h"
#import "Imageview.pbobjc.h"
#import "Button.pbobjc.h"
#import "Paint.pbobjc.h"
#import "Tabscreen.pbobjc.h"
#import "Stacknavigator.pbobjc.h"
#import "Switchview.pbobjc.h"
#import "Touch2.pbobjc.h"
#import "Resource.pbobjc.h"
#import "Color.pbobjc.h"
#import "Image.pbobjc.h"
#import "Textinput.pbobjc.h"
#import "Keyboard.pbobjc.h"
#import "Slider.pbobjc.h"
#import "ProgressView.pbobjc.h"
#import "SegmentView.pbobjc.h"
#import "Alert.pbobjc.h"
#import "Statusbar.pbobjc.h"

@interface UIColor (Matcha)
- (id)initWithProtobuf:(MatchaPBColor *)value;
- (MatchaPBColor *)protobuf;
@end

@interface NSAttributedString (Matcha)
- (id)initWithProtobuf:(MatchaPBStyledText *)value;
+ (NSDictionary *)attributesWithProtobuf:(MatchaPBTextStyle *)style;
- (MatchaPBStyledText *)protobuf;
@end

@interface UIFont (Matcha)
- (id)initWithProtobuf:(MatchaPBFont *)value;
- (MatchaPBFont *)protobuf;
@end

@interface UIImage (Matcha)
- (id)initWithProtobuf:(MatchaPBImage *)value;
- (id)initWithImageOrResourceProtobuf:(MatchaPBImageOrResource *)value;
- (MatchaPBImage *)protobuf;
@end

@interface MatchaViewPBLayoutPaintNode (Matcha)
@property (nonatomic, readonly) CGRect frame;
@property (nonatomic, readonly) MatchaColor matchaBackgroundColor;
@end

@interface MatchaLayoutPBRect (Matcha)
- (id)initWithCGRect:(CGRect)rect;
@property (nonatomic, readonly) CGRect toCGRect;
@end

@interface MatchaLayoutPBPoint (Matcha)
- (id)initWithCGPoint:(CGPoint)point;
- (id)initWithCGSize:(CGSize)size;
@property (nonatomic, readonly) CGPoint toCGPoint;
@property (nonatomic, readonly) CGSize toCGSize;
@end

@interface MatchaLayoutPBInsets (Matcha)
@property (nonatomic, readonly) UIEdgeInsets toUIEdgeInsets;
@end

@interface GPBTimestamp (Matcha)
- (id)initWithDate:(NSDate *)date;
@property (nonatomic, readonly) NSDate *toDate;
@end

CGColorRef MatchaCGColorWithProtobuf(MatchaPBColor *value);
UIKeyboardType MatchaKeyboardTypeWithProtobuf(MatchaKeyboardPBType t);
UIKeyboardAppearance MatchaKeyboardAppearanceWithProtobuf(MatchaKeyboardPBAppearance t);
UIReturnKeyType MatchaReturnTypeWithProtobuf(MatchaKeyboardPBReturnType t);
