#import "CameraView.h"
#import "Customview.pbobjc.h"
#import <AVFoundation/AVFoundation.h>

#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wdeprecated-declarations"

@interface CameraView ()
@property (nonatomic, strong) AVCaptureSession *captureSession;
@property (nonatomic, strong) AVCaptureDeviceInput *captureInput;
@property (nonatomic, strong) AVCaptureStillImageOutput *captureOutput;
@property (nonatomic, strong) AVCaptureVideoPreviewLayer *previewLayer;
@property (nonatomic, strong) UIView *buttonView;
@property (nonatomic, assign) BOOL frontCamera;
@end

@implementation CameraView

+ (void)load {
    [MatchaViewController registerView:@"gomatcha.io/matcha/examples/customview CameraView" block:^(MatchaViewNode *node){
        return [[CameraView alloc] initWithViewNode:node];
    }];
}

- (id)initWithViewNode:(MatchaViewNode *)viewNode {
    if ((self = [super initWithFrame:CGRectZero])) {
        self.viewNode = viewNode;
        
        // Capture session
        self.captureSession = [[AVCaptureSession alloc] init];
        self.captureSession.sessionPreset = AVCaptureSessionPresetPhoto;
        
        // Capture input
//        AVCaptureDevice *device = [AVCaptureDevice defaultDeviceWithMediaType:AVMediaTypeVideo];
//        NSError *error = nil;
//        self.captureInput = [AVCaptureDeviceInput deviceInputWithDevice:device error:&error];
//        if (input) {
//            [self.captureSession addInput:input];
//        } else {
//            NSLog(@"Couldn't create video capture device");
//        }
        [self.captureSession startRunning];
        
        // Capture output
        self.captureOutput = [[AVCaptureStillImageOutput alloc] init];
        self.captureOutput.outputSettings = @{AVVideoCodecKey: AVVideoCodecJPEG};
        if ([self.captureSession canAddOutput:self.captureOutput]) {
            [self.captureSession addOutput:self.captureOutput];
        }

        dispatch_async(dispatch_get_main_queue(), ^{
            self.previewLayer = [[AVCaptureVideoPreviewLayer alloc] initWithSession:self.captureSession];
            [self.layer addSublayer:self.previewLayer];
            self.previewLayer.frame = self.bounds;
        });
        
        self.buttonView = [[UIView alloc] init];
        self.buttonView.backgroundColor = [UIColor whiteColor];
        self.buttonView.layer.cornerRadius = 25;
        self.buttonView.clipsToBounds = true;
        [self.buttonView addGestureRecognizer:[[UITapGestureRecognizer alloc] initWithTarget:self action:@selector(tapAction)]];
        [self addSubview:self.buttonView];
    }
    return self;
}

- (void)dealloc {
    [self.captureSession stopRunning];
}

- (void)layoutSubviews {
    [super layoutSubviews];
    
    CGRect b = self.bounds;
    {
        CGRect f;
        f.size.width = 50;
        f.size.height = 50;
        f.origin.x = roundf(b.origin.x + b.size.width / 2 - f.size.width / 2);
        f.origin.y = roundf(CGRectGetMaxY(b) - 60);
        self.buttonView.frame = f;
        [self bringSubviewToFront:self.buttonView];
    }
}

- (void)setNativeState:(NSData *)nativeState {
    CustomViewProtoView *view = [CustomViewProtoView parseFromData:nativeState error:nil];
    
    // Set the camera input.
    if (view.frontCamera != self.frontCamera || self.captureInput == nil) {
        self.frontCamera = view.frontCamera;
        
        [self.captureSession beginConfiguration];
        
        if (self.captureInput != nil) {
            [self.captureSession removeInput:self.captureInput];
        }
        
        AVCaptureDevice *camera = nil;
        for (AVCaptureDevice *i in [AVCaptureDevice devicesWithMediaType:AVMediaTypeVideo]) {
            if ((view.frontCamera && i.position == AVCaptureDevicePositionFront) || (!view.frontCamera && i.position == AVCaptureDevicePositionBack)) {
                camera = i;
                break;
            }
        }
        
        self.captureInput = [[AVCaptureDeviceInput alloc] initWithDevice:camera error:nil];
        if (self.captureInput) {
            [self.captureSession addInput:self.captureInput];
        }
        
        [self.captureSession commitConfiguration];
    }
}

- (void)tapAction {
    AVCaptureConnection *connection = nil;
    for (AVCaptureConnection *i in self.captureOutput.connections){
        for (AVCaptureInputPort *port in [i inputPorts]) {
            if ([[port mediaType] isEqual:AVMediaTypeVideo]) {
                connection = i;
                break;
            }
        }
    }
    
    [self.captureOutput captureStillImageAsynchronouslyFromConnection:connection completionHandler:^(CMSampleBufferRef buf, NSError *error) {
        if (buf != NULL) {
            NSData *imageData = [AVCaptureStillImageOutput jpegStillImageNSDataRepresentation:buf];
            [self.viewNode call:@"OnCapture", [[MatchaGoValue alloc] initWithData:imageData], nil];
        }
    }];
}

@end

#pragma GCC diagnostic pop
