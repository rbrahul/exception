package exception

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
)

type ExceptionType string

const (
	UnknownErrorType    ExceptionType = "UnknownError"
	IndexErrorType      ExceptionType = "IndexError"
	RuntimeErrorType    ExceptionType = "RuntimeError"
	ValueErrorType      ExceptionType = "ValueError"
	NetworkErrorType    ExceptionType = "NetworkError"
	SyntaxErrorType     ExceptionType = "SyntaxError"
	PermissionErrorType ExceptionType = "PermissionError"
	TimeoutErrorType    ExceptionType = "TimeoutError"
	TypeErrorType       ExceptionType = "TypeError"
	AssertionErrorType  ExceptionType = "AssertionError"
	ConnectionErrorType ExceptionType = "ConnectionError"
	ReferenceErrorType  ExceptionType = "ReferenceError"
	EOFErrorType        ExceptionType = "EOFError"
	LookupErrorType     ExceptionType = "LookupError"
)

var exceptionErrorMap map[ExceptionType]string = map[ExceptionType]string{
	UnknownErrorType:    "Unknown Error",
	IndexErrorType:      "Index Error",
	ValueErrorType:      "Value Error",
	NetworkErrorType:    "Network Error",
	SyntaxErrorType:     "Syntax Error",
	PermissionErrorType: "Permission Error",
	TimeoutErrorType:    "Timeout Error",
	TypeErrorType:       "Type Error",
	AssertionErrorType:  "Assertion Error",
	ConnectionErrorType: "Connection Error",
	ReferenceErrorType:  "Reference Error",
	EOFErrorType:        "EOF Error",
	LookupErrorType:     "Lookup Error",
}

type Exception struct {
	Message    string
	Type       ExceptionType
	StackTrace string
}

//New is constructor which is used to create a new Exception
func New(exceptionType ExceptionType, args ...interface{}) *Exception {
	message, ok := exceptionErrorMap[exceptionType]
	if !ok {
		message = string(exceptionType)
	}
	if len(args) > 0 {
		message = fmt.Sprintf("%v", args[0])
	}
	return &Exception{
		Message: message,
		Type:    exceptionType,
	}
}

//AssertionError should be used to throw an exception indicating assertion failure. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func AssertionError(args ...interface{}) *Exception {
	return New(AssertionErrorType, args...)
}

//IndexError should be used to throw an exception indicating index is out of bound failure. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func IndexError(args ...interface{}) *Exception {
	return New(IndexErrorType, args...)
}

//ConnectionError should be used to throw an exception indicating connection related failure. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func ConnectionError(args ...interface{}) *Exception {
	return New(ConnectionErrorType, args...)
}

//EOFError should be used to throw an exception indicating end of file related issue. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func EOFError(args ...interface{}) *Exception {
	return New(EOFErrorType, args...)
}

//LookupError should be used to throw an exception indicating the key is unavailable in a map. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func LookupError(args ...interface{}) *Exception {
	return New(LookupErrorType, args...)
}

//NetworkError should be used to throw an exception indicating the network related issue. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func NetworkError(args ...interface{}) *Exception {
	return New(NetworkErrorType, args...)
}

//PermissionError should be used to throw an exception indicating the permission related issue. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func PermissionError(args ...interface{}) *Exception {
	return New(PermissionErrorType, args...)
}

//ReferenceError should be used to throw an exception indicating the reference related of issue. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func ReferenceError(args ...interface{}) *Exception {
	return New(ReferenceErrorType, args...)
}

//SyntaxError should be used to throw an exception indicating the reference related of issue. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func SyntaxError(args ...interface{}) *Exception {
	return New(SyntaxErrorType, args...)
}

//TypeError should be used to throw an exception indicating the type is not as expected. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func TypeError(args ...interface{}) *Exception {
	return New(TimeoutErrorType, args...)
}

//TimeoutError should be used to throw an exception indicating the Timeout related issue. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func TimeoutError(args ...interface{}) *Exception {
	return New(TimeoutErrorType, args...)
}

//ValueError should be used to throw an exception indicating the value is not in correct format. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func ValueError(args ...interface{}) *Exception {
	return New(ValueErrorType, args...)
}

