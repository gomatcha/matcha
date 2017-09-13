package io.gomatcha.matcha;

import android.content.Context;
import android.support.v7.widget.Toolbar;
import android.util.DisplayMetrics;
import android.util.Log;
import android.view.View;
import android.widget.RelativeLayout;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.app.R;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.android.PbStackView;

public class MatchaToolbarView extends MatchaChildView {
    Toolbar toolbar;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/android stackBarView", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaToolbarView(context, node);
            }
        });
    }

    public MatchaToolbarView(Context context, MatchaViewNode node) {
        super(context, node);

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
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
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbStackView.StackBar proto = buildNode.getBridgeValue().unpack(PbStackView.StackBar.class);
        } catch (InvalidProtocolBufferException e) {
        }
    }

    @Override
    protected void onMeasure(int widthMeasureSpec, int heightMeasureSpec) {
        int desiredHeight = (int)Math.ceil(56.0 * (float)getContext().getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT);
        if (MeasureSpec.getMode(heightMeasureSpec) == MeasureSpec.UNSPECIFIED ||
                (MeasureSpec.getMode(heightMeasureSpec) == MeasureSpec.AT_MOST && MeasureSpec.getSize(heightMeasureSpec) > desiredHeight)) {
            heightMeasureSpec = MeasureSpec.makeMeasureSpec(desiredHeight, MeasureSpec.EXACTLY);
        }
        super.onMeasure(widthMeasureSpec, heightMeasureSpec);
    }
}
