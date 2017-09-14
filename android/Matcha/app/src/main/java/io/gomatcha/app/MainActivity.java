package io.gomatcha.app;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.WindowManager;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.MatchaView;

public class MainActivity extends AppCompatActivity {
    
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        this.getSupportActionBar().hide();
        GoValue rootVC = GoValue.withFunc("gomatcha.io/matcha/examples/settings New").call("")[0];

        MatchaView v = new MatchaView(this, rootVC);
        setContentView(v);
    }
}
