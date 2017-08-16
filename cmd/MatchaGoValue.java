package matcha;

// public class MatchaJavaBridge {
//         static {
//                 System.loadLibrary("matcha");
//         }
// }

public class MatchaGoValue {
   static {
      System.loadLibrary("gojni");
   }
   
   private long goRef;
   
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
   
   private static native long matchaGoBool(boolean a);
   private static native long matchaGoLong(long a);
   private static native long matchaGoDouble(double a);
   private static native long matchaGoString(String a);
   private static native long matchaGoByteArray(byte[] v);
   
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
   
   private static native boolean matchaGoToBool(long a);
   private static native long matchaGoToLong(long a);
   private static native double matchaGoToDouble(long a);
   private static native String matchaGoToString(long a);
   private static native byte[] matchaGoToByteArray(long a);
   
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