package io.gomatcha.matcha;

import android.content.Context;
import android.text.Layout;
import android.util.AttributeSet;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;

public class MatchaLayout extends ViewGroup {
    public MatchaLayout(Context context) {
        super(context);
    }

    @Override
    protected void onMeasure (int widthMeasureSpec, int heightMeasureSpec) {
        int w = 0;
        int h = 0;
        int childCount = getChildCount();
        for(int i=0; i<childCount; i++) {
            View v = getChildAt(i);
            LayoutParams params = (LayoutParams)v.getLayoutParams();
            int measuredWidth = (int)(params.right - params.left);
            int measuredHeight = (int)(params.bottom - params.top);
            if (measuredWidth > w) {
                w = measuredWidth;
            }
            if (measuredHeight > h) {
                h = measuredHeight;
            }
        }
        w = resolveSize(w, widthMeasureSpec);
        h = resolveSize(h, heightMeasureSpec);
        setMeasuredDimension(w, h);
    }

    @Override
    protected void onLayout(boolean b, int left, int top, int right, int bottom) {
        int childCount = getChildCount();
        for(int i=0; i<childCount; i++) {
            View v = getChildAt(i);
            ViewGroup.LayoutParams params = v.getLayoutParams();
            LayoutParams matchaParams = (LayoutParams)params;
            if (matchaParams.full) {
                v.measure(MeasureSpec.makeMeasureSpec(right-left, MeasureSpec.AT_MOST), MeasureSpec.makeMeasureSpec(bottom-top, MeasureSpec.AT_MOST));
                v.layout(0, 0, right-left, bottom-top);
            } else {
                v.measure(MeasureSpec.makeMeasureSpec((int)(matchaParams.right-matchaParams.left), MeasureSpec.AT_MOST), MeasureSpec.makeMeasureSpec((int)(matchaParams.bottom-matchaParams.top), MeasureSpec.AT_MOST));
                v.layout((int) matchaParams.left, (int) matchaParams.top, (int) matchaParams.right, (int) matchaParams.bottom);
            }
        }
    }

    @Override
    protected ViewGroup.LayoutParams generateDefaultLayoutParams() {
        LayoutParams x = new LayoutParams();
        x.full = true;
        return x;
    }

    static class LayoutParams extends ViewGroup.LayoutParams {
        double left;
        double top;
        double right;
        double bottom;
        boolean full;
        public LayoutParams() {
            super(0, 0);
        }
    }
}
