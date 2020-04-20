package servlet

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewContainer(t *testing.T) {
	t.Run("create a new container", func(t *testing.T) {
		if container := NewContainer(); container == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_Container_Close(t *testing.T) {
	id := "id"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	container := NewContainer()

	entry := NewMockClosable(ctrl)
	entry.EXPECT().Close().Times(1)
	container.Add(id, func(c Container) interface{} { return entry })
	container.Get(id)

	t.Run("close on closable instatiated entries", func(t *testing.T) {
		if err := container.Close(); err != nil {
			t.Errorf("returned the (%v) error", err.Error())
		} else if container.Has(id) {
			t.Errorf("didn't removed the entry")
		}
	})
}

func Test_Container_Has(t *testing.T) {
	t.Run("return existence checks", func(t *testing.T) {
		scenarios := []struct {
			factories []struct {
				id      string
				factory Factory
			}
			request  string
			expected bool
		}{
			{ // return false if the container is empty
				factories: []struct {
					id      string
					factory Factory
				}{},
				request:  "id",
				expected: false,
			},
			{ // return false if the search entry is not present in the container
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "id1",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "id2",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "id3",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:  "id4",
				expected: false,
			},
			{ // return true if the search entry is present at the beginning of the container
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "id1",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "id2",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "id3",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:  "id1",
				expected: true,
			},
			{ // return true if the search entry is present at the middle of the container
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "id1",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "id2",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "id3",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:  "id2",
				expected: true,
			},
			{ // return true if the search entry is present at the end of the container
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "id1",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "id2",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "id3",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:  "id3",
				expected: true,
			},
		}

		for _, scn := range scenarios {
			container := NewContainer()

			for _, f := range scn.factories {
				container.Add(f.id, f.factory)
			}

			if result := container.Has(scn.request); result != scn.expected {
				t.Errorf("returned (%v) when requesting (%s), expected (%v)", result, scn.request, scn.expected)
			}
		}
	})
}

func Test_Container_Add(t *testing.T) {
	t.Run("error on register a nil factory", func(t *testing.T) {
		container := NewContainer()

		if err := container.Add("entry", nil); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "Invalid nil 'factory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register factories", func(t *testing.T) {
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
			{ // return the inserted entry if the container is empty
				factories: []struct {
					id      string
					factory Factory
				}{},
				request: struct {
					id      string
					factory Factory
				}{
					id:      "id1",
					factory: func(c Container) interface{} { return 1 },
				},
				expected: 1,
			},
			{ // return the inserted entry even if the container is not empty
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "id1",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "id2",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "id3",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request: struct {
					id      string
					factory Factory
				}{
					id:      "id4",
					factory: func(c Container) interface{} { return 4 },
				},
				expected: 4,
			},
			{ // override a previously inserted entry
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "id1",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "id2",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "id3",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request: struct {
					id      string
					factory Factory
				}{
					id:      "id2",
					factory: func(c Container) interface{} { return 4 },
				},
				expected: 4,
			},
		}

		for _, scn := range scenarios {
			container := NewContainer()

			for _, f := range scn.factories {
				if err := container.Add(f.id, f.factory); err != nil {
					t.Errorf("returned the (%v) error when inserting (%s)", err, f.id)
				}
			}

			if err := container.Add(scn.request.id, scn.request.factory); err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if value := container.Get(scn.request.id); value != scn.expected {
				t.Errorf("returned (%v) when requesting for inserted (%s), expected (%v)", value, scn.request.id, scn.expected)
			}
		}
	})
}

func Test_Container_Remove(t *testing.T) {
	id := "id"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	container := NewContainer()

	entry := NewMockClosable(ctrl)
	entry.EXPECT().Close().Times(1)
	container.Add(id, func(c Container) interface{} { return entry })
	container.Get(id)

	t.Run("remove a entry", func(t *testing.T) {
		container.Remove(id)
		if container.Has(id) {
			t.Errorf("didn't removed the requested entry")
		}
	})
}

func Test_Container_Get(t *testing.T) {
	t.Run("panic if doesn't exists", func(t *testing.T) {
		scenarios := []struct {
			factories []struct {
				id      string
				factory Factory
			}
			request       string
			expectedValue interface{}
			expectedError string
		}{
			{ // return nil if the container is empty
				factories: []struct {
					id      string
					factory Factory
				}{},
				request:       "id",
				expectedError: "Object 'id' not registed in the container",
			},
			{ // return nil if the search entry is not present in the container
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "id1",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "id2",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "id3",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:       "id4",
				expectedError: "Object 'id4' not registed in the container",
			},
		}

		for _, scn := range scenarios {
			container := NewContainer()

			for _, f := range scn.factories {
				container.Add(f.id, f.factory)
			}

			defer func() {
				if r := recover(); r == nil {
					t.Errorf("didn't panic when requesting (%s)", scn.request)
				} else if err := r.(error); err.Error() != scn.expectedError {
					t.Errorf("panic with the (%v) error when requesting (%s)", err, scn.request)
				}
			}()

			container.Get(scn.request)
		}
	})

	t.Run("return the intantiated entry", func(t *testing.T) {
		scenarios := []struct {
			factories []struct {
				id      string
				factory Factory
			}
			request       string
			expectedValue interface{}
			expectedError string
		}{
			{ // return the first entry value if the search entry is present at the beginning of the container
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "id1",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "id2",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "id3",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:       "id1",
				expectedValue: 1,
				expectedError: "",
			},
			{ // return middle entry value if the search entry is present at the middle of the container
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "id1",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "id2",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "id3",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:       "id2",
				expectedValue: 2,
				expectedError: "",
			},
			{ // return last entry value if the search entry is present at the end of the container
				factories: []struct {
					id      string
					factory Factory
				}{
					{
						id:      "id1",
						factory: func(c Container) interface{} { return 1 },
					},
					{
						id:      "id2",
						factory: func(c Container) interface{} { return 2 },
					},
					{
						id:      "id3",
						factory: func(c Container) interface{} { return 3 },
					},
				},
				request:       "id3",
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
				t.Errorf("returned (%v) when requesting (%s), expected (%v)", result, scn.request, scn.expectedValue)
			}
		}
	})

	t.Run("call factory once", func(t *testing.T) {
		call := 0
		id := "id"
		value := 1

		container := NewContainer()
		container.Add(id, func(c Container) interface{} { call = call + 1; return value })

		for range []int{1, 2, 3} {
			if result := container.Get(id); result != value {
				t.Errorf("returned (%v) when requesting for (%s), expected (%v)", result, id, value)
			} else if call != 1 {
				t.Errorf("factory was called (%v) times", call)
			}
		}
	})
}

func Test_Container(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	container := NewContainer()

	t.Run("don't instantiate a entry if not requested", func(t *testing.T) {
		entry := NewMockClosable(ctrl)
		entry.EXPECT().Close().Times(0)
		container.Add("id", func(c Container) interface{} { return entry })

		container.Close()
	})
}
