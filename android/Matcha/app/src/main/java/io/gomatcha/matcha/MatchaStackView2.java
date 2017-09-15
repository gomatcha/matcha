package io.gomatcha.matcha;

import android.content.Context;
import android.util.Log;
import android.view.View;
import android.widget.RelativeLayout;

import java.util.ArrayList;
import java.util.List;

public class MatchaStackView2 extends RelativeLayout {
    public MatchaStackView2(Context context) {
        super(context);
    }

    List<View> childViews = new ArrayList<View>();

    List<View> getChildViews() {
        return childViews;
    }

    void setChildViews(List<View> v) {
        if (childViews.size() > 0) {
            View prev = childViews.get(childViews.size()-1);
            this.removeView(prev);
        }

        childViews = v;

        if (childViews.size() > 0) {
            View top = childViews.get(childViews.size()-1);
            this.addView(top);
        }
    }
}
