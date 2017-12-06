package io.gomatcha.cameraview;

import android.content.Context;
import android.graphics.Color;
import android.graphics.Paint;
import android.graphics.drawable.ShapeDrawable;
import android.graphics.drawable.shapes.OvalShape;
import android.view.View;

public class CameraButton extends View {
    public CameraButton(Context context) {
        super(context);

        OvalShape oval = new OvalShape();
        oval.resize(50 , 50);

        ShapeDrawable drawable = new ShapeDrawable(oval);
        drawable.getPaint().setColor(Color.WHITE);
        drawable.getPaint().setStyle(Paint.Style.FILL);
        drawable.getPaint().setAntiAlias(true);
        drawable.getPaint().setFlags(Paint.ANTI_ALIAS_FLAG);
        setBackground(drawable);
    }
}
