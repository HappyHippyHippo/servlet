package config

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_NewConfig(t *testing.T) {
	t.Run("creates a new config", func(t *testing.T) {
		action := "Creating a new config object"

		config, err := NewConfig(60 * time.Second)

		if config == nil {
			t.Errorf("%s didn't return a valid reference to a new config", action)
		} else {
			config.Close()
		}

		if err != nil {
			t.Errorf("%s returned a unexpected error : %v", action, err)
		}
	})
}

func Test_Config_Close(t *testing.T) {
	t.Run("should propagate to registered sources", func(t *testing.T) {
		id1 := "source.1"
		priority1 := 0
		partial1 := partial{}

		id2 := "source.2"
		priority2 := 1
		partial2 := partial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Return(nil).Times(1)
		source1.EXPECT().Get("").Return(partial1).AnyTimes()

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Return(nil).Times(1)
		source2.EXPECT().Get("").Return(partial2).AnyTimes()

		config.AddSource(id1, priority1, source1)
		config.AddSource(id2, priority2, source2)
	})
}

func Test_Config_Has(t *testing.T) {
	t.Run("should correctly return the existence of the path", func(t *testing.T) {
		action := "Checking the existence of a path"

		scenarios := []struct {
			config   partial
			search   string
			expected bool
		}{
			{ // test the existence of a present path
				config:   partial{"node": "value"},
				search:   "node",
				expected: true,
			},
			{ // test the non-existence of a missing path
				config:   partial{"node": "value"},
				search:   "inexistent-node",
				expected: false,
			},
		}

		id := "source"
		priority := 0

		for _, scn := range scenarios {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			config, _ := NewConfig(60 * time.Second)
			defer config.Close()

			source := NewMockSource(ctrl)
			source.EXPECT().Close().Return(nil).Times(1)
			source.EXPECT().Get("").Return(scn.config).Times(1)

			config.AddSource(id, priority, source)

			if result := config.Has(scn.search); result != scn.expected {
				t.Errorf("%s didn't validated (%v) returning (%v), expected (%v)", action, scn.search, result, scn.expected)
			}
		}
	})
}

