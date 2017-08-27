package io.gomatcha.matcha;

import android.content.Context;
import android.content.res.Resources;
import android.graphics.Bitmap;
import android.util.Log;
import android.widget.ImageView;
import android.widget.RelativeLayout;

import com.google.protobuf.Any;
import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.matcha.pb.Pb;
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

            Pb.ImageOrResource imageOrResource = proto.getImage();
            if (imageOrResource.hasImage()) {
                Bitmap bitmap = Protobuf.newBitmap(imageOrResource.getImage());
                if (bitmap != null) {
                    view.setImageBitmap(bitmap);
                }
            } else {
                Resources res = this.getResources();
                int id = res.getIdentifier(imageOrResource.getPath(), "drawable", getContext().getPackageName());
                view.setImageResource(id);
            }

            Log.v("resize", "" + proto.getResizeMode());
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
