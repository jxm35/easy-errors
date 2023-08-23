package option

type Option[T any] interface {
	IsNone() bool
	IsPresent() bool
	Get() (T, bool)
	Value() T
	OrElse(other T) T
	Map(f func(T) T) Option[T]
	FlatMap(f func(T) (T, bool)) Option[T]
	FlatMapErr(f func(T) (T, error)) Option[T]
}

type optionImpl[T any] struct {
	val       T
	isPresent bool
}

func (o *optionImpl[T]) Get() (T, bool) {
	return o.val, o.isPresent
}

func (o *optionImpl[T]) Value() T {
	return o.val
}

func (o *optionImpl[T]) OrElse(other T) T {
	if o.isPresent {
		return o.val
	}
	return other
}

func (o *optionImpl[_]) IsPresent() bool {
	return o.isPresent
}

func (o *optionImpl[_]) IsNone() bool {
	return !o.isPresent
}

func (o *optionImpl[T]) Map(f func(T) T) Option[T] {
	if !o.isPresent {
		return o
	}
	return Some(f(o.val))
}

func (o *optionImpl[T]) FlatMap(f func(T) (T, bool)) Option[T] {
	if !o.isPresent {
		return o
	}
	res, success := f(o.val)
	if success {
		return Some(res)
	}
	return None[T]()
}

func (o *optionImpl[T]) FlatMapErr(f func(T) (T, error)) Option[T] {
	if !o.isPresent {
		return o
	}
	res, err := f(o.val)
	if err != nil {
		return None[T]()
	}
	return Some(res)
}

func Some[T any](val T) Option[T] {
	return &optionImpl[T]{
		val:       val,
		isPresent: true,
	}
}

// None returns an empty option
func None[T any]() Option[T] {
	return &optionImpl[T]{
		isPresent: false,
	}
}

func FromTupleOption[T any](val T, take bool) Option[T] {
	if take {
		return Some(val)
	}
	return None[T]()
}

func FromPointerOption[T any](val *T) Option[T] {
	if val == nil {
		return None[T]()
	}
	return Some(*val)
}

func MapOption[A, B any](opA Option[A], f func(A) B) Option[B] {
	if opA.IsNone() {
		return None[B]()
	}
	v := f(opA.Value())
	return Some(v)
}

func MapOptionTuple[A, B any](opA Option[A], f func(A) (B, bool)) Option[B] {
	if opA.IsNone() {
		return None[B]()
	}
	return FromTupleOption(f(opA.Value()))
}

func MapErrOption[A, B any](opA Option[A], f func(A) (B, error)) Option[B] {
	if opA.IsNone() {
		return None[B]()
	}
	v, err := f(opA.Value())
	if err != nil {
		return None[B]()
	}
	return Some(v)
}

func BindOption[A, B any](opA Option[A], f func(A) Option[B]) Option[B] {
	if opA.IsNone() {
		return None[B]()
	}
	return f(opA.Value())
}
