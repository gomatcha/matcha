package io.gomatcha.matcha;

import android.content.Context;
import android.widget.ImageView;
import android.widget.RelativeLayout;

import com.google.protobuf.Any;
import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.imageview.PbImageView;

public class MatchaImageView extends MatchaChildView {
    ImageView view;
    
    static {
        MatchaView.registerView("gomatcha.io/matcha/view/imageview", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaImageView(context, node);
            }
        });
    }
    
    public MatchaImageView(Context context, MatchaViewNode node) {
        super(context, node);

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
        view = new ImageView(context);
        addView(view, params);
    }
    
    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbImageView.View proto = buildNode.getBridgeValue().unpack(PbImageView.View.class);
        } catch (InvalidProtocolBufferException e) {
        }
    }
}
