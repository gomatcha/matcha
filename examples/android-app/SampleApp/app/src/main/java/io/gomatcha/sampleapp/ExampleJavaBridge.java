package io.gomatcha.sampleapp;

import io.gomatcha.bridge.Bridge;
import io.gomatcha.bridge.GoValue;

public class ExampleJavaBridge {
    static {
        ExampleJavaBridge b = new ExampleJavaBridge();
        Bridge.singleton().put("gomatcha.io/matcha/example/bridge", b);
    }

    public String callWithForeignValues(String param) {
        return String.format("Hello %s", param);
    }

    public GoValue callWithGoValues(GoValue param) {
        return new GoValue(String.format("Hello %s", param.toString()));
    }

    public String callGoFunctionWithForeignValues() {
        GoValue func = GoValue.withFunc("gomatcha.io/matcha/examples/bridge callWithForeignValues");
        return (String)func.call("", GoValue.WithObject("Ame"))[0].toObject();
    }

    public String callGoFunctionWithGoValues() {
        GoValue func = GoValue.withFunc("gomatcha.io/matcha/examples/bridge callWithGoValues");
        return func.call("", GoValue.WithString("Yuki"))[0].toString();
    }
}
