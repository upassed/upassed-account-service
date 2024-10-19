// code generated by protoc-gen-validate. DO NOT EDIT.
// source: teacher.proto

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

// define the regex for a UUID once up-front
var _teacher_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on TeacherCreateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TeacherCreateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TeacherCreateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TeacherCreateRequestMultiError, or nil if none found.
func (m *TeacherCreateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *TeacherCreateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetFirstName()); l < 2 || l > 20 {
		err := TeacherCreateRequestValidationError{
			field:  "FirstName",
			reason: "value length must be between 2 and 20 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetLastName()); l < 3 || l > 30 {
		err := TeacherCreateRequestValidationError{
			field:  "LastName",
			reason: "value length must be between 3 and 30 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetMiddleName()) > 20 {
		err := TeacherCreateRequestValidationError{
			field:  "MiddleName",
			reason: "value length must be at most 20 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateEmail(m.GetReportEmail()); err != nil {
		err = TeacherCreateRequestValidationError{
			field:  "ReportEmail",
			reason: "value must be a valid email address",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetUsername()); l < 5 || l > 20 {
		err := TeacherCreateRequestValidationError{
			field:  "Username",
			reason: "value length must be between 5 and 20 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return TeacherCreateRequestMultiError(errors)
	}

	return nil
}

func (m *TeacherCreateRequest) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

func (m *TeacherCreateRequest) _validateEmail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	addr = a.Address

	if len(addr) > 254 {
		return errors.New("email addresses cannot exceed 254 characters")
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return errors.New("email address local phrase cannot exceed 64 characters")
	}

	return m._validateHostname(parts[1])
}

// TeacherCreateRequestMultiError is an error wrapping multiple validation
// errors returned by TeacherCreateRequest.ValidateAll() if the designated
// constraints aren't met.
type TeacherCreateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TeacherCreateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TeacherCreateRequestMultiError) AllErrors() []error { return m }

// TeacherCreateRequestValidationError is the validation error returned by
// TeacherCreateRequest.Validate if the designated constraints aren't met.
type TeacherCreateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TeacherCreateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TeacherCreateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TeacherCreateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TeacherCreateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TeacherCreateRequestValidationError) ErrorName() string {
	return "TeacherCreateRequestValidationError"
}

// Error satisfies the builtin error interface
func (e TeacherCreateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTeacherCreateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TeacherCreateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TeacherCreateRequestValidationError{}

// Validate checks the field values on TeacherCreateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TeacherCreateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TeacherCreateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TeacherCreateResponseMultiError, or nil if none found.
func (m *TeacherCreateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *TeacherCreateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CreatedTeacherId

	if len(errors) > 0 {
		return TeacherCreateResponseMultiError(errors)
	}

	return nil
}

// TeacherCreateResponseMultiError is an error wrapping multiple validation
// errors returned by TeacherCreateResponse.ValidateAll() if the designated
// constraints aren't met.
type TeacherCreateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TeacherCreateResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TeacherCreateResponseMultiError) AllErrors() []error { return m }

// TeacherCreateResponseValidationError is the validation error returned by
// TeacherCreateResponse.Validate if the designated constraints aren't met.
type TeacherCreateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TeacherCreateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TeacherCreateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TeacherCreateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TeacherCreateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TeacherCreateResponseValidationError) ErrorName() string {
	return "TeacherCreateResponseValidationError"
}

// Error satisfies the builtin error interface
func (e TeacherCreateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTeacherCreateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TeacherCreateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TeacherCreateResponseValidationError{}

// Validate checks the field values on TeacherFindByIDRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TeacherFindByIDRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TeacherFindByIDRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TeacherFindByIDRequestMultiError, or nil if none found.
func (m *TeacherFindByIDRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *TeacherFindByIDRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetTeacherId()); err != nil {
		err = TeacherFindByIDRequestValidationError{
			field:  "TeacherId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return TeacherFindByIDRequestMultiError(errors)
	}

	return nil
}

func (m *TeacherFindByIDRequest) _validateUuid(uuid string) error {
	if matched := _teacher_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// TeacherFindByIDRequestMultiError is an error wrapping multiple validation
// errors returned by TeacherFindByIDRequest.ValidateAll() if the designated
// constraints aren't met.
type TeacherFindByIDRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TeacherFindByIDRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TeacherFindByIDRequestMultiError) AllErrors() []error { return m }

// TeacherFindByIDRequestValidationError is the validation error returned by
// TeacherFindByIDRequest.Validate if the designated constraints aren't met.
type TeacherFindByIDRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TeacherFindByIDRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TeacherFindByIDRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TeacherFindByIDRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TeacherFindByIDRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TeacherFindByIDRequestValidationError) ErrorName() string {
	return "TeacherFindByIDRequestValidationError"
}

// Error satisfies the builtin error interface
func (e TeacherFindByIDRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTeacherFindByIDRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TeacherFindByIDRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TeacherFindByIDRequestValidationError{}

// Validate checks the field values on TeacherFindByIDResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TeacherFindByIDResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TeacherFindByIDResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TeacherFindByIDResponseMultiError, or nil if none found.
func (m *TeacherFindByIDResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *TeacherFindByIDResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetTeacher()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TeacherFindByIDResponseValidationError{
					field:  "Teacher",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TeacherFindByIDResponseValidationError{
					field:  "Teacher",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTeacher()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TeacherFindByIDResponseValidationError{
				field:  "Teacher",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return TeacherFindByIDResponseMultiError(errors)
	}

	return nil
}

// TeacherFindByIDResponseMultiError is an error wrapping multiple validation
// errors returned by TeacherFindByIDResponse.ValidateAll() if the designated
// constraints aren't met.
type TeacherFindByIDResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TeacherFindByIDResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TeacherFindByIDResponseMultiError) AllErrors() []error { return m }

// TeacherFindByIDResponseValidationError is the validation error returned by
// TeacherFindByIDResponse.Validate if the designated constraints aren't met.
type TeacherFindByIDResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TeacherFindByIDResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TeacherFindByIDResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TeacherFindByIDResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TeacherFindByIDResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TeacherFindByIDResponseValidationError) ErrorName() string {
	return "TeacherFindByIDResponseValidationError"
}

// Error satisfies the builtin error interface
func (e TeacherFindByIDResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTeacherFindByIDResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TeacherFindByIDResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TeacherFindByIDResponseValidationError{}
