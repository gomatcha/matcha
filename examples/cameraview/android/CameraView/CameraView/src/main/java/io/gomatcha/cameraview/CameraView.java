package io.gomatcha.cameraview;

import android.content.Context;
import android.hardware.Camera;
import android.os.Handler;
import android.util.Log;
import android.view.Surface;
import android.view.View;
import android.view.WindowManager;
import android.widget.Button;

import com.google.gson.Gson;

import io.gomatcha.matcha.MatchaChildView;
import io.gomatcha.matcha.MatchaLayout;
import io.gomatcha.matcha.MatchaView;
import io.gomatcha.matcha.MatchaViewNode;
import io.gomatcha.matcha.Util;

import static android.content.Context.WINDOW_SERVICE;

public class CameraView extends MatchaChildView implements View.OnClickListener {
    Camera camera;
    Camera.CameraInfo cameraInfo;
    CameraPreview preview;
    MatchaViewNode viewNode;
    NativeState nativeState;
    CameraButton button;
    boolean visible;
    boolean frontCamera;

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

        preview = new CameraPreview(context);
        addView(preview);

        button = new CameraButton(context);
        button.setOnClickListener(this);
        addView(button);

        reloadCameraLayout();
    }

    class NativeState {
        private boolean frontCamera;
        NativeState() {}
    }

    @Override
    public void setNativeState(byte[] nativeState) {
        super.setNativeState(nativeState);

        Gson gson = new Gson();
        this.nativeState = gson.fromJson(new String(nativeState), NativeState.class);
        reloadCamera();
    }
    
    @Override
    public void onVisibilityChanged(View changedView, int visibility) {
        super.onVisibilityChanged(changedView, visibility);
        reloadCamera();
    }

    @Override
    public void onWindowFocusChanged(boolean hasFocus) {
        super.onWindowFocusChanged(hasFocus);

        // Delay a bit or the start preview fails.
        final Handler handler = new Handler();
        handler.postDelayed(new Runnable() {
            @Override
            public void run() {
                reloadCamera();
            }
        }, 100);
    }

    @Override
    public void onLayout(boolean changed, int left, int top, int right, int bottom) {
        reloadCameraLayout();
        super.onLayout(changed, left, top, right, bottom);
    }

    // OnClickListener

    public void onClick(View v) {
        Log.v("Matcha", "OnClick");
    }

    //

    void reloadCameraLayout() {
        MatchaLayout.LayoutParams params = new MatchaLayout.LayoutParams();
        if (camera != null) {
            float width = getWidth();
            float height = getHeight();
            
            Camera.Size cameraSize = camera.getParameters().getPreviewSize();
            float cameraHeight = cameraSize.height;
            float cameraWidth = cameraSize.width;
            int cameraOrientation = CameraView.cameraDisplayOrientation(getContext(), cameraInfo);
            if (cameraOrientation == 90 || cameraOrientation == 270) {
                float temp = cameraWidth;
                cameraWidth = cameraHeight;
                cameraHeight = temp;
            }
            
            if (cameraWidth/cameraHeight > width/height) {
                float ratio = width / cameraWidth;
                float previewHeight = height * ratio;

                params.left = 0;
                params.right = width;
                params.top = (height - previewHeight) / 2;
                params.bottom = params.top + previewHeight;
            } else {
                float ratio = height / cameraHeight;
                float previewWidth = width * ratio;

                params.top = 0;
                params.bottom = height;
                params.left = (width - previewWidth) / 2;
                params.right = params.left + previewWidth;
            }
        }
        preview.setLayoutParams(params);

        MatchaLayout.LayoutParams buttonParams = new MatchaLayout.LayoutParams();
        buttonParams.top = getHeight() - Util.dipToPixels(getContext(),60);
        buttonParams.left = (getWidth() - Util.dipToPixels(getContext(),50)) / 2;
        buttonParams.bottom = buttonParams.top + Util.dipToPixels(getContext(),50);
        buttonParams.right = buttonParams.left + Util.dipToPixels(getContext(),50);
        button.setLayoutParams(buttonParams);
    }

    void reloadCamera() {
        boolean visible = isShown() && hasWindowFocus();
        if (visible == this.visible && nativeState.frontCamera == frontCamera) {
            return;
        }
        this.visible = visible;
        this.frontCamera = nativeState.frontCamera;

        // Remove old camera
        if (camera != null) {
            preview.setCamera(null, null);
            camera.release();
            camera = null;
            cameraInfo = null;
        }

        // Add new camera
        if (visible && nativeState != null) {
            cameraInfo = new android.hardware.Camera.CameraInfo();;
            camera = getCameraInstance(nativeState.frontCamera, cameraInfo);
            preview.setCamera(camera, cameraInfo);
        }

        reloadCameraLayout();
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

    static int cameraDisplayOrientation(Context context, Camera.CameraInfo cameraInfo) {
        WindowManager manager = (WindowManager)context.getSystemService(WINDOW_SERVICE);
        int rotation = manager.getDefaultDisplay().getRotation();
        int degrees = 0;
        switch (rotation) {
            case Surface.ROTATION_0: degrees = 0; break;
            case Surface.ROTATION_90: degrees = 90; break;
            case Surface.ROTATION_180: degrees = 180; break;
            case Surface.ROTATION_270: degrees = 270; break;
        }
        int result;
        if (cameraInfo.facing == Camera.CameraInfo.CAMERA_FACING_FRONT) {
            result = (cameraInfo.orientation + degrees) % 360;
            result = (360 - result) % 360;  // compensate the mirror
        } else {  // back-facing
            result = (cameraInfo.orientation - degrees + 360) % 360;
        }
        return result;
    }
}
