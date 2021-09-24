package nuwa

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewValue(t *testing.T) {
	type testCase struct {
		desp   string
		typ    reflect.Type
		expect interface{}
		isErr  bool
	}

	vint := int(0)
	vtestCase := testCase{}
	pvint := &vint
	pvtestCase := &vtestCase
	testCases := []testCase{
		{
			desp:   "normal primitive type",
			typ:    reflect.TypeOf(int(1)),
			expect: int(0),
			isErr:  false,
		},
		{
			desp:   "normal struct type",
			typ:    reflect.TypeOf(testCase{}),
			expect: testCase{},
			isErr:  false,
		},
		{
			desp:   "normal ptr to primitive type",
			typ:    reflect.TypeOf(&vint),
			expect: &vint,
			isErr:  false,
		},
		{
			desp:   "normal ptr to struct type",
			typ:    reflect.TypeOf(&vtestCase),
			expect: &vtestCase,
			isErr:  false,
		},
		{
			desp:   "err ptr to primitive type",
			typ:    reflect.TypeOf(&pvint),
			expect: nil,
			isErr:  true,
		},
		{
			desp:   "normal ptr to struct type",
			typ:    reflect.TypeOf(&pvtestCase),
			expect: nil,
			isErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual, err := NewValue(tc.typ)
			if tc.isErr {
				g.Expect(err).To(HaveOccurred())
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual.Interface()).To(Equal(tc.expect))
		})
	}
}
