package io.gomatcha.customview;

import android.content.Context;
import android.support.v7.widget.SwitchCompat;
import android.util.DisplayMetrics;
import android.widget.CompoundButton;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.MatchaChildView;
import io.gomatcha.matcha.MatchaView;
import io.gomatcha.matcha.MatchaViewNode;
import io.gomatcha.matcha.proto.view.PbSwitchView;
import io.gomatcha.matcha.proto.view.PbView;

public class CustomView extends MatchaChildView {
    SwitchCompat view;
    boolean checked;

    static {
        MatchaView.registerView("github.com/overcyn/customview", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new CustomView(context, node);
            }
        });
    }

    public CustomView(Context context, MatchaViewNode node) {
        super(context, node);

        float ratio = (float)context.getResources().getDisplayMetrics().densityDpi / DisplayMetrics.DENSITY_DEFAULT;
        view = new SwitchCompat(context);
        view.setPadding(0, 0, (int)(7*ratio), 0);
        view.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
            public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
                if (isChecked != checked) {
                    checked = isChecked;
                    PbSwitchView.SwitchEvent event = PbSwitchView.SwitchEvent.newBuilder().setValue(isChecked).build();
                    CustomView.this.viewNode.rootView.call("OnChange", CustomView.this.viewNode.id, new GoValue(event.toByteArray()));
                }
            }
        });
        addView(view);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            PbSwitchView.SwitchView proto = buildNode.getBridgeValue().unpack(PbSwitchView.SwitchView.class);
            checked = proto.getValue();
            view.setChecked(proto.getValue());
            view.setEnabled(proto.getEnabled());
        } catch (InvalidProtocolBufferException e) {
        }
    }
}
