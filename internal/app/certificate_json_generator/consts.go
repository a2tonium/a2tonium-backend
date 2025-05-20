package certificate_json_generator

import (
	"errors"
	"path/filepath"
)

var (
	templatesPattern = filepath.Join("internal", "app", "script_generator", "templates", "*.gohtml")
)

var (
	ErrInvalidVusNumber   = errors.New("INVALID_VUS_NUMBER")
	ErrInvalidDuration    = errors.New("INVALID_DURATION")
	ErrInvalidDestination = errors.New("INVALID_DESTINATION")
	ErrInvalidHttpMethod  = errors.New("INVALID_HTTP_METHOD")
)
