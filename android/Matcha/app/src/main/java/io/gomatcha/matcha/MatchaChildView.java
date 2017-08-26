package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.widget.RelativeLayout;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.pb.view.PbView;

public class MatchaChildView extends RelativeLayout {
    MatchaViewNode viewNode;
    PbView.BuildNode buildNode;

    public MatchaChildView(Context context, MatchaViewNode node) {
        super(context);
        viewNode = node;
    }

    public void setNode(PbView.BuildNode buildNode) {
        this.buildNode = buildNode;
    }
}
