package io.gomatcha.sampleapp;

import io.gomatcha.bridge.Bridge;
import io.gomatcha.bridge.GoValue;

public class ExampleJavaBridge {
    static {
        ExampleJavaBridge b = new ExampleJavaBridge();
        Bridge.singleton().put("gomatcha.io/matcha/example", b);

        GoValue func1 = GoValue.withFunc("gomatcha.io/matcha/examples/bridge callWithGoValues");
        String str1 = func1.call("", new GoValue(123))[0].toString();

        GoValue func2 = GoValue.withFunc("gomatcha.io/matcha/examples/bridge callWithForeignValues");
        String str2 = (String)func2.call("", new GoValue(Long.valueOf(123)))[0].toObject();

    }

    public String callWithForeignValues(Long param) {
        return String.format("Hello %d", param);
    }

    public GoValue callWithGoValues(GoValue param) {
        return new GoValue(String.format("Hello %d", param.toLong()));
    }
}
