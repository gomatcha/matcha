package io.gomatcha.matcha;

import android.content.Context;
import android.view.View;
import android.view.ViewGroup;
import android.widget.RelativeLayout;

class MatchaToolbarWrapper extends RelativeLayout {
    public MatchaToolbarWrapper(Context context) {
        super(context);
    }

    protected MatchaToolbarView toolbarView;
    void setToolbarView(MatchaToolbarView v) {
        if (v == toolbarView) {
            return;
        }
        if (toolbarView != null) {
            this.removeView(toolbarView);
        }
        toolbarView = v;
        this.addView(toolbarView, new RelativeLayout.LayoutParams(MatchaLayout.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT));
    }
    protected View contentView;
    void setContentView(View v) {
        if (v == contentView) {
            return;
        }
        if (contentView != null) {
            this.removeView(contentView);
        }
        contentView = v;
        RelativeLayout.LayoutParams childViewParams = new RelativeLayout.LayoutParams(MatchaLayout.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.MATCH_PARENT);
        childViewParams.addRule(RelativeLayout.BELOW, toolbarView.getId());
        this.addView(contentView, childViewParams);
    }
}