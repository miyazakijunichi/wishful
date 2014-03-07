package useful

import (
	. "github.com/SimonRichardson/wishful/wishful"
)

type Promise struct {
	Fork func(resolve func(x AnyVal) AnyVal) AnyVal
}

func NewPromise(f func(resolve func(x AnyVal) AnyVal) AnyVal) Promise {
	return Promise{
		Fork: f,
	}
}

func (x Promise) Of(v AnyVal) Point {
	return Promise{func(resolve func(x AnyVal) AnyVal) AnyVal {
		return resolve(v)
	}}
}

func (x Promise) Ap(v Applicative) Applicative {
	return Promise{func(resolve func(x AnyVal) AnyVal) AnyVal {
		return x.Fork(func(f AnyVal) AnyVal {
			fun := v.(Functor)
			pro := fun.Map(func(x AnyVal) AnyVal {
				fun := NewFunction(f)
				res, _ := fun.Call(x)
				return res
			})
			return pro.(Promise).Fork(resolve)
		})
	}}
}

func (x Promise) Chain(f func(v AnyVal) Monad) Monad {
	return Promise{func(resolve func(x AnyVal) AnyVal) AnyVal {
		return x.Fork(func(a AnyVal) AnyVal {
			p := f(a).(Promise)
			return p.Fork(resolve)
		})
	}}
}

func (x Promise) Map(f func(v AnyVal) AnyVal) Functor {
	return Promise{func(resolve func(v AnyVal) AnyVal) AnyVal {
		return x.Fork(func(a AnyVal) AnyVal {
			return resolve(f(a))
		})
	}}
}

func (x Promise) Extract() AnyVal {
	return x.Fork(Identity)
}

func (x Promise) Extend(f func(p Comonad) AnyVal) Comonad {
	return x.Map(func(y AnyVal) AnyVal {
		fun := NewFunction(f)
		res, _ := fun.Call(x.Of(y))
		return res
	}).(Comonad)
}
