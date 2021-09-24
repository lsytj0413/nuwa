package property

import (
	"testing"

	"github.com/lsytj0413/nuwa/xerrors"
	. "github.com/onsi/gomega"
)

func TestSet(t *testing.T) {
	type testCase struct {
		desp   string
		key    string
		value  interface{}
		err    string
		expect map[string]string
	}
	testCases := []testCase{
		{
			desp:  "normal int8",
			key:   "k1",
			value: int8(1),
			err:   "",
			expect: map[string]string{
				"k1": "1",
			},
		},
		{
			desp:  "normal int16",
			key:   "k1",
			value: int16(0x1000),
			err:   "",
			expect: map[string]string{
				"k1": "4096",
			},
		},
		{
			desp:  "normal int32",
			key:   "k1",
			value: int32(0x10000000),
			err:   "",
			expect: map[string]string{
				"k1": "268435456",
			},
		},
		{
			desp:  "normal int64",
			key:   "k1",
			value: int64(0x1000000000000000),
			err:   "",
			expect: map[string]string{
				"k1": "1152921504606846976",
			},
		},
		{
			desp:  "normal int",
			key:   "k1",
			value: int(0x2000),
			err:   "",
			expect: map[string]string{
				"k1": "8192",
			},
		},
		{
			desp:  "normal uint8",
			key:   "k1",
			value: uint8(1),
			err:   "",
			expect: map[string]string{
				"k1": "1",
			},
		},
		{
			desp:  "normal uint16",
			key:   "k1",
			value: uint16(0x1000),
			err:   "",
			expect: map[string]string{
				"k1": "4096",
			},
		},
		{
			desp:  "normal uint32",
			key:   "k1",
			value: uint32(0x10000000),
			err:   "",
			expect: map[string]string{
				"k1": "268435456",
			},
		},
		{
			desp:  "normal uint64",
			key:   "k1",
			value: uint64(0x1000000000000000),
			err:   "",
			expect: map[string]string{
				"k1": "1152921504606846976",
			},
		},
		{
			desp:  "normal uint",
			key:   "k1",
			value: uint(0x2000),
			err:   "",
			expect: map[string]string{
				"k1": "8192",
			},
		},
		{
			desp:  "normal float32",
			key:   "k1",
			value: float32(1.1),
			err:   "",
			expect: map[string]string{
				"k1": "1.1",
			},
		},
		{
			desp:  "normal float64",
			key:   "k1",
			value: float64(1.2),
			err:   "",
			expect: map[string]string{
				"k1": "1.2",
			},
		},
		{
			desp:  "normal bool",
			key:   "k1",
			value: bool(true),
			err:   "",
			expect: map[string]string{
				"k1": "true",
			},
		},
		{
			desp:  "normal array",
			key:   "k1",
			value: [3]int{1, 2, 3},
			err:   "",
			expect: map[string]string{
				"k1[0]": "1",
				"k1[1]": "2",
				"k1[2]": "3",
			},
		},
		{
			desp:  "normal slice",
			key:   "k1",
			value: []int{1, 2, 3},
			err:   "",
			expect: map[string]string{
				"k1[0]": "1",
				"k1[1]": "2",
				"k1[2]": "3",
			},
		},
		{
			desp: "normal map",
			key:  "k1",
			value: map[string]int{
				"k1": 1,
				"k2": 2,
			},
			err: "",
			expect: map[string]string{
				"k1.k1": "1",
				"k1.k2": "2",
			},
		},
		{
			desp: "normal map with int key",
			key:  "k1",
			value: map[int]int{
				1: 1,
				2: 2,
			},
			err: "",
			expect: map[string]string{
				"k1.1": "1",
				"k1.2": "2",
			},
		},
		{
			desp:  "set value for struct failed",
			key:   "k1",
			value: testCase{},
			err:   "Cannot convert value to string",
		},
		{
			desp:  "set value for array struct failed",
			key:   "k1",
			value: [1]testCase{{}},
			err:   `Cannot set val for array/slice index's key 'k1\[0\]'`,
		},
		{
			desp:  "set value for slice struct failed",
			key:   "k1",
			value: []testCase{{}},
			err:   `Cannot set val for array/slice index's key 'k1\[0\]'`,
		},
		{
			desp: "set value for map struct failed",
			key:  "k1",
			value: map[string]testCase{
				"k1": {},
			},
			err: "Cannot set val for map's key 'k1.k1'",
		},
		{
			desp: "set value for map key struct failed",
			key:  "k1",
			value: map[*testCase]string{
				&testCase{}: "",
			},
			err: "Cannot convert map's key",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			p := NewProperties().(propertiesImpl)
			err := p.Set(tc.key, tc.value)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(map[string]string(p)).To(Equal(tc.expect))
		})
	}
}

