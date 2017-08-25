package io.gomatcha.matcha;

import android.graphics.Color;
import android.util.Log;
import android.view.View;
import android.widget.AbsoluteLayout;

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
            Log.v("", String.format("xx:%s, %s", buildNode, Arrays.toString(root.getLayoutPaintNodesMap().entrySet().toArray())));
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
                MatchaViewNode childNode = children.get(i);
                this.view.addView(childNode.view);
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

            for (int i = 0; i < layoutPaintNode.getChildOrderCount(); i++) {
                MatchaViewNode childNode = children.get(layoutPaintNode.getChildOrder(i));
                this.view.bringChildToFront(childNode.view); // TODO(KD): Can be done more performantly.
            }

            double maxX = layoutPaintNode.getMaxx();
            double maxY = layoutPaintNode.getMaxy();
            double minX = layoutPaintNode.getMinx();
            double minY = layoutPaintNode.getMiny();

            AbsoluteLayout.LayoutParams params = new AbsoluteLayout.LayoutParams((int)(maxX-minX), (int)(maxY-minY), (int)minX, (int)minY);
            this.view.setLayoutParams(params);
        }

        // Paint view
        if (layoutPaintNode != null & this.paintId != layoutPaintNode.getLayoutId()) {
            this.paintId = layoutPaintNode.getPaintId();

            PbPaint.Style paintStyle = layoutPaintNode.getPaintStyle();
            if (paintStyle.hasBackgroundColor()) {
                Pb.Color c = paintStyle.getBackgroundColor();
                this.view.setBackgroundColor(Color.argb((float)c.getAlpha()/65535, (float)c.getRed()/65535, (float)c.getGreen()/65535, (float)c.getBlue()/65535));
            } else {
                this.view.setBackgroundColor(Color.alpha(0));
            }
        }

        this.children = children;
    }
}
