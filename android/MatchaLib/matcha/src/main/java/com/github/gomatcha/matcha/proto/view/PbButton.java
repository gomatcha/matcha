// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: github.com/gomatcha/matcha/proto/view/button.proto

package com.github.gomatcha.matcha.proto.view;

public final class PbButton {
  private PbButton() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  public interface ButtonOrBuilder extends
      // @@protoc_insertion_point(interface_extends:matcha.view.Button)
      com.google.protobuf.MessageOrBuilder {

    /**
     * <code>string str = 1;</code>
     */
    java.lang.String getStr();
    /**
     * <code>string str = 1;</code>
     */
    com.google.protobuf.ByteString
        getStrBytes();

    /**
     * <code>bool enabled = 2;</code>
     */
    boolean getEnabled();

    /**
     * <code>.matcha.Color color = 3;</code>
     */
    boolean hasColor();
    /**
     * <code>.matcha.Color color = 3;</code>
     */
    com.github.gomatcha.matcha.proto.Proto.Color getColor();
    /**
     * <code>.matcha.Color color = 3;</code>
     */
    com.github.gomatcha.matcha.proto.Proto.ColorOrBuilder getColorOrBuilder();
  }
  /**
   * Protobuf type {@code matcha.view.Button}
   */
  public  static final class Button extends
      com.google.protobuf.GeneratedMessageV3 implements
      // @@protoc_insertion_point(message_implements:matcha.view.Button)
      ButtonOrBuilder {
  private static final long serialVersionUID = 0L;
    // Use Button.newBuilder() to construct.
    private Button(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
      super(builder);
    }
    private Button() {
      str_ = "";
      enabled_ = false;
    }

    @java.lang.Override
    public final com.google.protobuf.UnknownFieldSet
    getUnknownFields() {
      return this.unknownFields;
    }
    private Button(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      this();
      if (extensionRegistry == null) {
        throw new java.lang.NullPointerException();
      }
      int mutable_bitField0_ = 0;
      com.google.protobuf.UnknownFieldSet.Builder unknownFields =
          com.google.protobuf.UnknownFieldSet.newBuilder();
      try {
        boolean done = false;
        while (!done) {
          int tag = input.readTag();
          switch (tag) {
            case 0:
              done = true;
              break;
            case 10: {
              java.lang.String s = input.readStringRequireUtf8();

              str_ = s;
              break;
            }
            case 16: {

              enabled_ = input.readBool();
              break;
            }
            case 26: {
              com.github.gomatcha.matcha.proto.Proto.Color.Builder subBuilder = null;
              if (color_ != null) {
                subBuilder = color_.toBuilder();
              }
              color_ = input.readMessage(com.github.gomatcha.matcha.proto.Proto.Color.parser(), extensionRegistry);
              if (subBuilder != null) {
                subBuilder.mergeFrom(color_);
                color_ = subBuilder.buildPartial();
              }

              break;
            }
            default: {
              if (!parseUnknownFieldProto3(
                  input, unknownFields, extensionRegistry, tag)) {
                done = true;
              }
              break;
            }
          }
        }
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        throw e.setUnfinishedMessage(this);
      } catch (java.io.IOException e) {
        throw new com.google.protobuf.InvalidProtocolBufferException(
            e).setUnfinishedMessage(this);
      } finally {
        this.unknownFields = unknownFields.build();
        makeExtensionsImmutable();
      }
    }
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return com.github.gomatcha.matcha.proto.view.PbButton.internal_static_matcha_view_Button_descriptor;
    }

