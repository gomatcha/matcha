package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.graphics.PorterDuff;
import android.view.View;
import android.widget.Button;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.matcha.proto.view.PbButton;

class MatchaButton extends MatchaChildView {
    MatchaViewNode viewNode;
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
        super(context);
        viewNode = node;
        view = new Button(context);
        view.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View view) {
                MatchaButton.this.viewNode.call("OnPress");
            }
        });
        addView(view);
    }

    @Override
    public void setNativeState(byte[] nativeState) {
        super.setNativeState(nativeState);
        try {
            PbButton.Button proto = PbButton.Button.parseFrom(nativeState);
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
