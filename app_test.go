package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_NewApp(t *testing.T) {
	t.Run("instantiate a new app", func(t *testing.T) {
		if NewApp() == nil {
			t.Error("didn't returned a valid reference")
		}
	})

	t.Run("instantiate a app container", func(t *testing.T) {
		if NewApp().container == nil {
			t.Error("didn't created the app container")
		}
	})

	t.Run("instantiate a list of providers", func(t *testing.T) {
		if NewApp().providers == nil {
			t.Error("didn't created the list of providers")
		}
	})

	t.Run("flag the app has not booted", func(t *testing.T) {
		if NewApp().boot {
			t.Error("didn't flagged the app as not booted")
		}
	})
}

func Test_App_Container(t *testing.T) {
	t.Run("retrieve the stored container", func(t *testing.T) {
		a := NewApp()
		if a.Container() != a.container {
			t.Error("didn't returned the stored container")
		}
	})
}

func Test_App_Add(t *testing.T) {
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

		var a *App
		_ = a.Add(nil)
	})

	t.Run("nil provider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		a := NewApp()

		if err := a.Add(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'provider' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error registering provider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		a := NewApp()

		expectedError := "error"

		p := NewMockAppProvider(ctrl)
		p.EXPECT().Register(a.container).Return(fmt.Errorf(expectedError)).Times(1)

		if err := a.Add(p); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		} else if len(a.providers) != 0 {
			t.Error("stored the failing provider")
		}
	})

	t.Run("adding a valid provider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		a := NewApp()

		p := NewMockAppProvider(ctrl)
		p.EXPECT().Register(a.container).Return(nil).Times(1)

		if err := a.Add(p); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if a.providers[0] != p {
			t.Error("didn't stored the added provider")
		}
	})
}

func Test_App_Boot(t *testing.T) {
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

		var a *App
		_ = a.Boot()
	})

	t.Run("error on boot", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		a := NewApp()
		p := NewMockAppProvider(ctrl)
		p.EXPECT().Register(a.container).Return(nil).Times(1)
		p.EXPECT().Boot(a.container).Return(fmt.Errorf(expectedError)).Times(1)
		_ = a.Add(p)

		if err := a.Boot(); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("boot all providers only once", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		a := NewApp()
		p := NewMockAppProvider(ctrl)
		p.EXPECT().Register(a.container).Times(1)
		p.EXPECT().Boot(a.container).Times(1)
		_ = a.Add(p)

		_ = a.Boot()
		_ = a.Boot()

		if !a.boot {
			t.Error("didn't flagged the app as booted")
		}
	})
}
