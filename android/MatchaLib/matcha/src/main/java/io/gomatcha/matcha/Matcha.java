package io.gomatcha.matcha;

import android.content.Context;

public class Matcha {
    public static synchronized void configure(Context ctx) {
        JavaBridge.configure(ctx);
    }
}
