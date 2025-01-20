// Code generated by "github.com/zzztttkkk/lion/enums", DO NOT EDIT
// Code generated @ 1737355349

package vld

import "fmt"



func (ev ErrorKind) String() string {
	switch(ev){
		
		case ErrorKindMissingRequired : {
			return "MissingRequired"
		}
		case ErrorKindCustom : {
			return "Custom"
		}
		case ErrorKindIntLtMin : {
			return "IntLtMin"
		}
		case ErrorKindIntGtMax : {
			return "IntGtMax"
		}
		case ErrorKindIntNotInRange : {
			return "IntNotInRange"
		}
		case ErrorKindTimeTooEarly : {
			return "TimeTooEarly"
		}
		case ErrorKindTimeTooLate : {
			return "TimeTooLate"
		}
		case ErrorKindStringTooLong : {
			return "StringTooLong"
		}
		case ErrorKindStringTooShort : {
			return "StringTooShort"
		}
		case ErrorKindStringNotMatched : {
			return "StringNotMatched"
		}
		case ErrorKindStringNotInRanges : {
			return "StringNotInRanges"
		}
		case ErrorKindContainerSizeTooLarge : {
			return "ContainerSizeTooLarge"
		}
		case ErrorKindContainerSizeTooSmall : {
			return "ContainerSizeTooSmall"
		}
		case ErrorKindNilPointer : {
			return "NilPointer"
		}
		case ErrorKindNilSlice : {
			return "NilSlice"
		}
		case ErrorKindNilMap : {
			return "NilMap"
		}
		default: {
			panic(fmt.Errorf("vld.ErrorKind: unknown enum value, %d", ev))
		} 
	}
}


func init(){
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindMissingRequired)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindCustom)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindIntLtMin)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindIntGtMax)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindIntNotInRange)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindTimeTooEarly)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindTimeTooLate)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindStringTooLong)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindStringTooShort)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindStringNotMatched)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindStringNotInRanges)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindContainerSizeTooLarge)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindContainerSizeTooSmall)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindNilPointer)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindNilSlice)
	
	AllErrorKinds = append(AllErrorKinds, ErrorKindNilMap)
	
}




