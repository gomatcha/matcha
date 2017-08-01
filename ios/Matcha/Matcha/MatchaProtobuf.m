#import <Foundation/Foundation.h>
#import "MatchaProtobuf.h"

@implementation UIColor (Matcha)

- (id)initWithProtobuf:(MatchaPBColor *)value {
    if (value == nil) {
        return nil;
    }
    return [UIColor colorWithRed:((double)value.red)/0xffff green:((double)value.green)/0xffff blue:((double)value.blue)/0xffff alpha:((double)value.alpha)/0xffff];
}

- (MatchaPBColor *)protobuf {
    CGFloat red, green, blue, alpha;
    [self getRed:&red green:&green blue:&blue alpha:&alpha];
    
    MatchaPBColor *color = [[MatchaPBColor alloc] init];
    color.red = red*0xffff;
    color.green = green*0xffff;
    color.blue = blue*0xffff;
    color.alpha = alpha*0xffff;
    return color;
}

@end


CGColorRef MatchaCGColorWithProtobuf(MatchaPBColor *value) {
    if (value == nil) {
        return nil;
    }
    CGFloat colors[4] = {((double)value.red)/0xffff, ((double)value.green)/0xffff, ((double)value.blue)/0xffff, ((double)value.alpha)/0xffff};
    return CGColorCreate(CGColorSpaceCreateDeviceRGB(), colors);
}


@implementation NSAttributedString (Matcha)

- (id)initWithProtobuf:(MatchaPBStyledText *)value {
    NSString *string = value.text.text;
    NSDictionary *attributes = [NSAttributedString attributesWithProtobuf:value.style];
    return [[NSAttributedString alloc] initWithString:string attributes:attributes];
}

- (MatchaPBStyledText *)protobuf {
    MatchaPBText *text = [[MatchaPBText alloc] init];
    text.text = self.string;
    
    NSDictionary *attributes = nil;
    if (self.length != 0) {
        attributes = [self attributesAtIndex:0 effectiveRange:NULL];
    }
    
    MatchaPBStyledText *styledText = [[MatchaPBStyledText alloc] init];
    styledText.text = text;
    styledText.style = [NSAttributedString protobufWithAttributes:attributes];
    return styledText;
}

+ (NSDictionary *)attributesWithProtobuf:(MatchaPBTextStyle *)style {
    NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
    NSMutableDictionary *dictionary = [[NSMutableDictionary alloc] init];
    dictionary[NSParagraphStyleAttributeName] = paragraphStyle;
    
    NSTextAlignment alignment;
    switch (style.textAlignment) {
        case 0:
            alignment = NSTextAlignmentLeft;
            break;
        case 1:
            alignment = NSTextAlignmentRight;
            break;
        case 2:
            alignment = NSTextAlignmentCenter;
            break;
        case 3:
            alignment = NSTextAlignmentJustified;
            break;
        default:
            alignment = NSTextAlignmentLeft;
    }
    paragraphStyle.alignment = alignment;
    
    NSUnderlineStyle strikethroughStyle;
    switch (style.strikethroughStyle) {
        case 0:
            strikethroughStyle = NSUnderlineStyleNone;
            break;
        case 1:
            strikethroughStyle = NSUnderlineStyleSingle;
            break;
        case 2:
            strikethroughStyle = NSUnderlineStyleDouble;
            break;
        case 3:
            strikethroughStyle = NSUnderlineStyleThick;
            break;
        case 4:
            strikethroughStyle = NSUnderlinePatternDot;
            break;
        case 5:
            strikethroughStyle = NSUnderlinePatternDash;
            break;
        default:
            strikethroughStyle = NSUnderlineStyleNone;
    }
    dictionary[NSStrikethroughStyleAttributeName] = @(strikethroughStyle);
    
    dictionary[NSStrikethroughColorAttributeName] = [[UIColor alloc] initWithProtobuf:style.strikethroughColor];
    
    NSUnderlineStyle underlineStyle;
    switch (style.underlineStyle) {
        case 0:
            underlineStyle = NSUnderlineStyleNone;
            break;
        case 1:
            underlineStyle = NSUnderlineStyleSingle;
            break;
        case 2:
            underlineStyle = NSUnderlineStyleDouble;
            break;
        case 3:
            underlineStyle = NSUnderlineStyleThick;
            break;
        case 4:
            underlineStyle = NSUnderlinePatternDot;
            break;
        case 5:
            underlineStyle = NSUnderlinePatternDash;
            break;
        default:
            underlineStyle = NSUnderlineStyleNone;
    }
    dictionary[NSUnderlineStyleAttributeName] = @(underlineStyle);
    
    dictionary[NSUnderlineColorAttributeName] = [[UIColor alloc] initWithProtobuf:style.underlineColor];
    dictionary[NSFontAttributeName] = [[UIFont alloc] initWithProtobuf:style.font];
    dictionary[NSHyphenationFactorDocumentAttribute] = @(style.hyphenation);
    paragraphStyle.lineHeightMultiple = style.lineHeightMultiple;
    // TODO(KD): AttributeKeyMaxLines
    dictionary[NSForegroundColorAttributeName] = [[UIColor alloc] initWithProtobuf:style.textColor];
    // TODO(KD): AttributeKeyTextWrap
    // TODO(KD): AttributeKeyTruncation
    // TODO(KD): AttributeKeyTruncationString
    return dictionary;
}

