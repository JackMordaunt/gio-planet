/*
Package form implements primitives that reduce form boilerplate by allowing the caller to specify
their fields exactly once.

All values are processed via a chain of transformations that map text into a structured value, and
visa versa. Each transformation is encapsulated in a `form.Value` implementation, for instance a
`value.Int` will transform text into a Go integer and signal any errors that occur during that
transformation.

Forms are initialized once with all the fields via a call to `form.Load`.
Each field binds an input to a value.

By contention, value objects depend on pointer variables, this means you can simply point into a
predefined "model" struct. Once the form is submitted, the model will contain the validated values
ready to use. However this is only a convention, a value object can arbitrarily handle it's internal
state.

The following is an example of one way to use the form:

	type Person struct {
		Age    int
		Name   string
		Salary float64
	}

	type PersonForm struct {
		form.Form
		// Model contains the structured data.
		Model Person
		// Inputs contains the input state.
		Inputs struct {
			Age    component.TextField
			Name   component.TextField
			Salary component.TextField
		}
		init sync.Once
	}

	func (pf *PersonForm) Update() {
		pf.init.Do(func() {
			pf.Form.Load([]form.Field{
				{
					Value: value.Int{Value: &pf.Model.Age, Default: 18},
					Input: &pf.Inputs.Age,
				},
				{
					Value: value.Required{value.Text{Value: &pf.Model.Name, Default: ""}},
					Input: &pf.Inputs.Name,
				},
				{
					Value: value.Float{Value: &pf.Model.Salary, Default: 0},
					Input: &pf.Inputs.Salary,
				},
			})
		})
	}
*/
package form
