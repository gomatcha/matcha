package io.gomatcha.cameraview;

import android.content.Context;
import android.hardware.Camera;
import android.util.Log;
import android.view.SurfaceHolder;
import android.view.SurfaceView;

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
            int orientation = CameraView.cameraDisplayOrientation(getContext(), cameraInfo);
            mCamera.setDisplayOrientation(orientation);
            
            try {
                mCamera.setPreviewDisplay(mHolder);
                mCamera.startPreview();
            } catch (Exception e){
                Log.d("x", "Error starting camera preview: " + e.getMessage());
            }
        }
    }

    public void surfaceCreated(SurfaceHolder holder) {
        // no-op
    }

    public void surfaceDestroyed(SurfaceHolder holder) {
        // no-op
    }

    public void surfaceChanged(SurfaceHolder holder, int format, int w, int h) {
        // If your preview can change or rotate, take care of those events here.
        // Make sure to stop the preview before resizing or reformatting it.
        this.setCamera(mCamera, mCameraInfo);
    }
}
