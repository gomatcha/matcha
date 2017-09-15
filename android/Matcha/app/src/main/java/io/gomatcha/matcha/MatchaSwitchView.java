package io.gomatcha.matcha;

import android.content.Context;
import android.graphics.Color;
import android.util.Log;
import android.view.View;
import android.widget.CompoundButton;
import android.widget.RelativeLayout;
import android.widget.Switch;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.PbSwitchView;

public class MatchaSwitchView extends MatchaChildView {
    Switch view;
    boolean checked;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/switch", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaSwitchView(context, node);
            }
        });
    }

    public MatchaSwitchView(Context context, MatchaViewNode node) {
        super(context, node);

        view = new Switch(context);
        view.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
            public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
                if (isChecked != checked) {
                    checked = isChecked;
                    PbSwitchView.SwitchEvent event = PbSwitchView.SwitchEvent.newBuilder().setValue(isChecked).build();
                    MatchaSwitchView.this.viewNode.rootView.call("OnChange", MatchaSwitchView.this.viewNode.id, new GoValue(event.toByteArray()));
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
