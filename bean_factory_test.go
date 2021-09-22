package nuwa

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

type Bean1 struct {
	Bean2 *Bean2
	// V2    Bean2
}

func (b Bean1) F1() {

}

func (b *Bean1) F2() {

}

type Bean2 struct {
	V int
}

func TestGetBean(t *testing.T) {
	g := NewWithT(t)
	f := NewBeanFactory()
	var v int
	f.RegisterBeanDefinition("test", &beanDefinitionImpl{
		Typ: reflect.TypeOf(v),
	})

	actual, err := f.GetBean("test")
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(actual).To(Equal(int(0)))
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

func TestGetBeanStructPtr(t *testing.T) {
	g := NewWithT(t)
	f := NewBeanFactory()

	b1 := &Bean1{}
	b2 := &Bean2{}
	f.RegisterBeanDefinition("bean1", &beanDefinitionImpl{
		Typ: reflect.TypeOf(b1),
		fieldDescriptors: []FieldDescriptor{
			{
				FieldIndex: 0,
				Name:       "Bean2",
				Typ:        reflect.TypeOf(b2),
			},
			// {
			// 	FieldIndex: 1,
			// 	Name:       "Bean2",
			// 	Typ:        reflect.TypeOf(*b2),
			// 	Bean: &BeanFieldDescriptor{
			// 		Name: "bean3",
			// 	},
			// },
		},
	})
	f.RegisterBeanDefinition("bean2", &beanDefinitionImpl{
		Typ: reflect.TypeOf(b2),
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
	f.RegisterBeanDefinition("bean3", &beanDefinitionImpl{
		Typ: reflect.TypeOf(*b2),
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

	actual, err := f.GetBean("bean1")
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(actual).To(Equal(&Bean1{
		Bean2: &Bean2{
			V: 1,
		},
	}))
}

func printMethod(v interface{}) {
	t := reflect.TypeOf(v)

	fmt.Println(t.String())
	for i := 0; i < t.NumMethod(); i++ {
		fmt.Println(t.Method(i).Name)
	}
}

func TestT(t *testing.T) {
	b := Bean1{}
	bp := &Bean1{}

	printMethod(b)
	printMethod(bp)
}
