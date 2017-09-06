package io.gomatcha.matcha;

import android.graphics.Color;
import android.util.DisplayMetrics;
import android.util.Log;
import android.view.View;
import android.widget.AbsoluteLayout;
import android.widget.RelativeLayout;

import com.google.protobuf.InvalidProtocolBufferException;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

import io.gomatcha.matcha.pb.Pb;
import io.gomatcha.matcha.pb.paint.PbPaint;
import io.gomatcha.matcha.pb.touch.PbTouch;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.scrollview.PbScrollView;

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

        // Create view
        if (this.view == null) {
            this.view = MatchaView.createView(buildNode.getBridgeName(), rootView.getContext(), this);
        }
        RelativeLayout layout = this.view.getLayout();

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
                layout.addView(childNode.view);
            }
            for (long i : removedKeys) {
                MatchaViewNode childNode = this.children.get(i);
                layout.removeView(childNode.view);
            }

            // Update gesture recognizers... TODO(KD):
            com.google.protobuf.Any gestures = buildNode.getValuesMap().get("gomatcha.io/matcha/touch");
            if (gestures != null) {
                try {
                    PbTouch.RecognizerList proto = gestures.unpack(PbTouch.RecognizerList.class);
                    for (PbTouch.Recognizer i : proto.getRecognizersList()) {
                        String type = i.getRecognizer().getTypeUrl();
                        if (type.equals("type.googleapis.com/matcha.touch.TapRecognizer")) {
                            this.view.matchaGestureDetector.tapGesture = i.getRecognizer();
                        } else if (type.equals("type.googleapis.com/matcha.touch.PressRecognizer")) {
                            this.view.matchaGestureDetector.pressGesture = i.getRecognizer();
                        } else if (type.equals("type.googleapis.com/matcha.touch.ButtonRecognizer")) {
                            this.view.matchaGestureDetector.buttonGesture = i.getRecognizer();
                        }
                    }
                } catch (InvalidProtocolBufferException e) {
                }
            }
        }

        // Layout subviews
        if (layoutPaintNode != null && this.layoutId != layoutPaintNode.getLayoutId()) {
            this.layoutId = layoutPaintNode.getLayoutId();

            for (int i = 0; i < layoutPaintNode.getChildOrderCount(); i++) {
                MatchaViewNode childNode = children.get(layoutPaintNode.getChildOrder(i));
                layout.bringChildToFront(childNode.view); // TODO(KD): Can be done more performantly.
            }

            double ratio = (float)this.view.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
            double maxX = layoutPaintNode.getMaxx() * ratio;
            double maxY = layoutPaintNode.getMaxy() * ratio;
            double minX = layoutPaintNode.getMinx() * ratio;
            double minY = layoutPaintNode.getMiny() * ratio;

            if (this.parent == null) {
            // } else if (this.parent.view.getClass().isInstance(MatchaScrollView.class)) {
            } else {
                RelativeLayout.LayoutParams params = (RelativeLayout.LayoutParams)this.view.getLayoutParams();
                if (params == null) {
                    params = new RelativeLayout.LayoutParams(0, 0);
                }
                params.width = (int)(maxX-minX);
                params.height = (int)(maxY-minY);
                params.leftMargin = (int)minX;
                params.topMargin = (int)minY;
                this.view.setLayoutParams(params);
            }
        }

        // Paint view
        if (layoutPaintNode != null & this.paintId != layoutPaintNode.getLayoutId()) {
            this.paintId = layoutPaintNode.getPaintId();

            PbPaint.Style paintStyle = layoutPaintNode.getPaintStyle();
            if (paintStyle.hasBackgroundColor()) {
                Pb.Color c = paintStyle.getBackgroundColor();
                this.view.setBackgroundColor(Protobuf.newColor(c));
            } else {
                this.view.setBackgroundColor(Color.alpha(0));
            }
        }

        this.children = children;
    }
}
