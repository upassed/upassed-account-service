// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: model.proto

package client

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on GroupDTO with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GroupDTO) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GroupDTO with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GroupDTOMultiError, or nil
// if none found.
func (m *GroupDTO) ValidateAll() error {
	return m.validate(true)
}

func (m *GroupDTO) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for SpecializationCode

	// no validation rules for GroupNumber

	if len(errors) > 0 {
		return GroupDTOMultiError(errors)
	}

	return nil
}

// GroupDTOMultiError is an error wrapping multiple validation errors returned
// by GroupDTO.ValidateAll() if the designated constraints aren't met.
type GroupDTOMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GroupDTOMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GroupDTOMultiError) AllErrors() []error { return m }

// GroupDTOValidationError is the validation error returned by
// GroupDTO.Validate if the designated constraints aren't met.
type GroupDTOValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GroupDTOValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GroupDTOValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GroupDTOValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GroupDTOValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GroupDTOValidationError) ErrorName() string { return "GroupDTOValidationError" }

// Error satisfies the builtin error interface
func (e GroupDTOValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGroupDTO.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GroupDTOValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GroupDTOValidationError{}

// Validate checks the field values on StudentDTO with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *StudentDTO) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StudentDTO with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in StudentDTOMultiError, or
// nil if none found.
func (m *StudentDTO) ValidateAll() error {
	return m.validate(true)
}

func (m *StudentDTO) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for FirstName

	// no validation rules for LastName

	// no validation rules for MiddleName

	// no validation rules for EducationalEmail

	// no validation rules for Username

	if all {
		switch v := interface{}(m.GetGroup()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, StudentDTOValidationError{
					field:  "Group",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, StudentDTOValidationError{
					field:  "Group",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetGroup()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return StudentDTOValidationError{
				field:  "Group",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return StudentDTOMultiError(errors)
	}

	return nil
}

// StudentDTOMultiError is an error wrapping multiple validation errors
// returned by StudentDTO.ValidateAll() if the designated constraints aren't met.
type StudentDTOMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StudentDTOMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StudentDTOMultiError) AllErrors() []error { return m }

// StudentDTOValidationError is the validation error returned by
// StudentDTO.Validate if the designated constraints aren't met.
type StudentDTOValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StudentDTOValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StudentDTOValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StudentDTOValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StudentDTOValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StudentDTOValidationError) ErrorName() string { return "StudentDTOValidationError" }

// Error satisfies the builtin error interface
func (e StudentDTOValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStudentDTO.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StudentDTOValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StudentDTOValidationError{}

// Validate checks the field values on TeacherDTO with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *TeacherDTO) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TeacherDTO with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in TeacherDTOMultiError, or
// nil if none found.
func (m *TeacherDTO) ValidateAll() error {
	return m.validate(true)
}

func (m *TeacherDTO) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for FirstName

	// no validation rules for LastName

	// no validation rules for MiddleName

	// no validation rules for ReportEmail

	// no validation rules for Username

	if len(errors) > 0 {
		return TeacherDTOMultiError(errors)
	}

	return nil
}

// TeacherDTOMultiError is an error wrapping multiple validation errors
// returned by TeacherDTO.ValidateAll() if the designated constraints aren't met.
type TeacherDTOMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TeacherDTOMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TeacherDTOMultiError) AllErrors() []error { return m }

// TeacherDTOValidationError is the validation error returned by
// TeacherDTO.Validate if the designated constraints aren't met.
type TeacherDTOValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TeacherDTOValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TeacherDTOValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TeacherDTOValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TeacherDTOValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TeacherDTOValidationError) ErrorName() string { return "TeacherDTOValidationError" }

// Error satisfies the builtin error interface
func (e TeacherDTOValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTeacherDTO.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TeacherDTOValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TeacherDTOValidationError{}