+ (MatchaPBTextStyle *)protobufWithAttributes:(NSDictionary *)dictionary {
    MatchaPBTextStyle *style = [[MatchaPBTextStyle alloc] init];
    
    NSMutableParagraphStyle *paragraphStyle = dictionary[NSParagraphStyleAttributeName];
    if (paragraphStyle) {
        int alignment;
        switch (paragraphStyle.alignment) {
        case NSTextAlignmentLeft:
            alignment = 0;
            break;
        case NSTextAlignmentRight:
            alignment = 1;
            break;
        case NSTextAlignmentCenter:
            alignment = 2;
            break;
        case NSTextAlignmentJustified:
            alignment = 3;
            break;
        default:
            alignment = 0;
        }
        style.textAlignment = alignment;
    }
    
    if (dictionary[NSStrikethroughStyleAttributeName]) {
        int strikethroughStyle;
        switch (((NSNumber *)dictionary[NSStrikethroughStyleAttributeName]).integerValue) {
        case NSUnderlineStyleNone:
            strikethroughStyle = 0;
            break;
        case NSUnderlineStyleSingle:
            strikethroughStyle = 1;
            break;
        case NSUnderlineStyleDouble:
            strikethroughStyle = 2;
            break;
        case NSUnderlineStyleThick:
            strikethroughStyle = 3;
            break;
        case NSUnderlinePatternDot:
            strikethroughStyle = 4;
            break;
        case NSUnderlinePatternDash:
            strikethroughStyle = 5;
            break;
        default:
            strikethroughStyle = 0;
        }
        style.strikethroughStyle = strikethroughStyle;
    }
    
    if (dictionary[NSStrikethroughColorAttributeName]) {
        style.strikethroughColor = ((UIColor *)dictionary[NSStrikethroughColorAttributeName]).protobuf;
    }
    
    if (dictionary[NSUnderlineStyleAttributeName]) {
        int strikethroughStyle;
        switch (((NSNumber *)dictionary[NSUnderlineStyleAttributeName]).integerValue) {
        case NSUnderlineStyleNone:
            strikethroughStyle = 0;
            break;
        case NSUnderlineStyleSingle:
            strikethroughStyle = 1;
            break;
        case NSUnderlineStyleDouble:
            strikethroughStyle = 2;
            break;
        case NSUnderlineStyleThick:
            strikethroughStyle = 3;
            break;
        case NSUnderlinePatternDot:
            strikethroughStyle = 4;
            break;
        case NSUnderlinePatternDash:
            strikethroughStyle = 5;
            break;
        default:
            strikethroughStyle = 0;
        }
        style.underlineStyle = strikethroughStyle;
    }
    
    if (dictionary[NSUnderlineColorAttributeName]) {
        style.underlineColor = ((UIColor *)dictionary[NSUnderlineColorAttributeName]).protobuf;
    }
    
    if (dictionary[NSFontAttributeName]) {
        style.font = ((UIFont *)dictionary[NSFontAttributeName]).protobuf;
    }
    
    if (dictionary[NSHyphenationFactorDocumentAttribute]) {
        style.hyphenation = ((NSNumber *)dictionary[NSHyphenationFactorDocumentAttribute]).integerValue;
    }
    
    style.lineHeightMultiple = paragraphStyle.lineHeightMultiple;
    if (dictionary[NSForegroundColorAttributeName]) {
        style.textColor = ((UIColor *)dictionary[NSForegroundColorAttributeName]).protobuf;
    }  
    // TODO(KD): AttributeKeyTextWrap
    // TODO(KD): AttributeKeyTruncation
    // TODO(KD): AttributeKeyTruncationString
    return style;
}

