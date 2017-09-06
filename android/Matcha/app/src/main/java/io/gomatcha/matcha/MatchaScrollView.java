package io.gomatcha.matcha;

import android.content.Context;
import android.os.Build;
import android.util.Log;
import android.view.MotionEvent;
import android.view.View;
import android.widget.RelativeLayout;
import android.widget.ScrollView;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.scrollview.PbScrollView;

public class MatchaScrollView extends MatchaChildView {
    ScrollView view;
    RelativeLayout childView;

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
        view = new ScrollView(context);
        view.setFillViewport(true);

        addView(view, params);

        childView = new RelativeLayout(context);
        view.addView(childView);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbScrollView.View proto = buildNode.getBridgeValue().unpack(PbScrollView.View.class);
            view.setVerticalScrollBarEnabled(proto.getShowsVerticalScrollIndicator());
            view.setHorizontalScrollBarEnabled(proto.getShowsHorizontalScrollIndicator());

            if (!proto.getScrollEnabled()) {
                view.setOnTouchListener(new OnTouchListener() {
                    @Override
                    public boolean onTouch(View view, MotionEvent motionEvent) {
                        return false;
                    }
                });
            } else {
                view.setOnTouchListener(null);
            }
        } catch (InvalidProtocolBufferException e) {
        }
    }

    @Override
    public RelativeLayout getLayout() {
        return childView;
    }
}
