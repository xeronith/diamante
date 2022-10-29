package localization

// TODO: It's more performant with integer keys
type (
	Key      = string
	Value    = string
	Resource = map[Key]Value

	ILocalizer interface {
		Register(Resource)
		Get(Key) Value
	}
)
