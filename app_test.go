package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

/// ---------------------------------------------------------------------------
/// AppContainer
/// ---------------------------------------------------------------------------

func Test_NewAppContainer(t *testing.T) {
	c := NewAppContainer()
	defer c.Close()

	t.Run("instantiate a new app container", func(t *testing.T) {
		if c == nil {
			t.Error("didn't returned a valid reference")
		}
	})

	t.Run("instantiate a the factory map", func(t *testing.T) {
		if c.factories == nil {
			t.Error("didn't created the factories map")
		}
	})

	t.Run("instantiate a the loaded entries map", func(t *testing.T) {
		if c.entries == nil {
			t.Error("didn't created the loaded entries map")
		}
	})
}

func Test_AppContainer_Close(t *testing.T) {
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

		var c *AppContainer
		c.Close()
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "id"

	c := NewAppContainer()

	entry := NewMockClosable(ctrl)
	entry.EXPECT().Close().Times(1)
	c.factories[id] = func(_ *AppContainer) (interface{}, error) { return entry, nil }
	_, _ = c.Get(id)
	c.Close()

	t.Run("remove all entries", func(t *testing.T) {
		if _, ok := c.factories[id]; ok {
			t.Error("didn't removed the entry")
		}
	})
}

func Test_AppContainer_Has(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "id"

	c := NewAppContainer()
	defer c.Close()

	entry := NewMockClosable(ctrl)
	c.factories[id] = func(_ *AppContainer) (interface{}, error) { return entry, nil }

	t.Run("validate the entry existence", func(t *testing.T) {
		if !c.Has(id) {
			t.Error("didn't found the entry")
		}
	})
}

func Test_AppContainer_Add(t *testing.T) {
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

		var c *AppContainer
		_ = c.Add("id", nil)
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "id"

	c := NewAppContainer()
	defer c.Close()

	t.Run(" nil factory", func(t *testing.T) {
		if err := c.Add(id, nil); err == nil {
			t.Error("didn't returned the expected error")
		}
	})

	t.Run("adding a entry", func(t *testing.T) {
		entry := NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)

		if err := c.Add(id, func(_ *AppContainer) (interface{}, error) { return entry, nil }); err != nil {
			t.Errorf("returned the (%s) error", err)
		} else if _, ok := c.factories[id]; !ok {
			t.Error("didn't found the added factory")
		} else if e, _ := c.Get(id); !reflect.DeepEqual(e, entry) {
			t.Error("didn't stored the requested entry factory")
		}
	})

	t.Run("overriding a loaded entry", func(t *testing.T) {
		entry := NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)

		if err := c.Add(id, func(_ *AppContainer) (interface{}, error) { return entry, nil }); err != nil {
			t.Errorf("returned the (%s) error", err)
		} else if _, ok := c.factories[id]; !ok {
			t.Error("didn't found the added factory")
		} else if e, _ := c.Get(id); !reflect.DeepEqual(e, entry) {
			t.Error("didn't stored the requested entry factory")
		}
	})
}

func Test_AppContainer_Remove(t *testing.T) {
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

		var c *AppContainer
		c.Remove("id")
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "id"

	c := NewAppContainer()
	defer c.Close()

	t.Run("removing a non-loaded entry", func(t *testing.T) {
		_ = c.Add(id, func(_ *AppContainer) (interface{}, error) { return "value", nil })
		c.Remove(id)

		if _, ok := c.factories[id]; ok {
			t.Error("didn't removed the factory")
		} else if _, ok := c.entries[id]; ok {
			t.Error("didn't removed the loaded entry")
		}
	})

	t.Run("removing a loaded entry", func(t *testing.T) {
		entry := NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)
		_ = c.Add(id, func(_ *AppContainer) (interface{}, error) { return entry, nil })
		_, _ = c.Get(id)
		c.Remove(id)

		if _, ok := c.factories[id]; ok {
			t.Error("didn't removed the factory")
		} else if _, ok := c.entries[id]; ok {
			t.Error("didn't removed the loaded entry")
		}
	})
}

func Test_AppContainer_Get(t *testing.T) {
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

		var c *AppContainer
		_, _ = c.Get("id")
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "id"

	c := NewAppContainer()
	defer c.Close()

	t.Run("retrieving a non-registered entry", func(t *testing.T) {
		if e, err := c.Get("invalid_id"); e != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned a error")
		} else if err.Error() != "entry 'invalid_id' not registered in the container" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error while calling the factory", func(t *testing.T) {
		expectedError := "error message"
		_ = c.Add(id, func(_ *AppContainer) (interface{}, error) { return nil, fmt.Errorf(expectedError) })

		if e, err := c.Get(id); e != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	count := 0
	t.Run("retrieving a non-loaded entry", func(t *testing.T) {
		entry := NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)
		_ = c.Add(id, func(_ *AppContainer) (interface{}, error) { count = count + 1; return entry, nil })

		if e, err := c.Get(id); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if e == nil {
			t.Error("didn't returned a valid reference")
		}
	})

	t.Run("retrieving a loaded entry", func(t *testing.T) {
		if e, err := c.Get(id); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if e == nil {
			t.Error("didn't returned a valid reference")
		} else if count != 1 {
			t.Error("called the factory more than once")
		}
	})
}

/// ---------------------------------------------------------------------------
/// App
/// ---------------------------------------------------------------------------

func Test_NewApp(t *testing.T) {
	a := NewApp()

	t.Run("instantiate a new app", func(t *testing.T) {
		if a == nil {
			t.Error("didn't returned a valid reference")
		}
	})

	t.Run("instantiate a app container", func(t *testing.T) {
		if a.container == nil {
			t.Error("didn't created the app container")
		}
	})

	t.Run("instantiate a list of providers", func(t *testing.T) {
		if a.providers == nil {
			t.Error("didn't created the list of providers")
		}
	})

	t.Run("flag the app has not booted", func(t *testing.T) {
		if a.boot {
			t.Error("didn't flagged the app as not booted")
		}
	})
}

func Test_App_Container(t *testing.T) {
	a := NewApp()

	t.Run("retrieve the stored container", func(t *testing.T) {
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := NewApp()

	t.Run("nil provider", func(t *testing.T) {
		if err := a.Add(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'provider' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error registering provider", func(t *testing.T) {
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
