package io.gomatcha.bridge;

import java.util.Map;
import java.util.HashMap;
import android.util.Log;

public class Bridge {
    private Map<String, Object> mapTable = new HashMap<String, Object>();
    private static final Bridge instance = new Bridge();
    private Bridge() {
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
}