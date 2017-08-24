package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.view.View;

import java.lang.ref.WeakReference;
import java.util.ArrayList;

import io.gomatcha.bridge.GoValue;

public class MatchaView extends View {
    static ArrayList<WeakReference<MatchaView>> views = new ArrayList<WeakReference<MatchaView>>();
    GoValue goValue;
    long identifier;

    public MatchaView(Context context, GoValue v) {
        super(context);
        goValue = v;
        // identifier = v.call("Id", null)[0].toLong();
        setBackgroundColor(Color.RED);

        views.add(new WeakReference<MatchaView>(this));
    }

    @Override
    protected void finalize() {
        // goValue.call("Stop", null);
    }
}
