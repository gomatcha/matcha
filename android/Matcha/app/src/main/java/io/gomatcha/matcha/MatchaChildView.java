package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.graphics.PointF;
import android.support.v4.view.GestureDetectorCompat;
import android.util.DisplayMetrics;
import android.util.Log;
import android.view.GestureDetector;
import android.view.MotionEvent;
import android.view.View;
import android.view.ViewConfiguration;
import android.widget.RelativeLayout;

import com.google.protobuf.Duration;
import com.google.protobuf.InvalidProtocolBufferException;

import java.util.Date;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.pb.touch.PbTouch;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.slider.PbSlider;

public class MatchaChildView extends RelativeLayout {
    MatchaViewNode viewNode;
    PbView.BuildNode buildNode;
    MatchaGestureDetector matchaGestureDetector;
    GestureDetector gestureDetector;

    public MatchaChildView(Context context, MatchaViewNode node) {
        super(context);
        viewNode = node;
        this.setClipChildren(false);
        this.matchaGestureDetector = new MatchaGestureDetector();
        this.matchaGestureDetector.childView = this;
        this.matchaGestureDetector.context = context;
        this.gestureDetector = new GestureDetector(this.getContext(), this.matchaGestureDetector);

        this.setOnTouchListener(new View.OnTouchListener() {
            @Override
            public boolean onTouch(View v, MotionEvent event) {
                return MatchaChildView.this.gestureDetector.onTouchEvent(event);
            }
        });
    }

    public void setNode(PbView.BuildNode buildNode) {
        this.buildNode = buildNode;
    }

    public RelativeLayout getLayout() {
        return this;
    }

    class MatchaGestureDetector extends GestureDetector.SimpleOnGestureListener {
        MatchaChildView childView;
        com.google.protobuf.Any tapGesture;
        com.google.protobuf.Any pressGesture;
        com.google.protobuf.Any buttonGesture;
        Context context;

        @Override
        public boolean onDown(MotionEvent event) {
            return true;
        }

        @Override
        public void onLongPress(MotionEvent event) {
            try {
                if (pressGesture == null) {
                    return;
                }
                PbTouch.PressRecognizer proto = pressGesture.unpack(PbTouch.PressRecognizer.class);
                float ratio = (float)MatchaChildView.this.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;

                Duration duration = Protobuf.toProtobuf(ViewConfiguration.get(context).getLongPressTimeout());
                PbTouch.PressEvent e = PbTouch.PressEvent.newBuilder()
                        .setDuration(duration)
                        .setTimestamp(Protobuf.toProtobuf(new Date()))
                        .setPosition(Protobuf.toProtobuf(new PointF(event.getX() / ratio, event.getY() / ratio)))
                        .setKind(PbTouch.EventKind.EVENT_KIND_RECOGNIZED)
                        .build();

                childView.viewNode.rootView.call(String.format("%d", proto.getFuncId()), childView.viewNode.id, new GoValue(e.toByteArray()));
            } catch (InvalidProtocolBufferException e) {
            }
        }

        @Override
        public boolean onDoubleTap(MotionEvent event) {
            try {
                if (tapGesture == null) {
                    return false;
                }
                PbTouch.TapRecognizer proto = tapGesture.unpack(PbTouch.TapRecognizer.class);
                if (proto.getCount() != 2) {
                    return false;
                }
                float ratio = (float)MatchaChildView.this.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
                PbTouch.TapEvent e = PbTouch.TapEvent.newBuilder()
                        .setKind(PbTouch.EventKind.EVENT_KIND_RECOGNIZED)
                        .setTimestamp(Protobuf.toProtobuf(new Date()))
                        .setPosition(Protobuf.toProtobuf(new PointF(event.getX() / ratio, event.getY() / ratio)))
                        .build();

                childView.viewNode.rootView.call(String.format("%d", proto.getRecognizedFunc()), childView.viewNode.id, new GoValue(e.toByteArray()));
                return true;
            } catch (InvalidProtocolBufferException e) {
            }
            return false;
        }

        @Override
        public boolean onSingleTapUp(MotionEvent event) {
            try {
                if (tapGesture == null) {
                    return false;
                }
                PbTouch.TapRecognizer proto = tapGesture.unpack(PbTouch.TapRecognizer.class);
                if (proto.getCount() != 1) {
                    return false;
                }
                float ratio = (float)MatchaChildView.this.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
                PbTouch.TapEvent e = PbTouch.TapEvent.newBuilder()
                        .setKind(PbTouch.EventKind.EVENT_KIND_RECOGNIZED)
                        .setTimestamp(Protobuf.toProtobuf(new Date()))
                        .setPosition(Protobuf.toProtobuf(new PointF(event.getX() / ratio, event.getY() / ratio)))
                        .build();

                childView.viewNode.rootView.call(String.format("%d", proto.getRecognizedFunc()), childView.viewNode.id, new GoValue(e.toByteArray()));
                return true;
            } catch (InvalidProtocolBufferException e) {
            }
            return false;
        }
    }
}
