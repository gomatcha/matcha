package io.gomatcha.matcha;

import android.graphics.Color;
import android.util.DisplayMetrics;
import android.util.Log;
import android.view.View;
import android.widget.AbsoluteLayout;
import android.widget.RelativeLayout;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

import io.gomatcha.matcha.pb.Pb;
import io.gomatcha.matcha.pb.paint.PbPaint;
import io.gomatcha.matcha.pb.view.PbView;

public class MatchaViewNode extends Object {
    MatchaViewNode parent;
    MatchaView rootView;
    long id;
    long buildId;
    long layoutId;
    long paintId;
    Map<Long, MatchaViewNode> children = new HashMap<Long, MatchaViewNode>();
    MatchaChildView view;

    public MatchaViewNode(MatchaViewNode parent, MatchaView rootView, long id) {
        this.parent = parent;
        this.rootView = rootView;
        this.id = id;
    }

    public void setRoot(PbView.Root root) {
        PbView.LayoutPaintNode layoutPaintNode = root.getLayoutPaintNodesOrDefault(id, null);
        PbView.BuildNode buildNode = root.getBuildNodesOrDefault(id, null);

        if (this.view == null) {
            this.view = new MatchaChildView(rootView.getContext(), this);
        }


        // Build children
        Map<Long, MatchaViewNode> children = new HashMap<Long, MatchaViewNode>();
        ArrayList<Long> removedKeys = new ArrayList<Long>();
        ArrayList<Long> addedKeys = new ArrayList<Long>();
        ArrayList<Long> unmodifiedKeys = new ArrayList<Long>();
        if (buildNode != null && this.buildId != buildNode.getBuildId()) {
            for (Long i : this.children.keySet()) {
                if (!root.containsBuildNodes(i)) {
                    removedKeys.add(i);
                }
            }
            for (Long i : buildNode.getChildrenList()) {
                MatchaViewNode prevChild = this.children.get(i);
                if (prevChild == null) {
                    addedKeys.add(i);
                    MatchaViewNode child = new MatchaViewNode(this, this.rootView, i);
                    children.put(i, child);
                } else {
                    unmodifiedKeys.add(i);
                    children.put(i, prevChild);
                }
            }
        } else {
            // Log.v("", String.format("xx:%s, %s", buildNode, Arrays.toString(root.getLayoutPaintNodesMap().entrySet().toArray())));
            children = this.children;
        }

        // Update children
        for (MatchaViewNode i : children.values()) {
            i.setRoot(root);
        }

        if (buildNode != null && this.buildId != buildNode.getBuildId()) {
            this.buildId = buildNode.getBuildId();

            // Update the views with native values
            this.view.setNode(buildNode);

            // Add/remove subviews
            for (long i : addedKeys) {
                RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(300, 300);
                params.leftMargin = 100;
                params.topMargin = 100;
                
                MatchaViewNode childNode = children.get(i);
                this.view.addView(childNode.view, params);
                Log.v("addView", String.format("%s, %s", childNode.id, params));
            }
            for (long i : removedKeys) {
                MatchaViewNode childNode = children.get(i);
                this.view.removeView(childNode.view);
            }

            // Update gesture recognizers... TODO(KD):
        }

        // Layout subviews
        if (layoutPaintNode != null && this.layoutId != layoutPaintNode.getLayoutId()) {
            this.layoutId = layoutPaintNode.getLayoutId();

            // for (int i = 0; i < layoutPaintNode.getChildOrderCount(); i++) {
            //     MatchaViewNode childNode = children.get(layoutPaintNode.getChildOrder(i));
            //     this.view.bringChildToFront(childNode.view); // TODO(KD): Can be done more performantly.
            // }

            double maxX = layoutPaintNode.getMaxx() * this.view.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
            double maxY = layoutPaintNode.getMaxy() * this.view.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
            double minX = layoutPaintNode.getMinx() * this.view.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
            double minY = layoutPaintNode.getMiny() * this.view.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;

            if (this.parent != null) {
                RelativeLayout.LayoutParams params = (RelativeLayout.LayoutParams)this.view.getLayoutParams();
                if (params == null) {
                    params = new RelativeLayout.LayoutParams(0, 0);
                }
                params.width = (int)(maxX-minX);
                params.height = (int)(maxY-minY);
                params.leftMargin = (int)minX;
                params.topMargin = (int)minY;
                this.view.setLayoutParams(params);

                Log.v("setLayoutParams", String.format("%s, %s", this.id, params));
            }
        }

        // Paint view
        if (layoutPaintNode != null & this.paintId != layoutPaintNode.getLayoutId()) {
            this.paintId = layoutPaintNode.getPaintId();

            PbPaint.Style paintStyle = layoutPaintNode.getPaintStyle();
            if (paintStyle.hasBackgroundColor()) {
                Pb.Color c = paintStyle.getBackgroundColor();
                this.view.setBackgroundColor(Color.argb(c.getAlpha()*255/65535, c.getRed()*255/65535, c.getGreen()*255/65535, c.getBlue()*255/65535));
            } else {
                this.view.setBackgroundColor(Color.alpha(0));
            }
        }

        this.children = children;
    }
}
