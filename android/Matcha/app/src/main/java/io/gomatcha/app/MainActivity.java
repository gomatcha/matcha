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
        // getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_STATE_ALWAYS_VISIBLE);
        
        GoValue rootVC = GoValue.withFunc("gomatcha.io/matcha/examples/view NewTextView").call("")[0];

        MatchaView v = new MatchaView(this, rootVC);
        setContentView(v);
    }
}