//Throw is used to throw an exception. You can use some builtin Exceptions such as LookupError, PermissionError, NetworkError etc. which are offered by this library or you can also create your custom Exception and throw. It accepts error message as an optional argument. Otherwise the ExceptionType is used as default error message.
func Throw(exp *Exception) {
	message := exp.Message
	errorMsg := fmt.Sprintf("Message::%s||Exception::%s", message, exp.Type)
	panic(errorMsg)
}

// In accepts a variable number of ExceptionType as arguments. It creates a list ExceptionType that will be used for the mathcing of exception.
func In(exceptionTypes ...ExceptionType) []ExceptionType {
	return exceptionTypes
}

type catchblockEntry struct {
	Exceptions []ExceptionType
	Handler    func(arg *Exception)
}

type exceptionHandler struct {
	exception      *Exception
	tryHandler     func()
	catchHandlers  []catchblockEntry
	finallyHandler func()
}

//Try executes your code and finds if there is any panic or exception and passes the exception to catch block
func Try(cb func()) *exceptionHandler {
	resp := &exceptionHandler{exception: &Exception{Message: ""}, catchHandlers: []catchblockEntry{}, finallyHandler: nil}
	resp.tryHandler = cb
	return resp
}

//Catch gets executed if any panic or exception occurred inside Try. You can control the execution of any Catch block by passing a e.In() matcher which listens for certain Exception to be thrown. If you pass nil as first argument then the Catch block will be executed as default if there is no matching Catch block found.
func (c *exceptionHandler) Catch(exceptionTypes []ExceptionType, cb func(excep *Exception)) *exceptionHandler {
	c.catchHandlers = append(c.catchHandlers, catchblockEntry{Exceptions: exceptionTypes, Handler: cb})
	return c
}

//Finally is executed always even if the Try block Succeeds or Fails. But it won't be executed if there is a uncaught or unhandled Exception
func (c *exceptionHandler) Finally(cb func()) *exceptionHandler {
	c.finallyHandler = cb
	return c
}

//Run must be called to run the Try-Catch-Finally handlers. It should be always invoked at the end of the chain operations. It triggers the execution of Exception Handling flow.
func (c *exceptionHandler) Run() {
	c.executeTry()
	c.executeCatchHanlder()
	c.executeFinally()
}

func (c *exceptionHandler) executeTry() {
	defer func() {
		err := recover()
		if err != nil {
			value := reflect.ValueOf(err)
			errorMessage := fmt.Sprintf("%v", value)
			if !strings.Contains(errorMessage, "||") {
				errorMessage = fmt.Sprintf("Message::%s||Exception::%s", errorMessage, "RuntimeError")

			}
			c.exception = &Exception{
				Message:    errorMessage,
				StackTrace: string(debug.Stack()),
			}
		}
	}()
	c.tryHandler()
}

func (c *exceptionHandler) getExceptionType() string {
	messageItems := strings.Split(c.exception.Message, "||")
	if len(messageItems) > 0 {
		exceptionPart := messageItems[1]
		parts := strings.Split(exceptionPart, "::")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return ""
}

func (c *exceptionHandler) getMessage() string {
	messageItems := strings.Split(c.exception.Message, "||")
	if len(messageItems) > 0 {
		messagePart := messageItems[0]
		splitedMessage := strings.Split(messagePart, "::")
		if len(splitedMessage) > 1 {
			return splitedMessage[1]
		}
	}
	return ""
}

func (c *exceptionHandler) executeCatchHanlder() {
	if len(c.catchHandlers) == 0 {
		return
	}
	if c.exception != nil && len(c.exception.Message) > 0 {
		var catchHandlerExecuted bool
		var defaultHandler func(_ *Exception)
		for _, handler := range c.catchHandlers {
			if handler.Exceptions != nil && len(handler.Exceptions) > 0 {
				for _, exceptionType := range handler.Exceptions {
					exceptionTypePart := c.getExceptionType()
					if exceptionTypePart == string(exceptionType) {
						c.exception.Type = ExceptionType(exceptionTypePart)
						c.exception.Message = c.getMessage()
						handler.Handler(c.exception)
						catchHandlerExecuted = true
						return
					}
				}
			} else {
				defaultHandler = handler.Handler
			}
		}
		if !catchHandlerExecuted && defaultHandler != nil {
			c.exception.Type = ExceptionType(c.getExceptionType())
			c.exception.Message = c.getMessage()
			defaultHandler(c.exception)
		}
	}
}

func (c *exceptionHandler) executeFinally() {
	if c.finallyHandler != nil {
		c.finallyHandler()
	}
}
