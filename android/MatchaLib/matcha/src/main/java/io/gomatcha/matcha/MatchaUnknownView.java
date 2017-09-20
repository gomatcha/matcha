package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;

class MatchaUnknownView extends MatchaChildView {
    public MatchaUnknownView(Context context, MatchaViewNode node) {
        super(context, node);
    }

    @Override
    public void setBackgroundColor(int v) {
        super.setBackgroundColor(Color.RED);
    }
}
