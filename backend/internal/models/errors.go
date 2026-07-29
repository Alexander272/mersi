package models

import (
	"errors"
	"net/http"
)

type HTTPError interface {
	error
	Status() int
	Code() string
	Message() string
}

type DomainError struct {
	err     error
	status  int
	code    string
	message string
}

func (e *DomainError) Error() string {
	return e.err.Error()
}

func (e *DomainError) Unwrap() error {
	return e.err
}

func (e *DomainError) Status() int {
	return e.status
}

func (e *DomainError) Code() string {
	return e.code
}

func (e *DomainError) Message() string {
	return e.message
}

func NewDomainError(err error, status int, code string, message string) *DomainError {
	return &DomainError{err: err, status: status, code: code, message: message}
}

var (
	// 404 Not Found
	ErrNoRows    = NewDomainError(errors.New("row not found"), http.StatusNotFound, "NF001", "Запись не найдена")
	ErrNoChannel = NewDomainError(errors.New("channel not found"), http.StatusBadRequest, "NF002", "Канал для получения уведомлений не указан")

	// 400 Bad Request
	ErrNoResponsible = NewDomainError(errors.New("responsible not found"), http.StatusBadRequest, "BR001", "Ответственный не указан")
	ErrNoInstrument  = NewDomainError(errors.New("instrument not found"), http.StatusBadRequest, "BR002", "Инструмент не найден")
	ErrNotValid      = NewDomainError(errors.New("data is not valid"), http.StatusBadRequest, "BR003", "Отправлены некорректные данные")
	ErrInvalidInput  = NewDomainError(errors.New("invalid input"), http.StatusBadRequest, "BR007", "Невалидные данные")

	// 409 Conflict
	ErrEmployeeAlreadyExists      = NewDomainError(errors.New("employee already exists in this department"), http.StatusConflict, "AE003", "Сотрудник с таким именем уже существует в данном подразделении")
	ErrRepairAlreadyExists        = NewDomainError(errors.New("repair already exists for this period"), http.StatusConflict, "AE004", "Ремонт с таким периодом уже существует")
	ErrPreservationAlreadyExists  = NewDomainError(errors.New("preservation already exists for this date"), http.StatusConflict, "AE005", "Консервация с такой датой начала уже существует")
	ErrTransferToDepAlreadyExists = NewDomainError(errors.New("transfer to department already exists for this date"), http.StatusConflict, "AE006", "Передача в подразделение с такой датой уже существует")
	ErrTransferToSaveAlreadyExists = NewDomainError(errors.New("transfer to save already exists for this date"), http.StatusConflict, "AE007", "Передача на хранение с такой датой уже существует")
	ErrWriteOffAlreadyExists      = NewDomainError(errors.New("write-off already exists for this date"), http.StatusConflict, "AE008", "Списание с такой датой уже существует")
	ErrVerificationAlreadyExists  = NewDomainError(errors.New("verification already exists for this date"), http.StatusConflict, "AE009", "Поверка с такой датой уже существует")

	// 400 Bad Request (deletion constraints)
	ErrDeleteDepartmentHasInstrument = NewDomainError(errors.New("cannot delete department with instruments"), http.StatusBadRequest, "BR004", "Нельзя удалить подразделение у которого числятся инструменты")
	ErrDeleteEmployeeHasInstrument   = NewDomainError(errors.New("cannot delete employee with instruments"), http.StatusBadRequest, "BR005", "Нельзя удалить работника у которого числятся инструменты")
	ErrDeleteInstrumentAtHolder      = NewDomainError(errors.New("cannot delete instrument at holder"), http.StatusBadRequest, "BR006", "Не удалось удалить инструмент. Нельзя удалить инструмент находящийся у сотрудника.")

	// 400 Bad Request (location-specific)
	ErrNotResponsible       = NewDomainError(errors.New("user is not responsible"), http.StatusBadRequest, "BR008", "Вы не являетесь ответственным")
	ErrCannotMoveInstrument = NewDomainError(errors.New("cannot move instrument"), http.StatusBadRequest, "BR009", "Вы не можете переместить инструмент")
	ErrCannotConfirmReceipt = NewDomainError(errors.New("cannot confirm receipt"), http.StatusBadRequest, "BR010", "Вы не можете подтвердить получение инструментов")
	ErrInstrumentReceived   = NewDomainError(errors.New("instrument already received or not found"), http.StatusBadRequest, "BR011", "Инструмент уже получен или не найден")
	ErrSingleLocationDelete = NewDomainError(errors.New("location is last or not found"), http.StatusBadRequest, "BR012", "Запись не найдена или это единственное перемещение")

	// 400 Bad Request (file/document-specific)
	ErrNoFileInRequest = NewDomainError(errors.New("no file in request"), http.StatusBadRequest, "BR014", "Не удалось получить файлы")

	// 403 Forbidden
	ErrForbidden = NewDomainError(errors.New("forbidden"), http.StatusForbidden, "FR001", "Доступ запрещён")

	// 404 Not Found
	ErrFileNotFound = NewDomainError(errors.New("file not found"), http.StatusNotFound, "NF004", "Файл не найден")

	// 404 Not Found
	ErrEmployeeNotFound = NewDomainError(errors.New("employee not found"), http.StatusNotFound, "NF003", "Сотрудник не найден")

	// 500 Internal Server Error
	ErrInvalidUserType = NewDomainError(errors.New("invalid user type in context"), http.StatusInternalServerError, "S001", "Внутренняя ошибка сервера")

	// 401 Unauthorized
	ErrSessionEmpty = NewDomainError(errors.New("user session not found"), http.StatusUnauthorized, "AU001", "Сессия не найдена")
)
