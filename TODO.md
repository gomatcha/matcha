* Stackview animation
* Stackview button items
* Android vertical/horizontal scrollviews?

High:
* Android custom view example
* Better error logging for panics.

Medium:
* Crash on "Stop"
* Rebuild Instagram.
* More Touch Recognizers: Pan, Swipe, Pinch, Rotation
* Modal presentation.
* Picker
* TextField
* Rotation / orientation
* StackBar height / hidden, color
* Android image resource scale??
* Android scroll views scrollposition
* Better notavailable-view. lightgray with centered "Unknown".

Low:
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
* Guide.Insets, GUide.Transform? Layout.Insets(top, left, bottom, right)?
* Table behaviors
* Button disabled/highlighted state using anonymous structs?

Very Low:
* Building for iPhone 5 Simulator doesn't work.
* Add tests around core functionality. Store, etc.
* Button should fade when disabled.
* Statusbar color
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
* Change enabled to disabled.
* Change OnPress to OnActivate, OnSelect?
* StackBarButton or StackBarItem
* should options be pointer or struct receiver. Should statusbar be a pointer?
* Touch.OnTouch rename to OnEvent? OnMajorEvent? OnRecognize? 
* Rename view.Model.NativeValues to NativeOptions?
* Rename app package to a name thats less used?? application.LoadImage()??

Pro:
* Add preload, and prepreload stages
* Webview
* Debugging
* LocalStorage / Keychain / UserDefaults
* Cliboard
* Notifications
* Video / Sound / Microphone / Accelerometer
* ActionSheet
* CameraView
* MapView
* GPS
* Accessibility
* Android touch highlight?

Target Android 4.1 API 16 JellyBean