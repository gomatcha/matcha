// +build matcha,android

#include "go-foreign.h"
#include "java-foreign.h"
#include <jni.h>
#include <stdlib.h>
#include <android/log.h>
#include <stdint.h>
#include <string.h>

JavaVM *g_JavaVM;
jint g_JavaVersion;
jobject g_Tracker;
jclass g_TrackerClass;
JNIEnv *g_Env;
JNIEnv *g_UntrackEnv; // Untrack happens in a separate goroutine

jmethodID g_foreignNil;
jmethodID g_foreignIsNil;
jmethodID g_foreignBool;
jmethodID g_foreignToBool;
jmethodID g_foreignInt64;
jmethodID g_foreignToInt64;
jmethodID g_foreignFloat64;
jmethodID g_foreignToFloat64;
jmethodID g_foreignGoRef;
jmethodID g_foreignToGoRef;
jmethodID g_foreignString;
jmethodID g_foreignToString;
jmethodID g_foreignBytes;
jmethodID g_foreignToBytes;
jmethodID g_foreignArray;
jmethodID g_foreignToArray;
jmethodID g_foreignBridge;
jmethodID g_foreignCall;
jmethodID g_foreignUntrack;
jmethodID g_foreignTrackerCount;
jmethodID g_foreignPanic;

#define printf(...) __android_log_print(ANDROID_LOG_DEBUG, "TAG", __VA_ARGS__);

