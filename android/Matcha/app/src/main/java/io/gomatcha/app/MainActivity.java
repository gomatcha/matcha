package io.gomatcha.app;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import matcha.MatchaGoValue;
import com.google.protobuf.Any;
import com.google.protobuf.Descriptors;

public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        Descriptors.Descriptor a = Any.getDescriptor();

        MatchaGoValue f = MatchaGoValue.withFunc("TESTFUNC");
        MatchaGoValue[] v3 = f.call("", null);
        System.out.format("%s%n", v3[0].toString());
    }
}
