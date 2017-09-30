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

### Installation

Matcha requires the following components to be installed. Unfortunately only macOS is supported at this time.

* Go 1.8
* Xcode 8.3
* Android Studio 2.3
* NDK
* Protobuf 3.3 - ObjC and Java

To start, fetch the project and install the matcha command.

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

For Android simply open the sample Android Studio project and hit run!

    open -a /Applications/Android\ Studio.app $GOPATH/src/gomatcha.io/matcha/examples/android-app/SampleApp

You can try out other
[examples](https://github.com/gomatcha/matcha/tree/master/examples) by replacing
`"gomatcha.io/matcha/examples/settings New"` in AppDelegate.m and MainActivity.java
with the name of the example.


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
