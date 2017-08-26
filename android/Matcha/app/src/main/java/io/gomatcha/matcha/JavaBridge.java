package io.gomatcha.matcha;

import android.util.Log;
import android.view.Choreographer;

import com.google.protobuf.InvalidProtocolBufferException;

import java.lang.ref.WeakReference;

import io.gomatcha.bridge.*;
import io.gomatcha.matcha.pb.view.PbView;

public class JavaBridge {
    static Choreographer.FrameCallback callback;

    static {
        Bridge bridge = Bridge.singleton();
        bridge.put("", new JavaBridge());

        callback = new Choreographer.FrameCallback() {
            @Override
            public void doFrame(long frameTimeNanos) {
                GoValue.withFunc("gomatcha.io/matcha/animate screenUpdate").call("");
                Choreographer.getInstance().postFrameCallback(callback);
            }
        };
        Choreographer.getInstance().postFrameCallback(callback);
    }

    void updateViewWithProtobuf(Long id, byte[] protobuf) {
        Log.v("X", "updateViewWithProtobuf");
        for (WeakReference<MatchaView> i : MatchaView.views) {
            if (i.get().identifier == id) {
                try {
                    PbView.Root root = PbView.Root.parseFrom(protobuf);
                    i.get().update(root);
                } catch (InvalidProtocolBufferException e) {

                }
            }
        }
    }
}
