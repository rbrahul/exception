package exception

import "testing"

func TestTryGetsExecuted(t *testing.T) {
	var hasTryExecuted bool
	Try(func() {
		hasTryExecuted = true
	}).
		Catch(nil, func(excep *Exception) {}).
		Run()

	if hasTryExecuted == false {
		t.Fatal("Try was not executed")
	}
}

func TestCatchGetsExecutedForAllErrorsIfExceptionThrownInTry(t *testing.T) {
	var exceptionType ExceptionType
	var hasCaughtException bool
	Try(func() {
		Throw(ReferenceError("Dummy error text"))
	}).
		Catch(nil, func(excep *Exception) {
			hasCaughtException = true
			exceptionType = excep.Type
		}).Run()

	if hasCaughtException == false {
		t.Fatal("Try was not executed")
	}
	if exceptionType != ReferenceErrorType {
		t.Fatalf("Expecting %s but found %s", ReferenceErrorType, exceptionType)
	}
}

func TestFinallyGetsExecutedAlways(t *testing.T) {
	var hasFinallyExecuted bool
	Try(func() {
		Throw(NetworkError())
	}).
		Catch(nil, func(excep *Exception) {}).
		Finally(func() {
			hasFinallyExecuted = true
		}).
		Run()

	if hasFinallyExecuted == false {
		t.Fatal("Finally was not executed")
	}
}

func TestExpectedCatchBlockGetsExecutedForDefinedExceptionType(t *testing.T) {
	var thrownExceptionType ExceptionType
	var caughtFrom string
	Try(func() {
		user := map[string]string{
			"name": "John Doe",
		}
		_, ok := user["email"]
		if !ok {
			Throw(LookupError("Email doesn't exist"))
		}
	}).
		Catch(In(LookupErrorType), func(excep *Exception) {
			thrownExceptionType = excep.Type
			caughtFrom = "LookupErrorHandler"
		}).
		Catch(In(ReferenceErrorType, IndexErrorType), func(excep *Exception) {
			thrownExceptionType = excep.Type
			caughtFrom = "ReferenceErrorHandler"
		}).
		Catch(nil, func(excep *Exception) {
			thrownExceptionType = excep.Type
			caughtFrom = "DefaultExceptionHandler"
		}).Run()

	if thrownExceptionType != LookupErrorType {
		t.Fatalf("Expecting Exception type to be %s but found %s", LookupErrorType, thrownExceptionType)
	}

	if caughtFrom != "LookupErrorHandler" {
		t.Fatalf("Expecting Catch block to be %s but found %s", "LookupErrorHandler", caughtFrom)
	}
}

func TestDefaultCatchBlockGetsExecutedForUnmatchedException(t *testing.T) {
	var thrownExceptionType ExceptionType
	var caughtFrom string
	Try(func() {
		Throw(New(UnkownExceptionType, "Unkown Error"))
	}).
		Catch(In(LookupErrorType), func(excep *Exception) {
			thrownExceptionType = excep.Type
			caughtFrom = "LookupErrorHandler"
		}).
		Catch(In(ReferenceErrorType, IndexErrorType), func(excep *Exception) {
			thrownExceptionType = excep.Type
			caughtFrom = "ReferenceErrorHandler"
		}).
		Catch(nil, func(excep *Exception) {
			thrownExceptionType = excep.Type
			caughtFrom = "DefaultExceptionHandler"
		}).Run()

	if thrownExceptionType != UnkownExceptionType {
		t.Fatalf("Expecting Exception type to be %s but found %s", UnkownExceptionType, thrownExceptionType)
	}

	if caughtFrom != "DefaultExceptionHandler" {
		t.Fatalf("Expecting Catch block to be %s but found %s", "LookupErrorHandler", caughtFrom)
	}
}

func TestAllPanicGetsRecoveredWithinTryCatch(t *testing.T) {
	var panicAttack = "Something went very wrong!"
	var caughtMessage string
	Try(func() {
		panic(panicAttack)
	}).
		Catch(nil, func(excep *Exception) {
			caughtMessage = excep.Message
		}).
		Run()

	if caughtMessage != panicAttack {
		t.Fatal("Could not recover panic attack!")
	}
}

func TestNestedExceptionWasHandledAsExpected(t *testing.T) {
	var firstThrownExceptionType ExceptionType
	var secondThrownExceptionType ExceptionType
	var hasCaughtNestedException bool
	var CustomExceptionType ExceptionType = "CustomException"

	Try(func() {
		Throw(ReferenceError("Dummy error text"))
	}).
		Catch(nil, func(excep1 *Exception) {
			Try(func() {
				// trying to save the world
				Throw(New(CustomExceptionType, "Custome Error Message"))
			}).Catch(In(CustomExceptionType), func(excep2 *Exception) {
				firstThrownExceptionType = excep1.Type
				secondThrownExceptionType = excep2.Type
				hasCaughtNestedException = true
			}).Run()
		}).Run()

	if hasCaughtNestedException == false {
		t.Fatal("Nested exception was not handled")
	}

	if firstThrownExceptionType != ReferenceErrorType {
		t.Fatalf("Expecting first exception to be %s but found %s", ReferenceErrorType, firstThrownExceptionType)
	}

	if secondThrownExceptionType != CustomExceptionType {
		t.Fatalf("Expecting second exception to be %s but found %s", CustomExceptionType, secondThrownExceptionType)
	}
}
