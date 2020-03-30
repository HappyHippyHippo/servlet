package servlet

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewContainer(t *testing.T) {
	t.Run("create a new container", func(t *testing.T) {
		action := "Creating a new container"

		if container := NewContainer(); container == nil {
			t.Errorf("%s didn't return a valid reference to a new container", action)
		}
	})
}

func Test_Container_Close(t *testing.T) {
	t.Run("should call close on closable instatiated entries", func(t *testing.T) {
		action := "Closing a container"

		var entry *MockClosable

		id := "__dummy_entry__"
		factory := func(c Container) interface{} { return entry }

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		entry = NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)

		container := NewContainer()
		container.Add(id, factory)
		container.Get(id)

		if err := container.Close(); err != nil {
			t.Errorf("%s returned the unexpected error : %v", action, err.Error())
		}
		if container.Has(id) {
			t.Errorf("%s didn't removed the stored (%s) entry", action, id)
		}
	})
}

func Test_Container_Has(t *testing.T) {
	t.Run("should correctly return existence checks requests", func(t *testing.T) {
		action := "Checking for the existence of a registed entry in the container"

		scenarios := []struct {
			factories []struct {
				id      string
				factory Factory
			}
			request  string
			expected bool
		}{
			// return false if the container is empty
			{
				factories: []struct {
					id      string
					factory Factory
				}{},
				request:  "__dummy_entry__",
				expected: false,
			},
			// return false if the search entry is not present in the container
			{
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "__dummy_entry_1__",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "__dummy_entry_2__",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "__dummy_entry_3__",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:  "__dummy_entry_4__",
				expected: false,
			},
			// return true if the search entry is present at the beginning of the container
			{
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "__dummy_entry_1__",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "__dummy_entry_2__",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "__dummy_entry_3__",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:  "__dummy_entry_1__",
				expected: true,
			},
			// return true if the search entry is present at the middle of the container
			{
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "__dummy_entry_1__",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "__dummy_entry_2__",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "__dummy_entry_3__",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:  "__dummy_entry_2__",
				expected: true,
			},
			// return true if the search entry is present at the end of the container
			{
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "__dummy_entry_1__",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "__dummy_entry_2__",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "__dummy_entry_3__",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:  "__dummy_entry_3__",
				expected: true,
			},
		}

		for _, scn := range scenarios {
			container := NewContainer()
			for _, f := range scn.factories {
				container.Add(f.id, f.factory)
			}

			if result := container.Has(scn.request); result != scn.expected {
				t.Errorf("%s returned (%v) when requesting (%s), expected (%v)", action, result, scn.request, scn.expected)
			}
		}
	})
}

func Test_Container_Add(t *testing.T) {
	t.Run("should correctly register entries", func(t *testing.T) {
		action := "Adding a entry into the container"

		scenarios := []struct {
			factories []struct {
				id      string
				factory Factory
			}
			request struct {
				id      string
				factory Factory
			}
			expected interface{}
		}{
			// return the inserted entry if the container is empty
			{
				factories: []struct {
					id      string
					factory Factory
				}{},
				request: struct {
					id      string
					factory Factory
				}{
					id:      "__dummy_entry_1__",
					factory: func(c Container) interface{} { return 1 },
				},
				expected: 1,
			},
			// return the inserted entry even if the container is not empty
			{
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "__dummy_entry_1__",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "__dummy_entry_2__",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "__dummy_entry_3__",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request: struct {
					id      string
					factory Factory
				}{
					id:      "__dummy_entry_4__",
					factory: func(c Container) interface{} { return 4 },
				},
				expected: 4,
			},
			// override a previously inserted entry
			{
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "__dummy_entry_1__",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "__dummy_entry_2__",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "__dummy_entry_3__",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request: struct {
					id      string
					factory Factory
				}{
					id:      "__dummy_entry_2__",
					factory: func(c Container) interface{} { return 4 },
				},
				expected: 4,
			},
		}

		for _, scn := range scenarios {
			container := NewContainer()
			for _, f := range scn.factories {
				if err := container.Add(f.id, f.factory); err != nil {
					t.Errorf("%s returned unexpected error while inserting (%s) : %v", action, f.id, err)
				}
			}

			if result := container.Add(
				scn.request.id,
				scn.request.factory); result != nil {
				t.Errorf("%s returned an unexpected error : %v", action, result)
				continue
			}

			if check := container.Get(scn.request.id); check != scn.expected {
				t.Errorf("%s returned (%v) when requesting for inserted (%s), expected (%v)", action, check, scn.request.id, scn.expected)
			}
		}
	})

	t.Run("should return an error on registering a nil factory", func(t *testing.T) {
		action := "Adding a entry into the container with a nil factory function"

		expected := "Invalid nil 'factory' argument"

		container := NewContainer()

		if result := container.Add("entry", nil); result == nil {
			t.Errorf("%s returned nil, when expected an error instance", action)
		} else {
			if result.Error() != expected {
				t.Errorf("%s returned (%v) error message, expected (%v)", action, result.Error(), expected)
			}
		}
	})
}

