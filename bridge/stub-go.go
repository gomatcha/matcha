// +build !matcha

package bridge

// RegisterFunc registers a function that can be called from ObjC.
//
// Go:
//  func init() {
//      bridge.RegisterFunc("github.com/gomatcha/matcha/examples/simple Add", func(a, b int) int {
//          return a + b
//      })
//  }
// Objective-C:
//  MatchaGoValue *func = [[MatchaGoValue alloc] initWithFunc:@"github.com/gomatcha/matcha/examples/simple Add"];
//  MatchaGoValue *a = [[MatchaGoValue alloc] initWithInt:1];
//  MatchaGoValue *b = [[MatchaGoValue alloc] initWithInt:3];
//  MatchaGoValue *c = [func call:nil args:@[a, b]];
//  NSLog(@"1+3=%@", @(c.toLongLong));
func RegisterFunc(str string, f interface{}) {
}
