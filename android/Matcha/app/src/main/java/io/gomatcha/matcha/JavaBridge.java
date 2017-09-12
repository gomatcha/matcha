package io.gomatcha.matcha;

import android.app.AlertDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.content.res.Resources;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.PointF;
import android.graphics.drawable.Drawable;
import android.text.Layout;
import android.text.SpannableString;
import android.text.StaticLayout;
import android.text.TextPaint;
import android.util.Log;
import android.view.Choreographer;

import com.google.protobuf.InvalidProtocolBufferException;

import java.lang.ref.WeakReference;
import java.nio.ByteBuffer;
import java.util.List;

import io.gomatcha.app.R;
import io.gomatcha.bridge.*;
import io.gomatcha.matcha.pb.Pb;
import io.gomatcha.matcha.pb.layout.PbLayout;
import io.gomatcha.matcha.pb.text.PbText;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.PbAlert;

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

    GoValue sizeForStyledText(byte[] protobuf, Long maxLines) {
        try {
            PbText.SizeFunc sizeFunc = PbText.SizeFunc.parseFrom(protobuf);
            SpannableString str = Protobuf.newAttributedString(sizeFunc.getText());
            PointF minSize = Protobuf.newPoint(sizeFunc.getMinSize());
            PointF maxSize = Protobuf.newPoint(sizeFunc.getMaxSize());

            TextPaint textPaint = new TextPaint();
            // textPaint.getTextBounds(someText, 0, someText.length(), bounds);
            StaticLayout layout = new StaticLayout(str, textPaint, 1000, Layout.Alignment.ALIGN_NORMAL, 0, 0, false);
            int height = layout.getHeight();
            int width = layout.getWidth();
            int lines = layout.getLineCount();

            float maxWidth = 0;
            for (int i = 0; i < lines; i++) {
                float lineWidth = layout.getLineWidth(i);
                if (lineWidth > maxWidth) {
                    maxWidth = lineWidth;
                }
            }

            PbLayout.Point p = Protobuf.toProtobuf(new PointF(maxWidth+5f, height)); // TODO(KD): Why am I adding 5f?
            return new GoValue(p.toByteArray());
        } catch (InvalidProtocolBufferException e) {
            PbLayout.Point p = Protobuf.toProtobuf(new PointF(0, 0));
            return new GoValue(p.toByteArray());
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
    
    void displayAlert(byte[] protobuf) {
        try {
            final PbAlert.Alert alert = PbAlert.Alert.parseFrom(protobuf);
        
            AlertDialog.Builder builder = new AlertDialog.Builder(context);
            builder.setTitle(alert.getTitle());
            if (alert.getMessage().length() > 0) {
                builder.setMessage(alert.getMessage());
            }

            List<PbAlert.AlertButton> buttons = alert.getButtonsList();
            if (buttons.size() == 0) {
                builder.setPositiveButton("OK", new DialogInterface.OnClickListener() {
                   public void onClick(DialogInterface dialog, int id) {
                       // no-op?
                   }
                });
            }
            if (buttons.size() > 0) {
                builder.setPositiveButton(buttons.get(0).getTitle(), new DialogInterface.OnClickListener() {
                   public void onClick(DialogInterface dialog, int id) {
                       GoValue.withFunc("gomatcha.io/matcha/view/alert onPress").call("", new GoValue(alert.getId()), new GoValue(0));
                   }
                });
            }
            if (buttons.size() > 1) {
                builder.setNegativeButton(buttons.get(1).getTitle(), new DialogInterface.OnClickListener() {
                   public void onClick(DialogInterface dialog, int id) {
                       GoValue.withFunc("gomatcha.io/matcha/view/alert onPress").call("", new GoValue(alert.getId()), new GoValue(1));
                   }
                });
            }
            if (buttons.size() > 2) {
                builder.setNeutralButton(buttons.get(2).getTitle(), new DialogInterface.OnClickListener() {
                   public void onClick(DialogInterface dialog, int id) {
                       GoValue.withFunc("gomatcha.io/matcha/view/alert onPress").call("", new GoValue(alert.getId()), new GoValue(2));
                   }
                });
            }
            builder.setCancelable(false);
            builder.show();
        } catch (InvalidProtocolBufferException e) {
        }
    }
}
