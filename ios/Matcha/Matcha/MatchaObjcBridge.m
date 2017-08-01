#import "MatchaObjcBridge.h"
#import "MatchaBridge.h"
#import "MatchaNode.h"
#import "MatchaViewController.h"
#import "MatchaDeadlockLogger.h"
#import "MatchaProtobuf.h"

@implementation MatchaObjcBridge (Extensions)

- (void)configure {
    static dispatch_once_t sOnce = 0;
    dispatch_once(&sOnce, ^{
        [MatchaDeadlockLogger sharedLogger]; // Initialize
    
        static CADisplayLink *displayLink = nil;
        if (displayLink == nil) {
            displayLink = [CADisplayLink displayLinkWithTarget:self selector:@selector(screenUpdate)];
    //        displayLink.preferredFramesPerSecond = 1;
            [displayLink addToRunLoop:[NSRunLoop mainRunLoop] forMode:NSRunLoopCommonModes];
        }
        
        MatchaGoValue *screenScaleFunc = [[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/internal/device setScreenScale"];
        [screenScaleFunc call:nil args:@[[[MatchaGoValue alloc] initWithDouble:UIScreen.mainScreen.scale]]];
    });
}

- (MatchaGoValue *)sizeForAttributedString:(NSData *)protobuf maxLines:(int)maxLines {
    MatchaPBSizeFunc *func = [[MatchaPBSizeFunc alloc] initWithData:protobuf error:nil];
    
    NSAttributedString *attrStr = [[NSAttributedString alloc] initWithProtobuf:func.text];
    CGRect rect = [attrStr boundingRectWithSize:func.maxSize.toCGSize options:NSStringDrawingUsesLineFragmentOrigin|NSStringDrawingUsesFontLeading context:nil];
    
    UIFont *font = [attrStr attributesAtIndex:0 effectiveRange:NULL][NSFontAttributeName];
    CGFloat height = rect.size.height;
    if (maxLines > 0 && height > font.pointSize * maxLines) {
        height = font.pointSize * maxLines;
    }
    
    MatchaLayoutPBPoint *point = [[MatchaLayoutPBPoint alloc] initWithCGSize:CGSizeMake(ceil(rect.size.width), ceil(height))];
    return [[MatchaGoValue alloc] initWithData:point.data];
}

- (void)screenUpdate {
    static MatchaGoValue *updateFunc = nil;
    if (updateFunc == nil) {
        updateFunc = [[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/animate screenUpdate"];
    }
    [updateFunc call:nil args:nil];
}

- (void)updateId:(NSInteger)identifier withProtobuf:(NSData *)protobuf {
    MatchaViewPBRoot *pbroot = [[MatchaViewPBRoot alloc] initWithData:protobuf error:nil];
    MatchaNodeRoot *root = [[MatchaNodeRoot alloc] initWithProtobuf:pbroot];
    
    MatchaViewController *vc = [MatchaViewController viewControllerWithIdentifier:identifier];
    [vc update:root];
}

- (NSString *)assetsDir {
     return [[NSBundle mainBundle] resourcePath];
}

- (MatchaGoValue *)imageForResource:(NSString *)path {
    UIImage *image = [UIImage imageNamed:path];
    if (image == nil) {
        return nil;
    }
    NSData *data = UIImagePNGRepresentation(image);
    return [[MatchaGoValue alloc] initWithData:data];
}

- (MatchaGoValue *)propertiesForResource:(NSString *)path {
    UIImage *image = [UIImage imageNamed:path];
    if (image == nil) {
        return nil;
    }
    MatchaPBImageProperties *props = [[MatchaPBImageProperties alloc] init];
    props.width = ceil(image.size.width * image.scale);
    props.height = ceil(image.size.height * image.scale);
    props.scale = image.scale;
    return [[MatchaGoValue alloc] initWithData:props.data];
}

- (void)displayAlert:(NSData *)protobuf {
    MatchaAlertPBView *pbalert = [[MatchaAlertPBView alloc] initWithData:protobuf error:nil];
    UIAlertController *alert = [UIAlertController alertControllerWithTitle:pbalert.title message:pbalert.message preferredStyle:UIAlertControllerStyleAlert];
    for (NSInteger i = 0; i < pbalert.buttonsArray.count; i++) {
        MatchaAlertPBButton *button = pbalert.buttonsArray[i];
        UIAlertAction *action = [UIAlertAction actionWithTitle:button.title style:(UIAlertActionStyle)button.style handler:^(UIAlertAction *a){
            MatchaGoValue *onPress = [[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/view/alert onPress"];
            [onPress call:nil args:@[[[MatchaGoValue alloc] initWithLongLong:pbalert.id_p], [[MatchaGoValue alloc] initWithLongLong:i]]];
        }];
        [alert addAction:action];
    }
    [[UIApplication sharedApplication].keyWindow.rootViewController presentViewController:alert animated:YES completion:nil];
    
}

@end
