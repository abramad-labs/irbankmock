package pointers

func Deref[T any](p *T) (v T, ok bool) {
	if p != nil {
		return *p, true
	}
	return v, false
}

func DerefZero[T any](p *T) (v T) {
	if p != nil {
		return *p
	}
	return v
}
