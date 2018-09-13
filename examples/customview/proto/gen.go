package proto

//go:generate bash -c "( cd $GOPATH/src && protoc --go_out=. github.com/gomatcha/matcha/examples/customview/proto/*.proto )"
//go:generate bash -c "( cd $GOPATH/src && protoc --java_out=github.com/gomatcha/matcha/examples/customview/android/CustomViewLib/customview/src/main/java github.com/gomatcha/matcha/examples/customview/proto/*.proto )"
//go:generate bash -c "( cd $GOPATH/src && protoc --objc_out=github.com/gomatcha/matcha/examples/customview/ios/Protobuf github.com/gomatcha/matcha/examples/customview/proto/*.proto )"
