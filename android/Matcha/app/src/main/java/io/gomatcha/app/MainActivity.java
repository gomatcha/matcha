package io.gomatcha.app;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import io.gomatcha.bridge.GoValue;
import com.google.protobuf.Any;
import com.google.protobuf.Descriptors;

public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        Descriptors.Descriptor a = Any.getDescriptor();

        GoValue f = GoValue.withFunc("TESTFUNC");
        GoValue[] v3 = f.call("", null);
        System.out.format("wubalubadub %s%n", v3[0].toString());
    }
}
