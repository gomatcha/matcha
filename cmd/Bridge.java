package io.gomatcha.bridge;

import java.util.Map;
import java.util.HashMap;
import android.util.Log;

public class Bridge {
    private Map<String, Object> mapTable = new HashMap<String, Object>();
    private static final Bridge instance = new Bridge();
    private Bridge() {
        this.put("a", new testClass());
    }
    public static Bridge singleton() {
        return instance;
    }
    public void put(String j, Object v) {
        this.mapTable.put(j, v);
    }
    public Object get(String j) {
        return this.mapTable.get(j);
    }
    
    public class testClass {
        public int test() {
            Log.v("Bridge", "test");
            return 42;
        }
    }
}