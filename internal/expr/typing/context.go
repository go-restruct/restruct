package typing

// Context represents the context used for type checking and inference.
type Context struct {
	Self   Type
	Global map[string]Type
}

// Resolve resolves an identifier contextually.
func (context Context) Resolve(ident string) (Type, error) {
	if typ, ok := context.Global[ident]; ok {
		return typ, nil
	}

	return context.Self.Field(ident)
}
