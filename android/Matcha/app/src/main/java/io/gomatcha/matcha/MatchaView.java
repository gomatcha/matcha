package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.util.Log;
import android.view.View;

import java.lang.ref.WeakReference;
import java.util.ArrayList;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.JavaBridge;
import io.gomatcha.matcha.pb.view.PbView;

public class MatchaView extends View {
    static ArrayList<WeakReference<MatchaView>> views = new ArrayList<WeakReference<MatchaView>>();
    GoValue goValue;
    long identifier;
    MatchaViewNode node;

    static {
        new JavaBridge();
    }

    public MatchaView(Context context, GoValue v) {
        super(context);
        goValue = v;
        identifier = v.call("Id")[0].toLong();
        long viewid = v.call("ViewId")[0].toLong();
        node = new MatchaViewNode(null, this, viewid);
        setBackgroundColor(Color.RED);

        views.add(new WeakReference<MatchaView>(this));
    }

    void update(PbView.Root root) {
        node.setRoot(root);
    }

    @Override
    protected void onSizeChanged(int w, int h, int oldw, int oldh) {
        Log.v("", "onSizeChange"+ w  + "," + h);
        goValue.call("SetSize", new GoValue((double)w), new GoValue((double)h));

        GoValue.withFunc("gomatcha.io/matcha/animate screenUpdate").call("");
    }

    @Override
    protected void finalize() {
        goValue.call("Stop");
    }
}
