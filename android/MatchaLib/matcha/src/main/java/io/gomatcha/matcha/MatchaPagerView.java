package io.gomatcha.matcha;

import android.content.Context;
import android.support.v4.view.PagerAdapter;
import android.support.v4.view.ViewPager;
import android.view.View;
import android.view.ViewGroup;
import android.widget.RelativeLayout;

import com.google.protobuf.InvalidProtocolBufferException;

import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.android.PbPagerView;

public class MatchaPagerView extends MatchaChildView {
    SlidingTabLayout tabStrip;
    ViewPager viewPager;
    MatchaPagerAdapter pagerAdapter;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/android PagerView", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaPagerView(context, node);
            }
        });
    }

    public MatchaPagerView(Context context, MatchaViewNode node) {
        super(context, node);

        pagerAdapter = new MatchaPagerAdapter();

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
        try {
            PbPagerView.PagerView proto = buildNode.getBridgeValue().unpack(PbPagerView.PagerView.class);
            pagerAdapter.protoChildViews = proto.getChildViewsList();
            pagerAdapter.notifyDataSetChanged();
            tabStrip.setViewPager(viewPager);
        } catch (InvalidProtocolBufferException e) {
        }
    }

    @Override
    public boolean isContainerView() {
        return true;
    }

    @Override
    public void setChildViews(List<View> childViews) {
        pagerAdapter.childViews = childViews;
        pagerAdapter.notifyDataSetChanged();
        tabStrip.setViewPager(viewPager);
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

    public class MatchaPagerAdapter extends PagerAdapter {
        List<View> childViews;
        List<PbPagerView.PagerChildView> protoChildViews;

        @Override
        public int getCount() {
            if (childViews == null) {
                return 0;
            }
            return childViews.size();
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
            if (childViews == null) {
                View v = new View(container.getContext());
                container.addView(v);
                return v;
            }
            View v = childViews.get(position);
            container.addView(v);
            return v;
        }
        @Override
        public CharSequence getPageTitle(int position) {
            return protoChildViews.get(position).getTitle();
        }
    }
}
