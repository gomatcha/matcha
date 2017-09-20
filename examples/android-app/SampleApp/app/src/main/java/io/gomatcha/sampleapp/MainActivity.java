package io.gomatcha.sampleapp;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.MatchaView;

public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        getSupportActionBar().hide();

        GoValue rootView = GoValue.withFunc("gomatcha.io/matcha/examples/android NewStatusBarView").call("")[0];
        setContentView(new MatchaView(this, rootView));
    }
}
