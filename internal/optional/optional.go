package optional

type Optional[T any] struct {
	value  T
	exists bool
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value:  value,
		exists: true,
	}
}

func None[T any]() Optional[T] {
	return Optional[T]{
		exists: false,
	}
}

func (o Optional[T]) TryValue() (T, bool) {
	return o.value, o.exists
}

func (o Optional[T]) Value() T {
	return o.value
}

func (o Optional[T]) OrElse(defaultValue T) T {
	if o.exists {
		return o.value
	}
	return defaultValue
}

func (o Optional[T]) IsSome() bool {
	return o.exists
}

func (o Optional[T]) IsNone() bool {
	return !o.exists
}
