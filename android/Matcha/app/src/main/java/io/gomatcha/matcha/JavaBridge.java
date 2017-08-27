package io.gomatcha.matcha;

import android.content.Context;
import android.content.res.Resources;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.drawable.Drawable;
import android.util.Log;
import android.view.Choreographer;

import com.google.protobuf.InvalidProtocolBufferException;

import java.lang.ref.WeakReference;
import java.nio.ByteBuffer;

import io.gomatcha.app.R;
import io.gomatcha.bridge.*;
import io.gomatcha.matcha.pb.Pb;
import io.gomatcha.matcha.pb.view.PbView;

public class JavaBridge {
    static Choreographer.FrameCallback callback;
    static Context context;

    static synchronized void init(Context ctx) {
        if (context != null) {
            return;
        }
        context = ctx;

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

    GoValue getImageForResource(String path) {
        Resources res = context.getResources();
        int id = res.getIdentifier(path, "drawable", context.getPackageName());
        Bitmap bitmap = BitmapFactory.decodeResource(res, id);

        int size = bitmap.getRowBytes() * bitmap.getHeight();
        ByteBuffer byteBuffer = ByteBuffer.allocate(size);
        bitmap.copyPixelsToBuffer(byteBuffer);
        return new GoValue(byteBuffer.array());
    }

    GoValue getPropertiesForResource(String path) {
        Resources res = context.getResources();
        int id = res.getIdentifier(path, "drawable", context.getPackageName());

        BitmapFactory.Options dimensions = new BitmapFactory.Options();
        dimensions.inJustDecodeBounds = true;
        BitmapFactory.decodeResource(res, id, dimensions);
        int height = dimensions.outHeight;
        int width =  dimensions.outWidth;

        Pb.ImageProperties.Builder builder = Pb.ImageProperties.newBuilder();
        builder.setWidth(dimensions.outWidth);
        builder.setHeight(dimensions.outHeight);
        builder.setScale(1); // TODO(KD): Figure out which image density was selected. https://developer.android.com/guide/practices/screens_support.html

        return new GoValue(builder.build().toByteArray());
    }
}
