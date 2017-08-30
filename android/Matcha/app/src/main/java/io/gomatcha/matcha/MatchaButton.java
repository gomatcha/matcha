package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.graphics.PorterDuff;
import android.text.SpannableString;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.RelativeLayout;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.app.R;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.alert.PbAlert;
import io.gomatcha.matcha.pb.view.button.PbButton;

import static android.R.color.primary_text_dark;
import static android.R.color.primary_text_light;
import static android.R.color.secondary_text_dark;
import static android.R.color.secondary_text_dark_nodisable;
import static android.R.color.secondary_text_light;

public class MatchaButton extends MatchaChildView {
    Button view;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/button", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaButton(context, node);
            }
        });
    }

    public MatchaButton(Context context, MatchaViewNode node) {
        super(context, node);

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
        view = new Button(context);
        view.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View view) {
                MatchaButton.this.viewNode.rootView.call("OnPress", MatchaButton.this.viewNode.id);
            }
        });
        addView(view, params);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbButton.View proto = buildNode.getBridgeValue().unpack(PbButton.View.class);
            view.setEnabled(proto.getEnabled());
            view.setText(proto.getStr());

            if (proto.hasColor() && proto.getEnabled()) {
                int color = Protobuf.newColor(proto.getColor());
                view.getBackground().setColorFilter(color, PorterDuff.Mode.MULTIPLY);
                view.setTextColor(Color.WHITE);
            } else {
                // TODO(KD): reset background and text color
            }

        } catch (InvalidProtocolBufferException e) {
        }
    }
}
