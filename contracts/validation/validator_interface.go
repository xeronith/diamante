package validation

type IValidator interface {
	Validate(interface{}) error
}
