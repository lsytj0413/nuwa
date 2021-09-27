package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/lsytj0413/nuwa"
	"github.com/lsytj0413/nuwa/app"
)

type Runner1 struct {
}

func (r *Runner1) Run(ctx context.Context) {
	fmt.Println("Runner1.Run called")
}

type Runner2 struct {
}

func (r *Runner2) Run(ctx context.Context) {
	fmt.Println("Runner2.Run called")
}

func main() {
	ap := app.NewApplication()

	ap.RegisterBeanDefinition("runner1", nuwa.MustNewBeanDefinition(reflect.TypeOf((*Runner1)(nil))))
	ap.RegisterBeanDefinition("runner2", nuwa.MustNewBeanDefinition(reflect.TypeOf((*Runner2)(nil))))

	err := ap.Run()
	if err != nil {
		panic(err)
	}
}