func TestGet(t *testing.T) {
	type testCase struct {
		desp   string
		p      Properties
		key    string
		err    error
		expect string
	}
	testCases := []testCase{
		{
			desp: "normal get",
			p: propertiesImpl(map[string]string{
				"k1": "v1",
			}),
			key:    "k1",
			err:    nil,
			expect: "v1",
		},
		{
			desp: "key not found",
			p: propertiesImpl(map[string]string{
				"k1": "v1",
			}),
			key:    "k0",
			err:    xerrors.ErrNotFound,
			expect: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual, err := tc.p.Get(tc.key)
			if tc.err != nil {
				g.Expect(err).To(HaveOccurred())
				g.Expect(xerrors.Is(err, tc.err)).To(BeTrue())
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestRetrive(t *testing.T) {
	type testCase struct {
		desp   string
		p      Properties
		key    string
		i      interface{}
		err    string
		expect interface{}
	}
	testCases := []testCase{
		{
			desp: "normal retrive uint8",
			p: propertiesImpl(map[string]string{
				"k1": "1",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i uint8 = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i uint8 = 1
				return &i
			}(),
		},
		{
			desp: "normal retrive uint16",
			p: propertiesImpl(map[string]string{
				"k1": "4096",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i uint16 = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i uint16 = 4096
				return &i
			}(),
		},
		{
			desp: "normal retrive uint32",
			p: propertiesImpl(map[string]string{
				"k1": "268435456",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i uint32 = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i uint32 = 268435456
				return &i
			}(),
		},
		{
			desp: "normal retrive uint64",
			p: propertiesImpl(map[string]string{
				"k1": "1152921504606846976",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i uint64 = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i uint64 = 1152921504606846976
				return &i
			}(),
		},
		{
			desp: "normal retrive uint",
			p: propertiesImpl(map[string]string{
				"k1": "1152921504606846976",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i uint = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i uint = 1152921504606846976
				return &i
			}(),
		},
		{
			desp: "normal retrive int8",
			p: propertiesImpl(map[string]string{
				"k1": "1",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i int8 = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i int8 = 1
				return &i
			}(),
		},
		{
			desp: "normal retrive int16",
			p: propertiesImpl(map[string]string{
				"k1": "4096",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i int16 = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i int16 = 4096
				return &i
			}(),
		},
		{
			desp: "normal retrive int32",
			p: propertiesImpl(map[string]string{
				"k1": "268435456",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i int32 = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i int32 = 268435456
				return &i
			}(),
		},
		{
			desp: "normal retrive int64",
			p: propertiesImpl(map[string]string{
				"k1": "1152921504606846976",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i int64 = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i int64 = 1152921504606846976
				return &i
			}(),
		},
		{
			desp: "normal retrive int",
			p: propertiesImpl(map[string]string{
				"k1": "1152921504606846976",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i int = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i int = 1152921504606846976
				return &i
			}(),
		},
		{
			desp: "normal retrive float32",
			p: propertiesImpl(map[string]string{
				"k1": "1.1",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i float32 = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i float32 = 1.1
				return &i
			}(),
		},
		{
			desp: "normal retrive float64",
			p: propertiesImpl(map[string]string{
				"k1": "1.1",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i float64 = 0
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i float64 = 1.1
				return &i
			}(),
		},
		{
			desp: "normal retrive bool",
			p: propertiesImpl(map[string]string{
				"k1": "true",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i bool = false
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i bool = true
				return &i
			}(),
		},
		{
			desp: "normal retrive string",
			p: propertiesImpl(map[string]string{
				"k1": "true",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i string = ""
				return &i
			}(),
			err: "",
			expect: func() interface{} {
				var i string = "true"
				return &i
			}(),
		},
		{
			desp: "retrive not found key",
			p: propertiesImpl(map[string]string{
				"k1": "true",
				"k2": "2",
			}),
			key: "k3",
			i: func() interface{} {
				var i string = ""
				return &i
			}(),
			err:    "property with key='k3' not found",
			expect: nil,
		},
		{
			desp: "retrive unsupport target type",
			p: propertiesImpl(map[string]string{
				"k1": "true",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i testCase
				return &i
			}(),
			err:    "Cannot retrive value for key 'k1', unsupported target type",
			expect: nil,
		},
		{
			desp: "retrive cannot set value",
			p: propertiesImpl(map[string]string{
				"k1": "true",
				"k2": "2",
			}),
			key: "k1",
			i: func() interface{} {
				var i int
				return i
			}(),
			err:    "The 'int' cannot been set, it must setable",
			expect: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.p.Retrive(tc.key, tc.i)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.i).To(Equal(tc.expect))
		})
	}
}
