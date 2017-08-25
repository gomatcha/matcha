package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.view.View;
import android.widget.AbsoluteLayout;

import java.lang.ref.WeakReference;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.pb.view.PbView;

public class MatchaChildView extends AbsoluteLayout {
    MatchaViewNode viewNode;
    PbView.BuildNode buildNode;

    public MatchaChildView(Context context, MatchaViewNode v) {
        super(context);
        viewNode = v;
    }

    void setNode(PbView.BuildNode buildNode) {
        this.buildNode = buildNode;
    }
}
