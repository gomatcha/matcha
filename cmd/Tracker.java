package matcha;

import java.util.Map;
import java.util.HashMap;

public class Tracker {
    private static final Tracker instance = new Tracker();
    private Map<Long, Object> mapTable = new HashMap<Long, Object>();
    private long maxKey = 0;
    private Tracker() {
    }
    public static Tracker singleton() {
        return instance;
    }
    public synchronized long track(Object v) {
        if (v == null) {
            return 0;
        }
        this.maxKey += 1;
        this.mapTable.put(this.maxKey, v);
        return this.maxKey;
    }
    public synchronized void untrack(long v) {
        if (v == 0) {
            return;
        }
        if (this.mapTable.remove(v) == null) {
            throw new IllegalArgumentException("Tracker doesn't contain key");
        }
    }
    public synchronized Object get(long v) {
        if (v == 0) {
            return null;
        }
        
        Object a = this.mapTable.get(v);
        if (a == null) {
            throw new IllegalArgumentException("Tracker doesn't contain key");
        }
        return a;
    }
    
    public synchronized Object call(long v, String method, long args) {
        return null;
    }
    
    public synchronized long foreignBool(boolean v) {
        return track(v);
    }
    public synchronized boolean foreignToBool(long v) {
        boolean a = (Boolean)this.get(v);
        return a;
    }
    public synchronized long foreignInt64(long v) {
        return track(v);
    }
    public synchronized long foreignToInt64(long v) {
        long a = (Long)this.get(v);
        return a;
    }
    public synchronized long foreignFloat64(double v) {
        return track(v);
    }
    public synchronized double foreignToFloat64(long v) {
        double a = (Double)this.get(v);
        return a;
    }
    public synchronized long foreignGoRef(long v) {
        return track(new MatchaGoValue(v, false));
    }
    public synchronized long foreignToGoRef(long v) {
        return ((MatchaGoValue)this.get(v)).goRef;
    }
    public synchronized long foreignString(String v) {
        return track(v);
    }
    public synchronized String foreignToString(long v) {
        return (String)this.get(v);
    }
    public synchronized long foreignBytes(byte[] v) {
        return track(v);
    }
    public synchronized byte[] foreignToBytes(long v) {
        return (byte[])this.get(v);
    }
}