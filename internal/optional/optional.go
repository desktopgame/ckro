package optional

type Optional[T comparable] struct {
	value  T
	exists bool
}

func Some[T comparable](value T) Optional[T] {
	return Optional[T]{
		value:  value,
		exists: true,
	}
}

func None[T comparable]() Optional[T] {
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

func (o Optional[T]) IsSame(other Optional[T]) bool {
	if o.IsNone() != other.IsNone() {
		return false
	}
	if o.IsNone() == other.IsNone() {
		return true
	}
	return o.value == other.value
}

func (o Optional[T]) IsSome() bool {
	return o.exists
}

func (o Optional[T]) IsNone() bool {
	return !o.exists
}
