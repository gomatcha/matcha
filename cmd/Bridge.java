package matcha;

import android.util.Log;

public class Bridge {
    private static final Bridge instance = new Bridge();
    private Bridge() {
    }
    public static Bridge singleton() {
        return instance;
    }
    public int test() {
        Log.v("Bridge", "test");
        return 42;
    }
}