package io.gomatcha.matcha;

import android.content.Context;
import android.util.Log;
import android.widget.RelativeLayout;
import android.widget.SeekBar;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.slider.PbSlider;

public class MatchaSlider extends MatchaChildView {
    SeekBar view;
    double value;
    double maxValue;
    double minValue;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/slider", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaSlider(context, node);
            }
        });
    }

    public MatchaSlider(Context context, MatchaViewNode node) {
        super(context, node);

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
        view = new SeekBar(context);
        view.setMax(10000);
        view.setOnSeekBarChangeListener(new SeekBar.OnSeekBarChangeListener() {
            @Override
            public void onProgressChanged(SeekBar seekBar, int i, boolean b) {
                if (b && i != value ) {
                    double maxValue = MatchaSlider.this.maxValue;
                    double minValue = MatchaSlider.this.minValue;
                    PbSlider.Event proto = PbSlider.Event.newBuilder().setValue((double) i / 10000.0 * (maxValue - minValue) + minValue).build();
                    MatchaSlider.this.viewNode.rootView.call("OnValueChange", MatchaSlider.this.viewNode.id, new GoValue(proto.toByteArray()));
                }
            }
            @Override
            public void onStartTrackingTouch(SeekBar seekBar) {

            }
            @Override
            public void onStopTrackingTouch(SeekBar seekBar) {
                double maxValue = MatchaSlider.this.maxValue;
                double minValue = MatchaSlider.this.minValue;
                PbSlider.Event proto = PbSlider.Event.newBuilder().setValue((double)MatchaSlider.this.view.getProgress()/10000.0 * (maxValue - minValue) + minValue).build();
                MatchaSlider.this.viewNode.rootView.call("OnSubmit", MatchaSlider.this.viewNode.id, new GoValue(proto.toByteArray()));
            }
        });
        addView(view, params);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbSlider.View proto = buildNode.getBridgeValue().unpack(PbSlider.View.class);
            view.setEnabled(proto.getEnabled());
            view.setProgress((int)((proto.getValue()- proto.getMinValue())*10000.0/(proto.getMaxValue() - proto.getMinValue())));
            this.value = view.getProgress();
            this.maxValue = proto.getMaxValue();
            this.minValue = proto.getMinValue();
            Log.v("x", "+" + this.value);
        } catch (InvalidProtocolBufferException e) {
        }
    }
}
