// +build matcha,android

#include "go-foreign.h"
#include "java-foreign.h"
#include <jni.h>
#include <stdlib.h>
#include <android/log.h>
#include <stdint.h>
#include <string.h>

JavaVM *sJavaVM;
JNIEnv *sEnv;
jint sJavaVersion;
jobject sTracker;

#define printf(...) __android_log_print(ANDROID_LOG_DEBUG, "TAG", __VA_ARGS__);

FgnRef MatchaForeignBool(bool v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignBool", "(Z)J");
    return (*sEnv)->CallLongMethod(sEnv, sTracker, mid, v);
}

bool MatchaForeignToBool(FgnRef v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignToBool", "(J)Z");
    return (*sEnv)->CallBooleanMethod(sEnv, sTracker, mid, v);
}

FgnRef MatchaForeignInt64(int64_t v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignInt64", "(J)J");
    return (*sEnv)->CallLongMethod(sEnv, sTracker, mid, v);
}

int64_t MatchaForeignToInt64(FgnRef v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignToInt64", "(J)J");
    return (*sEnv)->CallLongMethod(sEnv, sTracker, mid, v);
}

FgnRef MatchaForeignFloat64(double v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignFloat64", "(D)J");
    return (*sEnv)->CallLongMethod(sEnv, sTracker, mid, v);
}

double MatchaForeignToFloat64(FgnRef v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignToFloat64", "(J)D");
    return (*sEnv)->CallDoubleMethod(sEnv, sTracker, mid, v);
}

FgnRef MatchaForeignGoRef(GoRef v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignGoRef", "(J)J");
    return (*sEnv)->CallLongMethod(sEnv, sTracker, mid, v);
}

GoRef MatchaForeignToGoRef(FgnRef v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignToGoRef", "(J)J");
    return (*sEnv)->CallLongMethod(sEnv, sTracker, mid, v);
}

FgnRef MatchaForeignString(CGoBuffer buf) {
    jstring jstrBuf = MatchaCGoBufferToString(sEnv, buf);
    
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignString", "(Ljava/lang/String;)J");
    long a = (*sEnv)->CallLongMethod(sEnv, sTracker, mid, jstrBuf);
    
    (*sEnv)->DeleteLocalRef(sEnv, jstrBuf);
    return a;
}

CGoBuffer MatchaForeignToString(FgnRef v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignToString", "(J)Ljava/lang/String;");
    jstring str = (jstring)(*sEnv)->CallObjectMethod(sEnv, sTracker, mid, v);
    
    return MatchaStringToCGoBuffer(sEnv, str);
}

FgnRef MatchaForeignBytes(CGoBuffer bytes) {
    jbyteArray array = MatchaCGoBufferToByteArray(sEnv, bytes);
    
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignBytes", "([B)J");
    long a = (*sEnv)->CallLongMethod(sEnv, sTracker, mid, array);
    
    (*sEnv)->DeleteLocalRef(sEnv, array);
    return a;
}

CGoBuffer MatchaForeignToBytes(FgnRef v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignToBytes", "(J)[B");
    jbyteArray str = (jbyteArray)(*sEnv)->CallObjectMethod(sEnv, sTracker, mid, v);
    
    return MatchaByteArrayToCGoBuffer(sEnv, str);
}

FgnRef MatchaForeignArray(CGoBuffer buf) {
    jbyteArray array = MatchaCGoBufferToJlongArray(sEnv, buf);
    
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignArray", "([J)J");
    long a = (*sEnv)->CallLongMethod(sEnv, sTracker, mid, array);
    
    (*sEnv)->DeleteLocalRef(sEnv, array);
    return a;
}

CGoBuffer MatchaForeignToArray(FgnRef v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignToArray", "(J)[J");
    jlongArray str = (jbyteArray)(*sEnv)->CallObjectMethod(sEnv, sTracker, mid, v);
    
    return MatchaJlongArrayToCGoBuffer(sEnv, str);
}

FgnRef MatchaForeignBridge(CGoBuffer str) {
    jstring *string = MatchaCGoBufferToString(sEnv, str);
    
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignBridge", "(Ljava/lang/String;)J");
    return (*sEnv)->CallLongMethod(sEnv, sTracker, mid, string);
}

// Call

FgnRef MatchaForeignCall(FgnRef v, CGoBuffer str, CGoBuffer args) {
    jstring method = MatchaCGoBufferToString(sEnv, str);
    jbyteArray array = MatchaCGoBufferToJlongArray(sEnv, args);
    
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignCall", "(JLjava/lang/String;[J)J");
    return (*sEnv)->CallLongMethod(sEnv, sTracker, mid, v, method, array);
}

// Tracker

FgnRef MatchaForeignTrack(jobject v) {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "track", "(Ljava/lang/Object;)J");
    return (*sEnv)->CallLongMethod(sEnv, sTracker, mid, v);
}

void MatchaForeignUntrack(FgnRef key) {
    JNIEnv *env = NULL;
    jint success = (*sJavaVM)->GetEnv(sJavaVM, (void **)&env, sJavaVersion);
    if (success == JNI_EDETACHED) {
        (*sJavaVM)->AttachCurrentThread(sJavaVM, &env, NULL);
    }

    jclass cls = (*env)->GetObjectClass(env, sTracker);
    jmethodID mid = (*env)->GetMethodID(env, cls, "untrack", "(J)V");
    (*env)->CallVoidMethod(env, sTracker, mid, key);
    (*env)->DeleteLocalRef(env, cls);
}

// Other

void MatchaForeignPanic() {
    jclass cls = (*sEnv)->GetObjectClass(sEnv, sTracker);
    jmethodID mid = (*sEnv)->GetMethodID(sEnv, cls, "foreignPanic", "()V");
    return (*sEnv)->CallVoidMethod(sEnv, sTracker, mid);
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
    
    jstring jstrBuf = (*sEnv)->NewStringUTF(sEnv, str);
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
