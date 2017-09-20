package io.gomatcha.customview;

import android.content.Context;
import android.widget.Switch;

import io.gomatcha.matcha.MatchaChildView;
import io.gomatcha.matcha.MatchaView;
import io.gomatcha.matcha.MatchaViewNode;

public class CustomView extends MatchaChildView {

    static {
        MatchaView.registerView("github.com/overcyn/customview", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new CustomView(context, node);
            }
        });
    }

    public CustomView(Context context, MatchaViewNode node) {
        super(context, node);

        Switch v = new Switch(context);
        this.addView(v);
    }
}
