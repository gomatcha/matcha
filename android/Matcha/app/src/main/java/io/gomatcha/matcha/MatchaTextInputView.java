package io.gomatcha.matcha;

import android.content.Context;
import android.text.SpannableString;
import android.widget.EditText;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.matcha.pb.text.PbText;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.textinput.PbTextInput;

public class MatchaTextInputView extends MatchaChildView {
    EditText view;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/textinput", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaTextView(context, node);
            }
        });
    }

    public MatchaTextInputView(Context context, MatchaViewNode node) {
        super(context, node);

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
        view = new EditText(context);
        addView(view, params);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbTextInput.View proto = buildNode.getBridgeValue().unpack(PbTextInput.View.class);
        } catch (InvalidProtocolBufferException e) {
        }
    }
}
