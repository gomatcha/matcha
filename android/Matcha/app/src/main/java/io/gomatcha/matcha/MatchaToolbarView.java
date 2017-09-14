package io.gomatcha.matcha;

import android.content.Context;
import android.support.v7.widget.Toolbar;
import android.util.DisplayMetrics;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.RelativeLayout;

import com.google.protobuf.InvalidProtocolBufferException;

import java.util.List;

import io.gomatcha.app.R;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.android.PbStackView;

public class MatchaToolbarView extends MatchaChildView {
    Toolbar toolbar;
    MatchaStackView stackView;

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
        toolbar.setId(MatchaPagerView.generateViewId());
        toolbar.setNavigationOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View view) {
                stackView.back();
            }
        });
        addView(toolbar, params);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbStackView.StackBar proto = buildNode.getBridgeValue().unpack(PbStackView.StackBar.class);

            toolbar.setTitle(proto.getTitle());
            toolbar.setSubtitle(proto.getSubtitle());
            if (proto.getBackButtonHidden()) {
                toolbar.setNavigationIcon(null);
            } else {
                toolbar.setNavigationIcon(R.drawable.abc_ic_ab_back_material);
            }

            List<PbStackView.StackBarItem> itemList = proto.getItemsList();
            Menu menu = toolbar.getMenu();
            menu.clear();
            for (int i = 0; i < itemList.size(); i++) {
                PbStackView.StackBarItem protoItem = itemList.get(i);
                MenuItem item = menu.add(0, Menu.FIRST + i, Menu.NONE, protoItem.getTitle());
                item.setShowAsAction(MenuItem.SHOW_AS_ACTION_ALWAYS);
            }
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
