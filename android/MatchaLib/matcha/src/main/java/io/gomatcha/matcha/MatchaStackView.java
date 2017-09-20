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

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.pb.view.PbSwitchView;
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

        stackView2 = new MatchaStackView2(context);
        addView(stackView2);
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

    ArrayList<MatchaToolbarWrapper> wrappers = new ArrayList<MatchaToolbarWrapper>();

    @Override
    public void setChildViews(List<View> childViews) {
        while (wrappers.size() < childViews.size()/2) {
            wrappers.add(new MatchaToolbarWrapper(this.getContext()));
        }
        while (wrappers.size() > childViews.size()/2) {
            wrappers.remove(wrappers.size()-1);
        }

        ArrayList<View> toolbarViews = new ArrayList<View>();
        for (int i = 0; i < childViews.size() / 2; i++) {
            MatchaToolbarWrapper wrapper = wrappers.get(i);

            MatchaToolbarView toolbar = (MatchaToolbarView)childViews.get(i*2);
            toolbar.setId(MatchaPagerView.generateViewId());
            toolbar.stackView = this;
            wrapper.setToolbarView(toolbar);

            View childView = childViews.get(i * 2 + 1);
            wrapper.setContentView(childView);

            toolbarViews.add(wrapper);
        }
        stackView2.setChildViews(toolbarViews);
    }

    public void back() {
        viewNode.rootView.call("OnBack", viewNode.id);
    }
}
