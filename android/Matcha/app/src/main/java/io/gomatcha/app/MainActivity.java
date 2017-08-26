package io.gomatcha.app;

import android.app.Fragment;
import android.app.FragmentTransaction;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.FrameLayout;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.MatchaView;

import com.google.protobuf.Any;
import com.google.protobuf.Descriptors;

public class MainActivity extends AppCompatActivity {
    
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        
        GoValue rootVC = GoValue.withFunc("gomatcha.io/matcha/examples/constraints New").call("")[0];

        MatchaView v = new MatchaView(this, rootVC);
        setContentView(v);
    }
}
