package exception

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
)

type ExceptionType string

const (
	UnkownExceptionType ExceptionType = "Unkown"
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
	UnkownExceptionType: "Unkown Exception",
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

func AssertionError(args ...interface{}) *Exception {
	return New(AssertionErrorType, args...)
}

func IndexError(args ...interface{}) *Exception {
	return New(IndexErrorType, args...)
}

func ConnectionError(args ...interface{}) *Exception {
	return New(ConnectionErrorType, args...)
}

func EOFError(args ...interface{}) *Exception {
	return New(EOFErrorType, args...)
}

func LookupError(args ...interface{}) *Exception {
	return New(LookupErrorType, args...)
}

func NetworkError(args ...interface{}) *Exception {
	return New(NetworkErrorType, args...)
}

func PermissionError(args ...interface{}) *Exception {
	return New(PermissionErrorType, args...)
}

func ReferenceError(args ...interface{}) *Exception {
	return New(ReferenceErrorType, args...)
}

func SyntaxError(args ...interface{}) *Exception {
	return New(SyntaxErrorType, args...)
}

func TypeError(args ...interface{}) *Exception {
	return New(TimeoutErrorType, args...)
}

func TimeoutError(args ...interface{}) *Exception {
	return New(TimeoutErrorType, args...)
}

func ValueError(args ...interface{}) *Exception {
	return New(ValueErrorType, args...)
}

func Throw(exp *Exception) {
	message := exp.Message
	errorMsg := fmt.Sprintf("Message::%s||Exception::%s", message, exp.Type)
	panic(errorMsg)
}

func In(exceptionTypes ...ExceptionType) []ExceptionType {
	return exceptionTypes
}

type Exception struct {
	Message    string
	Type       ExceptionType
	StackTrace string
}

type CatchblockEntry struct {
	Exceptions []ExceptionType
	Handler    func(arg *Exception)
}

type ExceptionHandler struct {
	exception      *Exception
	tryHandler     func()
	catchHandlers  []CatchblockEntry
	finallyHandler func()
}

func (c *ExceptionHandler) Catch(exceptionTypes []ExceptionType, cb func(excep *Exception)) *ExceptionHandler {
	c.catchHandlers = append(c.catchHandlers, CatchblockEntry{Exceptions: exceptionTypes, Handler: cb})
	return c
}

func (c *ExceptionHandler) Finally(cb func()) *ExceptionHandler {
	c.finallyHandler = cb
	return c
}

func (c *ExceptionHandler) Run() {
	c.executeTry()
	c.executeCatchHanlder()
	c.executeFinally()
}

func (c *ExceptionHandler) executeTry() {
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

func (c *ExceptionHandler) getExceptionType() string {
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

func (c *ExceptionHandler) getMessage() string {
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

func (c *ExceptionHandler) executeCatchHanlder() {
	if len(c.catchHandlers) == 0 {
		return
	}
	if c.exception != nil && len(c.exception.Message) > 0 {
		catchHandlerExecuted := false
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

func (c *ExceptionHandler) executeFinally() {
	if c.finallyHandler != nil {
		c.finallyHandler()
	}
}

func Try(cb func()) *ExceptionHandler {
	resp := &ExceptionHandler{exception: &Exception{Message: ""}, catchHandlers: []CatchblockEntry{}, finallyHandler: nil}
	resp.tryHandler = cb
	return resp
}
