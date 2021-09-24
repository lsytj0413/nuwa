package utils

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestIndirectToInterface(t *testing.T) {
	type testCase struct {
		desp   string
		v      interface{}
		expect interface{}
	}
	testCases := []testCase{
		{
			desp:   "normal value",
			v:      int(100),
			expect: int(100),
		},
		{
			desp:   "indirect reflect.Value",
			v:      reflect.ValueOf(int(100)),
			expect: int(100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual := IndirectToInterface(tc.v)
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestIndirectToValue(t *testing.T) {
	type testCase struct {
		desp   string
		v      interface{}
		expect interface{}
	}
	testCases := []testCase{
		{
			desp:   "normal value",
			v:      int(100),
			expect: int(100),
		},
		{
			desp:   "indirect reflect.Value",
			v:      reflect.ValueOf(int(100)),
			expect: int(100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual := IndirectToValue(tc.v)
			g.Expect(actual.Interface()).To(Equal(tc.expect))
		})
	}
}

func TestIndirectToSetableValue(t *testing.T) {
	type testStruct struct {
		V *int
	}

	type testCase struct {
		desp   string
		v      interface{}
		err    string
		expect interface{}
	}
	testCases := []testCase{
		{
			desp: "normal value",
			v: func() interface{} {
				var i int = 0
				return &i
			}(),
			expect: func() interface{} {
				var i int = 0
				return i
			}(),
		},
		{
			desp: "nil pointer",
			v: func() interface{} {
				st := &testStruct{}
				return reflect.ValueOf(st).Elem().Field(0)
			}(),
			expect: func() interface{} {
				var i int = 0
				return i
			}(),
		},
		{
			desp: "cannot set",
			v:    int(0),
			err:  "cannot been set, it must setable",
		},
		{
			desp: "cannot set pointer",
			v:    (*int)(nil),
			err:  "cannot been set, it must setable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual, err := IndirectToSetableValue(tc.v)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual.Interface()).To(Equal(tc.expect))
		})
	}
}

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