    @java.lang.Override
    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return com.github.gomatcha.matcha.proto.view.PbButton.internal_static_matcha_view_Button_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              com.github.gomatcha.matcha.proto.view.PbButton.Button.class, com.github.gomatcha.matcha.proto.view.PbButton.Button.Builder.class);
    }

    public static final int STR_FIELD_NUMBER = 1;
    private volatile java.lang.Object str_;
    /**
     * <code>string str = 1;</code>
     */
    public java.lang.String getStr() {
      java.lang.Object ref = str_;
      if (ref instanceof java.lang.String) {
        return (java.lang.String) ref;
      } else {
        com.google.protobuf.ByteString bs = 
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        str_ = s;
        return s;
      }
    }
    /**
     * <code>string str = 1;</code>
     */
    public com.google.protobuf.ByteString
        getStrBytes() {
      java.lang.Object ref = str_;
      if (ref instanceof java.lang.String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        str_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
    }

    public static final int ENABLED_FIELD_NUMBER = 2;
    private boolean enabled_;
    /**
     * <code>bool enabled = 2;</code>
     */
    public boolean getEnabled() {
      return enabled_;
    }

    public static final int COLOR_FIELD_NUMBER = 3;
    private com.github.gomatcha.matcha.proto.Proto.Color color_;
    /**
     * <code>.matcha.Color color = 3;</code>
     */
    public boolean hasColor() {
      return color_ != null;
    }
    /**
     * <code>.matcha.Color color = 3;</code>
     */
    public com.github.gomatcha.matcha.proto.Proto.Color getColor() {
      return color_ == null ? com.github.gomatcha.matcha.proto.Proto.Color.getDefaultInstance() : color_;
    }
    /**
     * <code>.matcha.Color color = 3;</code>
     */
    public com.github.gomatcha.matcha.proto.Proto.ColorOrBuilder getColorOrBuilder() {
      return getColor();
    }

    private byte memoizedIsInitialized = -1;
    @java.lang.Override
    public final boolean isInitialized() {
      byte isInitialized = memoizedIsInitialized;
      if (isInitialized == 1) return true;
      if (isInitialized == 0) return false;

      memoizedIsInitialized = 1;
      return true;
    }

    @java.lang.Override
    public void writeTo(com.google.protobuf.CodedOutputStream output)
                        throws java.io.IOException {
      if (!getStrBytes().isEmpty()) {
        com.google.protobuf.GeneratedMessageV3.writeString(output, 1, str_);
      }
      if (enabled_ != false) {
        output.writeBool(2, enabled_);
      }
      if (color_ != null) {
        output.writeMessage(3, getColor());
      }
      unknownFields.writeTo(output);
    }

    @java.lang.Override
    public int getSerializedSize() {
      int size = memoizedSize;
      if (size != -1) return size;

      size = 0;
      if (!getStrBytes().isEmpty()) {
        size += com.google.protobuf.GeneratedMessageV3.computeStringSize(1, str_);
      }
      if (enabled_ != false) {
        size += com.google.protobuf.CodedOutputStream
          .computeBoolSize(2, enabled_);
      }
      if (color_ != null) {
        size += com.google.protobuf.CodedOutputStream
          .computeMessageSize(3, getColor());
      }
      size += unknownFields.getSerializedSize();
      memoizedSize = size;
      return size;
    }

    @java.lang.Override
    public boolean equals(final java.lang.Object obj) {
      if (obj == this) {
       return true;
      }
      if (!(obj instanceof com.github.gomatcha.matcha.proto.view.PbButton.Button)) {
        return super.equals(obj);
      }
      com.github.gomatcha.matcha.proto.view.PbButton.Button other = (com.github.gomatcha.matcha.proto.view.PbButton.Button) obj;

      boolean result = true;
      result = result && getStr()
          .equals(other.getStr());
      result = result && (getEnabled()
          == other.getEnabled());
      result = result && (hasColor() == other.hasColor());
      if (hasColor()) {
        result = result && getColor()
            .equals(other.getColor());
      }
      result = result && unknownFields.equals(other.unknownFields);
      return result;
    }

    @java.lang.Override
    public int hashCode() {
      if (memoizedHashCode != 0) {
        return memoizedHashCode;
      }
      int hash = 41;
      hash = (19 * hash) + getDescriptor().hashCode();
      hash = (37 * hash) + STR_FIELD_NUMBER;
      hash = (53 * hash) + getStr().hashCode();
      hash = (37 * hash) + ENABLED_FIELD_NUMBER;
      hash = (53 * hash) + com.google.protobuf.Internal.hashBoolean(
          getEnabled());
      if (hasColor()) {
        hash = (37 * hash) + COLOR_FIELD_NUMBER;
        hash = (53 * hash) + getColor().hashCode();
      }
      hash = (29 * hash) + unknownFields.hashCode();
      memoizedHashCode = hash;
      return hash;
    }

    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseFrom(
        java.nio.ByteBuffer data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseFrom(
        java.nio.ByteBuffer data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseFrom(
        com.google.protobuf.ByteString data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseFrom(
        com.google.protobuf.ByteString data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseFrom(byte[] data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseFrom(
        byte[] data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input, extensionRegistry);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseDelimitedFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseDelimitedFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseFrom(
        com.google.protobuf.CodedInputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static com.github.gomatcha.matcha.proto.view.PbButton.Button parseFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input, extensionRegistry);
    }

    @java.lang.Override
    public Builder newBuilderForType() { return newBuilder(); }
    public static Builder newBuilder() {
      return DEFAULT_INSTANCE.toBuilder();
    }
    public static Builder newBuilder(com.github.gomatcha.matcha.proto.view.PbButton.Button prototype) {
      return DEFAULT_INSTANCE.toBuilder().mergeFrom(prototype);
    }
    @java.lang.Override
    public Builder toBuilder() {
      return this == DEFAULT_INSTANCE
          ? new Builder() : new Builder().mergeFrom(this);
    }

    @java.lang.Override
    protected Builder newBuilderForType(
        com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
      Builder builder = new Builder(parent);
      return builder;
    }
    /**
     * Protobuf type {@code matcha.view.Button}
     */
    public static final class Builder extends
        com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
        // @@protoc_insertion_point(builder_implements:matcha.view.Button)
        com.github.gomatcha.matcha.proto.view.PbButton.ButtonOrBuilder {
      public static final com.google.protobuf.Descriptors.Descriptor
          getDescriptor() {
        return com.github.gomatcha.matcha.proto.view.PbButton.internal_static_matcha_view_Button_descriptor;
      }

      @java.lang.Override
      protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
          internalGetFieldAccessorTable() {
        return com.github.gomatcha.matcha.proto.view.PbButton.internal_static_matcha_view_Button_fieldAccessorTable
            .ensureFieldAccessorsInitialized(
                com.github.gomatcha.matcha.proto.view.PbButton.Button.class, com.github.gomatcha.matcha.proto.view.PbButton.Button.Builder.class);
      }

      // Construct using com.github.gomatcha.matcha.proto.view.PbButton.Button.newBuilder()
      private Builder() {
        maybeForceBuilderInitialization();
      }

      private Builder(
          com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
        super(parent);
        maybeForceBuilderInitialization();
      }
      private void maybeForceBuilderInitialization() {
        if (com.google.protobuf.GeneratedMessageV3
                .alwaysUseFieldBuilders) {
        }
      }
      @java.lang.Override
      public Builder clear() {
        super.clear();
        str_ = "";

        enabled_ = false;

        if (colorBuilder_ == null) {
          color_ = null;
        } else {
          color_ = null;
          colorBuilder_ = null;
        }
        return this;
      }

      @java.lang.Override
      public com.google.protobuf.Descriptors.Descriptor
          getDescriptorForType() {
        return com.github.gomatcha.matcha.proto.view.PbButton.internal_static_matcha_view_Button_descriptor;
      }

      @java.lang.Override
      public com.github.gomatcha.matcha.proto.view.PbButton.Button getDefaultInstanceForType() {
        return com.github.gomatcha.matcha.proto.view.PbButton.Button.getDefaultInstance();
      }

      @java.lang.Override
      public com.github.gomatcha.matcha.proto.view.PbButton.Button build() {
        com.github.gomatcha.matcha.proto.view.PbButton.Button result = buildPartial();
        if (!result.isInitialized()) {
          throw newUninitializedMessageException(result);
        }
        return result;
      }

      @java.lang.Override
      public com.github.gomatcha.matcha.proto.view.PbButton.Button buildPartial() {
        com.github.gomatcha.matcha.proto.view.PbButton.Button result = new com.github.gomatcha.matcha.proto.view.PbButton.Button(this);
        result.str_ = str_;
        result.enabled_ = enabled_;
        if (colorBuilder_ == null) {
          result.color_ = color_;
        } else {
          result.color_ = colorBuilder_.build();
        }
        onBuilt();
        return result;
      }

      @java.lang.Override
      public Builder clone() {
        return (Builder) super.clone();
      }
      @java.lang.Override
      public Builder setField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          java.lang.Object value) {
        return (Builder) super.setField(field, value);
      }
      @java.lang.Override
      public Builder clearField(
          com.google.protobuf.Descriptors.FieldDescriptor field) {
        return (Builder) super.clearField(field);
      }
      @java.lang.Override
      public Builder clearOneof(
          com.google.protobuf.Descriptors.OneofDescriptor oneof) {
        return (Builder) super.clearOneof(oneof);
      }
      @java.lang.Override
      public Builder setRepeatedField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          int index, java.lang.Object value) {
        return (Builder) super.setRepeatedField(field, index, value);
      }
      @java.lang.Override
      public Builder addRepeatedField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          java.lang.Object value) {
        return (Builder) super.addRepeatedField(field, value);
      }
      @java.lang.Override
      public Builder mergeFrom(com.google.protobuf.Message other) {
        if (other instanceof com.github.gomatcha.matcha.proto.view.PbButton.Button) {
          return mergeFrom((com.github.gomatcha.matcha.proto.view.PbButton.Button)other);
        } else {
          super.mergeFrom(other);
          return this;
        }
      }

      public Builder mergeFrom(com.github.gomatcha.matcha.proto.view.PbButton.Button other) {
        if (other == com.github.gomatcha.matcha.proto.view.PbButton.Button.getDefaultInstance()) return this;
        if (!other.getStr().isEmpty()) {
          str_ = other.str_;
          onChanged();
        }
        if (other.getEnabled() != false) {
          setEnabled(other.getEnabled());
        }
        if (other.hasColor()) {
          mergeColor(other.getColor());
        }
        this.mergeUnknownFields(other.unknownFields);
        onChanged();
        return this;
      }

      @java.lang.Override
      public final boolean isInitialized() {
        return true;
      }

      @java.lang.Override
      public Builder mergeFrom(
          com.google.protobuf.CodedInputStream input,
          com.google.protobuf.ExtensionRegistryLite extensionRegistry)
          throws java.io.IOException {
        com.github.gomatcha.matcha.proto.view.PbButton.Button parsedMessage = null;
        try {
          parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
        } catch (com.google.protobuf.InvalidProtocolBufferException e) {
          parsedMessage = (com.github.gomatcha.matcha.proto.view.PbButton.Button) e.getUnfinishedMessage();
          throw e.unwrapIOException();
        } finally {
          if (parsedMessage != null) {
            mergeFrom(parsedMessage);
          }
        }
        return this;
      }

      private java.lang.Object str_ = "";
      /**
       * <code>string str = 1;</code>
       */
      public java.lang.String getStr() {
        java.lang.Object ref = str_;
        if (!(ref instanceof java.lang.String)) {
          com.google.protobuf.ByteString bs =
              (com.google.protobuf.ByteString) ref;
          java.lang.String s = bs.toStringUtf8();
          str_ = s;
          return s;
        } else {
          return (java.lang.String) ref;
        }
      }
      /**
       * <code>string str = 1;</code>
       */
      public com.google.protobuf.ByteString
          getStrBytes() {
        java.lang.Object ref = str_;
        if (ref instanceof String) {
          com.google.protobuf.ByteString b = 
              com.google.protobuf.ByteString.copyFromUtf8(
                  (java.lang.String) ref);
          str_ = b;
          return b;
        } else {
          return (com.google.protobuf.ByteString) ref;
        }
      }
      /**
       * <code>string str = 1;</code>
       */
      public Builder setStr(
          java.lang.String value) {
        if (value == null) {
    throw new NullPointerException();
  }
  
        str_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>string str = 1;</code>
       */
      public Builder clearStr() {
        
        str_ = getDefaultInstance().getStr();
        onChanged();
        return this;
      }
      /**
       * <code>string str = 1;</code>
       */
      public Builder setStrBytes(
          com.google.protobuf.ByteString value) {
        if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);
        
        str_ = value;
        onChanged();
        return this;
      }

      private boolean enabled_ ;
      /**
       * <code>bool enabled = 2;</code>
       */
      public boolean getEnabled() {
        return enabled_;
      }
      /**
       * <code>bool enabled = 2;</code>
       */
      public Builder setEnabled(boolean value) {
        
        enabled_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>bool enabled = 2;</code>
       */
      public Builder clearEnabled() {
        
        enabled_ = false;
        onChanged();
        return this;
      }

      private com.github.gomatcha.matcha.proto.Proto.Color color_ = null;
      private com.google.protobuf.SingleFieldBuilderV3<
          com.github.gomatcha.matcha.proto.Proto.Color, com.github.gomatcha.matcha.proto.Proto.Color.Builder, com.github.gomatcha.matcha.proto.Proto.ColorOrBuilder> colorBuilder_;
      /**
       * <code>.matcha.Color color = 3;</code>
       */
      public boolean hasColor() {
        return colorBuilder_ != null || color_ != null;
      }
      /**
       * <code>.matcha.Color color = 3;</code>
       */
      public com.github.gomatcha.matcha.proto.Proto.Color getColor() {
        if (colorBuilder_ == null) {
          return color_ == null ? com.github.gomatcha.matcha.proto.Proto.Color.getDefaultInstance() : color_;
        } else {
          return colorBuilder_.getMessage();
        }
      }
      /**
       * <code>.matcha.Color color = 3;</code>
       */
      public Builder setColor(com.github.gomatcha.matcha.proto.Proto.Color value) {
        if (colorBuilder_ == null) {
          if (value == null) {
            throw new NullPointerException();
          }
          color_ = value;
          onChanged();
        } else {
          colorBuilder_.setMessage(value);
        }

        return this;
      }
      /**
       * <code>.matcha.Color color = 3;</code>
       */
      public Builder setColor(
          com.github.gomatcha.matcha.proto.Proto.Color.Builder builderForValue) {
        if (colorBuilder_ == null) {
          color_ = builderForValue.build();
          onChanged();
        } else {
          colorBuilder_.setMessage(builderForValue.build());
        }

        return this;
      }
      /**
       * <code>.matcha.Color color = 3;</code>
       */
      public Builder mergeColor(com.github.gomatcha.matcha.proto.Proto.Color value) {
        if (colorBuilder_ == null) {
          if (color_ != null) {
            color_ =
              com.github.gomatcha.matcha.proto.Proto.Color.newBuilder(color_).mergeFrom(value).buildPartial();
          } else {
            color_ = value;
          }
          onChanged();
        } else {
          colorBuilder_.mergeFrom(value);
        }

        return this;
      }
      /**
       * <code>.matcha.Color color = 3;</code>
       */
      public Builder clearColor() {
        if (colorBuilder_ == null) {
          color_ = null;
          onChanged();
        } else {
          color_ = null;
          colorBuilder_ = null;
        }

        return this;
      }
      /**
       * <code>.matcha.Color color = 3;</code>
       */
      public com.github.gomatcha.matcha.proto.Proto.Color.Builder getColorBuilder() {
        
        onChanged();
        return getColorFieldBuilder().getBuilder();
      }
      /**
       * <code>.matcha.Color color = 3;</code>
       */
      public com.github.gomatcha.matcha.proto.Proto.ColorOrBuilder getColorOrBuilder() {
        if (colorBuilder_ != null) {
          return colorBuilder_.getMessageOrBuilder();
        } else {
          return color_ == null ?
              com.github.gomatcha.matcha.proto.Proto.Color.getDefaultInstance() : color_;
        }
      }
      /**
       * <code>.matcha.Color color = 3;</code>
       */
      private com.google.protobuf.SingleFieldBuilderV3<
          com.github.gomatcha.matcha.proto.Proto.Color, com.github.gomatcha.matcha.proto.Proto.Color.Builder, com.github.gomatcha.matcha.proto.Proto.ColorOrBuilder> 
          getColorFieldBuilder() {
        if (colorBuilder_ == null) {
          colorBuilder_ = new com.google.protobuf.SingleFieldBuilderV3<
              com.github.gomatcha.matcha.proto.Proto.Color, com.github.gomatcha.matcha.proto.Proto.Color.Builder, com.github.gomatcha.matcha.proto.Proto.ColorOrBuilder>(
                  getColor(),
                  getParentForChildren(),
                  isClean());
          color_ = null;
        }
        return colorBuilder_;
      }
      @java.lang.Override
      public final Builder setUnknownFields(
          final com.google.protobuf.UnknownFieldSet unknownFields) {
        return super.setUnknownFieldsProto3(unknownFields);
      }

      @java.lang.Override
      public final Builder mergeUnknownFields(
          final com.google.protobuf.UnknownFieldSet unknownFields) {
        return super.mergeUnknownFields(unknownFields);
      }


      // @@protoc_insertion_point(builder_scope:matcha.view.Button)
    }

    // @@protoc_insertion_point(class_scope:matcha.view.Button)
    private static final com.github.gomatcha.matcha.proto.view.PbButton.Button DEFAULT_INSTANCE;
    static {
      DEFAULT_INSTANCE = new com.github.gomatcha.matcha.proto.view.PbButton.Button();
    }

    public static com.github.gomatcha.matcha.proto.view.PbButton.Button getDefaultInstance() {
      return DEFAULT_INSTANCE;
    }

    private static final com.google.protobuf.Parser<Button>
        PARSER = new com.google.protobuf.AbstractParser<Button>() {
      @java.lang.Override
      public Button parsePartialFrom(
          com.google.protobuf.CodedInputStream input,
          com.google.protobuf.ExtensionRegistryLite extensionRegistry)
          throws com.google.protobuf.InvalidProtocolBufferException {
        return new Button(input, extensionRegistry);
      }
    };

    public static com.google.protobuf.Parser<Button> parser() {
      return PARSER;
    }

    @java.lang.Override
    public com.google.protobuf.Parser<Button> getParserForType() {
      return PARSER;
    }

    @java.lang.Override
    public com.github.gomatcha.matcha.proto.view.PbButton.Button getDefaultInstanceForType() {
      return DEFAULT_INSTANCE;
    }

  }

  private static final com.google.protobuf.Descriptors.Descriptor
    internal_static_matcha_view_Button_descriptor;
  private static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_matcha_view_Button_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n2github.com/gomatcha/matcha/proto/view/" +
      "button.proto\022\013matcha.view\032,github.com/go" +
      "matcha/matcha/proto/image.proto\"D\n\006Butto" +
      "n\022\013\n\003str\030\001 \001(\t\022\017\n\007enabled\030\002 \001(\010\022\034\n\005color" +
      "\030\003 \001(\0132\r.matcha.ColorBF\n%com.github.goma" +
      "tcha.matcha.proto.viewB\010PbButtonZ\004view\242\002" +
      "\014MatchaViewPbb\006proto3"
    };
    com.google.protobuf.Descriptors.FileDescriptor.InternalDescriptorAssigner assigner =
        new com.google.protobuf.Descriptors.FileDescriptor.    InternalDescriptorAssigner() {
          public com.google.protobuf.ExtensionRegistry assignDescriptors(
              com.google.protobuf.Descriptors.FileDescriptor root) {
            descriptor = root;
            return null;
          }
        };
    com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          com.github.gomatcha.matcha.proto.Proto.getDescriptor(),
        }, assigner);
    internal_static_matcha_view_Button_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_matcha_view_Button_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_matcha_view_Button_descriptor,
        new java.lang.String[] { "Str", "Enabled", "Color", });
    com.github.gomatcha.matcha.proto.Proto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
