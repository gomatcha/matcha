package io.gomatcha.matcha;

import android.content.Context;

class MatchaBasicView extends MatchaChildView {
    static {
        MatchaView.registerView("", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaBasicView(context, node);
            }
        });
    }

    public MatchaBasicView(Context context, MatchaViewNode node) {
        super(context, node);
    }
}
