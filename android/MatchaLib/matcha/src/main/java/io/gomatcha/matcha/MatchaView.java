package io.gomatcha.matcha;

import android.app.Activity;
import android.content.Context;
import android.os.Build;
import android.util.DisplayMetrics;
import android.view.KeyEvent;
import android.view.View;
import android.view.Window;
import android.view.WindowManager;
import android.widget.RelativeLayout;

import java.lang.ref.WeakReference;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.proto.view.PbView;
import io.gomatcha.matcha.proto.view.android.PbStatusBar;

public class MatchaView extends RelativeLayout {
    static ArrayList<WeakReference<MatchaView>> views = new ArrayList<WeakReference<MatchaView>>();
    GoValue goValue;
    long identifier;
    MatchaViewNode node;

    static {
        new JavaBridge();
    }

    public MatchaView(Context context, GoValue v2) {
        super(context);
        GoValue v = GoValue.withFunc("gomatcha.io/matcha/view NewRoot").call("", v2)[0];
        goValue = v;
        identifier = v.call("Id")[0].toLong();
        long viewid = v.call("ViewId")[0].toLong();
        node = new MatchaViewNode(null, this, viewid);
        setFocusable(true);
        setFocusableInTouchMode(true);

        views.add(new WeakReference<MatchaView>(this));

        // Initialize JavaBridge
        JavaBridge.init(context);
    }

    boolean loaded = false;

    void update(PbView.Root root) {
        node.setRoot(root);

        if (!loaded) {
            loaded = true;

            RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
            node.view.setLayoutParams(params);

            addView(node.view);
        }

        Map<String, com.google.protobuf.Any> map = root.getMiddlewareMap();
        com.google.protobuf.Any any = map.get("gomatcha.io/matcha/view/android statusbar");
        if (any != null) {
            try {
                PbStatusBar.StatusBar proto = any.unpack(PbStatusBar.StatusBar.class);
                int color = Protobuf.newColor(proto.getColor());
                Window window = ((Activity)getContext()).getWindow();
                window.clearFlags(WindowManager.LayoutParams.FLAG_TRANSLUCENT_STATUS);
                window.addFlags(WindowManager.LayoutParams.FLAG_DRAWS_SYSTEM_BAR_BACKGROUNDS);
                if (Build.VERSION.SDK_INT >= 21) {
                    if (!proto.getStyle()) {
                        int flags = this.getSystemUiVisibility();
                        flags |= View.SYSTEM_UI_FLAG_LIGHT_STATUS_BAR;
                        this.setSystemUiVisibility(flags);
                    } else {
                        int flags = this.getSystemUiVisibility();
                        flags &= ~View.SYSTEM_UI_FLAG_LIGHT_STATUS_BAR;
                        this.setSystemUiVisibility(flags);
                    }
                    window.setStatusBarColor(color);
                }
            } catch (com.google.protobuf.InvalidProtocolBufferException e) {
            }
        }
    }

    @Override
    protected void onSizeChanged(int w, int h, int oldw, int oldh) {
        final double width = (double)w / this.getResources().getDisplayMetrics().densityDpi * DisplayMetrics.DENSITY_DEFAULT;
        final double height = (double)h / this.getResources().getDisplayMetrics().densityDpi * DisplayMetrics.DENSITY_DEFAULT;

        this.post( new Runnable() {
            @Override
            public void run() {
                goValue.call("SetSize", new GoValue((double)width), new GoValue((double)height));
                GoValue.withFunc("gomatcha.io/matcha/animate screenUpdate").call("");
            }
        });
    }

    @Override
    protected void finalize() {
        goValue.call("Stop");
    }

    public GoValue[] call(String func, long viewId, GoValue... args) {
        GoValue[] args2 = new GoValue[]{new GoValue(func), new GoValue(viewId), new GoValue(args)};
        return this.goValue.call("Call", args2);
    }

    // View registry
    static Map<String, ViewFactory> viewRegistry = new HashMap<String, ViewFactory>();

    static {
        try {
            Class.forName("io.gomatcha.matcha.MatchaBasicView");
            Class.forName("io.gomatcha.matcha.MatchaImageView");
            Class.forName("io.gomatcha.matcha.MatchaTextView");
            Class.forName("io.gomatcha.matcha.MatchaTextInputView");
            Class.forName("io.gomatcha.matcha.MatchaSwitchView");
            Class.forName("io.gomatcha.matcha.MatchaButton");
            Class.forName("io.gomatcha.matcha.MatchaSlider");
            Class.forName("io.gomatcha.matcha.MatchaScrollView");
            Class.forName("io.gomatcha.matcha.MatchaStackView");
            Class.forName("io.gomatcha.matcha.MatchaPagerView");
            Class.forName("io.gomatcha.matcha.MatchaToolbarView");
        } catch (ClassNotFoundException e) {
            throw new RuntimeException(e);
        }
    }
    
    public interface ViewFactory {
        MatchaChildView createView(Context context, MatchaViewNode node);
    }
    
    public synchronized static void registerView(String name, ViewFactory factory) {
        viewRegistry.put(name, factory);
    }

    synchronized static MatchaChildView createView(String name, Context context, MatchaViewNode node) {
        ViewFactory factory = viewRegistry.get(name);
        if (factory == null) {
            return new MatchaUnknownView(context, node);
        }
        return factory.createView(context, node);
    }

    long downTime = 0;
    @Override
    public boolean dispatchKeyEventPreIme(KeyEvent event) {
        if (event.getKeyCode() == KeyEvent.KEYCODE_BACK) {
            KeyEvent.DispatcherState state = getKeyDispatcherState();
            if (state != null) {
                if (event.getAction() == KeyEvent.ACTION_DOWN && event.getRepeatCount() == 0) {
                    GoValue[] rlt = GoValue.withFunc("gomatcha.io/view/android StackBarCanBack").call("");
                    if (!rlt[0].toBool()) {
                        return super.dispatchKeyEventPreIme(event);
                    }
                    downTime = event.getDownTime();
                    state.startTracking(event, this);
                    return true;
                } else if (event.getAction() == KeyEvent.ACTION_UP && !event.isCanceled() && state.isTracking(event) && downTime == event.getDownTime()) {
                    GoValue.withFunc("gomatcha.io/view/android StackBarOnBack").call("");
                    return true;
                }
            }
        }
        return super.dispatchKeyEventPreIme(event);
    }

}
