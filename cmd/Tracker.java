package io.gomatcha.bridge;

import java.util.Map;
import java.util.HashMap;
import android.util.Log;
import java.lang.reflect.Method;
import java.lang.reflect.InvocationTargetException;
import java.util.Arrays;

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
    public synchronized long foreignBridge(String key) {
        Bridge bridge = Bridge.singleton();
        return track(bridge.get(key));
    }
    public synchronized long foreignCall(long v, String method, long args) {
        Object[] va = (Object[])this.get(args);
        int len = 0;
        if (va != null) {
            len = va.length;
        }
        Object[] vb = new Object[len];
        Class[] vc = new Class[len];
        for (int i = 0; i < len; i++) {
            Object e = va[i];
            vb[i] = e;
            vc[i] = e.getClass();
        }
        
        long test = 0;
        try {
            Object a = this.get(v);
            // Log.v("Bridge", String.format("foreignCall, %s, %s", a.getClass().getName(), Arrays.toString(vc)));
            Method m = a.getClass().getMethod(method, vc);
            Object rlt = m.invoke(a, vb);
            test = track(rlt);
        } catch (NoSuchMethodException e) {
            Log.v("Bridge", String.format("foreignCall, %d, %s, %d, %s", v, method, args, e.getCause()));
            throw new RuntimeException(e);
        } catch (IllegalAccessException e) {
            Log.v("Bridge", String.format("foreignCall, %d, %s, %d, %s", v, method, args, e.getCause()));
            throw new RuntimeException(e);
        } catch (InvocationTargetException e) {
            Log.v("Bridge", String.format("foreignCall, %d, %s, %d, %s", v, method, args, e.getCause()));
            throw new RuntimeException(e);
        }
        return test;
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
        Object a = this.get(v);
        if (a instanceof Integer) {
            return ((Integer)a).longValue();
        }
        return (Long)a;
    }
    public synchronized long foreignFloat64(double v) {
        return track(v);
    }
    public synchronized double foreignToFloat64(long v) {
        Object a = this.get(v);
        if (a instanceof Float) {
            return ((Float)a).doubleValue();
        }
        return (Double)a;
    }
    public synchronized long foreignGoRef(long v) {
        return track(new GoValue(v, false));
    }
    public synchronized long foreignToGoRef(long v) {
        return ((GoValue)this.get(v)).goRef;
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
    public synchronized long foreignArray(int v) {
        Object[] a = new Object[v];
        return track(a);
    }
    public synchronized void foreignArraySet(long v, long val, int idx) {
        Object[] a = (Object[])this.get(v);
        a[idx] = this.get(val);
    }
    public synchronized long foreignArrayAt(long v, int idx) {
        Object[] a = (Object[])this.get(v);
        return track(a[idx]);
    }
    public synchronized long foreignArrayLen(long v) {
        Object[] a = (Object[])this.get(v);
        return a.length;
    }
}
