package low

func CSel[T any](c bool, t, e T) T {
	if c {
		return t
	} else {
		return e
	}
}