func Test_Config_Get(t *testing.T) {
	t.Run("should correctly return value associated to the path", func(t *testing.T) {
		action := "Retrieving the stored config value of a path"

		scenarios := []struct {
			config   partial
			search   string
			expected interface{}
		}{
			{ // test the retrieving of a value of a present path
				config:   partial{"node": "value"},
				search:   "node",
				expected: "value",
			},
			{ // test the retrieving of a value of a missing path
				config:   partial{"node": "value"},
				search:   "inexistent-node",
				expected: nil,
			},
		}

		id := "source"
		priority := 0

		for _, scn := range scenarios {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			config, _ := NewConfig(60 * time.Second)
			defer config.Close()

			source := NewMockSource(ctrl)
			source.EXPECT().Close().Return(nil).Times(1)
			source.EXPECT().Get("").Return(scn.config).Times(1)

			config.AddSource(id, priority, source)

			if result := config.Get(scn.search); result != scn.expected {
				t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, scn.search, result, scn.expected)
			}
		}
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		p := partial{"node1": partial{"node2": 101}}
		path := "node3"
		expectedValue := 3

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(p).Times(1)

		config.AddSource(id, priority, source)

		if result := config.Get(path, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetBool(t *testing.T) {
	t.Run("should correctly return boolean value associated to the path", func(t *testing.T) {
		action := "Retrieving the boolean stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := true

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetBool(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a bool", func(t *testing.T) {
		action := "Retrieving the a non-boolean stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetBool(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := true

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetBool(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetInt(t *testing.T) {
	t.Run("should correctly return the int value associated to the path", func(t *testing.T) {
		action := "Retrieving the integer stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := 32

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetInt(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a int", func(t *testing.T) {
		action := "Retrieving the a non-integer stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetInt(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := 123

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetInt(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetInt8(t *testing.T) {
	t.Run("should correctly return the int8 value associated to the path", func(t *testing.T) {
		action := "Retrieving the int8 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value int8 = 32

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetInt8(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a int8", func(t *testing.T) {
		action := "Retrieving the a non-int8 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetInt8(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := int8(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetInt8(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetInt16(t *testing.T) {
	t.Run("should correctly return the int16 value associated to the path", func(t *testing.T) {
		action := "Retrieving the int16 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value int16 = 32

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetInt16(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a int16", func(t *testing.T) {
		action := "Retrieving the a non-int16 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetInt16(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := int16(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetInt16(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetInt32(t *testing.T) {
	t.Run("should correctly return the int32 value associated to the path", func(t *testing.T) {
		action := "Retrieving the int32 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value int32 = 32

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetInt32(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a int32", func(t *testing.T) {
		action := "Retrieving the a non-int32 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetInt32(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := int32(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetInt32(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetInt64(t *testing.T) {
	t.Run("should correctly return the int64 value associated to the path", func(t *testing.T) {
		action := "Retrieving the int64 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value int64 = 32

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetInt64(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a int16", func(t *testing.T) {
		action := "Retrieving the a non-int64 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetInt64(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := int64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetInt64(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetUInt(t *testing.T) {
	t.Run("should correctly return the uint value associated to the path", func(t *testing.T) {
		action := "Retrieving the uint stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value uint = 32

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetUInt(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a uint", func(t *testing.T) {
		action := "Retrieving the a non-uint stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetUInt(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := uint(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetUInt(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetUInt8(t *testing.T) {
	t.Run("should correctly return the uint8 value associated to the path", func(t *testing.T) {
		action := "Retrieving the uint8 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value uint8 = 32

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetUInt8(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a uint8", func(t *testing.T) {
		action := "Retrieving the a non-uint8 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetUInt8(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := uint8(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetUInt8(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetUInt16(t *testing.T) {
	t.Run("should correctly return the uint16 value associated to the path", func(t *testing.T) {
		action := "Retrieving the uint16 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value uint16 = 32

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetUInt16(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a uint16", func(t *testing.T) {
		action := "Retrieving the a non-uint16 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetUInt16(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := uint16(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetUInt16(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetUInt32(t *testing.T) {
	t.Run("should correctly return the uint32 value associated to the path", func(t *testing.T) {
		action := "Retrieving the uint32 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value uint32 = 32

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetUInt32(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a uint32", func(t *testing.T) {
		action := "Retrieving the a non-uint32 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetUInt32(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := uint32(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetUInt32(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetUInt64(t *testing.T) {
	t.Run("should correctly return the uint64 value associated to the path", func(t *testing.T) {
		action := "Retrieving the uint64 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value uint64 = 32

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetUInt64(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a uint64", func(t *testing.T) {
		action := "Retrieving the a non-uint64 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetUInt64(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := uint64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetUInt64(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetFloat32(t *testing.T) {
	t.Run("should correctly return the float32 value associated to the path", func(t *testing.T) {
		action := "Retrieving the float32 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value float32 = 32.0

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetFloat32(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a float32", func(t *testing.T) {
		action := "Retrieving the a non-float32 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetFloat32(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := float32(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetFloat32(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetFloat64(t *testing.T) {
	t.Run("should correctly return the float64 value associated to the path", func(t *testing.T) {
		action := "Retrieving the float64 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value float64 = 32.0

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetFloat64(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a float64", func(t *testing.T) {
		action := "Retrieving the a non-float64 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetFloat64(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := float64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetFloat64(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetComplex64(t *testing.T) {
	t.Run("should correctly return the complex64 value associated to the path", func(t *testing.T) {
		action := "Retrieving the complex64 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value complex64 = 32.0

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetComplex64(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a complex64", func(t *testing.T) {
		action := "Retrieving the a non-complex64 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetComplex64(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		var expectedValue complex64 = 1.0i

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetComplex64(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetComplex128(t *testing.T) {
	t.Run("should correctly return the complex128 value associated to the path", func(t *testing.T) {
		action := "Retrieving the complex128 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value complex128 = 32.0

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetComplex128(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a complex128", func(t *testing.T) {
		action := "Retrieving the a non-complex128 stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := "__invalid_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetComplex128(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		var expectedValue complex128 = 1.0i

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetComplex128(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetRune(t *testing.T) {
	t.Run("should correctly return the rune value associated to the path", func(t *testing.T) {
		action := "Retrieving the rune stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		var value rune = 'a'

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetRune(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a rune", func(t *testing.T) {
		action := "Retrieving the a non-rune stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := 123

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetRune(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := 'r'

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetRune(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_GetStrign(t *testing.T) {
	t.Run("should correctly return the string value associated to the path", func(t *testing.T) {
		action := "Retrieving the string stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := ""

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetString(path); result != value {
			t.Errorf("%s didn't retrieve the path (%v) value, returning %v, expected %v", action, path, result, value)
		}
	})

	t.Run("should panic if the stored value is not a string", func(t *testing.T) {
		action := "Retrieving the a non-string stored config value of a path"

		id := "source"
		priority := 0
		path := "node"
		value := 123

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{path: value}).Times(1)

		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, path)
			}
		}()

		config.GetString(path)
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		id := "source"
		priority := 0
		search := "node"
		expectedValue := "default"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)

		config.AddSource(id, priority, source)

		if result := config.GetString(search, expectedValue); result != expectedValue {
			t.Errorf("%s didn't retrieve the default value, returning %v, expected %v", action, result, expectedValue)
		}
	})
}

func Test_Config_HasSource(t *testing.T) {
	t.Run("should return true if the source is registed", func(t *testing.T) {
		action := "Checking the existence of a registered source"

		id := "source"
		priority := 0
		partial := partial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)

		config.AddSource(id, priority, source)

		if !config.HasSource(id) {
			t.Errorf("%s didn't correctly validated the existence of the source", action)
		}
	})

	t.Run("should return false if the source is not registed", func(t *testing.T) {
		action := "Checking the existence of a non-registered source"

		id := "source"
		priority := 0
		partial := partial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)

		config.AddSource(id, priority, source)

		if config.HasSource("source-inexistent") {
			t.Errorf("%s didn't correctly validated the non-existence of the source", action)
		}
	})
}

func Test_Config_AddSource(t *testing.T) {
	t.Run("should register a new source", func(t *testing.T) {
		action := "Registering a new source"

		id := "source"
		priority := 0
		partial := partial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)

		if err := config.AddSource(id, priority, source); err != nil {
			t.Errorf("%s returned the (%v) error while registering the source", action, err)
		}

		if !config.HasSource(id) {
			t.Errorf("%s didn't correctly validated the existence of the inserted source", action)
		}
	})

	t.Run("should return an error of not passing the source reference", func(t *testing.T) {
		action := "Registering a new source passing a nil reference"

		id := "source"
		priority := 0
		expected := "Invalid nil 'source' argument"

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		if err := config.AddSource(id, priority, nil); err == nil {
			t.Errorf("%s didn't returned an error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s returned error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return an error registring a duplicate id", func(t *testing.T) {
		action := "Registering a new source passing a existing id"

		id := "source"
		priority := 0
		partial := partial{}
		expected := "Duplicate source id : source"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)

		if err := config.AddSource(id, priority, source); err != nil {
			t.Errorf("%s returned the (%v) error while registering the source", action, err)
		}

		if err := config.AddSource(id, priority, source); err == nil {
			t.Errorf("%s didn't returned an error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s returned error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should override path if the inserted source have higher priority", func(t *testing.T) {
		action := "Inserting a higher priority source with overriden path"

		node := "node"
		id1 := "source.1"
		priority1 := 1
		value1 := "value1"
		partial1 := partial{node: value1}
		id2 := "source.2"
		priority2 := 2
		value2 := "value1"
		partial2 := partial{node: value2}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(partial1).AnyTimes()

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(partial2).AnyTimes()

		if err := config.AddSource(id1, priority1, source1); err != nil {
			t.Errorf("%s returned the (%v) error while registering the source", action, err)
		}
		if check := config.Get(node); check != value1 {
			t.Errorf("%s returned (%v) value after the insertion of the first source, expected (%v)", action, check, value1)
		}

		if err := config.AddSource(id2, priority2, source2); err != nil {
			t.Errorf("%s returned the (%v) error while registering the source", action, err)
		}
		if check := config.Get(node); check != value2 {
			t.Errorf("%s returned (%v) value after the insertion of the second source, expected (%v)", action, check, value2)
		}
	})

	t.Run("should not override path if the inserted source have lower priority", func(t *testing.T) {
		action := "Inserting a higher priority source with overriden path"

		node := "node"
		id1 := "source.1"
		priority1 := 2
		value1 := "value1"
		partial1 := partial{node: value1}
		id2 := "source.2"
		priority2 := 1
		value2 := "value1"
		partial2 := partial{node: value2}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(partial1).AnyTimes()

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(partial2).AnyTimes()

		if err := config.AddSource(id1, priority1, source1); err != nil {
			t.Errorf("%s returned the (%v) error while registering the source", action, err)
		}
		if check := config.Get(node); check != value1 {
			t.Errorf("%s returned (%v) value after the insertion of the first source, expected (%v)", action, check, value1)
		}

		if err := config.AddSource(id2, priority2, source2); err != nil {
			t.Errorf("%s returned the (%v) error while registering the source", action, err)
		}
		if check := config.Get(node); check != value1 {
			t.Errorf("%s returned (%v) value after the insertion of the second source, expected (%v)", action, check, value1)
		}
	})

	t.Run("should still be able to get not overriden paths of a inserted lower priority", func(t *testing.T) {
		action := "Inserting a lower priority source without overriden path"

		id1 := "source.1"
		priority1 := 2
		node1 := "node1"
		value1 := "value1"
		partial1 := partial{node1: value1}
		id2 := "source.2"
		priority2 := 1
		node2 := "node2"
		value2 := "value1"
		partial2 := partial{node2: value2}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(partial1).AnyTimes()

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(partial2).AnyTimes()

		if err := config.AddSource(id1, priority1, source1); err != nil {
			t.Errorf("%s returned the (%v) error while registering the source", action, err)
		}
		if err := config.AddSource(id2, priority2, source2); err != nil {
			t.Errorf("%s returned the (%v) error while registering the source", action, err)
		}
		if check := config.Get(node2); check != value2 {
			t.Errorf("%s returned (%v) value after the insertion of the second source, expected (%v)", action, check, value1)
		}
	})
}

func Test_Config_RemoveSource(t *testing.T) {
	t.Run("should unregister a proviuously registed source", func(t *testing.T) {
		action := "Unregistering a previously registed source"

		id1 := "source.1"
		priority1 := 0
		partial1 := partial{}
		id2 := "source.2"
		priority2 := 0
		partial2 := partial{}
		id3 := "source.3"
		priority3 := 0
		partial3 := partial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(partial1).AnyTimes()

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(partial2).AnyTimes()

		source3 := NewMockSource(ctrl)
		source3.EXPECT().Close().Times(1)
		source3.EXPECT().Get("").Return(partial3).AnyTimes()

		config.AddSource(id1, priority1, source1)
		config.AddSource(id2, priority2, source2)
		config.AddSource(id3, priority3, source3)

		config.RemoveSource(id2)

		if config.HasSource(id2) {
			t.Errorf("%s didn't correctly remove the inserted source", action)
		}
	})

	t.Run("should recover path overridden by the removed source", func(t *testing.T) {
		action := "Removing a higher priority source with overriden path"

		node := "node"
		id1 := "source.1"
		priority1 := 0
		value1 := "value.2"
		partial1 := partial{node: value1}
		id2 := "source.2"
		priority2 := 0
		value2 := "value.2"
		partial2 := partial{node: value2}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(partial1).AnyTimes()

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(partial2).AnyTimes()

		config.AddSource(id1, priority1, source1)
		config.AddSource(id2, priority2, source2)

		if check := config.Get(node); check != value2 {
			t.Errorf("%s returned (%v) value after the insertion of the second source, expected (%v)", action, check, value2)
		}

		config.RemoveSource(id2)

		if check := config.Get(node); check != value2 {
			t.Errorf("%s returned (%v) value after the removal of the high priority source, expected (%v)", action, check, value1)
		}
	})
}

func Test_Config_Source(t *testing.T) {
	t.Run("should return the registed source", func(t *testing.T) {
		action := "Retrieving a registered source"

		id := "source"
		priority := 0
		partial := partial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)

		config.AddSource(id, priority, source)

		result, err := config.Source(id)

		if err != nil {
			t.Errorf("%s returned the unexpected error : %v", action, err)
		}
		if result == nil {
			t.Errorf("%s didn't returned the expected source reference", action)
		} else {
			if !reflect.DeepEqual(result, source) {
				t.Errorf("%s returned (%v), expected (%v)", action, result, source)
			}
		}
	})

	t.Run("should return error if the source don't exists", func(t *testing.T) {
		action := "Retrieving a non-registered source"

		id := "source"
		priority := 0
		partial := partial{}
		requestID := "inexistent-source"
		expected := fmt.Sprintf("Source not found : %s", requestID)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)

		config.AddSource(id, priority, source)

		result, err := config.Source(requestID)

		if result != nil {
			t.Errorf("%s returned an unexpected reference to a source : %v", action, result)
		}
		if err == nil {
			t.Errorf("%s didn't returned the expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s returned (%s) error, expected (%s)", action, err.Error(), expected)
			}
		}
	})
}

func Test_Config_SourcePriority(t *testing.T) {
	t.Run("should return a error if the source was not found", func(t *testing.T) {
		action := "Updating a priority of a non-existing source"

		id := "inexistent-source"
		expected := fmt.Sprintf("Source not found : %s", id)

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		if err := config.SourcePriority(id, 0); err == nil {
			t.Errorf("%s didn't returned a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s returned the (%s) error, expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should update the priority of the source", func(t *testing.T) {
		action := "Lowering the priority of the higher prioritized source"

		node := "node"
		id1 := "source.1"
		priority1 := 1
		value1 := "value1"
		partial1 := partial{node: value1}
		id2 := "source.2"
		priority2 := 2
		value2 := "value1"
		partial2 := partial{node: value2}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(partial1).AnyTimes()

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(partial2).AnyTimes()

		config.AddSource(id1, priority1, source1)
		config.AddSource(id2, priority2, source2)

		if result := config.Get(node); result != value2 {
			t.Errorf("%s returned the (%v) prior the priority change, expected (%v)", action, result, value2)
		}
		if err := config.SourcePriority(id2, 0); err != nil {
			t.Errorf("%s returned an unexpected error : %v", action, err)
		}
		if result := config.Get(node); result != value1 {
			t.Errorf("%s returned the (%v) prior the priority change, expected (%v)", action, result, value1)
		}
	})
}

func Test_Config_HasObserver(t *testing.T) {
	t.Run("should check correctly the existence of a registed observer", func(t *testing.T) {
		action := "Checking the existence of a observer"

		scenarios := []struct {
			observers []string
			search    string
			expected  bool
		}{
			{ // Search a non-existing path in a empty list of observers
				observers: []string{},
				search:    "node1",
				expected:  false,
			},
			{ // Search a non-existing path in a non-empty list of observers
				observers: []string{"node1", "node2"},
				search:    "node3",
				expected:  false,
			},
			{ // Search a existing path in alist of observers
				observers: []string{"node1", "node2", "node3"},
				search:    "node2",
				expected:  true,
			},
		}

		for _, scn := range scenarios {
			config, _ := NewConfig(60 * time.Second)
			defer config.Close()

			for _, observer := range scn.observers {
				config.AddObserver(observer, func(old, new interface{}) {})
			}

			if check := config.HasObserver(scn.search); check != scn.expected {
				t.Errorf("%s returned (%v), expected (%v) when requesting for (%s), in (%v)", action, check, scn.expected, scn.search, scn.observers)
			}
		}
	})

	t.Run("should return nil if trying to register a nil callback", func(t *testing.T) {
		action := "Registering a new observer with a nil callback"

		observer := "path"
		expected := "Invalid nil 'callback' argument"

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		err := config.AddObserver(observer, nil)

		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})
}

func Test_Config_RemoveObserver(t *testing.T) {
	t.Run("should check correctly remove a registed observer", func(t *testing.T) {
		action := "Removing a registed observer"

		observer1 := "node.1"
		observer2 := "node.2"
		observer3 := "node.3"

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		config.AddObserver(observer1, func(old, new interface{}) {})
		config.AddObserver(observer2, func(old, new interface{}) {})
		config.AddObserver(observer3, func(old, new interface{}) {})

		config.RemoveObserver(observer2)

		if config.HasObserver(observer2) {
			t.Errorf("%s didn't correctly removed the observer", action)
		}
	})
}

func Test_Config(t *testing.T) {
	t.Run("should check for reload on observable sources after period time", func(t *testing.T) {
		id := "source"
		priority := 0

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(20 * time.Millisecond)
		defer config.Close()

		source := NewMockObservableSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial{}).AnyTimes()
		source.EXPECT().Reload().Return(false, nil).MinTimes(1)

		config.AddSource(id, priority, source)

		time.Sleep(60 * time.Millisecond)
	})

	t.Run("should not rebuild if the observable source does not notify changes", func(t *testing.T) {
		id := "source"
		priority := 0

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(20 * time.Millisecond)
		defer config.Close()

		source := NewMockObservableSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		source.EXPECT().Reload().Return(false, nil).MinTimes(1)

		config.AddSource(id, priority, source)

		time.Sleep(60 * time.Millisecond)
	})

	t.Run("should rebuild if the observable source notify changes", func(t *testing.T) {
		action := "Retrieving the observable source value after reload"

		id := "source"
		priority := 0
		node := "node"
		value := "value"
		partial := partial{node: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(20 * time.Millisecond)
		defer config.Close()

		source := NewMockObservableSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Reload().Return(true, nil).MinTimes(1)
		source.EXPECT().Get("").Return(partial).MinTimes(2)

		config.AddSource(id, priority, source)

		time.Sleep(60 * time.Millisecond)

		if check := config.Get(node); check != value {
			t.Errorf("%s returned (%v), expected (%v)", action, check, value)
		}
	})

	t.Run("should call observer callback function on config changes", func(t *testing.T) {
		action := "Expecting a observer callback execution on config change"

		id := "source"
		priority := 0
		node := "node"
		value := "value"
		check := false

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(20 * time.Millisecond)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial{node: value})

		config.AddObserver(node, func(old, new interface{}) {
			check = true

			if old != nil {
				t.Errorf("%s, callback called with (%v) as old value, expected nil", action, old)
			}
			if new != value {
				t.Errorf("%s, callback called with (%v) as old value, expected (%v)", action, new, value)
			}
		})

		config.AddSource(id, priority, source)

		if !check {
			t.Errorf("%s didn't actually called the callback", action)
		}
	})
}
