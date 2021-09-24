package nuwa

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

type BeanOnlyBeanField struct {
	B2 *BeanOnlyPropertyField
}

type BeanOnlyPropertyField struct {
	V    int
	VPtr *int
}

type BeanMixField struct {
	B2 *BeanOnlyPropertyField
	V  int
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
			desp: "normal get bean only property field",
			beanNames: []string{
				"bean2",
			},
			beanDefinitions: []BeanDefinition{
				&beanDefinitionImpl{
					Typ: reflect.TypeOf((*BeanOnlyPropertyField)(nil)),
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
			expect: &BeanOnlyPropertyField{
				V: 100,
				VPtr: func() *int {
					v := int(200)
					return &v
				}(),
			},
		},
		{
			desp: "normal get bean only bean field",
			beanNames: []string{
				"bean1",
				"bean2",
			},
			beanDefinitions: []BeanDefinition{
				&beanDefinitionImpl{
					Typ: reflect.TypeOf((*BeanOnlyBeanField)(nil)),
					fieldDescriptors: []FieldDescriptor{
						{
							FieldIndex: 0,
							Name:       "B2",
							Typ:        reflect.TypeOf((*BeanOnlyPropertyField)(nil)),
							Bean: &BeanFieldDescriptor{
								Name: "bean2",
							},
						},
					},
				},
				&beanDefinitionImpl{
					Typ: reflect.TypeOf((*BeanOnlyPropertyField)(nil)),
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
			beanName: "bean1",
			err:      "",
			expect: &BeanOnlyBeanField{
				B2: &BeanOnlyPropertyField{
					V: 100,
					VPtr: func() *int {
						v := int(200)
						return &v
					}(),
				},
			},
		},
		{
			desp: "normal get bean mix field",
			beanNames: []string{
				"bean2",
				"beanmix",
			},
			beanDefinitions: []BeanDefinition{
				&beanDefinitionImpl{
					Typ: reflect.TypeOf((*BeanOnlyPropertyField)(nil)),
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
				&beanDefinitionImpl{
					Typ: reflect.TypeOf((*BeanMixField)(nil)),
					fieldDescriptors: []FieldDescriptor{
						{
							FieldIndex: 0,
							Name:       "B2",
							Typ:        reflect.TypeOf((*BeanOnlyPropertyField)(nil)),
							Bean: &BeanFieldDescriptor{
								Name: "bean2",
							},
						},
						{
							FieldIndex: 1,
							Name:       "VPtr",
							Typ:        reflect.TypeOf(int(0)),
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
			beanName: "beanmix",
			err:      "",
			expect: &BeanMixField{
				B2: &BeanOnlyPropertyField{
					V: 100,
					VPtr: func() *int {
						v := int(200)
						return &v
					}(),
				},
				V: int(200),
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
