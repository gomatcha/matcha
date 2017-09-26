package io.gomatcha.matcha;

import android.content.Context;
import android.util.Log;
import android.widget.ImageView;

import com.google.protobuf.InvalidProtocolBufferException;
import com.makeramen.roundedimageview.RoundedImageView;

import io.gomatcha.matcha.proto.view.PbImageView;

class MatchaImageView extends MatchaChildView {
    RoundedImageView view;
    MatchaViewNode viewNode;
    
    static {
        MatchaView.registerView("gomatcha.io/matcha/view/imageview", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaImageView(context, node);
            }
        });
    }
    
    public MatchaImageView(Context context, MatchaViewNode node) {
        super(context);
        viewNode = node;

        view = new RoundedImageView(context);
        addView(view);
    }

    @Override
    public void setNativeState(byte[] nativeState) {
        super.setNativeState(nativeState);
        try {
            PbImageView.ImageView proto  = PbImageView.ImageView.parseFrom(nativeState);
            view.setImageDrawable(Protobuf.newDrawable(proto.getImage(), getContext()));

            switch (proto.getResizeMode()) {
                case FIT:
                    view.setScaleType(ImageView.ScaleType.FIT_XY);
                    view.setAdjustViewBounds(true);
                    break;
                case FILL:
                    view.setScaleType(ImageView.ScaleType.FIT_XY); // TODO(KD): not correct...
                    view.setAdjustViewBounds(true);
                    break;
                case STRETCH:
                    view.setScaleType(ImageView.ScaleType.FIT_XY);
                    view.setAdjustViewBounds(false);
                    break;
                case CENTER:
                    view.setScaleType(ImageView.ScaleType.CENTER);
                    view.setAdjustViewBounds(false);
                    break;
                case UNRECOGNIZED:
                    view.setScaleType(ImageView.ScaleType.FIT_XY);
                    view.setAdjustViewBounds(true);
                    break;
            }

            if (proto.hasTint()) {
                view.setColorFilter(Protobuf.newColor(proto.getTint()));
            }
        } catch (InvalidProtocolBufferException e) {
        }
    }
}
