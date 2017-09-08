package io.gomatcha.matcha;

import android.content.Context;
import android.support.v7.widget.Toolbar;
import android.util.Log;
import android.widget.CompoundButton;
import android.widget.RelativeLayout;
import android.widget.Switch;
import android.support.v7.app.ActionBar;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.stackview.PbStackView;
import io.gomatcha.matcha.pb.view.switchview.PbSwitchView;

public class MatchaStackView extends MatchaChildView {
    Toolbar toolbar;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/stacknav", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaStackView(context, node);
            }
        });
    }

    public MatchaStackView(Context context, MatchaViewNode node) {
        super(context, node);

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
        toolbar = new Toolbar(context);
        toolbar.setTitle("TEST");
        addView(toolbar, params);
        Log.v("x", "what");
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbStackView.View proto = buildNode.getBridgeValue().unpack(PbStackView.View.class);
        } catch (InvalidProtocolBufferException e) {
        }
    }
}
