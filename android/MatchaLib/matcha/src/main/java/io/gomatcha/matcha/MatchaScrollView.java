package io.gomatcha.matcha;

import android.content.Context;
import android.view.MotionEvent;
import android.view.View;
import android.widget.RelativeLayout;
import android.widget.ScrollView;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.matcha.proto.view.PbScrollView;

class MatchaScrollView extends MatchaChildView {
    ScrollView scrollView;
    MatchaLayout childView;
    MatchaViewNode viewNode;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/scrollview", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaScrollView(context, node);
            }
        });
    }

    public MatchaScrollView(Context context, MatchaViewNode node) {
        super(context);
        viewNode = node;
        this.setClipChildren(true);

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
        scrollView = new ScrollView(context);
        scrollView.setFillViewport(true);
        addView(scrollView);

        childView = new MatchaLayout(context);
        scrollView.addView(childView);
    }

    @Override
    public void setNativeState(byte[] nativeState) {
        super.setNativeState(nativeState);
        try {
            PbScrollView.ScrollView proto  = PbScrollView.ScrollView.parseFrom(nativeState);
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
}
