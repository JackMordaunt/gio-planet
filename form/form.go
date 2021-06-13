package form

// Value implements a bi-directional mapping between textual data and structured data. Value handles
// data validation, which is expresssed by the error return.
type Value interface {
	// To converts the value into textual form.
	To() (text string, err error)
	// From parses the value from textual form.
	From(text string) (err error)
	// Clear resets the value.
	Clear()
}

// Input objects display text, and handle user input events.
// Input will typically be implemented by a graphic widget, such as the
// `gioui.org/x/component#TextField`.
type Input interface {
	Text() string
	SetText(string)
	SetError(string)
	ClearError()
}

// Field binds a Value to an Input.
type Field struct {
	Value Value
	Input Input
}

// Validate the field by running the text through the Valuer.
// Precise validation logic is implemented by the Valuer.
// Returns a boolean indicating success.
func (field *Field) Validate() bool {
	err := field.Value.From(field.Input.Text())
	if err != nil {
		field.Input.SetError(err.Error())
	} else {
		field.Input.ClearError()
	}
	return err == nil
}

// Form exercises field bindings.
//
// There's two primary ways of using a Form:
// 1. Realtime Validation
// 2. Batch Validation
//
// Realtime validation, expressed by `Form.Validate`, processes all non-zero and changed fields.
// That is, they must have a value and that value must be different since the last validation.
//
// Batch validation, expressed by `Form.Submit`, processes _all_ fields including zero-value fields.
//
// The semantic difference is that an unsubmitted zero-value input is not in an errored state
// because the user hasn't input a value yet. If you attempt to submit that zero-value input then it
// submission error and the field is now in an errored state.
//
// Realtime validation is useful for providing fast feedback on input events. You can create a
// `form.Value` that maps to some complex data source. For example, you can run queries on the fly
// to figure out if an entry exists as the user is typing.
//
// Batch validation is useful for quickly testing if the whole form is valid before using the field
// data.
//
// `form.Submit` must be called to synchronize field values.
// If not called, the values stored in each `form.Value` could be different to what is displayed in
// the graphical input.
type Form struct {
	Fields []Field
	// cache contains the previous contents of each field to detect changes.
	cache []string
}

// Load values into inputs.
func (f *Form) Load(fields []Field) {
	if len(fields) > 0 {
		f.Fields = fields
	}
	f.cache = make([]string, len(fields))
	for ii, field := range f.Fields {
		if text, err := field.Value.To(); err != nil {
			field.Input.SetError(err.Error())
		} else {
			f.cache[ii] = field.Input.Text()
			field.Input.ClearError()
			field.Input.SetText(text)
		}
	}
}

// Submit batch validates the fields and returns a boolean indication success.
// If true, all the fields validated and you can use the data.
func (f *Form) Submit() (ok bool) {
	ok = true
	for _, field := range f.Fields {
		if !field.Validate() {
			ok = false
		}
	}
	return ok
}

// Validate form fields.
// Can be used per frame for realtime validation.
func (f *Form) Validate() {
	for ii, field := range f.Fields {
		var (
			text    = field.Input.Text()
			changed = f.cache[ii] != text
		)
		if changed {
			f.cache[ii] = text
			field.Validate()
		}
	}
}

func (f *Form) Clear() {
	for ii, field := range f.Fields {
		field.Value.Clear()
		text, _ := field.Value.To()
		field.Input.ClearError()
		field.Input.SetText(text)
		f.cache[ii] = text
	}
}
