package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.PointF;
import android.util.DisplayMetrics;
import android.view.MotionEvent;
import android.view.View;

import com.google.protobuf.InvalidProtocolBufferException;

import java.util.Date;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.pb.touch.PbTouch;

public class MatchaGestureRecognizer implements View.OnTouchListener {
    MatchaChildView childView;
    com.google.protobuf.Any tapGesture;
    com.google.protobuf.Any pressGesture;
    com.google.protobuf.Any buttonGesture;
    TapRecognizer tapRecognizer;
    PressRecognizer pressRecognizer;
    ButtonRecognizer buttonRecognizer;
    Result prevButtonResult;
    Context context;

    public enum State {
        POSSIBLE,
        CHANGED,
        FAILED,
        RECOGNIZED,
    }

    public void reload() {
        try {
            if (tapGesture != null && tapRecognizer != null) {
                PbTouch.TapRecognizer proto = tapGesture.unpack(PbTouch.TapRecognizer.class);
                tapRecognizer.recognizerId = (int)proto.getOnEvent();
                tapRecognizer.ratio = context.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
            }
            if (pressGesture != null && pressRecognizer != null) {
                PbTouch.PressRecognizer proto = pressGesture.unpack(PbTouch.PressRecognizer.class);
                pressRecognizer.recognizerId = (int)proto.getOnEvent();
                pressRecognizer.minDurationMillis = Protobuf.newMillis(proto.getMinDuration());
            }
            if (buttonGesture != null && buttonRecognizer != null) {
                PbTouch.ButtonRecognizer proto = buttonGesture.unpack(PbTouch.ButtonRecognizer.class);
                buttonRecognizer.recognizerId = (int)proto.getOnEvent();
            }
        } catch (InvalidProtocolBufferException e) {

        }
    }

    @Override
    public boolean onTouch(View view, MotionEvent event) {
        switch (event.getAction()) {
        case MotionEvent.ACTION_DOWN:
            if (tapRecognizer == null && pressRecognizer == null && buttonRecognizer == null) {
                try {
                    if (tapGesture != null) {
                        PbTouch.TapRecognizer proto = tapGesture.unpack(PbTouch.TapRecognizer.class);
                        tapRecognizer = new TapRecognizer();
                        tapRecognizer.recognizerId = (int)proto.getOnEvent();
                        tapRecognizer.ratio = context.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
                    }
                    if (pressGesture != null) {
                        PbTouch.PressRecognizer proto = pressGesture.unpack(PbTouch.PressRecognizer.class);
                        pressRecognizer = new PressRecognizer();
                        pressRecognizer.recognizerId = (int)proto.getOnEvent();
                        pressRecognizer.ratio = context.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
                        pressRecognizer.minDurationMillis = Protobuf.newMillis(proto.getMinDuration());
                    }
                    if (buttonGesture != null) {
                        PbTouch.ButtonRecognizer proto = buttonGesture.unpack(PbTouch.ButtonRecognizer.class);
                        buttonRecognizer = new ButtonRecognizer();
                        buttonRecognizer.recognizerId = (int)proto.getOnEvent();
                        buttonRecognizer.ratio = context.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
                        buttonRecognizer.width = view.getWidth();
                        buttonRecognizer.height = view.getHeight();
                    }
                } catch (InvalidProtocolBufferException e) {

                }
            }
        case MotionEvent.ACTION_MOVE:
        case MotionEvent.ACTION_UP:
        case MotionEvent.ACTION_OUTSIDE:
        case MotionEvent.ACTION_CANCEL:
            boolean handled = false;
            if (tapRecognizer != null) {
                handled = true;
                Result rlt = tapRecognizer.onEvent(event);
                childView.viewNode.rootView.call(String.format("gomatcha.io/matcha/touch %d", tapRecognizer.recognizerId), childView.viewNode.id, new GoValue(rlt.message.toByteArray()));
                if (rlt.state == State.FAILED || rlt.state == State.RECOGNIZED) {
                    tapRecognizer = null;
                }
            }
            if (pressRecognizer != null) {
                handled = true;
                Result rlt = pressRecognizer.onEvent(event);
                childView.viewNode.rootView.call(String.format("gomatcha.io/matcha/touch %d", pressRecognizer.recognizerId), childView.viewNode.id, new GoValue(rlt.message.toByteArray()));
                if (rlt.state == State.FAILED || rlt.state == State.RECOGNIZED) {
                    pressRecognizer = null;
                }
            }
            if (buttonRecognizer != null) {
                handled = true;
                Result rlt = buttonRecognizer.onEvent(event);
                if (rlt.state == State.POSSIBLE && prevButtonResult != null && prevButtonResult.state == State.POSSIBLE && ((PbTouch.ButtonEvent)prevButtonResult.message).getInside() == ((PbTouch.ButtonEvent)rlt.message).getInside()) {
                    // Skip message.
                } else {
                    childView.viewNode.rootView.call(String.format("gomatcha.io/matcha/touch %d", buttonRecognizer.recognizerId), childView.viewNode.id, new GoValue(rlt.message.toByteArray()));
                }
                prevButtonResult = rlt;
                if (rlt.state == State.FAILED || rlt.state == State.RECOGNIZED) {
                    buttonRecognizer = null;
                }
            }
            if (handled) {
                return true;
            }
        }
        return false;
    }

    class Result {
        State state;
        com.google.protobuf.GeneratedMessageV3 message;
        Result(State state, com.google.protobuf.GeneratedMessageV3 message) {
            this.state = state;
            this.message = message;
        }
    }

