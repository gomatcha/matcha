// +build matcha,android

#include <jni.h>
#include <stdlib.h>
#include <string.h>
#include "java-go.h"
#include "go-go.h"
#include "java-foreign.h"
#include <stdbool.h>
#include <stdio.h>
#include <android/log.h>

#define printf(...) __android_log_print(ANDROID_LOG_DEBUG, "TAG", __VA_ARGS__);

JNIEXPORT void JNICALL Java_io_gomatcha_bridge_GoValue_matchaInit(JNIEnv *env, jclass c, jobject tracker) {
    MatchaInit(env, tracker);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoForeign(JNIEnv *env, jclass c, jlong v) {
    return matchaGoForeign(v);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoBool(JNIEnv *env, jclass c, jboolean v) {
    return matchaGoBool(v);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoInt(JNIEnv *env, jclass c, jint v) {
    return matchaGoInt(v);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoLong(JNIEnv *env, jclass c, jlong v) {
    return matchaGoInt64(v);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoDouble(JNIEnv *env, jclass c, jdouble v) {
    return matchaGoFloat64(v);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoString(JNIEnv *env, jclass c, jstring v) {
    return matchaGoString(MatchaStringToCGoBuffer(env, v));
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoByteArray(JNIEnv *env, jclass c, jbyteArray v) {
    CGoBuffer cstr = MatchaByteArrayToCGoBuffer(env, v);
    return matchaGoBytes(cstr);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoArray(JNIEnv *env, jclass c, jlongArray v) {
    CGoBuffer buf = MatchaJlongArrayToCGoBuffer(env, v);
    return matchaGoArray(buf);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoFunc(JNIEnv *env, jclass c, jstring v) {
    return matchaGoFunc(MatchaStringToCGoBuffer(env, v));
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoType(JNIEnv *env, jclass c, jstring v) {
    return matchaGoType(MatchaStringToCGoBuffer(env, v));
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoToForeign(JNIEnv *env, jclass c, jlong v) {
    return matchaGoToForeign(v);
}

JNIEXPORT jboolean JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoToBool(JNIEnv *env, jclass c, jlong v) {
    return matchaGoToBool(v);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoToLong(JNIEnv *env, jclass c, jlong v) {
    return matchaGoToInt64(v);
}

JNIEXPORT jdouble JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoToDouble(JNIEnv *env, jclass c, jlong v) {
    return matchaGoToFloat64(v);
}

JNIEXPORT jstring JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoToString(JNIEnv *env, jclass c, jlong v) {
    CGoBuffer buf = matchaGoToString(v);
    return MatchaCGoBufferToString(env, buf);
}

JNIEXPORT jbyteArray JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoToByteArray(JNIEnv *env, jclass c, jlong v) {
    CGoBuffer buf = matchaGoToBytes(v);
    jbyteArray a = MatchaCGoBufferToByteArray(env, buf);
    return a;
}

JNIEXPORT jlongArray JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoToArray(JNIEnv *env, jclass c, jlong v) {
    CGoBuffer buf = matchaGoToArray(v);
    return MatchaCGoBufferToJlongArray(env, buf);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoElem(JNIEnv *env, jclass c, jlong v) {
    return matchaGoElem(v);
}

JNIEXPORT jboolean JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoIsNil(JNIEnv *env, jclass c, jlong v) {
    return matchaGoIsNil(v);
}

JNIEXPORT jboolean JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoEqual(JNIEnv *env, jclass c, jlong v, jlong v2) {
    return matchaGoEqual(v, v2);
}

JNIEXPORT jlongArray JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoCall(JNIEnv *env, jclass c, jlong v, jstring v2, jlongArray v3) {
    CGoBuffer args = MatchaJlongArrayToCGoBuffer(env, v3);
    CGoBuffer rlt = matchaGoCall(v, MatchaStringToCGoBuffer(env, v2), args);
    return MatchaCGoBufferToJlongArray(env, rlt);
}

JNIEXPORT jlong JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoField(JNIEnv *env, jclass c, jlong v, jstring v2) {
    return matchaGoField(v, MatchaStringToCGoBuffer(env, v2));
}

JNIEXPORT void JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoFieldSet(JNIEnv *env, jclass c, jlong v, jstring v2, jlong v3) {
    matchaGoFieldSet(v, MatchaStringToCGoBuffer(env, v2), v3);
}

JNIEXPORT void JNICALL Java_io_gomatcha_bridge_GoValue_matchaTestFunc(JNIEnv *env, jclass c) {
    matchaTestFunc();
}

JNIEXPORT void JNICALL Java_io_gomatcha_bridge_GoValue_matchaGoUntrack(JNIEnv *env, jclass c, jlong v) {
    matchaGoUntrack(v);
}
