package io.gomatcha.cameraview;

import android.Manifest;
import android.content.Context;
import android.content.pm.PackageManager;
import android.hardware.Camera;
import android.support.v4.app.ActivityCompat;
import android.support.v4.content.ContextCompat;
import android.support.v7.widget.SwitchCompat;
import android.util.Log;
import android.view.View;

import com.google.gson.Gson;

import io.gomatcha.matcha.MatchaChildView;
import io.gomatcha.matcha.MatchaView;
import io.gomatcha.matcha.MatchaViewNode;

public class CameraView extends MatchaChildView {
    Camera camera;
    CameraPreview preview;
    MatchaViewNode viewNode;
    boolean visible;
    NativeState nativeState;

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
        Log.v("x", "New camera");
        this.setClipChildren(true);
        viewNode = node;

        preview = new CameraPreview(context);
        addView(preview);
    }

    class NativeState {
        private boolean frontCamera;
        NativeState() {}
    }

    @Override
    public void setNativeState(byte[] nativeState) {
        super.setNativeState(nativeState);

        Log.v("x", "camera set native state");
        Gson gson = new Gson();
        this.nativeState = gson.fromJson(new String(nativeState), NativeState.class);
        reloadCamera();
        Log.v("x", "camera did set native state");
    }
    
    @Override
    public void onVisibilityChanged(View changedView, int visibility) {
        super.onVisibilityChanged(changedView, visibility);
        reloadCamera();
        Log.v("x", "OnVisibilityChanged" + visibility);
    }
    
    @Override
    public void onWindowVisibilityChanged(int visibility) {
        super.onWindowVisibilityChanged(visibility);
        reloadCamera();
        Log.v("x", "OnWindowVisibilityChanged" + visibility);
    }


    
    void reloadCamera() {
        Log.v("x", "reloadCamera" + visible + this.visible);

        boolean visible = isShown() && hasWindowFocus();
        // if (visible == this.visible) {
        //     return;
        // }
        // this.visible = visible;
        // Log.v("x", "camera!!!" + visible + nativeState.frontCamera);
        
        // Remove old camera
        if (camera != null) {
            preview.setCamera(null, null);
            camera.release();
            camera = null;
        }

        // Add new camera
        if (visible && nativeState != null) {
            Camera.CameraInfo info = new android.hardware.Camera.CameraInfo();;
            camera = getCameraInstance(nativeState.frontCamera, info);
            preview.setCamera(camera, info);
        }
    }
    
    public Camera getCameraInstance(boolean frontCamera, Camera.CameraInfo info) {
        //String[] permissions = {"android.permission.CAMERA"};
        //if (ContextCompat.checkSelfPermission(this.getContext(), Manifest.permission.CAMERA) != PackageManager.PERMISSION_GRANTED) {
        //    ActivityCompat.requestPermissions(this, permissions, MY_PERMISSIONS_REQUEST_CAMERA);
        //}

        Camera c = null;
        try {
            int cameraCount = 0;
            Camera.CameraInfo cameraInfo = new Camera.CameraInfo();
            cameraCount = Camera.getNumberOfCameras();
            for (int camIdx = 0; camIdx < cameraCount; camIdx++) {
                Camera.getCameraInfo(camIdx, cameraInfo);
                if (frontCamera && cameraInfo.facing == Camera.CameraInfo.CAMERA_FACING_FRONT) {
                    c = Camera.open(camIdx);
                    android.hardware.Camera.getCameraInfo(camIdx, info);
                    break;
                } else if (!frontCamera && cameraInfo.facing == Camera.CameraInfo.CAMERA_FACING_BACK) {
                    c = Camera.open(camIdx);
                    android.hardware.Camera.getCameraInfo(camIdx, info);
                    break;
                }
            }
        }
        catch (Exception e){
            Log.v("x", e.toString());
            // Camera is not available (in use or does not exist)
        }
        return c; // returns null if camera is unavailable
    }
}
