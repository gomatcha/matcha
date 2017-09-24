// +build !matcha

package bridge

import "reflect"

// RegisterType registers a go type that can be created from Objc.
//
// Go:
//  type struct Point {
//      X float64
//      Y float64
//  }
//  func (p Point) Print {
//      fmt.Printf("{%v, %v}", p.X, p.Y)
//  }
//  func init() {
//      bridge.RegisterType("gomatcha.io/matcha/layout.Point", reflect.TypeOf(Point{}))
//  }
// Objective-C:
//  MatchaGoValue *a = [[MatchaGoValue alloc] initWithType:@"gomatcha.io/matcha/layout.Point"];
//  a[@"X"] = [[MatchaGoValue alloc] initWithDouble:x];
//  a[@"Y"] = [[MatchaGoValue alloc] initWithDouble:y];
//  [a call:@"Print" args:nil];
func RegisterType(str string, t reflect.Type) {
}

// RegisterFunc registers a function that can be called from ObjC.
//
// Go:
//  func init() {
//      bridge.RegisterFunc("gomatcha.io/matcha/examples/simple Add", func(a, b int) int {
//          return a + b
//      })
//  }
// Objective-C:
//  MatchaGoValue *func = [[MatchaGoValue alloc] initWithFunc:@"gomatcha.io/matcha/examples/simple Add"];
//  MatchaGoValue *a = [[MatchaGoValue alloc] initWithInt:1];
//  MatchaGoValue *b = [[MatchaGoValue alloc] initWithInt:3];
//  MatchaGoValue *c = [func call:nil args:@[a, b]];
//  NSLog(@"1+3=%@", @(c.toLongLong));
func RegisterFunc(str string, f interface{}) {
}
