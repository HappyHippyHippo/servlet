package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_NewAppContainer(t *testing.T) {
	t.Run("correctly instantiate a new app container", func(t *testing.T) {
		c := NewAppContainer()
		defer c.Close()

		if c == nil {
			t.Error("didn't returned a valid reference")
		}
	})

	t.Run("new app container instantiate a the factory map", func(t *testing.T) {
		c := NewAppContainer()
		defer c.Close()

		if c.factories == nil {
			t.Error("didn't created the factories map")
		}
	})

	t.Run("new app container instantiate a the loaded entries map", func(t *testing.T) {
		c := NewAppContainer()
		defer c.Close()

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

	t.Run("remove all entries", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := "id"

		c := NewAppContainer()

		entry := NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)
		c.factories[id] = func(_ *AppContainer) (interface{}, error) {
			return entry, nil
		}

		_, _ = c.Get(id)
		c.Close()

		if _, ok := c.factories[id]; ok {
			t.Error("didn't removed the entry")
		}
	})
}

func Test_AppContainer_Has(t *testing.T) {
	t.Run("validate the entry existence", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := "id"

		c := NewAppContainer()
		defer c.Close()

		entry := NewMockClosable(ctrl)
		c.factories[id] = func(_ *AppContainer) (interface{}, error) {
			return entry, nil
		}

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

	t.Run(" nil factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := "id"

		c := NewAppContainer()
		defer c.Close()

		if err := c.Add(id, nil); err == nil {
			t.Error("didn't returned the expected error")
		}
	})

	t.Run("adding a entry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := "id"

		c := NewAppContainer()
		defer c.Close()

		entry := NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)

		if err := c.Add(id, func(_ *AppContainer) (interface{}, error) {
			return entry, nil
		}); err != nil {
			t.Errorf("returned the (%s) error", err)
		} else if _, ok := c.factories[id]; !ok {
			t.Error("didn't found the added factory")
		} else if e, _ := c.Get(id); !reflect.DeepEqual(e, entry) {
			t.Error("didn't stored the requested entry factory")
		}
	})

	t.Run("overriding a loaded entry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := "id"

		c := NewAppContainer()
		defer c.Close()

		entry1 := NewMockClosable(ctrl)
		entry1.EXPECT().Close().Times(1)
		entry2 := NewMockClosable(ctrl)
		entry2.EXPECT().Close().Times(1)

		if err := c.Add(id, func(_ *AppContainer) (interface{}, error) {
			return entry1, nil
		}); err != nil {
			t.Errorf("returned the (%s) error", err)
		} else if _, ok := c.factories[id]; !ok {
			t.Error("didn't found the added factory")
		} else if e, _ := c.Get(id); !reflect.DeepEqual(e, entry1) {
			t.Error("didn't stored the requested first entry factory")
		} else if err := c.Add(id, func(_ *AppContainer) (interface{}, error) {
			return entry2, nil
		}); err != nil {
			t.Errorf("returned the (%s) error", err)
		} else if _, ok := c.factories[id]; !ok {
			t.Error("didn't found the added factory")
		} else if e, _ := c.Get(id); !reflect.DeepEqual(e, entry2) {
			t.Error("didn't stored the requested second entry factory")
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

	t.Run("removing a non-loaded entry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := "id"

		c := NewAppContainer()
		defer c.Close()

		_ = c.Add(id, func(_ *AppContainer) (interface{}, error) {
			return "value", nil
		})
		c.Remove(id)

		if _, ok := c.factories[id]; ok {
			t.Error("didn't removed the factory")
		} else if _, ok := c.entries[id]; ok {
			t.Error("didn't removed the loaded entry")
		}
	})

	t.Run("removing a loaded entry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := "id"

		c := NewAppContainer()
		defer c.Close()

		entry := NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)
		_ = c.Add(id, func(_ *AppContainer) (interface{}, error) {
			return entry, nil
		})

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

	t.Run("retrieving a non-registered entry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := NewAppContainer()
		defer c.Close()

		if e, err := c.Get("invalid_id"); e != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned a error")
		} else if err.Error() != "entry 'invalid_id' not registered in the container" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error while calling the factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := "id"

		c := NewAppContainer()
		defer c.Close()

		expectedError := "error message"
		_ = c.Add(id, func(_ *AppContainer) (interface{}, error) {
			return nil, fmt.Errorf(expectedError)
		})

		if e, err := c.Get(id); e != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("retrieving a non-loaded entry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := "id"

		c := NewAppContainer()
		defer c.Close()

		count := 0
		entry := NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)
		_ = c.Add(id, func(_ *AppContainer) (interface{}, error) {
			count = count + 1
			return entry, nil
		})

		if e, err := c.Get(id); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if e == nil {
			t.Error("didn't returned a valid reference")
		}
	})

	t.Run("retrieving a loaded entry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := "id"

		c := NewAppContainer()
		defer c.Close()

		count := 0
		entry := NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)
		_ = c.Add(id, func(_ *AppContainer) (interface{}, error) {
			count = count + 1
			return entry, nil
		})

		if e, err := c.Get(id); err != nil {
			t.Errorf("returned the (%v) error on the first get call", err)
		} else if e == nil {
			t.Error("didn't returned a valid reference")
		} else if count != 1 {
			t.Error("called the factory more than once")
		} else if e, err := c.Get(id); err != nil {
			t.Errorf("returned the (%v) error on the second get call", err)
		} else if e == nil {
			t.Error("didn't returned a valid reference")
		} else if count != 1 {
			t.Error("called the factory more than once")
		}
	})
}
