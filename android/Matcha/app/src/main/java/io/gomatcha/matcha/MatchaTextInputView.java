package io.gomatcha.matcha;

import android.accessibilityservice.AccessibilityService;
import android.app.Activity;
import android.content.Context;
import android.os.Handler;
import android.os.SystemClock;
import android.text.Editable;
import android.text.InputType;
import android.text.SpannableString;
import android.text.SpannableStringBuilder;
import android.text.TextWatcher;
import android.util.Log;
import android.view.Gravity;
import android.view.KeyEvent;
import android.view.MotionEvent;
import android.view.View;
import android.view.Window;
import android.view.WindowManager;
import android.view.inputmethod.EditorInfo;
import android.view.inputmethod.InputMethodManager;
import android.widget.EditText;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.google.protobuf.InvalidProtocolBufferException;

import io.gomatcha.bridge.GoValue;
import io.gomatcha.matcha.pb.text.PbText;
import io.gomatcha.matcha.pb.view.PbView;
import io.gomatcha.matcha.pb.view.slider.PbSlider;
import io.gomatcha.matcha.pb.view.textinput.PbTextInput;

import static android.widget.TextView.BufferType.EDITABLE;

public class MatchaTextInputView extends MatchaChildView {
    EditText view;
    boolean editing;
    boolean focused;

    static {
        MatchaView.registerView("gomatcha.io/matcha/view/textinput", new MatchaView.ViewFactory() {
            @Override
            public MatchaChildView createView(Context context, MatchaViewNode node) {
                return new MatchaTextInputView(context, node);
            }
        });
    }

