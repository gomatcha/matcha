package io.gomatcha.matcha;

import android.content.Context;
import android.text.SpannableString;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.matcha.pb.text.PbText;
import io.gomatcha.matcha.pb.view.PbView;

public class MatchaTextView extends MatchaChildView {
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

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
        view = new TextView(context);
        addView(view, params);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbText.StyledText proto = buildNode.getBridgeValue().unpack(PbText.StyledText.class);
            SpannableString str = Protobuf.newAttributedString(proto);
            view.setText(str);
        } catch (InvalidProtocolBufferException e) {
        }
    }
}
