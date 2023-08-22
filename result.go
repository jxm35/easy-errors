package easy_errors

type Result[T any] interface {
	GetError() error
	Get() (T, error)
	Value() T
}

type resultImpl[T any] struct {
	val T
	err error
}

func (r *resultImpl[T]) GetError() error {
	return r.err
}
func (r *resultImpl[T]) Get() (T, error) {
	return r.val, r.err
}
func (r *resultImpl[T]) Value() T {
	return r.val
}

func Ok[T any](val T) Result[T] {
	return &resultImpl[T]{
		val: val,
		err: nil,
	}
}
func Err[T any](err error) Result[T] {
	return &resultImpl[T]{
		err: err,
	}
}

func Map[A, B any](resA Result[A], f func(A) B) Result[B] {
	if err := resA.GetError(); err != nil {
		return Err[B](err)
	}
	v := f(resA.Value())
	return Ok(v)
}

func MapErr[A, B any](resA Result[A], f func(A) (B, error)) Result[B] {
	if err := resA.GetError(); err != nil {
		return Err[B](err)
	}
	v, err := f(resA.Value())
	if err != nil {
		return Err[B](err)
	}
	return Ok(v)
}

func Bind[A, B any](resA Result[A], f func(A) Result[B]) Result[B] {
	if err := resA.GetError(); err != nil {
		return Err[B](err)
	}
	return f(resA.Value())
}
