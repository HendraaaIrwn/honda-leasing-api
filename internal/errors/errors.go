package errs

import "errors"

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
	ErrInvalidSort       = errors.New("invalid sort parameters")
	ErrInvalidSearch     = errors.New("search name cannot be empty")
	ErrInvalidEmail      = errors.New("email already exist")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidDecision   = errors.New("invalid workflow decision")

	// when create
	ErrCreateUser = errors.New("error when create user")
	ErrAssignRole = errors.New("error when assigned user role")
	ErrTrxUser    = errors.New("transaction errors")

	// leasing workflow
	ErrMotorUnitNotReady       = errors.New("motor unit is not ready")
	ErrInvalidStatusTransition = errors.New("invalid contract status transition")
	ErrContractNotDraft        = errors.New("contract must be in draft status")
	ErrContractNotApproved     = errors.New("contract must be in approved status")
	ErrDPOutOfRange            = errors.New("down payment is outside allowed product range")
	ErrInvalidPaymentAmount    = errors.New("invalid payment amount")
)
