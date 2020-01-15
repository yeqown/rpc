package rpc

import (
	"log"
	"reflect"
	"unicode"
	"unicode/utf8"
)

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	// sync.Mutex // protects counters
	// numCalls   uint
}

type service struct {
	name   string
	rcvr   reflect.Value
	typ    reflect.Type
	method map[string]*methodType
}

func (s *service) call(mtype *methodType, argv, replyv reflect.Value) error {
	function := mtype.method.Func
	returnValues := function.Call([]reflect.Value{s.rcvr, argv, replyv})
	DebugF("[service.call] argv: %v, replyv: %v", argv, replyv)
	if i := returnValues[0].Interface(); i != nil {
		return i.(error)
	}
	return nil
}

// type isExported
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// suitableMethods get all method of registering-type
// into a map[string]*methodType
func suitableMethods(typ reflect.Type) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		if mt := suitableMethod(method); mt != nil {
			methods[method.Name] = mt
		}
	}
	return methods
}

func suitableMethodWtihName(typ reflect.Type, methodName string) *methodType {
	if method, ex := typ.MethodByName(methodName); ex {
		return suitableMethod(method)
	}
	log.Println("rpc.RegisterName: has no such method")
	return nil
}

func suitableMethod(method reflect.Method) *methodType {
	mtype := method.Type
	mname := method.Name

	// Method must be exported.
	if method.PkgPath != "" {
		return nil
	}
	// Method needs three ins: receiver, *args, *reply.
	if mtype.NumIn() != 3 {
		log.Printf("rpc.Register: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
		return nil
	}
	// First arg need not be a pointer.
	argType := mtype.In(1)
	if !isExportedOrBuiltinType(argType) {
		log.Printf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
		return nil
	}
	// Second arg must be a pointer.
	replyType := mtype.In(2)
	if replyType.Kind() != reflect.Ptr {
		log.Printf("rpc.Register: reply type of method %q is not a pointer: %q\n", mname, replyType)
		return nil
	}
	// Reply type must be exported.
	if !isExportedOrBuiltinType(replyType) {
		log.Printf("rpc.Register: reply type of method %q is not exported: %q\n", mname, replyType)
		return nil
	}
	// Method needs one out.
	if mtype.NumOut() != 1 {
		log.Printf("rpc.Register: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
		return nil
	}
	// The return type of the method must be error.
	if returnType := mtype.Out(0); returnType != typeOfError {
		log.Printf("rpc.Register: return type of method %q is %q, must be error\n", mname, returnType)
		return nil
	}
	return &methodType{method: method, ArgType: argType, ReplyType: replyType}
}
