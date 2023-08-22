package results

import "errors"

type Option[T any] interface {
	IsNone() bool
	Get() (T, error)
	Value() T
}

type optionImpl[T any] struct {
	val    T
	isNone bool
}

func (o *optionImpl[T]) Get() (T, error) {
	if o.isNone {
		return o.val, errors.New("option is None")
	}
	return o.val, nil
}

func (o *optionImpl[T]) Value() T {
	return o.val
}

func (o *optionImpl[_]) IsNone() bool {
	return o.isNone
}

func Some[T any](val T) Option[T] {
	return &optionImpl[T]{
		val: val,
	}
}

// None returns an empty option
func None[T any]() Option[T] {
	return &optionImpl[T]{
		isNone: true,
	}
}

func MapOption[A, B any](opA Option[A], f func(A) B) Option[B] {
	if opA.IsNone() {
		return None[B]()
	}
	v := f(opA.Value())
	return Some(v)
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
