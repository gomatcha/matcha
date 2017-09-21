package io.gomatcha.matcha;

import android.content.Context;
import android.view.View;

import java.util.List;

import io.gomatcha.matcha.proto.view.PbView;

public class MatchaChildView extends MatchaLayout {
    public MatchaViewNode viewNode;
    public PbView.BuildNode buildNode;
    MatchaGestureRecognizer matchaGestureRecognizer;

    public MatchaChildView(Context context, MatchaViewNode node) {
        super(context);
        final Context ctx = context;
        viewNode = node;
        this.setClipChildren(false);
        this.setClipToPadding(false);
        this.matchaGestureRecognizer = new MatchaGestureRecognizer();
        this.matchaGestureRecognizer.childView = this;
        this.matchaGestureRecognizer.context = context;
        this.setOnTouchListener(this.matchaGestureRecognizer);
    }

    public void setNode(PbView.BuildNode buildNode) {
        this.buildNode = buildNode;
    }

    public boolean isContainerView() {
        return false;
    }

    public void setChildViews(List<View> childViews) {
        // no-op
    }

    public MatchaLayout getLayout() {
        return this;
    }
}
