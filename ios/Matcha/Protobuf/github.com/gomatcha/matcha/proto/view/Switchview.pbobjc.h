// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: github.com/gomatcha/matcha/proto/view/switchview.proto

// This CPP symbol can be defined to use imports that match up to the framework
// imports needed when using CocoaPods.
#if !defined(GPB_USE_PROTOBUF_FRAMEWORK_IMPORTS)
 #define GPB_USE_PROTOBUF_FRAMEWORK_IMPORTS 0
#endif

#if GPB_USE_PROTOBUF_FRAMEWORK_IMPORTS
 #import <Protobuf/GPBProtocolBuffers.h>
#else
 #import "GPBProtocolBuffers.h"
#endif

#if GOOGLE_PROTOBUF_OBJC_VERSION < 30002
#error This file was generated by a newer version of protoc which is incompatible with your Protocol Buffer library sources.
#endif
#if 30002 < GOOGLE_PROTOBUF_OBJC_MIN_SUPPORTED_VERSION
#error This file was generated by an older version of protoc which is incompatible with your Protocol Buffer library sources.
#endif

// @@protoc_insertion_point(imports)

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"

CF_EXTERN_C_BEGIN

NS_ASSUME_NONNULL_BEGIN

#pragma mark - MatchaViewPbSwitchviewRoot

/**
 * Exposes the extension registry for this file.
 *
 * The base class provides:
 * @code
 *   + (GPBExtensionRegistry *)extensionRegistry;
 * @endcode
 * which is a @c GPBExtensionRegistry that includes all the extensions defined by
 * this file and all files that it depends on.
 **/
@interface MatchaViewPbSwitchviewRoot : GPBRootObject
@end

#pragma mark - MatchaViewPbSwitchView

typedef GPB_ENUM(MatchaViewPbSwitchView_FieldNumber) {
  MatchaViewPbSwitchView_FieldNumber_Value = 1,
  MatchaViewPbSwitchView_FieldNumber_Enabled = 2,
};

@interface MatchaViewPbSwitchView : GPBMessage

@property(nonatomic, readwrite) BOOL value;

@property(nonatomic, readwrite) BOOL enabled;

@end

#pragma mark - MatchaViewPbSwitchEvent

typedef GPB_ENUM(MatchaViewPbSwitchEvent_FieldNumber) {
  MatchaViewPbSwitchEvent_FieldNumber_Value = 1,
};

@interface MatchaViewPbSwitchEvent : GPBMessage

@property(nonatomic, readwrite) BOOL value;

@end

NS_ASSUME_NONNULL_END

CF_EXTERN_C_END

#pragma clang diagnostic pop

// @@protoc_insertion_point(global_scope)
