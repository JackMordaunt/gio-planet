package form

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"
)

// Global state for convenient random string generation.
var (
	random = RandomStringer{}
)

// RandomStringer generates random strings.
type RandomStringer struct {
	Src  *rand.Rand
	Buf  []byte
	init sync.Once
}

func (rs *RandomStringer) String() string {
	rs.init.Do(func() {
		rs.Src = rand.New(rand.NewSource(10))
		rs.Buf = make([]byte, 10)
	})
	if _, err := rs.Src.Read(rs.Buf); err != nil {
		panic(fmt.Errorf("generating random string: %w", err))
	}
	return string(rs.Buf)
}

// MockInput counts method calls.
type MockInput struct {
	Calls MockInputCalls
}

type MockInputCalls struct {
	SetText    uint
	Text       uint
	SetError   uint
	ClearError uint
}

func (mi *MockInput) Text() string {
	mi.Calls.Text++
	return ""
}
func (mi *MockInput) SetText(t string) {
	mi.Calls.SetText++
}

func (mi *MockInput) SetError(err string) {
	mi.Calls.SetError++
}

func (mi *MockInput) ClearError() {
	mi.Calls.ClearError++
}

func (mi MockInput) String() string {
	return fmt.Sprintf("%+v", mi.Calls)
}

// StableInput simulates an input that doesn't change.
type StableInput struct {
	MockInput
}

func (si StableInput) Text() string {
	return "stable"
}

// VolatileInput simulates an input that changes constantly.
type VolatileInput struct {
	MockInput
}

func (vi *VolatileInput) Text() string {
	return random.String()
}

// MockValue counts method calls.
type MockValue struct {
	Calls MockValueCalls
}

type MockValueCalls struct {
	To    uint
	From  uint
	Clear uint
}

func (mv *MockValue) From(text string) error {
	mv.Calls.From++
	return nil
}
func (mv *MockValue) To() (string, error) {
	mv.Calls.To++
	return "", nil
}
func (mv *MockValue) Clear() {
	mv.Calls.Clear++
}

func (mv MockValue) String() string {
	return fmt.Sprintf("%+v", mv.Calls)
}

// FailingValue simulates a value that fails.
type FailingValue struct {
	MockValue
}

func (fv *FailingValue) From(string) error {
	return fmt.Errorf("failed")
}

func (fv *FailingValue) To() (string, error) {
	return "", fmt.Errorf("failed")
}

// TestFormFieldError ensures that inputs get signalled with an error when the corresponding value
// fails to validate.
func TestForm(t *testing.T) {
	for _, tt := range []struct {
		Label  string
		Fields []Field
		Want   []Field
	}{
		{
			Label: "set error when validation fails",
			Fields: []Field{
				{
					Value: &FailingValue{},
					Input: &MockInput{},
				},
			},
			Want: []Field{
				{
					Value: &FailingValue{},
					Input: &MockInput{
						Calls: MockInputCalls{
							Text:     1,
							SetError: 1,
						},
					},
				},
			},
		},
		{
			Label: "clear error when validation succeeds",
			Fields: []Field{
				{
					Value: &MockValue{},
					Input: &MockInput{},
				},
			},
			Want: []Field{
				{
					Value: &MockValue{
						Calls: MockValueCalls{
							To: 1,
						},
					},
					Input: &MockInput{
						Calls: MockInputCalls{
							Text:       2,
							ClearError: 1,
							SetText:    1,
						},
					},
				},
			},
		},
		{
			Label: "value is validated only when the input changes",
			Fields: []Field{
				// Changing field: "From" called once.
				{
					Value: &MockValue{},
					Input: &VolatileInput{},
				},
				// Stable field: "From" called nonce.
				{
					Value: &MockValue{},
					Input: &StableInput{},
				},
			},
			Want: []Field{
				{
					Value: &MockValue{
						Calls: MockValueCalls{
							To:   1,
							From: 1,
						},
					},
					Input: &VolatileInput{
						MockInput{
							Calls: MockInputCalls{
								Text:       0,
								SetText:    1,
								ClearError: 2,
							},
						},
					},
				},
				{
					Value: &MockValue{
						Calls: MockValueCalls{
							To: 1,
						},
					},
					Input: &StableInput{
						MockInput{
							Calls: MockInputCalls{
								SetText:    1,
								ClearError: 1,
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tt.Label, func(t *testing.T) {
			var form Form
			form.Load(tt.Fields)
			form.Validate()
			for ii := 0; ii < len(tt.Fields); ii++ {
				var (
					got  = tt.Fields[ii]
					want = tt.Want[ii]
				)
				if !reflect.DeepEqual(want, got) {
					t.Fatalf("\nfield: %d \n want: %+v \n  got: %+v\n", ii, want, got)
				}
			}
		})
	}
}
