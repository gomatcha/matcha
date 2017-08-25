package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.widget.RelativeLayout;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.pb.view.PbView;

public class MatchaChildView extends RelativeLayout {
    MatchaViewNode viewNode;
    PbView.BuildNode buildNode;

    public MatchaChildView(Context context, MatchaViewNode v) {
        super(context);
        viewNode = v;
        
        this.setBackgroundColor(Color.GREEN);
    }

    void setNode(PbView.BuildNode buildNode) {
        this.buildNode = buildNode;
    }
}
