package io.gomatcha.matcha;

import android.graphics.Bitmap;
import android.graphics.Color;
import android.util.Log;

import java.nio.ByteBuffer;

import io.gomatcha.matcha.pb.Pb;

public class Protobuf {
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
}
