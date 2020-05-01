package servlet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"os"
	"testing"
)

/// ---------------------------------------------------------------------------
/// GinAppParams
/// ---------------------------------------------------------------------------

func Test_NewGinAppParams(t *testing.T) {
	t.Run("no env override", func(t *testing.T) {
		p := NewGinAppParams()
		if p.EngineID != ContainerGinEngineID {
			t.Errorf("stored the '%s' gin engine container id", p.EngineID)
		}
	})

	t.Run("with env override", func(t *testing.T) {
		id := "test_id"

		_ = os.Setenv(EnvContainerGinEngineID, id)
		defer func() { _ = os.Setenv(EnvContainerGinEngineID, "") }()
		p := NewGinAppParams()

		if p.EngineID != id {
			t.Errorf("stored the '%s' gin engine container id", p.EngineID)
		}
	})
}

/// ---------------------------------------------------------------------------
/// GinApp
/// ---------------------------------------------------------------------------

func Test_NewGinApp(t *testing.T) {
	t.Run("without params", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		if a := NewGinApp(nil); a == nil {
			t.Error("didn't returned a valid reference")
		} else if a.container == nil {
			t.Error("didn't created the app container")
		} else if a.providers == nil {
			t.Error("didn't created the list of providers")
		} else if a.boot {
			t.Error("didn't flagged the app as not booted")
		} else if a.engine == nil {
			t.Error("didn't instantiate the gin engine")
		} else if _, ok := a.container.factories[ContainerGinEngineID]; !ok {
			t.Error("didn't registered the gin engine in the default id")
		} else if e, _ := a.container.Get(ContainerGinEngineID); e != a.engine {
			t.Error("didn't registered in the container the correct engine retrieving factory")
		}
	})

	t.Run("with defined params", func(t *testing.T) {
		id := "test_id"
		p := NewGinAppParams()
		p.EngineID = id

		gin.SetMode(gin.ReleaseMode)
		if a := NewGinApp(p); a == nil {
			t.Error("didn't returned a valid reference")
		} else if a.container == nil {
			t.Error("didn't created the app container")
		} else if a.providers == nil {
			t.Error("didn't created the list of providers")
		} else if a.boot {
			t.Error("didn't flagged the app as not booted")
		} else if a.engine == nil {
			t.Error("didn't instantiate the gin engine")
		} else if _, ok := a.container.factories[id]; !ok {
			t.Error("didn't registered the gin engine in the defined id")
		} else if e, _ := a.container.Get(id); e != a.engine {
			t.Error("didn't registered in the container the correct engine retrieving factory")
		}
	})
}

func Test_GinApp_Engine(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	a := NewGinApp(nil)

	t.Run("retrieve the gin engine", func(t *testing.T) {
		if a.Engine() != a.engine {
			t.Error("didn't returned the gin engine")
		}
	})
}

func Test_GinApp_Run(t *testing.T) {
	t.Run("nil pointer receiver", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("didn't panic")
			} else {
				switch e := r.(type) {
				case error:
					if e.Error() != "nil pointer receiver" {
						t.Errorf("panic with the (%v) error", e)
					}
				default:
					t.Error("didn't panic with an error")
				}
			}
		}()

		var a *GinApp
		_ = a.Run()
	})

	t.Run("error on boot", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.ReleaseMode)
		app := NewGinApp(nil)

		provider := NewMockAppProvider(ctrl)
		provider.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		provider.EXPECT().Boot(gomock.Any()).Return(fmt.Errorf("error")).Times(1)
		_ = app.Add(provider)

		if err := app.Run(); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("run the application", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		port := "port_value"
		expectedError := fmt.Errorf("error")

		gin.SetMode(gin.ReleaseMode)
		app := NewGinApp(nil)

		e := NewMockGinEngine(ctrl)
		e.EXPECT().Run(port).Return(expectedError).Times(1)
		app.engine = e

		if err := app.Run(port); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError.Error() {
			t.Errorf("returned the (%v) error", err)
		} else if !app.boot {
			t.Error("didn't booted the application")
		}
	})
}
