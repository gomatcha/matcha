package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.util.DisplayMetrics;
import android.util.Log;
import android.view.View;
import android.widget.RelativeLayout;

import java.lang.ref.WeakReference;
import java.util.ArrayList;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.JavaBridge;
import io.gomatcha.matcha.pb.view.PbView;

public class MatchaView extends RelativeLayout {
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

    boolean loaded = false;

    void update(PbView.Root root) {
        node.setRoot(root);

        if (!loaded) {
            loaded = true;

            RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
            node.view.setLayoutParams(params);

            addView(node.view);
        }
    }

    @Override
    protected void onSizeChanged(int w, int h, int oldw, int oldh) {
        final double width = (double)w / this.getResources().getDisplayMetrics().densityDpi * DisplayMetrics.DENSITY_DEFAULT;
        final double height = (double)h / this.getResources().getDisplayMetrics().densityDpi * DisplayMetrics.DENSITY_DEFAULT;
        Log.v("OnSizeChange", width  + "," + height +  "." + this.getResources().getDisplayMetrics().densityDpi + "." + DisplayMetrics.DENSITY_DEFAULT);

        this.post( new Runnable() {
            @Override
            public void run() {
                goValue.call("SetSize", new GoValue((double)width), new GoValue((double)height));
                GoValue.withFunc("gomatcha.io/matcha/animate screenUpdate").call("");
            }
        });
    }

    @Override
    protected void finalize() {
        goValue.call("Stop");
    }
}
