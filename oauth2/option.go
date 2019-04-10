package oauth2

type EmptyOption struct{}

func (EmptyOption) apply(*AdalHandler) (err error) { return }

type funcOption struct {
	f func(*AdalHandler) error
}

func (fo *funcOption) apply(h *AdalHandler) error {

	return fo.f(h)
}

func newFuncOption(f func(*AdalHandler) error) *funcOption {
	return &funcOption{
		f: f,
	}
}
