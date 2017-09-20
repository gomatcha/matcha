* MailChimp thing.
* BLog thing.
* Better images, readme docs.
* Android custom view example
* Consulting.

High:
* Non-image resources.
    * Assets directory that gets merged from any folder containing /assets and importing "gomatcha.io/matcha". What to do about images? _1x.png _2x.png
* Android double tap.
* Rotation / orientation

Medium:
* Better notavailable-view. lightgray with centered "Unknown".
* Android scroll views scrollposition
* Move images loading to native side.
* Webview
* Stackview button items
* Android vertical/horizontal scrollviews?
* Crash on "Stop"
* Rebuild Instagram.
* More Touch Recognizers: Pan, Swipe, Pinch, Rotation
* Modal presentation.
* Picker
* StackBar height / hidden, color

Low:
* Better error logging for panics.
* Not getting start and inprogress events for UITapGestureRecognizer.
* UIButtonGestureRecognizer only sends a event if inside/outside changes. Its faster but less generic?
* Delay any ui events while viewcontroller is updating??
* Optimize middleware so they aren't called on every view.
* Improve function call performance.
* Switching quickly between navigation item causes visual glitch. 2 quick backs.
* How to respond to memory pressure?
* Have matcha flag that generates a new xcodeproj for easy setup.
* Examples. Start rebuild a few apps. Pintrest, Instagram, Settings, Slack
* Table ScrollBehaviors, Table Direction
* Custom painters.
* Compile a list of things that should be easy to do and implement them. Button activation cancelled by vertical scrolling but not horizontal, Pinch to zoom, Highlighting a view and dragging outside of it and back in., Horizontal swipe on tableview to show delete button, Touch driven animations. AKA swipe back to navigate.
* Guide.Insets, Guide.Transform? Layout.Insets(top, left, bottom, right)?
* Table behaviors
* Button disabled/highlighted state using anonymous structs?

Very Low:
* Building for iPhone 5 Simulator doesn't work.
* Add tests around core functionality. Store, etc.
* Button should fade when disabled.
* Automatically insert copyright notice.
* Text selection.
* Localization
* View 3d transforms.
* GridView
* Debug constraints.
* Collect native resources into assets.
* Animations: Spring, Delay, Batch, Reverse, Decay, Repeat
* Rework Slider.FloatNotifier to use comm.Float64Value and give it a better name InOutValue?
* Flexbox
* Strikethrough doesn't work.

Refactors
* Change OnPress to OnActivate, OnSelect?
* Touch.OnTouch rename to OnEvent? OnMajorEvent? OnRecognize? 
* Touch package to Pointer?

Documentation
* animate
* comm
* layout
* layout/absolute
* layout/constraint
* layout/full
* layout/table
* paint
* text
* touch
* view
* view/ios
* view/android

Upgrade:
* Add preload, and prepreload stages :
* Webview : 
* Debugging : 
* LocalStorage / Keychain / UserDefaults : pref, iospref, osxpref, andpref, securepref?
* Cliboard : clipboard
* Notifications : notification, iosnotification, andnotification
* Video / Sound / Microphone / Accelerometer : zview, zios, zandroid
* ActionSheet
* CameraView
* MapView
* GPS
* Accessibility : access, 
* Android touch highlight?

Target Android 4.1 API 16 JellyBean