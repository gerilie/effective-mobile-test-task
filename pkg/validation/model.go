package validation

// Resp represents a map of field names to their corresponding validation error messages.
type Resp map[string]string

// ValidationPrefix is a constant string used as a prefix for validation error messages.
const ValidationPrefix = "The field must be"
