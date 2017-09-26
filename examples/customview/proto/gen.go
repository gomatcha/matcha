package proto

//go:generate bash -c "( cd $GOPATH/src && protoc --go_out=. gomatcha.io/matcha/examples/customview/proto/*.proto )"
//go:generate bash -c "( cd $GOPATH/src && protoc --java_out=gomatcha.io/matcha/examples/customview/android/CustomViewLib/customview/src/main/java gomatcha.io/matcha/examples/customview/proto/*.proto )"
//go:generate bash -c "( cd $GOPATH/src && protoc --objc_out=gomatcha.io/matcha/examples/customview/ios/Protobuf gomatcha.io/matcha/examples/customview/proto/*.proto )"