    public void showKeyboardWithFocus(Activity a) {
        try {
            view.requestFocus();
            InputMethodManager imm = (InputMethodManager) a.getSystemService(Context.INPUT_METHOD_SERVICE);
            imm.showSoftInput(view, InputMethodManager.SHOW_IMPLICIT);
            a.getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_STATE_VISIBLE);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public MatchaTextInputView(Context context, MatchaViewNode node) {
        super(context, node);
        final Context ctx = context;

        RelativeLayout.LayoutParams params = new RelativeLayout.LayoutParams(RelativeLayout.LayoutParams.MATCH_PARENT, RelativeLayout.LayoutParams.MATCH_PARENT);
        view = new EditText(context);
        view.setGravity(Gravity.TOP);
        view.setOnEditorActionListener(new TextView.OnEditorActionListener() {
            @Override
            public boolean onEditorAction(TextView textView, int i, KeyEvent keyEvent) {
                boolean handled = false;
                if (i == EditorInfo.IME_ACTION_DONE) {
                    Log.v("x", "submit");
                    MatchaTextInputView.this.viewNode.rootView.call("OnSubmit", MatchaTextInputView.this.viewNode.id);
                    handled = true;
                }
                return handled;
            }
        });
        view.addTextChangedListener(new TextWatcher() {
            @Override
            public void beforeTextChanged(CharSequence charSequence, int i, int i1, int i2) {
                // no-op
            }
            @Override
            public void onTextChanged(CharSequence charSequence, int i, int i1, int i2) {
                if (!editing) {
                    PbText.StyledText styledText = Protobuf.toProtobuf((SpannableStringBuilder) charSequence);
                    PbTextInput.Event proto = PbTextInput.Event.newBuilder().setStyledText(styledText).build();
                    MatchaTextInputView.this.viewNode.rootView.call("OnTextChange", MatchaTextInputView.this.viewNode.id, new GoValue(proto.toByteArray()));
                }
            }
            @Override
            public void afterTextChanged(Editable editable) {
                // no-op
            }
        });
        view.setOnFocusChangeListener(new OnFocusChangeListener() {
            @Override
            public void onFocusChange(View view, boolean b) {
                if (!editing) {
                    PbTextInput.FocusEvent proto = PbTextInput.FocusEvent.newBuilder().setFocused(b).build();
                    MatchaTextInputView.this.viewNode.rootView.call("OnFocus", MatchaTextInputView.this.viewNode.id, new GoValue(proto.toByteArray()));
                }
            }
        });
        addView(view, params);
    }

    @Override
    public void setNode(PbView.BuildNode buildNode) {
        super.setNode(buildNode);
        try {
            editing = true;
            PbTextInput.View proto = buildNode.getBridgeValue().unpack(PbTextInput.View.class);
            SpannableString str = Protobuf.newAttributedString(proto.getStyledText());
            if (!str.toString().equals(view.getText().toString())) {
                view.setText(str, TextView.BufferType.SPANNABLE);
            }
            if (view.hasFocus() && !proto.getFocused()) {
                view.clearFocus();
                InputMethodManager imm = (InputMethodManager) MatchaTextInputView.this.getContext().getSystemService(Context.INPUT_METHOD_SERVICE);
                imm.hideSoftInputFromWindow(((Activity)this.getContext()).getCurrentFocus().getWindowToken(), InputMethodManager.HIDE_NOT_ALWAYS);
            } else if (!view.hasFocus() && proto.getFocused()) {
                new Handler().postDelayed(new Runnable() {
                    public void run() {
                        editing = true;
                        view.requestFocus();
                        InputMethodManager imm = (InputMethodManager) MatchaTextInputView.this.getContext().getSystemService(Context.INPUT_METHOD_SERVICE);
                        imm.showSoftInput(view, 0);
                        editing = false;
                    }
                }, 100);
            }
            int inputType;
            switch (proto.getKeyboardType()) {
                case DEFAULT_TYPE:
                    inputType = InputType.TYPE_CLASS_TEXT;
                break;
                case NUMBER_TYPE:
                    inputType = InputType.TYPE_CLASS_NUMBER;
                break;
                case NUMBER_PUNCTUATION_TYPE:
                    inputType = InputType.TYPE_CLASS_NUMBER;
                break;
                case DECIMAL_TYPE:
                    inputType = InputType.TYPE_CLASS_NUMBER;
                break;
                case PHONE_TYPE:
                    inputType = InputType.TYPE_CLASS_PHONE;
                break;
                case ASCII_TYPE:
                    inputType = InputType.TYPE_CLASS_TEXT;
                break;
                case EMAIL_TYPE:
                    inputType = InputType.TYPE_CLASS_TEXT | InputType.TYPE_TEXT_VARIATION_EMAIL_ADDRESS;
                break;
                case URL_TYPE:
                    inputType = InputType.TYPE_CLASS_TEXT | InputType.TYPE_TEXT_VARIATION_URI;
                break;
                case WEB_SEARCH_TYPE:
                    inputType = InputType.TYPE_CLASS_TEXT;
                break;
                case NAME_PHONE_TYPE:
                    inputType = InputType.TYPE_CLASS_TEXT;
                break;
                default:
                    inputType = InputType.TYPE_CLASS_TEXT;
            }
            if (proto.getSecureTextEntry()) {
                inputType |= InputType.TYPE_TEXT_VARIATION_PASSWORD;
            }

            int imeOptions = 0;
            switch (proto.getKeyboardReturnType()) {
                case DEFAULT_RETURN_TYPE:
                break;
                case GO_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_GO;
                break;
                case GOOGLE_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_GO;
                break;
                case JOIN_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_DONE;
                break;
                case NEXT_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_NEXT;
                break;
                case ROUTE_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_DONE;
                break;
                case SEARCH_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_SEARCH;
                break;
                case SEND_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_SEND;
                break;
                case YAHOO_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_DONE;
                break;
                case DONE_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_DONE;
                break;
                case EMERGENCY_CALL_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_DONE;
                break;
                case CONTINUE_RETURN_TYPE:
                    imeOptions = EditorInfo.IME_ACTION_DONE;
                break;
                default:
            }
            view.setImeOptions(imeOptions);
            view.setSingleLine(proto.getMaxLines() == 1);

            view.setHint(Protobuf.newAttributedString(proto.getPlaceholderText()));
            focused = proto.getFocused();
            editing = false;
            
        } catch (InvalidProtocolBufferException e) {
        }
    }
}
