# form

[![Go Reference](https://pkg.go.dev/badge/git.sr.ht/~jackmordaunt/gio-planet/form.svg)](https://pkg.go.dev/git.sr.ht/~jackmordaunt/gio-planet/form)

Primitives for specifying forms.

The fundamental idea explored here is in separating form concerns.
The three core abstractions are the `Form`, the `Value` and the `Input`.
`Value` and `Input` combine together into a `Field`.

- `Input` abstracts the interactible graphical widget that accepts input.
- `Value` abstracts the transformation of text to structured data.
- `Field` maps an `Input` to a `Value`.
- `Form` provides a consistent api for handling a group of fields, namely batch and realtime validation.

## Status

This repo is experimental, and does not have a stable interface.

Ideally form initialization would be fully declarative, where the zero value is directly usable.
However it's not clear how to achieve such an api.

Some potential avenues to explore:

- code generation: take a struct definition and generate the form code.
- reflection: introspect a struct definition and generate `form.Field`s on the fly.

## Use

Since the api is experimental, there is no "idiomatic usage".

Fundamentally, a form binds inputs to values and there are numerous ways to compose it.

This is an example of one way to use the package: one-time initialization that you can call in an
update function once per frame.

```go
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
```
