High:
* Take screenshots / videos of a running app.
* Delay any ui events while viwecontroller is updating.
* Scroll is forcing a relayout.

Medium:
* Rotation / orientation
* Rebuild settings app, Todo App.
* Multiple view controllers.
* Picker
* TextField
* Optimize middleware so they aren't called on every view.

Low:
* Skip "ctx view.Context, key string, " paramater on views we know are top level?
* Improve function call performance.
* Switching quickly between navigation item causes visual glitch. 2 quick backs.
* How to respond to memory pressure?
* Have matcha flag that generates a new xcodeproj for easy setup.
* Add tests around core functionality. Store, etc.
* Examples. Start rebuild a few apps. Pintrest, Instagram, Settings, Slack
* Modal presentation
* Asset catalog
* StackBar height / hidden, color
* More Touch Recognizers: Pan, Swipe, Pinch, EdgePan, Rotation
* Table ScrollBehaviors, Table Direction
* Custom painters.
* Compile a list of things that should be easy to do and implement them. Button activation cancelled by vertical scrolling but not horizontal, Pinch to zoom, Highlighting a view and dragging outside of it and back in., Horizontal swipe on tableview to show delete button, Touch driven animations. AKA swipe back to navigate.
* Building for iPhone 5 Simulator doesn't work.

Very Low:
* Statusbar color
* Automatically insert copyright notice.
* StyledText
* Text selection.
* Localization
* View 3d transforms.
* GridView
* Add preload, and prepreload stages
* Debug constraints.
* Collect native resources into assets.
* Animations: Spring, Delay, Batch, Reverse, Decay, Repeat
* Rework Slider.FloatNotifier to use comm.Float64Value and give it a better name InOutValue?
* Flexbox

Bugs:

Documentation:
* cmd
* comm
* store
* docs
* env
* examples
* layout
* view 
    * stackscreen
    * tabscrceen
    * textinput
    * textview

Pro:
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