@end

@implementation UIImage (Matcha)

- (id)initWithProtobuf:(MatchaPBImage *)value {
    CIImage *image = [CIImage imageWithBitmapData:value.data_p bytesPerRow:(size_t)value.stride size:CGSizeMake(value.width, value.height) format:kCIFormatRGBA8 colorSpace:CGColorSpaceCreateDeviceRGB()];
    return [self initWithCIImage:image];
}

- (id)initWithImageOrResourceProtobuf:(MatchaPBImageOrResource *)value {
    if (value.hasImage) {
        return [self initWithProtobuf:value.image];
    } else {
        return [UIImage imageNamed:value.path];
    }
}

- (MatchaPBImage *)protobuf {
    return nil;
}

@end

@implementation MatchaViewPBLayoutPaintNode (Matcha)

- (CGRect)frame {
    return CGRectMake(self.minx, self.miny, self.maxx - self.minx, self.maxy - self.miny);
}

- (MatchaColor) matchaBackgroundColor {
    MatchaPBColor *c = self.paintStyle.backgroundColor;
    MatchaColor color = {
        .red = c.red,
        .green = c.green,
        .blue = c.blue,
        .alpha = c.alpha,
    };
    return color;
}

@end

@implementation UIFont (Matcha)

- (id)initWithProtobuf:(MatchaPBFont *)value {
    NSMutableDictionary *attr = [NSMutableDictionary dictionary];
    attr[UIFontDescriptorFamilyAttribute] = value.family;
    attr[UIFontDescriptorFaceAttribute] = value.face;
    attr[UIFontDescriptorSizeAttribute] = @(value.size);
    
    UIFontDescriptor *desc = [[UIFontDescriptor alloc] initWithFontAttributes:attr];
    UIFont *font = [UIFont fontWithDescriptor:desc size:0];
    return font;
}

- (MatchaPBFont *)protobuf {
    NSDictionary *attr = self.fontDescriptor.fontAttributes;
    
    MatchaPBFont *font = [[MatchaPBFont alloc] init];
    font.family = attr[UIFontDescriptorFamilyAttribute];
    font.face = attr[UIFontDescriptorFaceAttribute];
    font.size = ((NSNumber *)attr[UIFontDescriptorSizeAttribute]).doubleValue;
    return font;
}

@end

@implementation MatchaLayoutPBRect (Matcha)

- (id)initWithCGRect:(CGRect)rect {
    if ((self = [super init])) {
        self.min = [[MatchaLayoutPBPoint alloc] initWithCGPoint:CGPointMake(rect.origin.x, rect.origin.y)];
        self.max = [[MatchaLayoutPBPoint alloc] initWithCGPoint:CGPointMake(rect.origin.x + rect.size.width, rect.origin.y + rect.size.height)];
    }
    return self;
}

- (CGRect)toCGRect {
    CGPoint min = self.min.toCGPoint;
    CGPoint max = self.max.toCGPoint;
    return CGRectMake(min.x, min.y, max.x - min.x, max.y - min.y);
}

@end
@implementation MatchaLayoutPBPoint (Matcha)

- (id)initWithCGPoint:(CGPoint)point {
    if ((self = [super init])) {
        self.x = point.x;
        self.y = point.y;
    }
    return self;
}

- (CGPoint)toCGPoint {
    return CGPointMake(self.x, self.y);
}

- (id)initWithCGSize:(CGSize)size {
    if ((self = [super init])) {
        self.x = size.width;
        self.y = size.height;
    }
    return self;
}

- (CGSize)toCGSize {
    return CGSizeMake(self.x, self.y);
}

@end

@implementation MatchaLayoutPBInsets (Matcha)
- (UIEdgeInsets)toUIEdgeInsets {
    return UIEdgeInsetsMake(self.top, self.left, self.bottom, self.right);
}
@end

@implementation GPBTimestamp (Matcha)
- (id)initWithDate:(NSDate *)date {
    if ((self = [super init])) {
        double integral;
        double fractional = modf([date timeIntervalSince1970], &integral);
        
        self.seconds = integral;
        self.nanos = fractional * NSEC_PER_SEC;
    }
    return self;
}

