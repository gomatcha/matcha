package io.gomatcha.cameraview;

import android.content.Context;
import android.hardware.Camera;
import android.util.Log;
import android.view.Surface;
import android.view.SurfaceHolder;
import android.view.SurfaceView;
import android.view.WindowManager;

import java.io.IOException;
import java.util.List;

import static android.content.Context.WINDOW_SERVICE;

public class CameraPreview extends SurfaceView implements SurfaceHolder.Callback {
    private SurfaceHolder mHolder;
    private Camera mCamera;
    private Camera.CameraInfo mCameraInfo;

    public CameraPreview(Context context) {
        super(context);

        // Install a SurfaceHolder.Callback so we get notified when the
        // underlying surface is created and destroyed.
        mHolder = getHolder();
        mHolder.addCallback(this);
        // deprecated setting, but required on Android versions prior to 3.0
        mHolder.setType(SurfaceHolder.SURFACE_TYPE_PUSH_BUFFERS);
    }

    public void setCamera(Camera camera, Camera.CameraInfo cameraInfo) {
        if (mCamera != null) {
            try {
                mCamera.stopPreview();
                mCamera = null;
            } catch (Exception e){
                // ignore: tried to stop a non-existent preview
            }
        }
        
        mCamera = camera;
        mCameraInfo = cameraInfo;

        if (mHolder.getSurface() != null && camera != null) {
            // Set Rotation
            WindowManager manager = (WindowManager) getContext().getSystemService(WINDOW_SERVICE);
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
            camera.setDisplayOrientation(result);
            
            Camera.Parameters params = camera.getParameters();
            List<Camera.Size> supportedPreviewSizes = params.getSupportedPreviewSizes();
            Camera.Size camPreviewSize = getOptimalPreviewSize(supportedPreviewSizes, getWidth() , getHeight());
            params.setPreviewSize(camPreviewSize.width, camPreviewSize.height);
            //camera.setParameters(params);

            
            try {
                Log.v("x", "StartCamera" + camPreviewSize.width + camPreviewSize.height + supportedPreviewSizes.toString());
                mCamera.setPreviewDisplay(mHolder);
                mCamera.startPreview();
            } catch (Exception e){
                Log.d("x", "Error starting camera preview: " + e.getMessage());
            }
        }
    }

    public void surfaceCreated(SurfaceHolder holder) {
        Log.v("x", "surfacecreated");
        //this.setCamera(mCamera, mCameraInfo);
    }

    public void surfaceDestroyed(SurfaceHolder holder) {
        // empty. Take care of releasing the Camera preview in your activity.
        Log.v("x", "surfacedestroyed");
    }

    public void surfaceChanged(SurfaceHolder holder, int format, int w, int h) {
        Log.v("x", "surfacechanged");
        // If your preview can change or rotate, take care of those events here.
        // Make sure to stop the preview before resizing or reformatting it.
        this.setCamera(mCamera, mCameraInfo);
    }
    
    private Camera.Size getOptimalPreviewSize(List<Camera.Size> sizes, int w, int h) {
        final double ASPECT_TOLERANCE = 0.1;
        double targetRatio=(double)h / w;

        if (sizes == null) return null;

        Camera.Size optimalSize = null;
        double minDiff = Double.MAX_VALUE;

        int targetHeight = h;

        for (Camera.Size size : sizes) {
            double ratio = (double) size.width / size.height;
            if (Math.abs(ratio - targetRatio) > ASPECT_TOLERANCE) continue;
            if (Math.abs(size.height - targetHeight) < minDiff) {
                optimalSize = size;
                minDiff = Math.abs(size.height - targetHeight);
            }
        }

        if (optimalSize == null) {
            minDiff = Double.MAX_VALUE;
            for (Camera.Size size : sizes) {
                if (Math.abs(size.height - targetHeight) < minDiff) {
                    optimalSize = size;
                    minDiff = Math.abs(size.height - targetHeight);
                }
            }
        }
        return optimalSize;
    }
}
