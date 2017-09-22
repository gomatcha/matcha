package io.gomatcha.matcha;

import android.content.Context;
import android.view.MotionEvent;
import android.view.View;
import android.widget.RelativeLayout;
import android.widget.ScrollView;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.matcha.proto.view.PbScrollView;
import io.gomatcha.matcha.proto.view.PbView;

class MatchaScrollView extends MatchaChildView {
    ScrollView scrollView;
    MatchaLayout childView;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/scrollview", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaScrollView(context, node);
            }
        });
    }

    public MatchaScrollView(Context context, MatchaViewNode node) {
        super(context, node);
        this.setClipChildren(true);

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
        scrollView = new ScrollView(context);
        scrollView.setFillViewport(true);
        addView(scrollView);

        childView = new MatchaLayout(context);
        scrollView.addView(childView);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbScrollView.ScrollView proto = buildNode.getBridgeValue().unpack(PbScrollView.ScrollView.class);
            scrollView.setVerticalScrollBarEnabled(proto.getShowsVerticalScrollIndicator());
            scrollView.setHorizontalScrollBarEnabled(proto.getShowsHorizontalScrollIndicator());

            if (!proto.getScrollEnabled()) {
                scrollView.setOnTouchListener(new OnTouchListener() {
                    @Override
                    public boolean onTouch(View view, MotionEvent motionEvent) {
                        return false;
                    }
                });
            } else {
                scrollView.setOnTouchListener(null);
            }
        } catch (InvalidProtocolBufferException e) {
        }
    }

    @Override
    public MatchaLayout getLayout() {
        return childView;
    }
}
