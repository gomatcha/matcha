package io.gomatcha.matcha;

import android.content.Context;
import android.support.v7.widget.Toolbar;
import android.util.Log;
import android.view.View;

import com.google.protobuf.InvalidProtocolBufferException;

import java.util.ArrayList;
import java.util.List;

import io.gomatcha.matcha.proto.view.android.PbStackView;

class MatchaStackView extends MatchaChildView {
    MatchaViewNode viewNode;
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
        super(context);
        viewNode = node;

        stackView2 = new MatchaStackView2(context);
        addView(stackView2);
    }

    @Override
    public void setNativeState(byte[] nativeState) {
        super.setNativeState(nativeState);
        try {
            PbStackView.StackView proto  = PbStackView.StackView.parseFrom(nativeState);
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
            if (toolbar.getId() < 0) {
                toolbar.setId(MatchaPagerView.generateViewId());
            }
            toolbar.stackView = this;
            wrapper.setToolbarView(toolbar);

            View childView = childViews.get(i * 2 + 1);
            wrapper.setContentView(childView);

            toolbarViews.add(wrapper);
        }
        stackView2.setChildViews(toolbarViews);
    }

    public void back() {
        viewNode.call("OnBack");
    }
}
