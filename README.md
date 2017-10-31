# Matcha - iOS and Android apps in Go

[gomatcha.io](https://gomatcha.io)

Matcha is in early development! There are many rough edges and APIs may still
change. Please file issues for any bugs you find.

### What is Matcha?

Matcha is a package for building iOS and Android applications and frameworks in
Go. Matcha provides a UI component library similar to ReactNative and exposes
bindings to Objective-C and Java code through reflection. The library also
provides Go APIs for common app tasks.

### Examples

[![settings-example](docs/settings.gif)](https://github.com/gomatcha/matcha/tree/master/examples/settings)
[![insta-example](docs/insta.gif)](https://github.com/gomatcha/matcha/tree/master/examples/insta)
[![todo-example](docs/todo.gif)](https://github.com/gomatcha/matcha/tree/master/examples/todo)

### Installation - macOS

Matcha requires the following components to be installed.

* Go 1.8+
* Xcode 8.3+
* Android Studio 2.3+ (with SDK 26, NDK and Android Support)

Start by installing Xcode and Android Studio. Instructions can be found at
https://developer.apple.com/download/ and
https://developer.android.com/studio/install.html.

Open Android Studio's SDK Manager and under the SDK Platforms tab, install
the Android 8 Platform (API 27). And in the SDK Tools tab, install NDK and the
Android Support Repository.

Configure the ANDROID_HOME enviromental variable to point to the Android SDK.
This is usually located at ~/Library/Android/sdk.

    export ANDROID_HOME=<SDK location>

Fetch the project and install the matcha command.

    go get gomatcha.io/matcha/...

We build the Go standard library for iOS and Android with the following command.
This may take awhile. The output is installed at `$GOPATH/pkg/matcha`. If your
path doesn't contain $GOPATH/bin, you may need to replace these calls with
`$GOPATH/bin/matcha`.

    matcha init

Now build the example project. The output is installed at `$GOPATH/src/gomatcha.io/matcha/ios/MatchaBridge/MatchaBridge/MatchaBridge.a` and `$GOPATH/src/gomatcha.io/matcha/android/matchabridge.aar`.

    matcha build gomatcha.io/matcha/examples

We can now open the sample iOS project.

    open $GOPATH/src/gomatcha.io/matcha/examples/ios-app/SampleApp.xcworkspace
    
Set the Development Team in Xcode under General > Signing and select `SampleApp` in
the target dropdown in the upper right. Then run the App!

For Android, you may need to run the following to allow Android Studio to read
your GOPATH (https://stackoverflow.com/a/14285335). This also must be done on
reboot.

    launchctl setenv GOPATH $GOPATH

Then simply open the sample Android Studio project and hit run!

    open -a /Applications/Android\ Studio.app $GOPATH/src/gomatcha.io/matcha/examples/android-app/SampleApp

### Installation - Linux

Matcha requires the following components to be installed. iOS builds are not
supported on Linux.

* Go 1.8+
* Android Studio 2.3+ (with SDK 26, NDK and Android Support Library)

Start by installing Android Studio. Instructions can be found at
https://developer.android.com/studio/install.html. 

Open Android Studio's SDK Manager and under the SDK Platforms tab, install
the Android 8 Platform (API 27). And in the SDK Tools tab, install NDK and the
Android Support Repository.

Configure the ANDROID_HOME enviromental variable to point to the Android SDK.
This is often located at ~/Android/Sdk depending on your install.

    export ANDROID_HOME=<SDK location>

Modify your PATH to include the Java compiler if it does not already. javac can
be found at android-studio/jre/bin.

    export PATH=${PATH}:<Java Compiler location>

Fetch the project and install the matcha command.

    go get gomatcha.io/matcha/...

We build the Go standard library for Android with the following command.
This may take awhile. The output is installed at `$GOPATH/pkg/matcha`. If your
path doesn't contain $GOPATH/bin, you may need to replace these calls with
`$GOPATH/bin/matcha`.

    matcha init

Now build the example project. The output is installed at `$GOPATH/src/gomatcha.io/matcha/android/matchabridge.aar`.

    matcha build gomatcha.io/matcha/examples

Now open the sample Android Studio project and hit run!

    <Android Studio location>/bin/studio.sh $GOPATH/src/gomatcha.io/matcha/examples/android-app/SampleApp

### Installation - Windows

Matcha requires the following components to be installed. iOS builds are not
supported on Windows.

* Go 1.8+
* Android Studio 2.3+ (with SDK 26, NDK and Android Support Library)

Start by installing Android Studio. Instructions can be found at
https://developer.android.com/studio/install.html. 

Open Android Studio's SDK Manager and under the SDK Platforms tab, install
the Android 8 Platform (API 27). And in the SDK Tools tab, install NDK and the
Android Support Repository.

Configure the ANDROID_HOME enviromental variable to point to the Android SDK.
This is often located at ~/Android/Sdk depending on your install.

    export ANDROID_HOME=<SDK location>

Modify your PATH to include the Java compiler if it does not already. javac can
be found at android-studio/jre/bin.

    export PATH=${PATH}:<Java Compiler location>

Fetch the project and install the matcha command.

    go get gomatcha.io/matcha/...

We build the Go standard library for Android with the following command.
This may take awhile. The output is installed at `$GOPATH/pkg/matcha`. If your
path doesn't contain $GOPATH/bin, you may need to replace these calls with
`$GOPATH/bin/matcha`.

    matcha init

Now build the example project. The output is installed at `$GOPATH/src/gomatcha.io/matcha/android/matchabridge.aar`.

    matcha build gomatcha.io/matcha/examples

Now open the sample Android Studio project and hit run!

    <Android Studio location>/bin/studio.sh $GOPATH/src/gomatcha.io/matcha/examples/android-app/SampleApp



<h3>Try it out!</h3>
<ul>
    <li><a href="https://gomatcha.io/guide/installation/">Install</a> the project</li>
    <li>Read the <a href="https://gomatcha.io/guide/getting-started/">Getting Started</a> guide</li>
    <li>Go through some <a href="https://github.com/gomatcha/matcha/tree/master/examples">examples</a></li>
    <li>Learn the <a href="https://gomatcha.io/guide/concepts/">basic concepts</a></li>
</ul>
<h3>Contact us</h3>
<ul>
    <li>Join the Gophers <a href="https://gophers.slack.com/messages/matcha">Slack</a> channel</li>
    <li>Tweet <a href="http://twitter.com/gomatchaio">@gomatcha.io</a> on Twitter</li>
    <li>Star us on <a href="https://github.com/gomatcha/matcha">GitHub</a></li>
    <li><a href="mailto:kevin@gomatcha.io">Email</a> the team</li>
</ul>
