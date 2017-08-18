package matcha;

import matcha.Bridge;

public class MatchaGoValue {
   static {
      System.loadLibrary("gojni");
      
      matchaInit(Tracker.singleton(), Bridge.singleton());
   }
   
   private static native void matchaInit(Object tracker, Object bridge);
   
   public long goRef;
   
   private MatchaGoValue(long goref, boolean empty) {
      this.goRef = goref;
   }
   public MatchaGoValue(boolean v) {
      this(matchaGoBool(v), false);
   }
   public MatchaGoValue(long v) {
      this(matchaGoLong(v), false);
   }
   public MatchaGoValue(double v) {
      this(matchaGoDouble(v), false);
   }
   public MatchaGoValue(String v) {
      this(matchaGoString(v), false);
   }
   public MatchaGoValue(byte[] v) {
      this(matchaGoByteArray(v), false);
   }
   public MatchaGoValue(MatchaGoValue[] v) {
      this(makeGoArray(v), false);
   }
   public static MatchaGoValue withFunc(String v) {
      return new MatchaGoValue(matchaGoFunc(v), false);
   }
   public static MatchaGoValue withType(String v) {
      return new MatchaGoValue(matchaGoType(v), false);
   }
   private static long makeGoArray(MatchaGoValue[] v) {
      long[] array = new long[v.length];
      for (int i = 0; i < v.length; i++) {
         array[i] = v[i].goRef;
      }
      return matchaGoArray(array);
   }
   
   private static native long matchaGoBool(boolean a);
   private static native long matchaGoLong(long a);
   private static native long matchaGoDouble(double a);
   private static native long matchaGoString(String a);
   private static native long matchaGoByteArray(byte[] v);
   private static native long matchaGoArray(long[] v);
   private static native long matchaGoFunc(String a);
   private static native long matchaGoType(String a);
   
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
   public MatchaGoValue[] toArray() {
      long[] array = matchaGoToArray(this.goRef);
      
      MatchaGoValue[] array2 = new MatchaGoValue[array.length];
      for (int i = 0; i < array.length; i++) {
         array2[i] = new MatchaGoValue(array[i], false);
      }
      return array2;
   }
   
   private static native boolean matchaGoToBool(long a);
   private static native long matchaGoToLong(long a);
   private static native double matchaGoToDouble(long a);
   private static native String matchaGoToString(long a);
   private static native byte[] matchaGoToByteArray(long a);
   private static native long[] matchaGoToArray(long a);
   
   public MatchaGoValue elem() {
      return new MatchaGoValue(matchaGoElem(this.goRef), false);
   }
   
   public boolean isNil() {
      return matchaGoIsNil(this.goRef);
   }
   
   public boolean equals(MatchaGoValue v) {
      return matchaGoEqual(this.goRef, v.goRef);
   }
   
   public MatchaGoValue[] call(String v, MatchaGoValue[] v2) {
      if (v2 == null) {
         v2 = new MatchaGoValue[0];
      }
      MatchaGoValue x = new MatchaGoValue(v2);
      long goRef = matchaGoCall(this.goRef, v, x.goRef);
      return new MatchaGoValue(goRef, false).toArray();
   }
   
   public MatchaGoValue field(String v) {
      return new MatchaGoValue(matchaGoField(this.goRef, v), false);
   }
   
   public void setField(String a, MatchaGoValue v) {
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
// - (id)initWithGoRef:(GoRef)ref;
// - (id)initWithBool:(BOOL)v;
// - (id)initWithInt:(int)v;
// - (id)initWithLongLong:(long long)v;
// - (id)initWithUnsignedLongLong:(unsigned long long)v;
// - (id)initWithDouble:(double)v;
// - (id)initWithString:(NSString *)v;
// - (id)initWithData:(NSData *)v;
// - (id)initWithArray:(NSArray<MatchaGoValue *> *)v;
// - (id)initWithType:(NSString *)typeName;
// - (id)initWithFunc:(NSString *)funcName;
// @property (nonatomic, readonly) GoRef ref;
// - (BOOL)toBool;
// - (long long)toLongLong;
// - (unsigned long long)toUnsignedLongLong;
// - (double)toDouble;
// - (NSString *)toString;
// - (NSData *)toData;
// - (NSArray *)toArray;
// - (NSMapTable *)toMapTable;
// // - (NSDictionary *)toDictionary;
// - (BOOL)isNil;
// - (BOOL)isEqual:(MatchaGoValue *)value;
// - (MatchaGoValue *)elem;
// - (NSArray<MatchaGoValue *> *)call:(NSString *)method args:(NSArray<MatchaGoValue *> *)args; // pass in nil for the method o call a closure.
// - (MatchaGoValue *)field:(NSString *)name;
// - (void)setField:(NSString *)name value:(MatchaGoValue *)value;
// - (MatchaGoValue *)objectForKeyedSubscript:(NSString *)key;
// - (void)setObject:(MatchaGoValue *)object forKeyedSubscript:(NSString *)key;
}