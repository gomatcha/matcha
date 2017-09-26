package io.gomatcha.bridge;

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
   public static GoValue withFunc(String v) {
      return new GoValue(matchaGoFunc(v), false);
   }
   public static GoValue withType(String v) {
      return new GoValue(matchaGoType(v), false);
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
   private static native long matchaGoLong(long a);
   private static native long matchaGoDouble(double a);
   private static native long matchaGoString(String a);
   private static native long matchaGoByteArray(byte[] v);
   private static native long matchaGoArray(long[] v);
   private static native long matchaGoFunc(String a);
   private static native long matchaGoType(String a);
   
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
   
   private static native long matchaGoToForeign(long a);
   private static native boolean matchaGoToBool(long a);
   private static native long matchaGoToLong(long a);
   private static native double matchaGoToDouble(long a);
   private static native String matchaGoToString(long a);
   private static native byte[] matchaGoToByteArray(long a);
   private static native long[] matchaGoToArray(long a);
   
   public GoValue elem() {
      return new GoValue(matchaGoElem(this.goRef), false);
   }
   
   public boolean isNil() {
      return matchaGoIsNil(this.goRef);
   }
   
   public boolean equals(GoValue v) {
      return matchaGoEqual(this.goRef, v.goRef);
   }
   
   public GoValue[] call(String v, GoValue...v2) {
      if (v2 == null) {
         v2 = new GoValue[0];
      }
      GoValue x = new GoValue(v2);
      long goRef = matchaGoCall(this.goRef, v, x.goRef);
      return new GoValue(goRef, false).toArray();
   }
   
   public GoValue field(String v) {
      return new GoValue(matchaGoField(this.goRef, v), false);
   }
   
   public void setField(String a, GoValue v) {
      matchaGoFieldSet(this.goRef, a, v.goRef);
   }
   
   private static native long matchaGoElem(long a);
   private static native boolean matchaGoIsNil(long a);
   private static native boolean matchaGoEqual(long a, long b);
   private static native long matchaGoCall(long a, String b, long c);
   private static native long matchaGoField(long a, String b);
   private static native void matchaGoFieldSet(long a, String b, long c);

   private static native void matchaGoUntrack(long a);
   
   protected void finalize() throws Throwable {
      matchaGoUntrack(this.goRef);
   }
}