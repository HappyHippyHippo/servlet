package servlet

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func Test_NewApplication(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	parameters := NewDefaultParameters()
	app := NewApplication(parameters)

	t.Run("create a application", func(t *testing.T) {
		if app == nil {
			t.Errorf("didn't return a valid reference")
		}
	})

	t.Run("register a gin engine", func(t *testing.T) {
		switch app.GetContainer().Get(parameters.EngineID).(type) {
		case (Engine):
		default:
			t.Errorf("didn't register a gin engine with the (%s) id", parameters.EngineID)
		}
	})
}

func Test_Application_Boot(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := NewApplication(NewDefaultParameters())

	provider1 := NewMockProvider(ctrl)
	provider1.EXPECT().Register(app.GetContainer()).Times(1)
	provider1.EXPECT().Boot(app.GetContainer()).Times(1)
	app.AddProvider(provider1)

	provider2 := NewMockProvider(ctrl)
	provider2.EXPECT().Register(app.GetContainer()).Times(1)
	provider2.EXPECT().Boot(app.GetContainer()).Times(1)
	app.AddProvider(provider2)

	t.Run("call the boot method on added providers", func(t *testing.T) {
		app.Boot()
	})

	t.Run("call providers boot method only once", func(t *testing.T) {
		app.Boot()
		app.Boot()
		app.Boot()
	})
}

func Test_Application_Engine(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	app := NewApplication(NewDefaultParameters())

	t.Run("return a valid gin engine", func(t *testing.T) {
		if engine := app.Engine(); engine == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_Application_GetContainer(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	app := NewApplication(NewDefaultParameters())

	t.Run("return a valid container", func(t *testing.T) {
		if container := app.GetContainer(); container == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_Application_SetContainer(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	app := NewApplication(NewDefaultParameters())

	t.Run("error when register a nil container", func(t *testing.T) {
		if err := app.SetContainer(nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'container' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register a new container", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)

		if err := app.SetContainer(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if container != app.GetContainer() {
			t.Errorf("didn't stored the assigned container")
		}
	})
}

func Test_Application_AddProvider(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	app := NewApplication(NewDefaultParameters())

	t.Run("error if a nil reference is given", func(t *testing.T) {
		if err := app.AddProvider(nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'provider' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the provider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewMockProvider(ctrl)
		provider.EXPECT().Register(app.GetContainer()).Times(1)

		if err := app.AddProvider(provider); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_Application_Run(t *testing.T) {
	args := []string{"arg1", "arg2"}
	expected := []interface{}{"arg1", "arg2"}

	t.Run("redirect to gin engine", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(NewDefaultParameters()).(*application)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		engine := NewMockEngine(ctrl)
		engine.EXPECT().Run(expected...).Return(nil).Times(1)
		app.engine = engine

		app.Run(args...)
	})

	t.Run("boot if not yet booted", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(NewDefaultParameters()).(*application)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		engine := NewMockEngine(ctrl)
		engine.EXPECT().Run(expected...).Return(nil).Times(1)
		app.engine = engine

		provider := NewMockProvider(ctrl)
		provider.EXPECT().Register(app.GetContainer()).Times(1)
		provider.EXPECT().Boot(app.GetContainer()).Times(1)
		app.AddProvider(provider)

		app.Run(args...)
	})
}
