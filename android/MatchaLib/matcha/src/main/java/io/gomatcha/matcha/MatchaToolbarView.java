package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.PorterDuff;
import android.graphics.drawable.Drawable;
import android.support.v7.widget.Toolbar;
import android.util.DisplayMetrics;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;

import com.google.protobuf.InvalidProtocolBufferException;

import java.util.List;

import io.gomatcha.matcha.proto.view.android.PbStackView;

class MatchaToolbarView extends MatchaChildView {
    Toolbar toolbar;
    MatchaStackView stackView;
    MatchaViewNode viewNode;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/android stackBarView", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaToolbarView(context, node);
            }
        });
    }

    public MatchaToolbarView(Context context, MatchaViewNode node) {
        super(context);
        viewNode = node;

        toolbar = new Toolbar(context);
        toolbar.setId(MatchaPagerView.generateViewId());
        toolbar.setNavigationOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View view) {
                stackView.back();
            }
        });
        if (android.os.Build.VERSION.SDK_INT >= 21){
            float ratio = (float)context.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
            this.setElevation(4*ratio);
        }
        addView(toolbar);
    }

    @Override
    public void setNativeState(byte[] nativeState) {
        super.setNativeState(nativeState);
        try {
            PbStackView.StackBar proto  = PbStackView.StackBar.parseFrom(nativeState);

            if (proto.hasStyledTitle()) {
                toolbar.setTitle(Protobuf.newAttributedString(proto.getStyledTitle()));
            } else {
                toolbar.setTitle(proto.getTitle());
            }
            if (proto.hasStyledSubtitle()) {
                toolbar.setSubtitle(Protobuf.newAttributedString(proto.getStyledSubtitle()));
            } else {
                toolbar.setSubtitle(proto.getSubtitle());
            }
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
                final String onPressFunc = protoItem.getOnPressFunc();

                MenuItem item = menu.add(0, Menu.FIRST + i, Menu.NONE, protoItem.getTitle());
                item.setShowAsAction(MenuItem.SHOW_AS_ACTION_ALWAYS);
                item.setOnMenuItemClickListener(new MenuItem.OnMenuItemClickListener() {
                    @Override
                    public boolean onMenuItemClick(MenuItem menuItem) {
                        MatchaToolbarView.this.viewNode.call(onPressFunc);
                        return true;
                    }
                });
                item.setEnabled(!protoItem.getDisabled());

                if (protoItem.hasIcon()) {
                    Drawable icon = Protobuf.newDrawable(protoItem.getIcon(), getContext());
                    if (protoItem.hasIconTint()) {
                        icon.setColorFilter(Protobuf.newColor(protoItem.getIconTint()), PorterDuff.Mode.SRC_ATOP);
                    }
                    item.setIcon(icon);
                }
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
