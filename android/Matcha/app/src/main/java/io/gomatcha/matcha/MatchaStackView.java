package io.gomatcha.matcha;

import android.content.Context;
import android.support.v7.widget.Toolbar;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;
import android.view.ViewParent;
import android.widget.RelativeLayout;

import com.google.protobuf.InvalidProtocolBufferException;

import java.util.ArrayList;
import java.util.List;

import io.gomatcha.app.R;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.android.PbStackView;

public class MatchaStackView extends MatchaChildView {
    Toolbar toolbar;
    MatchaStackView2 stackView2;

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

        RelativeLayout.LayoutParams contentParams = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        stackView2 = new MatchaStackView2(context);
        stackView2.setBackgroundColor(0xff00ffff);
        addView(stackView2, contentParams);
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
        ArrayList<View> toolbarViews = new ArrayList<View>();
        for (int i = 0; i < childViews.size() / 2; i++) {
            RelativeLayout layout = new RelativeLayout(getContext());

            RelativeLayout.LayoutParams toolbarParams = new RelativeLayout.LayoutParams(LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT);
            MatchaToolbarView toolbar = (MatchaToolbarView)childViews.get(i*2);
            ViewParent toolbarParent = toolbar.getParent();
            if (toolbarParent != null) {
                ((RelativeLayout)toolbarParent).removeView(toolbar);
            }
            toolbar.setId(MatchaPagerView.generateViewId());
            layout.addView(toolbar, toolbarParams);

            RelativeLayout.LayoutParams childViewParams = new RelativeLayout.LayoutParams(LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.MATCH_PARENT);
            childViewParams.addRule(RelativeLayout.BELOW, toolbar.getId());
            View childView = childViews.get(i * 2 + 1);
            ViewParent childParent = childView.getParent();
            if (childParent != null) {
                ((RelativeLayout)childParent).removeView(childView);
            }
            layout.addView(childView, childViewParams);

            toolbarViews.add(layout);
        }
        stackView2.setChildViews(toolbarViews);
    }
}
