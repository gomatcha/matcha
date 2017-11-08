// +build matcha,android

#ifndef JAVA_FOREIGN_H
#define JAVA_FOREIGN_H

#include <jni.h>

// Init
void MatchaInit(JNIEnv *env, jobject tracker);

// Utilities
CGoBuffer MatchaStringToCGoBuffer(JNIEnv *env, jstring v); // return buffer needs to be released.
jstring MatchaCGoBufferToString(JNIEnv *env, CGoBuffer buf); // releases buf, jstring needs to be released.
CGoBuffer MatchaByteArrayToCGoBuffer(JNIEnv *env, jbyteArray v); // returned buffer needs to be released.
jbyteArray MatchaCGoBufferToByteArray(JNIEnv *env, CGoBuffer buf); // releases buf
CGoBuffer MatchaJlongArrayToCGoBuffer(JNIEnv *env, jlongArray v); // returned buffer needs to be released.
jlongArray MatchaCGoBufferToJlongArray(JNIEnv *env, CGoBuffer buf); // releases buf

#endif //JAVA_FOREIGN_H