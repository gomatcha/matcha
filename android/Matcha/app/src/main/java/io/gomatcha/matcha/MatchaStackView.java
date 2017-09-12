package io.gomatcha.matcha;

import android.content.Context;
import android.support.v7.widget.Toolbar;
import android.util.Log;
import android.view.View;
import android.widget.RelativeLayout;

import com.google.protobuf.InvalidProtocolBufferException;

import java.util.List;

import io.gomatcha.app.R;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.android.PbStackView;

public class MatchaStackView extends MatchaChildView {
    Toolbar toolbar;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/android StackView", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaStackView(context, node);
            }
        });
    }

    public MatchaStackView(Context context, MatchaViewNode node) {
        super(context, node);

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, LayoutParams.WRAP_CONTENT);
        toolbar = new Toolbar(context);
        toolbar.setTitle("TEST");
        toolbar.setId(MatchaPagerView.generateViewId());
        toolbar.setBackgroundColor(0xffff0000);
        toolbar.setNavigationIcon(R.drawable.abc_ic_ab_back_material);
        toolbar.setNavigationOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View view) {
                Log.v("x", "onClick");
            }
        });

        addView(toolbar, params);

        RelativeLayout.LayoutParams contentParams = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentParams.addRule(RelativeLayout.BELOW, toolbar.getId());
        View view = new View(context);
        view.setBackgroundColor(0xff00ffff);
        addView(view, contentParams);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbStackView.StackView proto = buildNode.getBridgeValue().unpack(PbStackView.StackView.class);
        } catch (InvalidProtocolBufferException e) {
        }
    }

    @Override
    public boolean isContainerView() {
        return true;
    }

    @Override
    public void setChildViews(List<View> childViews) {

    }
}
