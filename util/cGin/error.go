package cGin

// CustomError custom error
type CustomError struct {
	HTTPCode int
	Code     int
	Message  string
}

func (cErr CustomError) Error() string {
	return cErr.Message
}

// type CheckCustomErrorFunc func(*CustomError) bool

// // CheckCustomError compare input err's chain matches CustomError
// // if err's chain is contains CustomError, check the CustomError with checkFunc
// func CheckCustomError(err error, checkFuncs ...CheckCustomErrorFunc) bool {
// 	valid, cErr := RecursivelyUnwrap(err)
// 	if !valid {
// 		return false
// 	}

// 	for _, checkFunc := range checkFuncs {
// 		if checkFuncValid := checkFunc(cErr); !checkFuncValid {
// 			return false
// 		}
// 	}
// 	return true
// }

// // Is implements errors interface Is(error) bool
// func (cErr CustomError) Is(target error) bool {
// 	causeErr := errors.Cause(target)
// 	if _cErr, ok := causeErr.(*CustomError); ok {
// 		return cErr.Code == _cErr.Code
// 	}

// 	if _cErr, ok := causeErr.(CustomError); ok {
// 		return cErr.Code == _cErr.Code
// 	}

// 	return false
// }

// // As implements errors interface As(interface) bool (go 1.15)
// func (cErr CustomError) As(target interface{}) bool {
// 	if _, ok := target.(CustomError); ok {
// 		return ok
// 	}

// 	if _, ok := target.(*CustomError); ok {
// 		return ok
// 	}

// 	return false
// }

// /* RecursivelyUnwrap
//  * check the err's type in err's chain is contains CustomError
//  * if true, unwrap error until the CustomError return
//  * if false, return nil
//  */
// func RecursivelyUnwrap(err error) (bool, *CustomError) {
// 	return recursivelyUnwrap(err, false, 0)

// }
// func recursivelyUnwrap(err error, skipCheckType bool, unwrapTimes uint) (bool, *CustomError) {
// 	if !skipCheckType && !errors.As(err, &CustomError{}) {
// 		return false, nil
// 	}

// 	var unwrapErr error
// 	unwrapErr = err
// 	if unwrapTimes != 0 {
// 		unwrapErr = errors.Unwrap(err)
// 	}

// 	cusErr, ok := unwrapErr.(CustomError)
// 	if ok {
// 		return true, &cusErr
// 	}

// 	ptrCusErr, ok := unwrapErr.(*CustomError)
// 	if ok {
// 		return true, ptrCusErr
// 	}
// 	unwrapTimes++
// 	return recursivelyUnwrap(unwrapErr, true, unwrapTimes)
// }
