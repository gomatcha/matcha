package io.gomatcha.matcha;

import android.content.Context;
import android.content.res.Resources;
import android.graphics.Bitmap;
import android.graphics.Color;
import android.graphics.Paint;
import android.graphics.PointF;
import android.graphics.Typeface;
import android.graphics.drawable.BitmapDrawable;
import android.graphics.drawable.Drawable;
import android.text.Layout;
import android.text.SpannableString;
import android.text.SpannableStringBuilder;
import android.text.style.AbsoluteSizeSpan;
import android.text.style.AlignmentSpan;
import android.text.style.ForegroundColorSpan;
import android.text.style.LineHeightSpan;
import android.text.style.StrikethroughSpan;
import android.text.style.StyleSpan;
import android.text.style.TypefaceSpan;
import android.text.style.UnderlineSpan;
import android.util.Log;

import com.google.protobuf.Duration;
import com.google.protobuf.Timestamp;

import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

import io.gomatcha.matcha.pb.Pb;
import io.gomatcha.matcha.pb.keyboard.PbKeyboard;
import io.gomatcha.matcha.pb.layout.PbLayout;
import io.gomatcha.matcha.pb.text.PbText;

import static android.text.Spanned.SPAN_INCLUSIVE_EXCLUSIVE;
import static android.text.Spanned.SPAN_INCLUSIVE_INCLUSIVE;

public class Protobuf {
    public static long newMillis(Duration d) {
        return d.getSeconds() * 1000 + d.getNanos() / 1000000;
    }

    public static Duration toProtobuf(long millis) {
        return Duration.newBuilder()
                .setNanos((int)(millis % 1000 * 1000000))
                .setSeconds(millis/1000)
                .build();
    }

    public static Date newDate(Timestamp t) {
        return new Date(t.getSeconds() * 1000 + t.getNanos() / 1000000);
    }

    public static Timestamp toProtobuf(Date d) {
        long millis = d.getTime();
        return Timestamp.newBuilder().setSeconds(millis/1000).setNanos((int)(millis % 1000 * 1000000)).build();
    }

    public static Drawable newDrawable(Pb.ImageOrResource res, Context ctx) {
        if (res.hasImage()) {
            Bitmap bitmap = Protobuf.newBitmap(res.getImage());
            if (bitmap != null) {
                return new BitmapDrawable(ctx.getResources(), bitmap);
            }
            return null;
        } else {
            Resources resources = ctx.getResources();
            int id = resources.getIdentifier(res.getPath(), "drawable", ctx.getPackageName());
            return resources.getDrawable(id);
        }
    }

    public static Bitmap newBitmap(Pb.Image image) {
        byte[] buf2 = new byte[(int)image.getWidth()*(int)image.getHeight()*4];
        image.getData().copyTo(buf2, 0);
        ByteBuffer buf = ByteBuffer.wrap(buf2);
        // TODO(KD): Why doesn't the below work?
        // ByteBuffer buf = image.getData().asReadOnlyByteBuffer();

        Bitmap bitmap = Bitmap.createBitmap((int)image.getWidth(), (int)image.getHeight(), Bitmap.Config.ARGB_8888);
        bitmap.copyPixelsFromBuffer(buf);
        return bitmap;
    }
    
    public static Bitmap newBitmap(Pb.ImageOrResource image) {
        if (!image.hasImage()) {
            return null;
        }
        return newBitmap(image.getImage());
    }

    public static int newColor(Pb.Color c) {
        return Color.argb(c.getAlpha()*255/65535, c.getRed()*255/65535, c.getGreen()*255/65535, c.getBlue()*255/65535);
    }
    
    public static PointF newPoint(PbLayout.Point pt) {
        return new PointF((float)pt.getX(), (float)pt.getY());
    }
    
    public static PbLayout.Point toProtobuf(PointF pt) {
        PbLayout.Point.Builder builder = PbLayout.Point.newBuilder();
        builder.setX(pt.x);
        builder.setY(pt.y);
        return builder.build();
    }
    