- (NSDate *)toDate {
    NSTimeInterval interval = (NSTimeInterval)self.seconds + (NSTimeInterval)self.nanos / NSEC_PER_SEC;
    return [NSDate dateWithTimeIntervalSince1970:interval];
}

@end

UIKeyboardType MatchaKeyboardTypeWithProtobuf(MatchaKeyboardPBType a) {
    UIKeyboardType t = UIKeyboardTypeDefault;
    switch (a) {
    case MatchaKeyboardPBType_GPBUnrecognizedEnumeratorValue:
    case MatchaKeyboardPBType_DefaultType: {
        t = UIKeyboardTypeDefault;
        break;
    }
    case MatchaKeyboardPBType_NumberType: {
        t = UIKeyboardTypeNumberPad;
        break;
    }
    case MatchaKeyboardPBType_NumberPunctuationType: {
        t = UIKeyboardTypeNumbersAndPunctuation;
        break;
    }
    case MatchaKeyboardPBType_DecimalType: {
        t = UIKeyboardTypeDecimalPad;
        break;
    }
    case MatchaKeyboardPBType_PhoneType: {
        t = UIKeyboardTypePhonePad;
        break;
    }
    case MatchaKeyboardPBType_AsciiType: {
        t = UIKeyboardTypeASCIICapable;
        break;
    }
    case MatchaKeyboardPBType_EmailType: {
        t = UIKeyboardTypeEmailAddress;
        break;
    }
    case MatchaKeyboardPBType_URLType: {
        t = UIKeyboardTypeURL;
        break;
    }
    case MatchaKeyboardPBType_WebSearchType: {
        t = UIKeyboardTypeWebSearch;
        break;
    }
    case MatchaKeyboardPBType_NamePhoneType: {
        t = UIKeyboardTypeNamePhonePad;
        break;
    }
    }
    return t;
}

UIKeyboardAppearance MatchaKeyboardAppearanceWithProtobuf(MatchaKeyboardPBAppearance a) {
    UIKeyboardAppearance t = UIKeyboardAppearanceDefault;
    switch (a) {
    case MatchaKeyboardPBAppearance_GPBUnrecognizedEnumeratorValue:
    case MatchaKeyboardPBAppearance_DefaultAppearance: {
        t = UIKeyboardAppearanceDefault;
        break;
    }
    case MatchaKeyboardPBAppearance_LightAppearance: {
        t = UIKeyboardAppearanceLight;
        break;
    }
    case MatchaKeyboardPBAppearance_DarkAppearance: {
        t = UIKeyboardAppearanceDark;
        break;
    }
    }
    return t;
}

UIReturnKeyType MatchaReturnTypeWithProtobuf(MatchaKeyboardPBReturnType a) {
    UIReturnKeyType t = UIReturnKeyDefault;
    switch (a) {
    case MatchaKeyboardPBReturnType_GPBUnrecognizedEnumeratorValue:
    case MatchaKeyboardPBReturnType_DefaultReturnType: {
        t = UIReturnKeyDefault;
        break;
    }
    case MatchaKeyboardPBReturnType_GoReturnType: {
        t = UIReturnKeyGo;
        break;
    }
    case MatchaKeyboardPBReturnType_GoogleReturnType: {
        t = UIReturnKeyGoogle;
        break;
    }
    case MatchaKeyboardPBReturnType_JoinReturnType: {
        t = UIReturnKeyJoin;
        break;
    }
    case MatchaKeyboardPBReturnType_NextReturnType: {
        t = UIReturnKeyNext;
        break;
    }
    case MatchaKeyboardPBReturnType_RouteReturnType: {
        t = UIReturnKeyRoute;
        break;
    }
    case MatchaKeyboardPBReturnType_SearchReturnType: {
        t = UIReturnKeySearch;
        break;
    }
    case MatchaKeyboardPBReturnType_SendReturnType: {
        t = UIReturnKeySend;
        break;
    }
    case MatchaKeyboardPBReturnType_YahooReturnType: {
        t = UIReturnKeyYahoo;
        break;
    }
    case MatchaKeyboardPBReturnType_DoneReturnType: {
        t = UIReturnKeyDone;
        break;
    }
    case MatchaKeyboardPBReturnType_EmergencyCallReturnType: {
        t = UIReturnKeyEmergencyCall;
        break;
    }
    case MatchaKeyboardPBReturnType_ContinueReturnType: {
        t = UIReturnKeyContinue;
        break;
    }
    }
    return t;
}
