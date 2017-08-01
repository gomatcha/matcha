#import "MatchaTextInput.h"
#import "MatchaProtobuf.h"
#import "MatchaViewController.h"
#import "UITextView+Placeholder.h"

@implementation MatchaTextInput

+ (void)load {
    MatchaRegisterView(@"gomatcha.io/matcha/view/textinput", ^(MatchaViewNode *node){
        return [[MatchaTextInput alloc] initWithViewNode:node];
    });
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        self.delegate = self;
        self.scrollEnabled = false;
    }
    return self;
}

- (void)setNode:(MatchaBuildNode *)value {
    _node = value;
    GPBAny *state = value.nativeViewState;
    NSError *error = nil;
    MatchaTextInputPBView *view = (id)[state unpackMessageClass:[MatchaTextInputPBView class] error:&error];
    
    NSDictionary *attributes = [NSAttributedString attributesWithProtobuf:view.styledText.style];
    self.font = attributes[NSFontAttributeName];
    self.textColor = attributes[NSForegroundColorAttributeName];
    
    NSAttributedString *attrString = [[NSAttributedString alloc] initWithProtobuf:view.styledText];
    if (![attrString.string isEqual:self.attributedText.string]) { // TODO(KD): Better comparison.
        self.attributedText = attrString;
    }
    
    NSAttributedString *attrPlaceholder = [[NSAttributedString alloc] initWithProtobuf:view.placeholderText];
    if (![attrPlaceholder.string isEqual:self.attributedPlaceholder]) { // TODO(KD): Better comparison.
        self.attributedPlaceholder = attrPlaceholder;
    }
    
    self.attrStr2 = attrString;
    self.hasFocus = view.focused;
    self.keyboardType = MatchaKeyboardTypeWithProtobuf(view.keyboardType);
    self.keyboardAppearance = MatchaKeyboardAppearanceWithProtobuf(view.keyboardAppearance);
    self.returnKeyType = MatchaReturnTypeWithProtobuf(view.keyboardReturnType);
    self.multiline = view.multiline;
    self.secureTextEntry = view.secureTextEntry;
    
    if (self.hasFocus && !self.isFirstResponder) {
        [self becomeFirstResponder];
    } else if (!self.hasFocus && self.isFirstResponder) {
        [self resignFirstResponder];
    }
}

- (void)textViewDidChange:(UITextView *)textView {
    if ([self.attributedText isEqual:self.attrStr2] || self.attributedText == self.attrStr2) {
        return;
    }
    self.attrStr2 = self.attributedText;
    MatchaTextInputPBEvent *event = [[MatchaTextInputPBEvent alloc] init];
    event.styledText = self.attributedText.protobuf;
    
    NSData *data = [event data];
    MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:data];
    
    [self.viewNode.rootVC call:@"OnTextChange" viewId:self.node.identifier.longLongValue args:@[value]];
}

- (BOOL)textView:(UITextView *)textView shouldChangeTextInRange:(NSRange)range replacementText:(NSString *)text {
    if ([text isEqualToString:@"\n"]) {
        [self.viewNode.rootVC call:@"OnSubmit" viewId:self.node.identifier.longLongValue args:nil];
        return NO;
    }
    return YES;
}

- (void)textViewDidBeginEditing:(UITextView *)textView {
    [self focusDidChange];
}

- (void)textViewDidEndEditing:(UITextView *)textView {
    [self focusDidChange];
}

- (void)focusDidChange {
    if ((self.hasFocus && !self.isFirstResponder) || (!self.hasFocus && self.isFirstResponder)) {
        MatchaTextInputPBFocusEvent *event = [[MatchaTextInputPBFocusEvent alloc] init];
        event.focused = self.isFirstResponder;
        
        MatchaGoValue *value = [[MatchaGoValue alloc] initWithData:event.data];
        [self.viewNode.rootVC call:@"OnFocus" viewId:self.node.identifier.longLongValue args:@[value]];
    }
}

@end
