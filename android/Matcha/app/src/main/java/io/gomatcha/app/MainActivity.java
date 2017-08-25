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
        
        //GoValue f = GoValue.withFunc("TESTFUNC");
        //GoValue[] v3 = f.call("", null);
        //System.out.format("wubalubadub %s%n", v3[0].toString());

        GoValue rootVC = GoValue.withFunc("gomatcha.io/matcha/examples/constraints New").call("", null)[0];

        MatchaView v = new MatchaView(this, rootVC);
        setContentView(v);
    }

    
    // @Override
    // public View onCreateView(LayoutInflater inflater, ViewGroup parent, Bundle savedInstanceState) {
    //     MatchaView v = new MatchaView(parent);
    //     return v;
    //     // getActivity().getActionBar().setDisplayHomeAsUpEnabled(true);
    //     // View rootView = inflater.inflate(R.layout.fragment_single_image, parent, false);
    //     // ImageView imageView = (ImageView)rootView.findViewById(R.id.currentImage);
    //     // imageView.setImageBitmap(currentImage);
    //     // return rootView;
    // }
}
