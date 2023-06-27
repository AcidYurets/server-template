package err_const

import "fmt"

// Константные ошибки

// Const тип используемый для константных ошибок, позволяет избегать возможных мутаций значений ошибок.
// Не рекомендуется использовать их для создания ошибок в рамках бизнес-логики.
type Const string

func (e Const) Error() string {
	return string(e)
}

const (
	ErrJsonUnMarshal          = Const("не удалось декодировать JSON")
	ErrJsonMarshal            = Const("не удалось упаковать данные в JSON")
	ErrDatabaseRecordNotFound = Const("запись не найдена")
	ErrMissingUser            = Const("пользователь не указан")

	ErrIdMissing    = Const("не указан Id объекта")
	ErrUuidMissing  = Const("не указан Uuid объекта")
	ErrIdValidate   = Const("неверно указан Id объекта")
	ErrUuidValidate = Const("неверно указан Uuid объекта")

	ErrMissingToken = Const("не указан токен авторизации")
)

type PanicError struct {
	Message string
	Detail  string
}

func FromPanic(panicInstance interface{}) *PanicError {
	// Ошибка, полученная из паники
	var panicMsg string

	switch msg := panicInstance.(type) {
	case string:
		panicMsg = msg
	case error:
		panicMsg = msg.Error()

	default:
		panicMsg = fmt.Sprint(panicInstance)
	}

	return &PanicError{
		Message: "внутрення ошибка системы",
		Detail:  panicMsg,
	}
}

func (r *PanicError) Error() string {
	return r.Message
}
