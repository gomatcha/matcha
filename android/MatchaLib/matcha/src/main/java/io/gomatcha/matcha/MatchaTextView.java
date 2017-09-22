package io.gomatcha.matcha;

import android.content.Context;
import android.text.SpannableString;
import android.widget.TextView;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.matcha.proto.text.PbText;

class MatchaTextView extends MatchaChildView {
    TextView view;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/textview", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaTextView(context, node);
            }
        });
    }

    public MatchaTextView(Context context, MatchaViewNode node) {
        super(context, node);

        view = new TextView(context);
        addView(view);
    }

    @Override
    public void setNativeState(byte[] nativeState) {
        super.setNativeState(nativeState);
        try {
            PbText.StyledText proto  = PbText.StyledText.parseFrom(nativeState);
            SpannableString str = Protobuf.newAttributedString(proto);
            view.setText(str);
        } catch (InvalidProtocolBufferException e) {
        }
    }
}
