package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_NewConfig(t *testing.T) {
	t.Run("creates a new config", func(t *testing.T) {
		if config, err := NewConfig(60 * time.Second); config == nil {
			t.Errorf("didn't return a valid reference")
		} else {
			config.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

func Test_Config_Close(t *testing.T) {
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

	t.Run("propagate close to sources", func(t *testing.T) {
		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Return(nil).Times(1)
		source1.EXPECT().Get("").Return(partial1).AnyTimes()
		config.AddSource(id1, priority1, source1)

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Return(nil).Times(1)
		source2.EXPECT().Get("").Return(partial2).AnyTimes()
		config.AddSource(id2, priority2, source2)
	})
}

func Test_Config_Has(t *testing.T) {
	id := "source"
	priority := 0

	t.Run("return the existence of the path", func(t *testing.T) {
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
				t.Errorf("returned (%v) when expected (%v)", result, scn.expected)
			}
		}
	})
}

func Test_Config_Get(t *testing.T) {
	id := "source"
	priority := 0
	p := partial{"node1": partial{"node2": 101}}
	search := "node3"
	defValue := 3

	t.Run("return path value", func(t *testing.T) {
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
				t.Errorf("returned (%v) when expected (%v)", result, scn.expected)
			}
		}
	})

	t.Run("return default if path was not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(p).Times(1)
		config.AddSource(id, priority, source)

		if result := config.Get(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetBool(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	value := true
	invalidValue := "true"
	defValue := true

	t.Run("return boolean value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetBool(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a bool", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetBool(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetBool(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetInt(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	value := 12
	invalidValue := "12"
	defValue := 34

	t.Run("return int value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetInt(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a int", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetInt(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetInt(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetInt8(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value int8 = 12
	invalidValue := "12"
	var defValue int8 = 34

	t.Run("return int8 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetInt8(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a int8", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetInt8(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetInt8(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetInt16(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value int16 = 12
	invalidValue := "12"
	var defValue int16 = 34

	t.Run("return int16 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetInt16(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a int16", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetInt16(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetInt16(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetInt32(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value int32 = 12
	invalidValue := "12"
	var defValue int32 = 34

	t.Run("return int32 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetInt32(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a int32", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetInt32(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetInt32(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetInt64(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value int64 = 12
	invalidValue := "12"
	var defValue int64 = 34

	t.Run("return int64 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetInt64(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a int16", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetInt64(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetInt64(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetUInt(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value uint = 12
	invalidValue := "12"
	var defValue uint = 34

	t.Run("return uint value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetUInt(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a uint", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetUInt(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetUInt(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetUInt8(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value uint8 = 12
	invalidValue := "12"
	var defValue uint8 = 34

	t.Run("return uint8 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetUInt8(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a uint8", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetUInt8(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetUInt8(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetUInt16(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value uint16 = 12
	invalidValue := "12"
	var defValue uint16 = 34

	t.Run("return uint16 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetUInt16(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a uint16", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetUInt16(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetUInt16(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetUInt32(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value uint32 = 12
	invalidValue := "12"
	var defValue uint32 = 34

	t.Run("return uint32 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetUInt32(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a uint32", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetUInt32(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetUInt32(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetUInt64(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value uint64 = 12
	invalidValue := "12"
	var defValue uint64 = 34

	t.Run("return uint64 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetUInt64(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a uint64", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetUInt64(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetUInt64(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetFloat32(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value float32 = 12.0
	invalidValue := "12.0"
	var defValue float32 = 34.0

	t.Run("return float32 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetFloat32(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a float32", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetFloat32(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetFloat32(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetFloat64(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value float64 = 12.0
	invalidValue := "12.0"
	var defValue float64 = 34.0

	t.Run("return float64 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetFloat64(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a float64", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetFloat64(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetFloat64(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetComplex64(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value complex64 = 12.0
	invalidValue := "12.0"
	var defValue complex64 = 34.0

	t.Run("return complex64 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetComplex64(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a complex64", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetComplex64(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetComplex64(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetComplex128(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value complex128 = 12.0
	invalidValue := "12.0"
	var defValue complex128 = 34.0

	t.Run("return complex128 value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetComplex128(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a complex128", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetComplex128(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetComplex128(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetRune(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value rune = 'r'
	invalidValue := "12.0"
	var defValue rune = 'w'

	t.Run("return rune value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetRune(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a rune", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetRune(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetRune(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetString(t *testing.T) {
	id := "source"
	priority := 0
	search := "node"
	var value string = "value"
	invalidValue := 12.0
	var defValue string = "default value"

	t.Run("return string value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: value}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetString(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{search: invalidValue}).Times(1)
		config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		config.GetString(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial{}).Times(1)
		config.AddSource(id, priority, source)

		if result := config.GetString(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_HasSource(t *testing.T) {
	id := "source"
	priority := 0
	partial := partial{}

	t.Run("validate if the source is registed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)
		config.AddSource(id, priority, source)

		if !config.HasSource(id) {
			t.Errorf("returned false")
		}
	})

	t.Run("invalidate if the source is not registed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Return(nil).Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)
		config.AddSource(id, priority, source)

		if config.HasSource("source-inexistent") {
			t.Errorf("returned true")
		}
	})
}

func Test_Config_AddSource(t *testing.T) {
	id := "source"
	priority := 1
	p := partial{}

	t.Run("error of not passing a source", func(t *testing.T) {
		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		if err := config.AddSource(id, priority, nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'source' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register a new source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(p).Times(1)

		if err := config.AddSource(id, priority, source); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !config.HasSource(id) {
			t.Errorf("didn't stored the source")
		}
	})

	t.Run("error registring a duplicate id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(p).Times(1)
		config.AddSource(id, priority, source)

		if err := config.AddSource(id, priority, source); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Duplicate source id : source" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	node := "node"
	id1 := "source.1"
	value1 := "value1"
	p1 := partial{node: value1}
	id2 := "source.2"
	value2 := "value2"
	p2 := partial{node: value2}

	t.Run("override path if the insert have higher priority", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(p1).AnyTimes()
		config.AddSource(id1, priority, source1)

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(p2).AnyTimes()
		config.AddSource(id2, priority+1, source2)

		if check := config.Get(node); check != value2 {
			t.Errorf("returned the (%v) value", check)
		}
	})

	t.Run("do not override path if the insert have lower priority", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(p1).AnyTimes()
		config.AddSource(id1, priority, source1)

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(p2).AnyTimes()
		config.AddSource(id2, priority-1, source2)

		if check := config.Get(node); check != value1 {
			t.Errorf("returned the (%v) value", check)
		}
	})

	extendedNode := "extendedNode"
	extendedValue := "extraValue"
	extendedPartial := p2.merge(partial{extendedNode: extendedValue})

	t.Run("still be able to get not overriden paths of a inserted lower priority", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(p1).AnyTimes()
		config.AddSource(id1, priority, source1)

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(extendedPartial).AnyTimes()
		config.AddSource(id2, priority-1, source2)

		if check := config.Get(extendedNode); check != extendedValue {
			t.Errorf("returned the (%v) value", check)
		}
	})
}

func Test_Config_RemoveSource(t *testing.T) {
	id1 := "source.1"
	priority1 := 0
	node := "node"
	value1 := "value.1"
	p1 := partial{node: value1}
	id2 := "source.2"
	priority2 := 0
	value2 := "value.2"
	p2 := partial{node: value2}
	id3 := "source.3"
	priority3 := 0
	value3 := "value.3"
	p3 := partial{node: value3}

	t.Run("unregister a proviuously registed source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(p1).AnyTimes()
		config.AddSource(id1, priority1, source1)

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(p2).AnyTimes()
		config.AddSource(id2, priority2, source2)

		source3 := NewMockSource(ctrl)
		source3.EXPECT().Close().Times(1)
		source3.EXPECT().Get("").Return(p3).AnyTimes()
		config.AddSource(id3, priority3, source3)

		config.RemoveSource(id2)

		if config.HasSource(id2) {
			t.Errorf("didn't remove the source")
		}
	})

	t.Run("recover path overridden by the removed source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(p1).AnyTimes()
		config.AddSource(id1, priority1, source1)

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(p2).AnyTimes()
		config.AddSource(id2, priority2, source2)

		source3 := NewMockSource(ctrl)
		source3.EXPECT().Close().Times(1)
		source3.EXPECT().Get("").Return(p3).AnyTimes()
		config.AddSource(id3, priority3, source3)

		config.RemoveSource(id3)

		if value := config.Get(node); value != value2 {
			t.Errorf("returned (%v) value", value)
		}
	})
}

func Test_Config_Source(t *testing.T) {
	id := "source"
	priority := 0
	partial := partial{}
	inexistingID := "inexistent-source"
	expectedError := "Source not found : inexistent-source"

	t.Run("error if the source don't exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)
		config.AddSource(id, priority, source)

		if result, err := config.Source(inexistingID); result != nil {
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("return the registed source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)
		config.AddSource(id, priority, source)

		if result, err := config.Source(id); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if result == nil {
			t.Errorf("returned nil")
		} else if !reflect.DeepEqual(result, source) {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_SourcePriority(t *testing.T) {
	inexistingID := "inexistent-source"
	expectedError := "Source not found : inexistent-source"

	t.Run("error if the source was not found", func(t *testing.T) {
		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		if err := config.SourcePriority(inexistingID, 0); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	node := "node"
	id1 := "source.1"
	priority1 := 1
	value1 := "value1"
	partial1 := partial{node: value1}
	id2 := "source.2"
	priority2 := 2
	value2 := "value1"
	partial2 := partial{node: value2}

	t.Run("update the priority of the source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(partial1).AnyTimes()
		config.AddSource(id1, priority1, source1)

		source2 := NewMockSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(partial2).AnyTimes()
		config.AddSource(id2, priority2, source2)

		if result := config.Get(node); result != value2 {
			t.Errorf("returned the (%v) value prior the change", result)
		}
		if err := config.SourcePriority(id2, 0); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
		if result := config.Get(node); result != value1 {
			t.Errorf("returned the (%v) value after the change", result)
		}
	})
}

func Test_Config_HasObserver(t *testing.T) {
	t.Run("check the existence of a observer", func(t *testing.T) {
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
				t.Errorf("returned (%v)", check)
			}
		}
	})
}

func Test_Config_AddObserver(t *testing.T) {
	search := "path"
	expectedError := "Invalid nil 'callback' argument"

	t.Run("error on registering a nil callback", func(t *testing.T) {
		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		if err := config.AddObserver(search, nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("return the (%v) error", err)
		}
	})
}

func Test_Config_RemoveObserver(t *testing.T) {
	observer1 := "node.1"
	observer2 := "node.2"
	observer3 := "node.3"

	t.Run("remove a registed observer", func(t *testing.T) {
		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		config.AddObserver(observer1, func(old, new interface{}) {})
		config.AddObserver(observer2, func(old, new interface{}) {})
		config.AddObserver(observer3, func(old, new interface{}) {})
		config.RemoveObserver(observer2)

		if config.HasObserver(observer2) {
			t.Errorf("didn't removed the observer")
		}
	})
}

func Test_Config(t *testing.T) {
	id := "source"
	priority := 0
	node := "node"
	value := "value"
	partial := partial{node: value}

	t.Run("reload on observable sources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(20 * time.Millisecond)
		defer config.Close()

		source := NewMockObservableSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)
		source.EXPECT().Reload().Return(false, nil).MinTimes(1)
		config.AddSource(id, priority, source)

		time.Sleep(60 * time.Millisecond)
	})

	t.Run("rebuild if the observable source notify changes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(20 * time.Millisecond)
		defer config.Close()

		source := NewMockObservableSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).MinTimes(2)
		source.EXPECT().Reload().Return(true, nil).MinTimes(1)
		config.AddSource(id, priority, source)

		time.Sleep(60 * time.Millisecond)

		if check := config.Get(node); check != value {
			t.Errorf("returned (%v)", check)
		}
	})

	check := false

	t.Run("should call observer callback function on config changes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(20 * time.Millisecond)
		defer config.Close()

		config.AddObserver(node, func(old, new interface{}) {
			check = true

			if old != nil {
				t.Errorf("callback called with (%v) as old value", old)
			}
			if new != value {
				t.Errorf("callback called with (%v) as new value", new)
			}
		})

		source := NewMockSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial)
		config.AddSource(id, priority, source)

		if !check {
			t.Errorf("didn't actually called the callback")
		}
	})
}
