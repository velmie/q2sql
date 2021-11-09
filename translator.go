package q2sql

// Translator translates given strings from one format to another
// for example "fieldName" could be translated to "field_name" or "table.field_name"
type Translator func(in []string) (out []string, err error)

// MapTranslator creates format translator function based on the passed map
func MapTranslator(m map[string]string) Translator {
	return func(in []string) (out []string, err error) {
		out = make([]string, len(in))
		for i := 0; i < len(in); i++ {
			entry, ok := m[in[i]]
			if !ok {
				return nil, &TranslationError{
					Entry:   in[i],
					Message: "translation is not found",
				}
			}
			out[i] = entry
		}
		return out, nil
	}
}
