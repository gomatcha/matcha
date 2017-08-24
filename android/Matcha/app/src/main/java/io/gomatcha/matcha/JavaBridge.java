package io.gomatcha.matcha;

import io.gomatcha.bridge.*;

public class JavaBridge {
    static {
        Bridge bridge = Bridge.singleton();
        bridge.put("", new JavaBridge());
    }

    void updateViewWithProtobuf(long id, byte[] protobuf) {
        System.out.format("updateViewWithProtobuf");
    }
}
