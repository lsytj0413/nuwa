package nuwa

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

type Bean2 struct {
	V    int
	VPtr *int
}

func TestGetBean(t *testing.T) {
	type testCase struct {
		desp            string
		beanNames       []string
		beanDefinitions []BeanDefinition
		valNames        []string
		valValues       []interface{}
		beanName        string
		err             string
		expect          interface{}
	}
	testCases := []testCase{
		{
			desp: "normal get bean",
			beanNames: []string{
				"bean2",
			},
			beanDefinitions: []BeanDefinition{
				&beanDefinitionImpl{
					Typ: reflect.TypeOf((*Bean2)(nil)),
					fieldDescriptors: []FieldDescriptor{
						{
							FieldIndex: 0,
							Name:       "V",
							Typ:        reflect.TypeOf(int(0)),
							Property: &PropertyFieldDescriptor{
								Name: "val",
							},
						},
						{
							FieldIndex: 1,
							Name:       "VPtr",
							Typ:        reflect.TypeOf((*int)(nil)),
							Property: &PropertyFieldDescriptor{
								Name: "valPtr",
							},
						},
					},
				},
			},
			valNames: []string{
				"val",
				"valPtr",
			},
			valValues: func() []interface{} {
				v := "100"
				valPtr := "200"
				return []interface{}{v, valPtr}
			}(),
			beanName: "bean2",
			err:      "",
			expect: &Bean2{
				V: 100,
				VPtr: func() *int {
					v := int(200)
					return &v
				}(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			f := NewBeanFactory()
			for i := range tc.beanNames {
				err := f.RegisterBeanDefinition(tc.beanNames[i], tc.beanDefinitions[i])
				g.Expect(err).ToNot(HaveOccurred())
			}
			for i := range tc.valNames {
				err := f.Set(tc.valNames[i], tc.valValues[i])
				g.Expect(err).ToNot(HaveOccurred())
			}

			actual, err := f.GetBean(tc.beanName)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestGetBeanWithProperty(t *testing.T) {
	g := NewWithT(t)
	f := NewBeanFactory()

	f.RegisterBeanDefinition("bean2", &beanDefinitionImpl{
		Typ: reflect.TypeOf((*Bean2)(nil)),
		fieldDescriptors: []FieldDescriptor{
			{
				FieldIndex: 0,
				Name:       "V",
				Typ:        reflect.TypeOf(int(0)),
				Property: &PropertyFieldDescriptor{
					Name: "val",
				},
			},
		},
	})
	f.Set("val", 1)

	actual, err := f.GetBean("bean2")
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(actual).To(Equal(&Bean2{
		V: 1,
	}))
}
