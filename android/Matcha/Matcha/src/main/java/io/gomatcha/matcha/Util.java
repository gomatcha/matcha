package io.gomatcha.matcha;

import android.content.Context;
import android.util.Log;
import android.util.TypedValue;

public class Util {
    public static int dipToPixels(Context context, float dip) {
        return (int)TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, dip, context.getResources().getDisplayMetrics());
    }

    static void log(String msg) {
        Log.v("Matcha", msg);
    }
}
