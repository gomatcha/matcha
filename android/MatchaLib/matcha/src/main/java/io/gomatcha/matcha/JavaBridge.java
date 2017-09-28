package io.gomatcha.matcha;

import android.app.AlertDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.res.Resources;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.PointF;
import android.graphics.drawable.Drawable;
import android.net.Uri;
import android.text.SpannableString;
import android.util.DisplayMetrics;
import android.util.Log;
import android.view.Choreographer;
import android.view.Surface;
import android.view.View;
import android.view.WindowManager;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.google.protobuf.InvalidProtocolBufferException;

import java.lang.ref.WeakReference;
import java.nio.ByteBuffer;
import java.util.HashMap;
import java.util.List;

import io.gomatcha.bridge.Bridge;
import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.proto.Proto;
import io.gomatcha.matcha.proto.layout.PbLayout;
import io.gomatcha.matcha.proto.text.PbText;
import io.gomatcha.matcha.proto.view.PbAlert;
import io.gomatcha.matcha.proto.view.PbView;

public class JavaBridge {
    static JavaBridge javaBridge;
    static Choreographer.FrameCallback callback;
    static Context context;
    static TextView textView;
    static HashMap<Long, WeakReference<MatchaView>> viewMap = new HashMap<Long, WeakReference<MatchaView>>();

    static synchronized void init(Context ctx) {
        if (context != null) {
            return;
        }
        context = ctx;
        textView = new TextView(context);
        callback = new Choreographer.FrameCallback() {
            @Override
            public void doFrame(long frameTimeNanos) {
                GoValue.withFunc("gomatcha.io/matcha/animate screenUpdate").call("");
                Choreographer.getInstance().postFrameCallback(callback);
            }
        };
        Choreographer.getInstance().postFrameCallback(callback);
        javaBridge = new JavaBridge();
        Bridge.singleton().put("", javaBridge);
    }

    public boolean updateViewWithProtobuf(Long id, byte[] protobuf) {
        WeakReference<MatchaView> v = viewMap.get(id);
        if (v == null || v.get() == null) {
            viewMap.remove(id);
            return false;
        }
        try {
            v.get().update(PbView.Root.parseFrom(protobuf));
        } catch (InvalidProtocolBufferException e) {
        }
        return true;
    }

    public GoValue sizeForStyledText(byte[] protobuf, Long maxLines) {
        try {
            PbText.SizeFunc sizeFunc = PbText.SizeFunc.parseFrom(protobuf);
            SpannableString str = Protobuf.newAttributedString(sizeFunc.getText());
            PointF minSize = Protobuf.newPoint(sizeFunc.getMinSize());
            PointF maxSize = Protobuf.newPoint(sizeFunc.getMaxSize());

            float ratio = (float)context.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
            int widthMeasureSpec = View.MeasureSpec.makeMeasureSpec((int)(maxSize.x*ratio), View.MeasureSpec.AT_MOST);
            int heightMeasureSpec = View.MeasureSpec.makeMeasureSpec((int)(maxSize.y*ratio), View.MeasureSpec.AT_MOST);

            if (maxLines == 0) {
                maxLines = (long)99999999;
            }

            textView.setText(str);
            textView.setMaxLines(maxLines.intValue());
            textView.measure(widthMeasureSpec, heightMeasureSpec);
            textView.setLayoutParams(new RelativeLayout.LayoutParams(0, 0)); // We need this or setText throws a null pointer exception.

            PointF calculatedSize = new PointF();
            calculatedSize.x = (float)textView.getMeasuredWidth() / ratio + 1;
            calculatedSize.y = (float)textView.getMeasuredHeight() / ratio;
            PbLayout.Point p = Protobuf.toProtobuf(calculatedSize);

            return new GoValue(p.toByteArray());
        } catch (InvalidProtocolBufferException e) {
            Log.v("x", "exception" + e);
            PbLayout.Point p = Protobuf.toProtobuf(new PointF(0, 0));
            return new GoValue(p.toByteArray());
        }
    }

    public GoValue getImageForResource(String path) {
        Resources res = context.getResources();
        int id = res.getIdentifier(path, "drawable", context.getPackageName());
        Bitmap bitmap = BitmapFactory.decodeResource(res, id);

        int size = bitmap.getRowBytes() * bitmap.getHeight();
        ByteBuffer byteBuffer = ByteBuffer.allocate(size);
        bitmap.copyPixelsToBuffer(byteBuffer);
        return new GoValue(byteBuffer.array());
    }

    public GoValue getPropertiesForResource(String path) {
        Resources res = context.getResources();
        int id = res.getIdentifier(path, "drawable", context.getPackageName());

        Drawable drawable = res.getDrawableForDensity(id, DisplayMetrics.DENSITY_MEDIUM);
        Proto.ImageProperties.Builder builder = Proto.ImageProperties.newBuilder()
                .setWidth(drawable.getMinimumWidth())
                .setHeight(drawable.getMinimumHeight())
                .setScale(1); // TODO(KD): Figure out which image density was selected. https://developer.android.com/guide/practices/screens_support.html

        return new GoValue(builder.build().toByteArray());
    }

    public boolean openURL(String url) {
        Intent browserIntent = new Intent(Intent.ACTION_VIEW, Uri.parse("http://www.google.com"));
        context.startActivity(browserIntent);
        return true;
    }

    public int orientation() {
        int rotation = ((WindowManager) context.getSystemService(Context.WINDOW_SERVICE)).getDefaultDisplay().getOrientation();
        switch (rotation) {
        case Surface.ROTATION_0:
            return 0;
        case Surface.ROTATION_90:
            return 3;
        case Surface.ROTATION_180:
            return 1;
        case Surface.ROTATION_270:
            return 2;
        default:
            return 0;
        }
    }

    void didChangeOrientation() {
        Log.v("x", "DidChangeOrientation");
        GoValue.withFunc("gomatcha.io/matcha/application SetOrientation").call("", new GoValue(orientation()));
    }

    public void displayAlert(byte[] protobuf) {
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