    class TapRecognizer {
        long millis;
        float x;
        float y;
        int recognizerId;
        float ratio;

        Result onEvent(MotionEvent event) {
            double distance = Math.sqrt(Math.pow((x - event.getX()), 2) + Math.pow((y - event.getY()), 2)) / ratio;
            double duration = System.currentTimeMillis() - millis;
            PbTouch.TapEvent.Builder e = PbTouch.TapEvent.newBuilder()
                .setTimestamp(Protobuf.toProtobuf(new Date()))
                .setPosition(Protobuf.toProtobuf(new PointF(event.getX() / ratio, event.getY() / ratio)));

            switch (event.getAction()) {
            case MotionEvent.ACTION_DOWN:
                if (event.getPointerCount() > 1) {
                    return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
                }
                millis = System.currentTimeMillis();
                x = event.getX();
                y = event.getY();
                return new Result(State.POSSIBLE, e.setKind(PbTouch.EventKind.EVENT_KIND_POSSIBLE).build());
            case MotionEvent.ACTION_MOVE: {
                if (distance > 10 || duration > 1000*0.75) {
                    return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
                }
                return new Result(State.POSSIBLE, e.setKind(PbTouch.EventKind.EVENT_KIND_POSSIBLE).build());
            }
            case MotionEvent.ACTION_UP:
                if (distance > 10 || duration > 1000*0.75) {
                    return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
                }
                return new Result(State.RECOGNIZED, e.setKind(PbTouch.EventKind.EVENT_KIND_RECOGNIZED).build());
            case MotionEvent.ACTION_OUTSIDE:
            case MotionEvent.ACTION_CANCEL:
                return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
            }
            return null;
        }
    }
    class PressRecognizer {
        long minDurationMillis;
        long millis;
        float x;
        float y;
        int recognizerId;
        float ratio;
        Result onEvent(MotionEvent event) {
            double distance = Math.sqrt(Math.pow((x - event.getX()), 2) + Math.pow((y - event.getY()), 2)) / ratio;
            double duration = System.currentTimeMillis() - millis;
            PbTouch.PressEvent.Builder e = PbTouch.PressEvent.newBuilder()
                    .setTimestamp(Protobuf.toProtobuf(new Date()))
                    .setPosition(Protobuf.toProtobuf(new PointF(event.getX() / ratio, event.getY() / ratio)))
                    .setDuration(Protobuf.toProtobuf((long) duration));

            switch (event.getAction()) {
                case MotionEvent.ACTION_DOWN:
                    if (event.getPointerCount() > 1) {
                        return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
                    }
                    millis = System.currentTimeMillis();
                    x = event.getX();
                    y = event.getY();
                    return new Result(State.POSSIBLE, e.setKind(PbTouch.EventKind.EVENT_KIND_POSSIBLE).setDuration(Protobuf.toProtobuf(0)).build());
                case MotionEvent.ACTION_MOVE:
                    if (duration < minDurationMillis) {
                        if (distance > 10) {
                            return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
                        } else {
                            return new Result(State.POSSIBLE, e.setKind(PbTouch.EventKind.EVENT_KIND_POSSIBLE).build());
                        }
                    }
                    return new Result(State.CHANGED, e.setKind(PbTouch.EventKind.EVENT_KIND_CHANGED).build());
                case MotionEvent.ACTION_UP:
                    if (duration < minDurationMillis) {
                        return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
                    }
                    return new Result(State.RECOGNIZED, e.setKind(PbTouch.EventKind.EVENT_KIND_RECOGNIZED).build());
                case MotionEvent.ACTION_OUTSIDE:
                case MotionEvent.ACTION_CANCEL:
                    return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
            }
            return null;
        }
    }
    class ButtonRecognizer {
        long millis;
        float x;
        float y;
        float width;
        float height;
        int recognizerId;
        float ratio;
        Result onEvent(MotionEvent event) {
            double distance = Math.sqrt(Math.pow((x - event.getX()), 2) + Math.pow((y - event.getY()), 2)) / ratio;
            double duration = System.currentTimeMillis() - millis;
            boolean inside = event.getX() >= 0 && event.getY() >= 0 && event.getX() <= width && event.getY() <= height;
            PbTouch.ButtonEvent.Builder e = PbTouch.ButtonEvent.newBuilder()
                    .setTimestamp(Protobuf.toProtobuf(new Date()))
                    .setInside(inside);
                    //.setPosition(Protobuf.toProtobuf(new PointF(event.getX() / ratio, event.getY() / ratio)));

            switch (event.getAction()) {
                case MotionEvent.ACTION_DOWN:
                    if (event.getPointerCount() > 1) {
                        return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
                    }
                    millis = System.currentTimeMillis();
                    x = event.getX();
                    y = event.getY();
                    return new Result(State.POSSIBLE, e.setKind(PbTouch.EventKind.EVENT_KIND_POSSIBLE).build());
                case MotionEvent.ACTION_MOVE:
                    return new Result(State.POSSIBLE, e.setKind(PbTouch.EventKind.EVENT_KIND_POSSIBLE).build());
                case MotionEvent.ACTION_UP:
                    if (!inside) {
                        return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
                    }
                    return new Result(State.RECOGNIZED, e.setKind(PbTouch.EventKind.EVENT_KIND_RECOGNIZED).build());
                case MotionEvent.ACTION_OUTSIDE:
                case MotionEvent.ACTION_CANCEL:
                    return new Result(State.FAILED, e.setKind(PbTouch.EventKind.EVENT_KIND_FAILED).build());
            }
            return null;
        }
    }
}
