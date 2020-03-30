package servlet

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func Test_NewApplication(t *testing.T) {
	t.Run("create a application with the default parameters", func(t *testing.T) {
		action := "Creating a new application"

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(nil)

		if app == nil {
			t.Errorf("%s didn't return a valid reference to a new application", action)
		}
	})

	t.Run("should register a gin engine into the internal container with the default id", func(t *testing.T) {
		action := "Creating a new application without parameters"

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(nil)

		container := app.GetContainer()
		engine := container.Get(ContainerEngineID)

		switch engine.(type) {
		case (Engine):
		default:
			t.Errorf("%s didn't register a gin engine in the container with the %s id", action, ContainerEngineID)
		}
	})

	t.Run("should register a gin engine into the internal container with a defined id", func(t *testing.T) {
		action := "Creating a new application without parameters"

		id := "__dummy_id__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		parameters := NewMockApplicationParameters(ctrl)
		parameters.EXPECT().GetEngineID().Return(id).Times(1)

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(parameters)

		container := app.GetContainer()
		engine := container.Get(id)

		switch engine.(type) {
		case (Engine):
		default:
			t.Errorf("%s didn't register a gin engine in the container with the %s id", action, id)
		}
	})
}

func Test_Application_Boot(t *testing.T) {
	t.Run("should call the boot method on added providers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(nil)

		provider1 := NewMockProvider(ctrl)
		provider1.EXPECT().Register(app.GetContainer()).Times(1)
		provider1.EXPECT().Boot(app.GetContainer()).Times(1)

		provider2 := NewMockProvider(ctrl)
		provider2.EXPECT().Register(app.GetContainer()).Times(1)
		provider2.EXPECT().Boot(app.GetContainer()).Times(1)

		app.AddProvider(provider1)
		app.AddProvider(provider2)

		app.Boot()
	})

	t.Run("should call the boot methods only once", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(nil)

		provider1 := NewMockProvider(ctrl)
		provider1.EXPECT().Register(app.GetContainer()).Times(1)
		provider1.EXPECT().Boot(app.GetContainer()).Times(1)

		provider2 := NewMockProvider(ctrl)
		provider2.EXPECT().Register(app.GetContainer()).Times(1)
		provider2.EXPECT().Boot(app.GetContainer()).Times(1)

		app.AddProvider(provider1)
		app.AddProvider(provider2)

		app.Boot()
		app.Boot()
		app.Boot()
	})
}

func Test_Application_Run(t *testing.T) {
	t.Run("should redirect to the application engine", func(t *testing.T) {
		args := []string{"__dummy_arg_1__", "__dummy_arg_2__"}
		expected := []interface{}{"__dummy_arg_1__", "__dummy_arg_2__"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		engine := NewMockEngine(ctrl)
		engine.EXPECT().Run(expected...).Return(nil).Times(1)

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(nil).(*application)
		app.engine = engine

		app.Run(args...)
	})

	t.Run("should boot if not yet started", func(t *testing.T) {
		args := []string{"__dummy_arg_1__", "__dummy_arg_2__"}
		expected := []interface{}{"__dummy_arg_1__", "__dummy_arg_2__"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		engine := NewMockEngine(ctrl)
		engine.EXPECT().Run(expected...).Return(nil).Times(1)

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(nil).(*application)
		app.engine = engine

		provider := NewMockProvider(ctrl)
		provider.EXPECT().Register(app.GetContainer()).Times(1)
		provider.EXPECT().Boot(app.GetContainer()).Times(1)

		app.AddProvider(provider)

		app.Run(args...)
	})
}

func Test_Application_Engine(t *testing.T) {
	t.Run("should return a valid gin engine instance", func(t *testing.T) {
		action := "Retrieving the gin engine container"

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(nil)

		if engine := app.Engine(); engine == nil {
			t.Errorf("%s didn't return a valid reference to a gin engine", action)
		}
	})
}

func Test_Application_GetContainer(t *testing.T) {
	t.Run("should return a valid container instance", func(t *testing.T) {
		action := "Retrieving the application container"

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(nil)

		if container := app.GetContainer(); container == nil {
			t.Errorf("%s didn't return a valid reference to a container", action)
		}
	})
}

func Test_Application_SetContainer(t *testing.T) {
	t.Run("should return a error when trying to register a nil container", func(t *testing.T) {
		action := "Registering a nil container"

		expectedError := "Invalid nil 'container' argument"

		app := NewApplication(nil)

		check := app.SetContainer(nil)
		if check == nil {
			t.Errorf("%s didn't return the expected error", action)
		} else {
			if check.Error() != expectedError {
				t.Errorf("%s returned the (%s) error, expected (%s)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should correctly register a new container", func(t *testing.T) {
		action := "Registering a new container"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)

		app := NewApplication(nil)

		if err := app.SetContainer(container); err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if container != app.GetContainer() {
			t.Errorf("%s didn't stored the assigned container", action)
		}
	})
}

func Test_Application_AddProvider(t *testing.T) {
	t.Run("should return an errors if a nil reference is given", func(t *testing.T) {
		action := "Adding a invalid nil reference to a application provider"

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(nil)

		if err := app.AddProvider(nil); err == nil {
			t.Errorf("%s didn't return the expected error", action)
		}
	})

	t.Run("should properly register the given providers", func(t *testing.T) {
		action := "Adding a valid application provider"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.ReleaseMode)
		app := NewApplication(nil)

		provider := NewMockProvider(ctrl)
		provider.EXPECT().Register(app.GetContainer()).Times(1)

		if err := app.AddProvider(provider); err != nil {
			t.Errorf("%s return the unexpected error (%v)", action, err)
		}
	})
}
