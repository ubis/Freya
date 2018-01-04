// Based on https://github.com/cenkalti/rpc2
// Package rpc provides bi-directional RPC client and server
package rpc

import (
	"io"
	"net"
	"reflect"
	"share/event"
	"share/log"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// Precompute the reflect type for error.  Can't use error directly
// because Typeof takes an empty interface value.  This is annoying.
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfClient = reflect.TypeOf((*Client)(nil))

type Server struct {
	handlers map[string]*handler

	Port int
}

type handler struct {
	fn        reflect.Value
	argType   reflect.Type
	replyType reflect.Type
}

// Initializes RPC Server
func (s *Server) Init() {
	log.Info("Initializing RPC Server...")
	s.handlers = make(map[string]*handler)
}

// Registers a new RPC call function
func (s *Server) Register(method string, handlerFunc interface{}) {
	log.Infof("Registered RPC packet `%s`", method)
	addHandler(s.handlers, method, handlerFunc)
}

// Starts up RPC server
func (s *Server) Run() {
	var listen, err = net.Listen("tcp", ":"+strconv.Itoa(s.Port))
	defer listen.Close()

	if err != nil {
		log.Fatal("Error starting RPC Server:", err.Error())
	}

	log.Info("Listening on " + listen.Addr().String() + "...")
	s.accept(listen)
}

// Adds a new RPC call function handler
func addHandler(handlers map[string]*handler, mname string, handlerFunc interface{}) {
	if _, ok := handlers[mname]; ok {
		panic("rpc: multiple registrations for " + mname)
	}

	method := reflect.ValueOf(handlerFunc)
	mtype := method.Type()
	// Method needs three ins: *client, *args, *reply.
	if mtype.NumIn() != 3 {
		log.Panic("method", mname, "has wrong number of ins:", mtype.NumIn())
	}
	// First arg must be a pointer to rpc2.Client.
	clientType := mtype.In(0)
	if clientType.Kind() != reflect.Ptr {
		log.Panic("method", mname, "client type not a pointer:", clientType)
	}
	if clientType != typeOfClient {
		log.Panic("method", mname, "first argument", clientType.String(), "not *rpc2.Client")
	}
	// Second arg need not be a pointer.
	argType := mtype.In(1)
	if !isExportedOrBuiltinType(argType) {
		log.Panic(mname, "argument type not exported:", argType)
	}
	// Third arg must be a pointer.
	replyType := mtype.In(2)
	if replyType.Kind() != reflect.Ptr {
		log.Panic("method", mname, "reply type not a pointer:", replyType)
	}
	// Reply type must be exported.
	if !isExportedOrBuiltinType(replyType) {
		log.Panic("method", mname, "reply type not exported:", replyType)
	}
	// Method needs one out.
	if mtype.NumOut() != 1 {
		log.Panic("method", mname, "has wrong number of outs:", mtype.NumOut())
	}
	// The return type of the method must be error.
	if returnType := mtype.Out(0); returnType != typeOfError {
		log.Panic("method", mname, "returns", returnType.String(), "not error")
	}
	handlers[mname] = &handler{
		fn:        method,
		argType:   argType,
		replyType: replyType,
	}
}

// Checks if this type is exported or a builtin
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// Checks if name starts up with uppercase
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection. Invokes it in a go statement.
func (s *Server) accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatal("rpc.Serve: accept:", err.Error())
		}
		go s.serveConn(conn, conn.RemoteAddr().String())
	}
}

// Runs the server on a single connection. Triggers SyncConnectEvent
func (s *Server) serveConn(conn io.ReadWriteCloser, endpnt string) {
	var codec = newGobCodec(conn)
	defer codec.Close()

	// client also handles the incoming connections
	var client = NewClientWithCodec(codec)

	client.server = true
	client.handlers = s.handlers
	client.endpnt = endpnt
	client.connected = true

	event.Trigger(event.SyncConnectEvent, client)
	client.Run()
}
