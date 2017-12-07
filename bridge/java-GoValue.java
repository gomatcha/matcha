package io.gomatcha.bridge;

import android.util.Log;
import java.util.Arrays;

public class GoValue {
   static {
      System.loadLibrary("gojni");
      
      matchaInit(Tracker.singleton());
   }
   
   private static native void matchaInit(Object tracker);
   
   protected long goRef;
   
   protected GoValue(long goref, boolean empty) {
      this.goRef = goref;
   }
   public GoValue(Object v) {
      this(matchaGoForeign(Tracker.singleton().track(v)), false);
   }
   public GoValue(boolean v) {
      this(matchaGoBool(v), false);
   }
   public GoValue(int v) {
      this(matchaGoInt(v), false);
   }
   public GoValue(long v) {
      this(matchaGoLong(v), false);
   }
   public GoValue(double v) {
      this(matchaGoDouble(v), false);
   }
   public GoValue(String v) {
      this(matchaGoString(v), false);
   }
   public GoValue(byte[] v) {
      this(matchaGoByteArray(v), false);
   }
   public GoValue(GoValue[] v) {
      this(makeGoArray(v), false);
   }
   public static GoValue WithBoolean(boolean v) {
      return new GoValue(matchaGoBool(v), false);
   }
   public static GoValue WithInt(int v) {
      return new GoValue(matchaGoInt(v), false);
   }
   public static GoValue WithLong(long v) {
      return new GoValue(matchaGoLong(v), false);
   }
   public static GoValue WithDouble(double v) {
      return new GoValue(matchaGoDouble(v), false);
   }
   public static GoValue WithString(String v) {
      return new GoValue(matchaGoString(v), false);
   }
   public static GoValue WithByteArray(byte[] v) {
      return new GoValue(matchaGoByteArray(v), false);
   }
   public static GoValue WithArray(GoValue[] v) {
      return new GoValue(makeGoArray(v), false);
   }
   public static GoValue WithObject(Object v) {
      return new GoValue(matchaGoForeign(Tracker.singleton().track(v)), false);
   }
   public static GoValue withFunc(String v) {
      return new GoValue(matchaGoFunc(v), false);
   }
   private static long makeGoArray(GoValue[] v) {
      long[] array = new long[v.length];
      for (int i = 0; i < v.length; i++) {
         array[i] = v[i].goRef;
      }
      return matchaGoArray(array);
   }
   
   private static native long matchaGoForeign(long a);
   private static native long matchaGoBool(boolean a);
   private static native long matchaGoInt(int a);
   private static native long matchaGoLong(long a);
   private static native long matchaGoDouble(double a);
   private static native long matchaGoString(String a);
   private static native long matchaGoByteArray(byte[] v);
   private static native long matchaGoArray(long[] v);
   private static native long matchaGoFunc(String a);
   
   public boolean isNil() {
      return matchaGoIsNil(this.goRef);
   }
   
   public Object toObject() {
      long foreignRef = matchaGoToForeign(this.goRef);
      return Tracker.singleton().get(foreignRef);
   }
   public boolean toBool() {
      return matchaGoToBool(this.goRef);
   }
   public long toLong() {
      return matchaGoToLong(this.goRef);
   }
   public double toDouble() {
      return matchaGoToDouble(this.goRef);
   }
   public String toString() {
      return matchaGoToString(this.goRef);
   }
   public byte[] toByteArray() {
      return matchaGoToByteArray(this.goRef);
   }
   public GoValue[] toArray() {
      long[] array = matchaGoToArray(this.goRef);
      
      GoValue[] array2 = new GoValue[array.length];
      for (int i = 0; i < array.length; i++) {
         array2[i] = new GoValue(array[i], false);
      }
      return array2;
   }
   
   private static native boolean matchaGoIsNil(long a);
   private static native long matchaGoToForeign(long a);
   private static native boolean matchaGoToBool(long a);
   private static native long matchaGoToLong(long a);
   private static native double matchaGoToDouble(long a);
   private static native String matchaGoToString(long a);
   private static native byte[] matchaGoToByteArray(long a);
   private static native long[] matchaGoToArray(long a);
   
   public GoValue[] call(String v, GoValue...v2) {
      if (v2 == null) {
         v2 = new GoValue[0];
      }
      long[] args = new long[v2.length];
      for (int i = 0; i < v2.length; i++) {
         args[i] = v2[i].goRef;
      }
      
      long[] refs = matchaGoCall(this.goRef, v, args);
      
      GoValue[] array2 = new GoValue[refs.length];
      for (int i = 0; i < refs.length; i++) {
         array2[i] = new GoValue(refs[i], false);
      }
      return array2;
   }
   
   public static void testFunc() {
      matchaTestFunc();
   }
   
   protected void finalize() throws Throwable {
      matchaGoUntrack(this.goRef);
      super.finalize();
   }
   
   private static native long[] matchaGoCall(long a, String b, long[] c);
   private static native void matchaTestFunc();
   private static native void matchaGoUntrack(long a);
}