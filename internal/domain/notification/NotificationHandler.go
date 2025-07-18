package notification

import (
	"bytes"
	"fmt"
	"sync"
	templateTx "text/template"
)

type NotificationData struct {
	Entity string
	Data   string
}

type DomainError struct {
	Message string
	Code    string
	Cause   error
}

func (e *DomainError) Error() string {
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

var templateCache = struct {
	sync.RWMutex
	templates map[string]*templateTx.Template
}{
	templates: make(map[string]*templateTx.Template),
}

func getOrCreateTemplate(templateText string) *templateTx.Template {
	templateCache.RLock()
	tmpl, exists := templateCache.templates[templateText]
	templateCache.RUnlock()

	if exists {
		return tmpl
	}

	templateCache.Lock()
	defer templateCache.Unlock()

	if tmpl, exists = templateCache.templates[templateText]; exists {
		return tmpl
	}

	tmpl = templateTx.Must(templateTx.New("Notific").Parse(templateText))
	templateCache.templates[templateText] = tmpl
	return tmpl
}

func CreateCustomNotification(template string, entity string, data interface{}) error {
	var buffer bytes.Buffer
	var errorMessage string

	switch v := data.(type) {
	case error:
		errorMessage = v.Error()
	case string:
		errorMessage = v
	case nil:
		errorMessage = ""
	default:
		errorMessage = fmt.Sprintf("%v", v)
	}

	notificationData := NotificationData{
		Entity: entity,
		Data:   errorMessage,
	}

	tmpl := getOrCreateTemplate(template)
	if err := tmpl.Execute(&buffer, notificationData); err != nil {
		return &DomainError{
			Message: "Erro ao processar template de notificação",
			Code:    "TEMPLATE_ERROR",
			Cause:   err,
		}
	}

	var cause error
	if err, ok := data.(error); ok {
		cause = err
	}

	return &DomainError{
		Message: buffer.String(),
		Code:    getErrorCode(template),
		Cause:   cause,
	}
}

func CreateSimpleNotification(template string, data error) error {
	var buffer bytes.Buffer

	notificationData := NotificationData{
		Entity: "",
		Data:   data.Error(),
	}

	if err := getOrCreateTemplate(template).Execute(&buffer, notificationData); err != nil {
		return &DomainError{
			Message: "Erro ao processar template de notificação",
			Code:    "TEMPLATE_ERROR",
			Cause:   err,
		}
	}

	return &DomainError{
		Message: buffer.String(),
		Code:    getErrorCode(template),
		Cause:   data,
	}
}

func CreateNotification(template string) error {
	var buffer bytes.Buffer

	if err := getOrCreateTemplate(template).Execute(&buffer, ""); err != nil {
		return err
	}

	return &DomainError{
		Message: buffer.String(),
		Code:    getErrorCode(template),
	}
}

func getErrorCode(template string) string {
	switch template {
	case NotFound:
		return "NOT_FOUND"
	case InvalidData:
		return "INVALID_DATA"
	case InvalidMethod:
		return "INVALID_METHOD"
	case ScanErrorRepository:
		return "SCAN_ERROR"
	case FindErrorRepository:
		return "FIND_ERROR"
	case FindAllErrorRepository:
		return "FIND_ALL_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}