func Test_Container_Remove(t *testing.T) {
	t.Run("should correctly remove a entry", func(t *testing.T) {
		action := "Removing a entry from the container"

		var entry *MockClosable

		id := "__dummy_entry__"
		factory := func(c Container) interface{} {
			return entry
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		entry = NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(1)

		container := NewContainer()
		container.Add(id, factory)
		container.Get(id)
		container.Remove(id)

		if container.Has(id) {
			t.Errorf("%s didn't removed the requested entry (%s)", action, id)
		}
	})
}

func Test_Container_Get(t *testing.T) {
	t.Run("should correctly panic if does not exists", func(t *testing.T) {
		action := "Retrieving a registed entry not in the container"

		scenarios := []struct {
			factories []struct {
				id      string
				factory Factory
			}
			request       string
			expectedValue interface{}
			expectedError string
		}{
			// return nil if the container is empty
			{
				factories: []struct {
					id      string
					factory Factory
				}{},
				request:       "__dummy_entry__",
				expectedError: "Object '__dummy_entry__' not registed in the container",
			},
			// return nil if the search entry is not present in the container
			{
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "__dummy_entry_1__",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "__dummy_entry_2__",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "__dummy_entry_3__",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:       "__dummy_entry_4__",
				expectedError: "Object '__dummy_entry_4__' not registed in the container",
			},
		}

		for _, scn := range scenarios {
			container := NewContainer()
			for _, f := range scn.factories {
				container.Add(f.id, f.factory)
			}

			defer func() {
				if r := recover(); r == nil {
					t.Errorf("%s when requesting (%s) did not panic", action, scn.request)
				} else {
					if check := r.(error).Error(); check != scn.expectedError {
						t.Errorf("%s when requesting (%s) did not panic with the expected (%s) value : %s ", action, scn.request, scn.expectedError, check)
					}
				}
			}()

			container.Get(scn.request)
		}
	})

	t.Run("should correctly return entries requests", func(t *testing.T) {
		action := "Retrieving a registed entry in the container"

		scenarios := []struct {
			factories []struct {
				id      string
				factory Factory
			}
			request       string
			expectedValue interface{}
			expectedError string
		}{
			// return the first entry value if the search entry is present at the beginning of the container
			{
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "__dummy_entry_1__",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "__dummy_entry_2__",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "__dummy_entry_3__",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:       "__dummy_entry_1__",
				expectedValue: 1,
				expectedError: "",
			},
			// return middle entry value if the search entry is present at the middle of the container
			{
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "__dummy_entry_1__",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "__dummy_entry_2__",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "__dummy_entry_3__",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:       "__dummy_entry_2__",
				expectedValue: 2,
				expectedError: "",
			},
			// return last entry value if the search entry is present at the end of the container
			{
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "__dummy_entry_1__",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "__dummy_entry_2__",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "__dummy_entry_3__",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:       "__dummy_entry_3__",
				expectedValue: 3,
				expectedError: "",
			},
		}

		for _, scn := range scenarios {
			container := NewContainer()
			for _, f := range scn.factories {
				container.Add(f.id, f.factory)
			}

			if result := container.Get(scn.request); result != scn.expectedValue {
				t.Errorf("%s returned (%v) when requesting (%s), expected (%v)", action, result, scn.request, scn.expectedValue)
			}
		}
	})

	t.Run("should only call the factory of a entry once", func(t *testing.T) {
		action := "Retrieving a registed entry from the container multiple times"

		call := 0
		id := "__dummy_entry__"
		value := 1
		factory := func(c Container) interface{} {
			call = call + 1
			return value
		}

		container := NewContainer()
		container.Add(id, factory)

		for range []int{1, 2, 3} {
			if result := container.Get(id); result != value {
				t.Errorf("%s returned (%v) when requesting for (%s), expected (%v)", action, result, id, value)
			}

			if call != 1 {
				t.Errorf("%s factory was called (%v) times, expected once", action, call)
			}
		}
	})
}

func Test_Container(t *testing.T) {
	t.Run("should not instantiate a entry if not requested", func(t *testing.T) {
		var entry *MockClosable

		id := "__dummy_entry__"
		factory := func(c Container) interface{} {
			return entry
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		entry = NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(0)

		container := NewContainer()
		container.Add(id, factory)
		container.Close()
	})
}
