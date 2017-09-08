package io.gomatcha.matcha;

import android.content.Context;
import android.support.v4.view.PagerAdapter;
import android.support.v4.view.PagerTabStrip;
import android.support.v4.view.PagerTitleStrip;
import android.support.v4.view.ViewPager;
import android.util.Log;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.widget.RelativeLayout;

import java.util.concurrent.atomic.AtomicInteger;

import io.gomatcha.matcha.pb.view.PbView;

public class MatchaSwipeView extends MatchaChildView {
    SlidingTabLayout tabStrip;
    ViewPager viewPager;
    PagerAdapter pagerAdapter;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/android SwipeView", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaSwipeView(context, node);
            }
        });
    }

    public MatchaSwipeView(Context context, MatchaViewNode node) {
        super(context, node);

        pagerAdapter = new PagerAdapter() {
            @Override
            public int getCount() {
                return 2;
            }
            @Override
            public boolean isViewFromObject(View view, Object object) {
                return object == view;
            }
            @Override
            public void destroyItem(ViewGroup container, int position, Object object) {
                container.removeView((View)object);
            }
            @Override
            public Object instantiateItem(ViewGroup container, int position) {
                RelativeLayout v = new RelativeLayout(container.getContext());
                if (position == 0) {
                    v.setBackgroundColor(0xffffffff);
                }
                container.addView(v, 0);
                return v;
            }
            @Override
            public CharSequence getPageTitle(int position) {
                switch (position) {
                    case 0:
                        return "Tab One";
                    case 1:
                        return "Tab Two";
                    case 2:
                        return "Tab Three";
                }
                return null;
            }
        };

        //ViewPager.LayoutParams tabParams = new ViewPager.LayoutParams();
        //tabParams.width = ViewPager.LayoutParams.MATCH_PARENT;
        //tabParams.height = ViewPager.LayoutParams.WRAP_CONTENT;
        RelativeLayout.LayoutParams tabParams = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, LayoutParams.WRAP_CONTENT);
        tabStrip = new SlidingTabLayout(context);
        tabStrip.setId(generateViewId());
        tabStrip.setBackgroundColor(0xff00ffff);
        addView(tabStrip, tabParams);

        RelativeLayout.LayoutParams contentParams = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentParams.addRule(RelativeLayout.BELOW, tabStrip.getId());
        viewPager = new ViewPager(context);
        viewPager.setId(generateViewId());
        viewPager.setBackgroundColor(0xff0000ff);
        viewPager.setAdapter(pagerAdapter);
        addView(viewPager, contentParams);

        tabStrip.setDistributeEvenly(true);
        tabStrip.setViewPager(viewPager);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        // try {
            //PbStackView.View proto = buildNode.getBridgeValue().unpack(PbStackView.View.class);
        // } catch (InvalidProtocolBufferException e) {
        // }
    }

    private static final AtomicInteger sNextGeneratedId = new AtomicInteger(1);
    public static int generateViewId() {
        for (;;) {
            final int result = sNextGeneratedId.get();
            // aapt-generated IDs have the high byte nonzero; clamp to the range under that.
            int newValue = result + 1;
            if (newValue > 0x00FFFFFF) newValue = 1; // Roll over to 1, not 0.
            if (sNextGeneratedId.compareAndSet(result, newValue)) {
                return result;
            }
        }
    }

}
