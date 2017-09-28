package io.gomatcha.sampleapp;

import android.content.res.Configuration;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.MatchaView;

public class MainActivity extends AppCompatActivity {
    MatchaView view;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        getSupportActionBar().hide();

        if (view == null) {
            GoValue rootView = GoValue.withFunc("gomatcha.io/matcha/examples/settings New").call("")[0];
            view = new MatchaView(this, rootView);
        }
        setContentView(view);
    }

    static {
        try {
            // Load in example classes.
            Class.forName("io.gomatcha.customview.CustomView");
            Class.forName("io.gomatcha.sampleapp.ExampleJavaBridge");
        } catch (ClassNotFoundException e) {
            throw new RuntimeException(e);
        }
    }
}
