package io.gomatcha.matcha;

import android.content.Context;
import android.view.View;

import java.util.List;

public class MatchaChildView extends MatchaLayout {
    public MatchaViewNode viewNode;

    public MatchaChildView(Context context, MatchaViewNode node) {
        super(context);
        viewNode = node;
    }

    public void setNativeState(byte[] nativeState) {
        //no-op
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
