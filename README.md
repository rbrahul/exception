## Go Try Catch Exception Handler
By design, Go doesn't offer any mechanism for Exception handling. But Programmers from different backgrounds like Java, C++, Php might find it sceptical. Exception handling with *Try Catch Finally* is well adapted in all the modern languages. To ease the pain, this library offers utility functions for Exception Handling, which will help programmers to write Go code with *Try-Catch-Finally* approach.

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

## Documentation will be updated soon!
