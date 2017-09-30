package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.graphics.Typeface;
import android.graphics.drawable.ColorDrawable;
import android.graphics.drawable.Drawable;
import android.util.Log;
import android.view.Gravity;
import android.widget.TextView;

class MatchaUnknownView extends MatchaChildView {
    public MatchaUnknownView(Context context, MatchaViewNode node) {
        super(context);

        TextView tv = new TextView(context);
        tv.setText("Unknown View");
        tv.setTextColor(Color.WHITE);
        tv.setTypeface(null, Typeface.BOLD);
        tv.setTextSize(13);
        tv.setGravity(Gravity.CENTER);
        addView(tv);
    }

    @Override
    public void setBackground(Drawable d) {
        super.setBackground(new ColorDrawable(Color.RED));
    }
}
