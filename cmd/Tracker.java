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
}