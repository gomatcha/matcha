package io.gomatcha.cameraview;

import android.Manifest;
import android.content.Context;
import android.content.pm.PackageManager;
import android.hardware.Camera;
import android.support.v4.app.ActivityCompat;
import android.support.v4.content.ContextCompat;
import android.support.v7.widget.SwitchCompat;
import android.util.Log;
import android.widget.SeekBar;

import io.gomatcha.matcha.MatchaChildView;
import io.gomatcha.matcha.MatchaView;
import io.gomatcha.matcha.MatchaViewNode;

public class CameraView extends MatchaChildView {
    Camera camera;
    CameraPreview preview;
    MatchaViewNode viewNode;

    static {
        MatchaView.registerView("gomatcha.io/matcha/examples/cameraview CameraView", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new CameraView(context, node);
            }
        });
    }

    public CameraView(Context context, MatchaViewNode node) {
        super(context);
        this.setClipChildren(true);
        viewNode = node;

        camera = this.getCameraInstance();
        preview = new CameraPreview(context, camera);
        addView(preview);
        
        // addView(new SeekBar(context));
    }

    @Override
    public void setNativeState(byte[] nativeState) {
        super.setNativeState(nativeState);
    }
    
    public Camera getCameraInstance() {
        //String[] permissions = {"android.permission.CAMERA"};
        //if (ContextCompat.checkSelfPermission(this.getContext(), Manifest.permission.CAMERA) != PackageManager.PERMISSION_GRANTED) {
        //    ActivityCompat.requestPermissions(this, permissions, MY_PERMISSIONS_REQUEST_CAMERA);
        //}


        Camera c = null;
        try {
            c = Camera.open(); // attempt to get a Camera instance
        }
        catch (Exception e){
            Log.v("x", e.toString());
            // Camera is not available (in use or does not exist)
        }
        return c; // returns null if camera is unavailable
    }
}
