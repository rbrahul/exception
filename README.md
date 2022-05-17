## Go Try Catch Exception Handler
By design, Go doesn't offer any mechanism for Exception handling. But Programmers from different backgrounds like Java, C++, Php can be sceptical about the decision. Exception handling with *Try Catch Finally* is well adapted in all the modern languages. To ease the pain, this library offers utility functions for Exception Handling, which will help programmers to write Go code with *Try-Catch-Finally* approach.

### This how you can throw Exception and handle within Catch:

```go
 import(
     e "github.com/rbrahul/exception"
 )

...
    e.Try(func() {
        data := getValue() // get me the value from Allions
        if data != 100 {
		    e.Throw(e.AssertionError("Expected value is not same as 100"))
        }
	})
    .Catch(e.In(e.AssertionErrorType, e.ValueErrorType), func(excep *Exception) {
        fmt.Println("Message:",excep.Message)
        fmt.Println("Exception Type:",excep.Type)
        fmt.Println("Here is the Stack Trace:",excep.StackTrace)
    })
    .Catch(nil, func(excep *Exception) {
        fmt.Println("I'll be executed as fallback:",excep.Message)
    })
    .Finally(func() {
		fmt.Println("I have been executing always to clean the world!")
	})
    .Run()
...
```

### Throwing a custom exception

You have to define a exception variable with ExceptionType.

```go
const SomethingWentWrongError  e.ExceptionType = "SomethingWentWrongError"
```

Now you have to initialize and throw your exception via e.New constructor. You can pass a proper error message as optional argument.

```go
    e.Try(func() {
        e.Throw(e.New(SomethingWentWrongError, "Something went worng!"))
	})
    .Catch(e.In(SomethingWentWrongError), func(excep *Exception) {
        fmt.Println("Message:",excep.Message)
        fmt.Println("Exception Type:",excep.Type)
    })
    .Finally(func() {
		fmt.Println("I'm Gonna fix it!")
	})
    .Run()
```

### You can wrap any panic with try-catch and recover it elegently

```go
    e.Try(func() {
        panic("I'm gonna panic but don't worry")
	})
    .Catch(nil, func(excep *Exception) {
        fmt.Println("I knew you are gonna catch me :p", excep.Message)
    })
    .Run()
```