    public static SpannableString newAttributedString(PbText.StyledText st) {
        SpannableString str = new SpannableString(st.getText().getText());
        List<PbText.TextStyle> styles = st.getStylesList();
        for (int i = 0; i < styles.size(); i++) {
            PbText.TextStyle style = styles.get(i);
            int start = (int)style.getIndex();
            int end = 0;
            if (i < styles.size() - 1) {
                end = (int)styles.get(i+1).getIndex();
            } else {
                end = str.length();
            }
            ArrayList<Object> spans = newSpanArrayList(style);
            for (Object j : spans) {
                str.setSpan(j, start, end, SPAN_INCLUSIVE_INCLUSIVE);
            }
        }
        return str;
    }

    public static PbText.StyledText toProtobuf(SpannableStringBuilder str) {
        PbText.Text.Builder textBuilder = PbText.Text.newBuilder().setText(str.toString());

        PbText.StyledText.Builder builder = PbText.StyledText.newBuilder();
        builder.setText(textBuilder.build());
        return builder.build();
    }
    
    public static ArrayList<Object> newSpanArrayList(PbText.TextStyle textStyle) {
        ArrayList<Object> arrayList = new ArrayList<Object>();
        
        Object span;
        switch (textStyle.getTextAlignment()) {
            case TEXT_ALIGNMENT_LEFT:
                span = new AlignmentSpan.Standard(Layout.Alignment.ALIGN_NORMAL);
                break;
            case TEXT_ALIGNMENT_RIGHT:
                span = new AlignmentSpan.Standard(Layout.Alignment.ALIGN_OPPOSITE);
                break;
            case TEXT_ALIGNMENT_CENTER:
                span = new AlignmentSpan.Standard(Layout.Alignment.ALIGN_CENTER);
                break;
            case TEXT_ALIGNMENT_JUSTIFIED:
                span = new AlignmentSpan.Standard(Layout.Alignment.ALIGN_NORMAL);
                break;
            default:
                span = new AlignmentSpan.Standard(Layout.Alignment.ALIGN_NORMAL);
                break;
        }
        if (span != null) {
            arrayList.add(span);
        }
        
        span = null;
        switch (textStyle.getStrikethroughStyle()) {
            case STRIKETHROUGH_STYLE_NONE:
                break;
            case STRIKETHROUGH_STYLE_SINGLE:
            case STRIKETHROUGH_STYLE_DOUBLE:
            case STRIKETHROUGH_STYLE_THICK:
            case STRIKETHROUGH_STYLE_DOTTED:
            case STRIKETHROUGH_STYLE_DASHED:
            default:
                span = new StrikethroughSpan();
                break;
        }
        if (span != null) {
            arrayList.add(span);
        }
        
        // TODO(KD): Strikethrough color...
        // TODO(KD): Underline color...
        
        span = null;
        switch (textStyle.getUnderlineStyle()) {
            case UNDRELINE_STYLE_NONE:
                break;
            case UNDRELINE_STYLE_SINGLE:
            case UNDRELINE_STYLE_DOUBLE:
            case UNDRELINE_STYLE_THICK:
            case UNDRELINE_STYLE_DOTTED:
            case UNDRELINE_STYLE_DASHED:
            default:
                span = new UnderlineSpan();
                break;
        }
        if (span != null) {
            arrayList.add(span);
        }

        PbText.Font font = textStyle.getFont();
        String fontName = font.getFamily();
        if (fontName.endsWith("-bold")) {
            fontName = fontName.substring(0, fontName.length() - 5);
            arrayList.add(new StyleSpan(Typeface.BOLD));
        } else if (fontName.endsWith("-italic")) {
            fontName = fontName.substring(0, fontName.length() - 7);
            arrayList.add(new StyleSpan(Typeface.ITALIC));
        } else if (fontName.endsWith("-bolditalic")) {
            fontName = fontName.substring(0, fontName.length() - 11);
            arrayList.add(new StyleSpan(Typeface.BOLD_ITALIC));
        }
        span = new TypefaceSpan(font.getFamily());
        arrayList.add(span);

        span = new AbsoluteSizeSpan((int)font.getSize(), true);
        arrayList.add(span);

        // span = new LineHeightSpan();

        int color = newColor(textStyle.getTextColor());
        span = new ForegroundColorSpan(color);
        arrayList.add(span);

        // TODO(KD): AttributeKeyTextWrap
        // TODO(KD): AttributeKeyTruncation
        // TODO(KD): AttributeKeyTruncationString

        return arrayList;
    }
}
