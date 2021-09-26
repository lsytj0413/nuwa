package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/lsytj0413/nuwa"
)

// Application is the interface for app
type Application interface {
	Run() error
	Shutdown()

	nuwa.BeanFactory
}

// AppRunner is the interface for runner
type AppRunner interface {
	// Run will been called when the application is ready.
	Run(ctx context.Context)
}

type nuwaApplication struct {
	exitChan chan struct{}

	nuwa.BeanFactory
}

// NewApplication return the application
func NewApplication() Application {
	return &nuwaApplication{
		exitChan:    make(chan struct{}),
		BeanFactory: nuwa.NewBeanFactory(),
	}
}

func (a *nuwaApplication) Run() error {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		sig := <-ch
		a.shutdownWithMessage(fmt.Sprintf("Receive signal: %v", sig))
	}()

	// Prepare the application
	runners := []AppRunner{}
	err := a.RetriveBeans(&runners)
	if err != nil {
		return err
	}

	for _, r := range runners {
		r.Run(context.TODO())
	}

	<-a.exitChan
	return nil
}

func (a *nuwaApplication) Shutdown() {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	fnName := "<unknown>"
	if fn != nil {
		fnName = fn.Name()
	}
	caller := fmt.Sprintf("%s:%v %s", file, line, fnName)

	a.shutdownWithMessage(fmt.Sprintf("Shutdown from %v", caller))
}

func (a *nuwaApplication) shutdownWithMessage(msg string) {
	fmt.Fprintf(os.Stdout, "Application exit: %v", msg)
	select {
	case <-a.exitChan:
	default:
		close(a.exitChan)
	}
}
