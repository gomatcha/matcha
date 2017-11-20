package io.gomatcha.matchaexample;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.Matcha;
import io.gomatcha.matcha.MatchaView;

public class MainActivity extends AppCompatActivity {
    MatchaView view;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        getSupportActionBar().hide();

        if (view == null) {
            Matcha.configure(this);

            GoValue rootView = GoValue.withFunc("gomatcha.io/matcha/examples NewExamplesView").call("")[0];
            view = new MatchaView(this, rootView);
        }
        setContentView(view);
    }

    static {
        try {
            // Load in example classes.
            Class.forName("io.gomatcha.cameraview.CameraView");
            Class.forName("io.gomatcha.sampleapp.ExampleJavaBridge");
        } catch (ClassNotFoundException e) {
            throw new RuntimeException(e);
        }
    }
}
