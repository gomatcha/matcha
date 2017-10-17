package io.gomatcha.matcha;

import android.graphics.Color;
import android.graphics.PointF;
import android.graphics.drawable.GradientDrawable;
import android.util.DisplayMetrics;
import android.util.Log;
import android.view.View;
import android.widget.ScrollView;
import android.widget.HorizontalScrollView;

import com.google.protobuf.InvalidProtocolBufferException;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.proto.paint.PbPaint;
import io.gomatcha.matcha.proto.pointer.PbPointer;
import io.gomatcha.matcha.proto.view.PbView;

public class MatchaViewNode extends Object {
    MatchaViewNode parent;
    public MatchaView rootView;
    private long id;
    long buildId;
    long layoutId;
    long paintId;
    Map<Long, MatchaViewNode> children = new HashMap<Long, MatchaViewNode>();
    ArrayList<MatchaViewNode> childList = new ArrayList<MatchaViewNode>();
    MatchaChildView view;

    MatchaViewNode(MatchaViewNode parent, MatchaView rootView, long id) {
        this.parent = parent;
        this.rootView = rootView;
        this.id = id;
    }

    public void call(String func, GoValue... args) {
        if (this.rootView.updating) {
            return;
        }
        this.rootView.call(func, this.id, args);
    }

    void setRoot(PbView.Root root) {
        PbView.LayoutPaintNode layoutPaintNode = root.getLayoutPaintNodesOrDefault(id, null);
        PbView.BuildNode buildNode = root.getBuildNodesOrDefault(id, null);

        // Create view
        if (this.view == null) {
            this.view = MatchaView.createView(buildNode.getBridgeName(), rootView.getContext(), this);
            this.view.matchaGestureRecognizer.viewNode = this;
        }
        
        MatchaLayout layout = this.view;
        if (this.view instanceof MatchaScrollView) {
            layout = ((MatchaScrollView)this.view).childView;
        }

        // Build children
        Map<Long, MatchaViewNode> children = new HashMap<Long, MatchaViewNode>();
        ArrayList<MatchaViewNode> childList = new ArrayList<MatchaViewNode>();
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
                    childList.add(child);
                    children.put(i, child);
                } else {
                    unmodifiedKeys.add(i);
                    childList.add(prevChild);
                    children.put(i, prevChild);
                }
            }
        } else {
            children = this.children;
            childList = this.childList;
        }

        // Update children
        for (MatchaViewNode i : children.values()) {
            i.setRoot(root);
        }

        if (buildNode != null && this.buildId != buildNode.getBuildId()) {
            this.buildId = buildNode.getBuildId();

            // Update the views with native values
            this.view.setNativeState(buildNode.getBridgeValue().toByteArray());

            // Add/remove subviews
            if (this.view.isContainerView()) {
                ArrayList<View> childViews = new ArrayList<View>();
                for (MatchaViewNode i : childList) {
                    childViews.add(i.view);
                }
                this.view.setChildViews(childViews);
            } else {
                for (long i : addedKeys) {
                    MatchaViewNode childNode = children.get(i);
                    layout.addView(childNode.view);
                }
                for (long i : removedKeys) {
                    MatchaViewNode childNode = this.children.get(i);
                    layout.removeView(childNode.view);
                }
            }

            // Update gesture recognizers... TODO(KD):
            com.google.protobuf.ByteString gestures = buildNode.getValuesMap().get("gomatcha.io/matcha/touch");
            if (gestures != null) {
                try {
                    PbPointer.RecognizerList proto = PbPointer.RecognizerList.parseFrom(gestures);
                    for (PbPointer.Recognizer i : proto.getRecognizersList()) {
                        String type = i.getRecognizer().getTypeUrl();
                        if (type.equals("type.googleapis.com/matcha.pointer.TapRecognizer")) {
                            this.view.matchaGestureRecognizer.tapGesture = i.getRecognizer();
                        } else if (type.equals("type.googleapis.com/matcha.pointer.PressRecognizer")) {
                            this.view.matchaGestureRecognizer.pressGesture = i.getRecognizer();
                        } else if (type.equals("type.googleapis.com/matcha.pointer.ButtonRecognizer")) {
                            this.view.matchaGestureRecognizer.buttonGesture = i.getRecognizer();
                        }
                    }
                    this.view.matchaGestureRecognizer.reload();
                    this.view.setClickable(proto.getRecognizersCount() > 0);
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
            } else if (this.parent.view.isContainerView()) {
                // Let containers do their own layout.
            } else if (this.parent.view instanceof MatchaScrollView) {
                // Translate the scrollview's contentView offset into a ScrollX and ScrollY.
                MatchaScrollView matchaScrollView = (MatchaScrollView)this.parent.view;
                double offsetX = -minX;
                double offsetY = -minY;
                minX = 0;
                minY = 0;
                maxX += offsetX;
                maxY += offsetY;

                if (!matchaScrollView.horizontal) {
                    ScrollView scrollView = matchaScrollView.scrollView;
                    scrollView.setScrollX((int)offsetX);
                    scrollView.setScrollY((int)offsetY);
                } else {
                    HorizontalScrollView scrollView = matchaScrollView.hScrollView;
                    scrollView.setScrollX((int)offsetX);
                    scrollView.setScrollY((int)offsetY);
                }
                matchaScrollView.matchaX = (int)offsetX;
                matchaScrollView.matchaY = (int)offsetY;

                MatchaLayout.LayoutParams params = (MatchaLayout.LayoutParams)this.view.getLayoutParams();
                if (params == null) {
                    params = new MatchaLayout.LayoutParams();
                }
                params.left = minX;
                params.top = minY;
                params.right = maxX;
                params.bottom = maxY;
                this.view.setLayoutParams(params);
            } else {
                MatchaLayout.LayoutParams params = (MatchaLayout.LayoutParams)this.view.getLayoutParams();
                if (params == null) {
                    params = new MatchaLayout.LayoutParams();
                }
                params.left = minX;
                params.top = minY;
                params.right = maxX;
                params.bottom = maxY;
                this.view.setLayoutParams(params);
            }
        }

        // Paint scrollView
        if (layoutPaintNode != null & this.paintId != layoutPaintNode.getPaintId()) {
            this.paintId = layoutPaintNode.getPaintId();

            PbPaint.Style paintStyle = layoutPaintNode.getPaintStyle();

            double ratio = (float)this.view.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
            GradientDrawable gd = new GradientDrawable();

            double cornerRadius = paintStyle.getCornerRadius();
            gd.setCornerRadius((float)(cornerRadius * ratio));

            if (paintStyle.getHasBorderColor()) {
                int c = Protobuf.newColor(paintStyle.getBorderColorRed(), paintStyle.getBorderColorGreen(), paintStyle.getBorderColorBlue(), paintStyle.getBorderColorAlpha());
                gd.setStroke((int)(paintStyle.getBorderWidth() * ratio), c);
            } else {
                gd.setStroke(0, 0);
            }

            if (this.view instanceof MatchaImageView) {
                ((MatchaImageView)this.view).view.setCornerRadius((float)(cornerRadius*ratio));
                //((MatchaImageView)this.view).view.setBorderColor(Protobuf.newColor(paintStyle.getBorderColor()));
                ((MatchaImageView)this.view).view.setBorderWidth((float)(paintStyle.getBorderWidth()*ratio));
            }

            if (paintStyle.getHasBackgroundColor()) {
                int c = Protobuf.newColor(paintStyle.getBackgroundColorRed(), paintStyle.getBackgroundColorGreen(), paintStyle.getBackgroundColorBlue(), paintStyle.getBackgroundColorAlpha());
                gd.setColor(c);
            } else {
                gd.setColor(Color.alpha(0));
            }
            this.view.setBackground(gd);

            this.view.setAlpha((float)(1.0 - paintStyle.getTransparency()));
        }

        this.children = children;
    }
}
