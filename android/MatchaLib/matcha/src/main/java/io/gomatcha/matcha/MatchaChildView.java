package io.gomatcha.matcha;

import android.content.Context;
import android.view.View;

import java.util.List;

public class MatchaChildView extends MatchaLayout {
    public MatchaChildView(Context context) {
        super(context);
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
}