void MatchaInit(JNIEnv *env, jobject tracker) {
    (*env)->GetJavaVM(env, &g_JavaVM);
    g_JavaVersion = (*env)->GetVersion(env);
    g_Env = env;
    g_Tracker = (*env)->NewGlobalRef(env, tracker);
    
    jclass trackerClass = (*env)->GetObjectClass(env, tracker);
    g_TrackerClass = (*env)->NewGlobalRef(env, trackerClass);
    (*env)->DeleteLocalRef(env, trackerClass);
    
    g_foreignNil = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignNil", "()J");
    g_foreignIsNil = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignIsNil", "(J)Z");
    g_foreignBool = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignBool", "(Z)J");
    g_foreignToBool = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignToBool", "(J)Z");
    g_foreignInt64 = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignInt64", "(J)J");
    g_foreignToInt64 = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignToInt64", "(J)J");
    g_foreignFloat64 = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignFloat64", "(D)J");
    g_foreignToFloat64 = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignToFloat64", "(J)D");
    g_foreignGoRef = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignGoRef", "(J)J");
    g_foreignToGoRef = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignToGoRef", "(J)J");
    g_foreignString = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignString", "(Ljava/lang/String;)J");
    g_foreignToString = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignToString", "(J)Ljava/lang/String;");
    g_foreignBytes = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignBytes", "([B)J");
    g_foreignToBytes = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignToBytes", "(J)[B");
    g_foreignArray = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignArray", "([J)J");
    g_foreignToArray = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignToArray", "(J)[J");
    g_foreignBridge = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignBridge", "(Ljava/lang/String;)J");
    g_foreignCall = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignCall", "(JLjava/lang/String;[J)J");
    g_foreignUntrack = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "untrack", "(J)V");
    g_foreignTrackerCount = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "trackerCount", "()J");
    g_foreignPanic = (*g_Env)->GetMethodID(g_Env, g_TrackerClass, "foreignPanic", "()V");
}

// Foreign functions

FgnRef MatchaForeignNil() {
    return (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignNil);
}

bool MatchaForeignIsNil(FgnRef v) {
    return (*g_Env)->CallBooleanMethod(g_Env, g_Tracker, g_foreignIsNil, v);
}

FgnRef MatchaForeignBool(bool v) {
    return (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignBool, v);
}

bool MatchaForeignToBool(FgnRef v) {
    return (*g_Env)->CallBooleanMethod(g_Env, g_Tracker, g_foreignToBool, v);
}

FgnRef MatchaForeignInt64(int64_t v) {
    return (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignInt64, v);
}

int64_t MatchaForeignToInt64(FgnRef v) {
    return (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignToInt64, v);
}

FgnRef MatchaForeignFloat64(double v) {
    return (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignFloat64, v);
}

double MatchaForeignToFloat64(FgnRef v) {
    return (*g_Env)->CallDoubleMethod(g_Env, g_Tracker, g_foreignToFloat64, v);
}

FgnRef MatchaForeignGoRef(GoRef v) {
    return (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignGoRef, v);
}

GoRef MatchaForeignToGoRef(FgnRef v) {
    return (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignToGoRef, v);
}

FgnRef MatchaForeignString(CGoBuffer buf) {
    jstring jstrBuf = MatchaCGoBufferToString(g_Env, buf);
    long a = (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignString, jstrBuf);
    (*g_Env)->DeleteLocalRef(g_Env, jstrBuf);
    return a;
}

CGoBuffer MatchaForeignToString(FgnRef v) {
    jstring str = (jstring)(*g_Env)->CallObjectMethod(g_Env, g_Tracker, g_foreignToString, v);
    CGoBuffer a = MatchaStringToCGoBuffer(g_Env, str);
    (*g_Env)->DeleteLocalRef(g_Env, str);
    return a;
}

FgnRef MatchaForeignBytes(CGoBuffer bytes) {
    jbyteArray array = MatchaCGoBufferToByteArray(g_Env, bytes);
    long a = (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignBytes, array);
    (*g_Env)->DeleteLocalRef(g_Env, array);
    return a;
}

CGoBuffer MatchaForeignToBytes(FgnRef v) {
    jbyteArray str = (jbyteArray)(*g_Env)->CallObjectMethod(g_Env, g_Tracker, g_foreignToBytes, v);
    CGoBuffer a = MatchaByteArrayToCGoBuffer(g_Env, str);
    (*g_Env)->DeleteLocalRef(g_Env, str);
    return a;
}

FgnRef MatchaForeignArray(CGoBuffer buf) {
    jbyteArray array = MatchaCGoBufferToJlongArray(g_Env, buf);
    long a = (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignArray, array);
    (*g_Env)->DeleteLocalRef(g_Env, array);
    return a;
}

CGoBuffer MatchaForeignToArray(FgnRef v) {
    jlongArray str = (jbyteArray)(*g_Env)->CallObjectMethod(g_Env, g_Tracker, g_foreignToArray, v);
    CGoBuffer a = MatchaJlongArrayToCGoBuffer(g_Env, str);
    (*g_Env)->DeleteLocalRef(g_Env, str);
    return a;
}

FgnRef MatchaForeignBridge(CGoBuffer str) {
    jstring string = MatchaCGoBufferToString(g_Env, str);
    long a = (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignBridge, string);
    (*g_Env)->DeleteLocalRef(g_Env, string);
    return a;
}

FgnRef MatchaForeignCall(FgnRef v, CGoBuffer str, CGoBuffer args) {
    jstring method = MatchaCGoBufferToString(g_Env, str);
    jbyteArray array = MatchaCGoBufferToJlongArray(g_Env, args);
    long a = (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignCall, v, method, array);
    (*g_Env)->DeleteLocalRef(g_Env, method);
    (*g_Env)->DeleteLocalRef(g_Env, array);
    return a;
}

void MatchaForeignUntrack(FgnRef key) {
    // Garbage collection happens on a unknown background thread. We must get the correct JNIEnv.
    if (g_UntrackEnv == NULL) {
        jint success = (*g_JavaVM)->GetEnv(g_JavaVM, (void **)&g_UntrackEnv, g_JavaVersion);
        if (success == JNI_EDETACHED) {
            (*g_JavaVM)->AttachCurrentThread(g_JavaVM, &g_UntrackEnv, NULL);
        }
    }
    
    (*g_UntrackEnv)->CallVoidMethod(g_UntrackEnv, g_Tracker, g_foreignUntrack, key);
}

int64_t MatchaForeignTrackerCount() {
    return (*g_Env)->CallLongMethod(g_Env, g_Tracker, g_foreignTrackerCount);
}

void MatchaForeignPanic() {
    return (*g_Env)->CallVoidMethod(g_Env, g_Tracker, g_foreignPanic);
}

// Utilities

CGoBuffer MatchaStringToCGoBuffer(JNIEnv *env, jstring v) {
    const char *nativeString = (*env)->GetStringUTFChars(env, v, 0);
    
    int len = strlen(nativeString);
    char *buf = (char *)malloc(len);
    strncpy(buf, nativeString, len);
    
    (*env)->ReleaseStringUTFChars(env, v, nativeString);
   
    CGoBuffer cstr;
    cstr.ptr = buf;
    cstr.len = len;
    return cstr;
}

jstring MatchaCGoBufferToString(JNIEnv *env, CGoBuffer buf) {
    char *str = malloc(buf.len+1);
    strncpy(str, buf.ptr, buf.len);
    str[buf.len] = '\0';
    
    jstring jstrBuf = (*g_Env)->NewStringUTF(g_Env, str);
    free(buf.ptr);
    free(str);
    return jstrBuf;
}

CGoBuffer MatchaByteArrayToCGoBuffer(JNIEnv *env, jbyteArray v) {
    int len = (*env)->GetArrayLength(env, v);
    if (len == 0) {
        return (CGoBuffer){0};
    }
    char *buf = (char *)malloc(len);
    (*env)->GetByteArrayRegion(env, v, 0, len, (jbyte *)buf);
  
    CGoBuffer cstr;
    cstr.ptr = buf;
    cstr.len = len;
    return cstr;
}

jbyteArray MatchaCGoBufferToByteArray(JNIEnv *env, CGoBuffer buf) {
    jbyteArray array = (*env)->NewByteArray(env, buf.len);
    (*env)->SetByteArrayRegion(env, array, 0, buf.len, buf.ptr);
    free(buf.ptr);
    return array;
}

jlongArray MatchaCGoBufferToJlongArray(JNIEnv *env, CGoBuffer buf) {
    int len = buf.len/8;
    jlongArray array = (*env)->NewLongArray(env, len);
    jlong *arr = (*env)->GetLongArrayElements(env, array, NULL);
    char *data = buf.ptr;
    for (int i = 0; i < len; i++) {
        GoRef ref = 0;
        memcpy(&ref, data, 8);
        arr[i] = ref;
        data += 8;
    }
    
    (*env)->ReleaseLongArrayElements(env, array, arr, 0);
    free(buf.ptr);
    return array;
}

CGoBuffer MatchaJlongArrayToCGoBuffer(JNIEnv *env, jlongArray v) {
    int len = (*env)->GetArrayLength(env, v);
    if (len == 0) {
        return (CGoBuffer){0};
    }
    
    char *buf = (char *)malloc(len * 8);
    char *data = buf;
    jlong *arr = (*env)->GetLongArrayElements(env, v, 0);
    for (int i = 0; i < len; i++) {
        int64_t ref = arr[i];
        memcpy(data+i*8, &ref, 8);
    }
    (*env)->ReleaseLongArrayElements(env, v, arr, 0);
    
    CGoBuffer cstr;
    cstr.ptr = buf;
    cstr.len = len * 8;
    return cstr;
}
