package servlet

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
	"time"
)

func Test_NewConfig(t *testing.T) {
	t.Run("new config without reload", func(t *testing.T) {
		if config, err := NewConfig(0 * time.Second); config == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			defer config.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if config.mutex == nil {
				t.Error("didn't instantiate the access mutex")
			} else if config.sources == nil {
				t.Error("didn't instantiate the sources storing array")
			} else if config.observers == nil {
				t.Error("didn't instantiate the observers storing array")
			} else if config.loader != nil {
				t.Error("instantiated the sources reload trigger")
			}
		}
	})

	t.Run("new config with reload", func(t *testing.T) {
		if config, err := NewConfig(60 * time.Second); config == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			defer config.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if config.mutex == nil {
				t.Error("didn't instantiate the access mutex")
			} else if config.sources == nil {
				t.Error("didn't instantiate the sources storing array")
			} else if config.observers == nil {
				t.Error("didn't instantiate the observers storing array")
			} else if config.loader == nil {
				t.Error("didn't instantiate the sources reload trigger")
			}
		}
	})
}

func Test_Config_Close(t *testing.T) {
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

		var config *Config
		config.Close()
	})

	t.Run("propagate close to sources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)

		id1 := "source.1"
		priority1 := 0
		source1 := NewMockConfigSource(ctrl)
		source1.EXPECT().Get("").Return(ConfigPartial{}).AnyTimes()
		source1.EXPECT().Close().Times(1)
		_ = config.AddSource(id1, priority1, source1)

		id2 := "source.2"
		priority2 := 1
		source2 := NewMockConfigSource(ctrl)
		source2.EXPECT().Get("").Return(ConfigPartial{}).AnyTimes()
		source2.EXPECT().Close().Times(1)
		_ = config.AddSource(id2, priority2, source2)

		config.Close()
	})
}

func Test_Config_Has(t *testing.T) {
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

		var config *Config
		config.Has("path")
	})

	t.Run("return the existence of the path", func(t *testing.T) {
		scenarios := []struct {
			config   ConfigPartial
			search   string
			expected bool
		}{
			{ // test the existence of a present path
				config:   ConfigPartial{"node": "value"},
				search:   "node",
				expected: true,
			},
			{ // test the non-existence of a missing path
				config:   ConfigPartial{"node": "value"},
				search:   "invalid-node",
				expected: false,
			},
		}

		for _, scn := range scenarios {
			id := "source"
			priority := 0

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			config, _ := NewConfig(60 * time.Second)
			defer config.Close()

			source := NewMockConfigSource(ctrl)
			source.EXPECT().Close().Times(1)
			source.EXPECT().Get("").Return(scn.config).Times(1)
			_ = config.AddSource(id, priority, source)

			if result := config.Has(scn.search); result != scn.expected {
				t.Errorf("returned (%v) when expected (%v)", result, scn.expected)
			}
		}
	})
}

func Test_Config_Get(t *testing.T) {
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

		var config *Config
		config.Get("path")
	})

	t.Run("return path value", func(t *testing.T) {
		scenarios := []struct {
			config   ConfigPartial
			search   string
			expected interface{}
		}{
			{ // test the retrieving of a value of a present path
				config:   ConfigPartial{"node": "value"},
				search:   "node",
				expected: "value",
			},
			{ // test the retrieving of a value of a missing path
				config:   ConfigPartial{"node": "value"},
				search:   "invalid-node",
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

			source := NewMockConfigSource(ctrl)
			source.EXPECT().Close().Times(1)
			source.EXPECT().Get("").Return(scn.config).Times(1)
			_ = config.AddSource(id, priority, source)

			if result := config.Get(scn.search); result != scn.expected {
				t.Errorf("returned (%v) when expected (%v)", result, scn.expected)
			}
		}
	})

	t.Run("return default if path was not found", func(t *testing.T) {
		id := "source"
		priority := 0
		p := ConfigPartial{"node1": ConfigPartial{"node2": 101}}
		search := "node3"
		defValue := 3

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(p).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.Get(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetBool(t *testing.T) {
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

		var config *Config
		config.GetBool("path")
	})

	t.Run("return boolean value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: true}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetBool(search); !result {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a bool", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "true"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetBool(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetBool(search, true); !result {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetInt(t *testing.T) {
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

		var config *Config
		config.GetInt("path")
	})

	t.Run("return int value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := 123

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetInt(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a int", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetInt(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := 123

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetInt(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetInt8(t *testing.T) {
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

		var config *Config
		config.GetInt8("path")
	})

	t.Run("return int8 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := int8(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetInt8(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a int8", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetInt8(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := int8(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetInt8(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetInt16(t *testing.T) {
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

		var config *Config
		config.GetInt16("path")
	})

	t.Run("return int16 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := int16(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetInt16(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a int16", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetInt16(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := int16(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetInt16(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetInt32(t *testing.T) {
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

		var config *Config
		config.GetInt32("path")
	})

	t.Run("return int32 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := int32(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetInt32(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a int32", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetInt32(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := int32(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetInt32(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetInt64(t *testing.T) {
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

		var config *Config
		config.GetInt64("path")
	})

	t.Run("return int64 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := int64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetInt64(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a int64", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetInt64(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := int64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetInt64(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetUInt(t *testing.T) {
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

		var config *Config
		config.GetUInt("path")
	})

	t.Run("return uint value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := uint(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetUInt(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a uint", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetUInt(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := uint(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetUInt(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetUInt8(t *testing.T) {
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

		var config *Config
		config.GetUInt8("path")
	})

	t.Run("return uint8 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := uint8(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetUInt8(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a uint8", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetUInt8(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := uint8(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetUInt8(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetUInt16(t *testing.T) {
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

		var config *Config
		config.GetUInt16("path")
	})

	t.Run("return int16 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := uint16(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetUInt16(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a uint16", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetUInt16(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := uint16(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetUInt16(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetUInt32(t *testing.T) {
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

		var config *Config
		config.GetUInt32("path")
	})

	t.Run("return uint32 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := uint32(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetUInt32(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a uint32", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetUInt32(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := uint32(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetUInt32(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetUInt64(t *testing.T) {
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

		var config *Config
		config.GetUInt64("path")
	})

	t.Run("return uint64 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := uint64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetUInt64(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a uint64", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetUInt64(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := uint64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetUInt64(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetFloat32(t *testing.T) {
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

		var config *Config
		config.GetFloat32("path")
	})

	t.Run("return float32 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := float32(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetFloat32(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a float32", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetFloat32(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := float32(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetFloat32(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetFloat64(t *testing.T) {
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

		var config *Config
		config.GetFloat64("path")
	})

	t.Run("return float64 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := float64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetFloat64(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a float64", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetFloat64(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := float64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetFloat64(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetComplex64(t *testing.T) {
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

		var config *Config
		config.GetComplex64("path")
	})

	t.Run("return complex64 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := complex64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetComplex64(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a complex64", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetComplex64(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := complex64(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetComplex64(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetComplex128(t *testing.T) {
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

		var config *Config
		config.GetComplex128("path")
	})

	t.Run("return complex128 value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := complex128(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetComplex128(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a complex128", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetComplex128(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := complex128(123)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetComplex128(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetRune(t *testing.T) {
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

		var config *Config
		config.GetRune("path")
	})

	t.Run("return rune value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := 'r'

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetRune(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a rune", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetRune(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := 'r'

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetRune(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_GetString(t *testing.T) {
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

		var config *Config
		config.GetString("path")
	})

	t.Run("return string value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		value := "value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: value}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetString(search); result != value {
			t.Errorf("returned (%v)", result)
		}
	})

	t.Run("panic if the stored value is not a string", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		invalidValue := 123

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{search: invalidValue}).Times(1)
		_ = config.AddSource(id, priority, source)

		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()

		config.GetString(search)
	})

	t.Run("return the default value", func(t *testing.T) {
		id := "source"
		priority := 0
		search := "node"
		defValue := "value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(ConfigPartial{}).Times(1)
		_ = config.AddSource(id, priority, source)

		if result := config.GetString(search, defValue); result != defValue {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_HasSource(t *testing.T) {
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

		var config *Config
		_ = config.HasSource("path")
	})

	t.Run("validate if the source is registered", func(t *testing.T) {
		id := "source"
		priority := 0
		partial := ConfigPartial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)
		_ = config.AddSource(id, priority, source)

		if !config.HasSource(id) {
			t.Error("returned false")
		}
	})

	t.Run("invalidate if the source is not registered", func(t *testing.T) {
		id := "source"
		priority := 0
		partial := ConfigPartial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)
		_ = config.AddSource(id, priority, source)

		if config.HasSource("invalid source id") {
			t.Error("returned true")
		}
	})
}

func Test_Config_AddSource(t *testing.T) {
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

		var config *Config
		_ = config.AddSource("id", 1, nil)
	})

	t.Run("nil source", func(t *testing.T) {
		id := "source"
		priority := 1

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		if err := config.AddSource(id, priority, nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'source' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register a new source", func(t *testing.T) {
		id := "source"
		priority := 1
		p := ConfigPartial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(p).Times(1)

		if err := config.AddSource(id, priority, source); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !config.HasSource(id) {
			t.Error("didn't stored the source")
		}
	})

	t.Run("duplicate id", func(t *testing.T) {
		id := "source"
		priority := 1
		p := ConfigPartial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(p).Times(1)
		_ = config.AddSource(id, priority, source)

		if err := config.AddSource(id, priority, source); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "duplicate source id : source" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("override path if the insert have higher priority", func(t *testing.T) {
		node := "node"
		priority := 1

		id1 := "source.1"
		value1 := "value1"
		p1 := ConfigPartial{node: value1}

		id2 := "source.2"
		value2 := "value2"
		p2 := ConfigPartial{node: value2}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockConfigSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(p1).AnyTimes()
		_ = config.AddSource(id1, priority, source1)

		source2 := NewMockConfigSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(p2).AnyTimes()
		_ = config.AddSource(id2, priority+1, source2)

		if check := config.Get(node); check != value2 {
			t.Errorf("returned the (%v) value", check)
		}
	})

	t.Run("do not override path if the insert have lower priority", func(t *testing.T) {
		node := "node"
		priority := 1

		id1 := "source.1"
		value1 := "value1"
		p1 := ConfigPartial{node: value1}

		id2 := "source.2"
		value2 := "value2"
		p2 := ConfigPartial{node: value2}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockConfigSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(p1).AnyTimes()
		_ = config.AddSource(id1, priority, source1)

		source2 := NewMockConfigSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(p2).AnyTimes()
		_ = config.AddSource(id2, priority-1, source2)

		if check := config.Get(node); check != value1 {
			t.Errorf("returned the (%v) value", check)
		}
	})

	t.Run("still be able to get not overridden paths of a inserted lower priority", func(t *testing.T) {
		node := "node"
		priority := 1

		id1 := "source.1"
		value1 := "value1"
		p1 := ConfigPartial{node: value1}

		id2 := "source.2"
		value2 := "value2"
		p2 := ConfigPartial{node: value2}

		extendedNode := "extendedNode"
		extendedValue := "extraValue"
		extendedPartial := p2.merge(ConfigPartial{extendedNode: extendedValue})

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockConfigSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(p1).AnyTimes()
		_ = config.AddSource(id1, priority, source1)

		source2 := NewMockConfigSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(extendedPartial).AnyTimes()
		_ = config.AddSource(id2, priority-1, source2)

		if check := config.Get(extendedNode); check != extendedValue {
			t.Errorf("returned the (%v) value", check)
		}
	})
}

func Test_Config_RemoveSource(t *testing.T) {
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

		var config *Config
		config.RemoveSource("id")
	})

	t.Run("unregister a previously registered source", func(t *testing.T) {
		id1 := "source.1"
		priority1 := 0
		node := "node"
		value1 := "value.1"
		p1 := ConfigPartial{node: value1}

		id2 := "source.2"
		priority2 := 0
		value2 := "value.2"
		p2 := ConfigPartial{node: value2}

		id3 := "source.3"
		priority3 := 0
		value3 := "value.3"
		p3 := ConfigPartial{node: value3}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockConfigSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(p1).AnyTimes()
		_ = config.AddSource(id1, priority1, source1)

		source2 := NewMockConfigSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(p2).AnyTimes()
		_ = config.AddSource(id2, priority2, source2)

		source3 := NewMockConfigSource(ctrl)
		source3.EXPECT().Close().Times(1)
		source3.EXPECT().Get("").Return(p3).AnyTimes()
		_ = config.AddSource(id3, priority3, source3)

		config.RemoveSource(id2)

		if config.HasSource(id2) {
			t.Error("didn't remove the source")
		}
	})

	t.Run("recover path overridden by the removed source", func(t *testing.T) {
		id1 := "source.1"
		priority1 := 0
		node := "node"
		value1 := "value.1"
		p1 := ConfigPartial{node: value1}

		id2 := "source.2"
		priority2 := 0
		value2 := "value.2"
		p2 := ConfigPartial{node: value2}

		id3 := "source.3"
		priority3 := 0
		value3 := "value.3"
		p3 := ConfigPartial{node: value3}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockConfigSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(p1).AnyTimes()
		_ = config.AddSource(id1, priority1, source1)

		source2 := NewMockConfigSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(p2).AnyTimes()
		_ = config.AddSource(id2, priority2, source2)

		source3 := NewMockConfigSource(ctrl)
		source3.EXPECT().Close().Times(1)
		source3.EXPECT().Get("").Return(p3).AnyTimes()
		_ = config.AddSource(id3, priority3, source3)

		config.RemoveSource(id3)

		if value := config.Get(node); value != value2 {
			t.Errorf("returned (%v) value", value)
		}
	})
}

func Test_Config_Source(t *testing.T) {
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

		var config *Config
		_, _ = config.Source("id")
	})

	t.Run("error if the source don't exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		if result, err := config.Source("invalid id"); result != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "source not found : invalid id" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("return the registered source", func(t *testing.T) {
		id := "source"
		priority := 0
		partial := ConfigPartial{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)
		_ = config.AddSource(id, priority, source)

		if result, err := config.Source(id); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if result == nil {
			t.Error("returned nil")
		} else if !reflect.DeepEqual(result, source) {
			t.Errorf("returned (%v)", result)
		}
	})
}

func Test_Config_SourcePriority(t *testing.T) {
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

		var config *Config
		_ = config.SourcePriority("id", 0)
	})

	t.Run("error if the source was not found", func(t *testing.T) {
		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		if err := config.SourcePriority("invalid id", 0); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "source not found : invalid id" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("update the priority of the source", func(t *testing.T) {
		node := "node"
		id1 := "source.1"
		priority1 := 1
		value1 := "value1"
		partial1 := ConfigPartial{node: value1}

		id2 := "source.2"
		priority2 := 2
		value2 := "value1"
		partial2 := ConfigPartial{node: value2}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		source1 := NewMockConfigSource(ctrl)
		source1.EXPECT().Close().Times(1)
		source1.EXPECT().Get("").Return(partial1).AnyTimes()
		_ = config.AddSource(id1, priority1, source1)

		source2 := NewMockConfigSource(ctrl)
		source2.EXPECT().Close().Times(1)
		source2.EXPECT().Get("").Return(partial2).AnyTimes()
		_ = config.AddSource(id2, priority2, source2)

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

		var config *Config
		_ = config.HasObserver("id")
	})

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
			{ // Search a existing path in a list of observers
				observers: []string{"node1", "node2", "node3"},
				search:    "node2",
				expected:  true,
			},
		}

		for _, scn := range scenarios {
			config, _ := NewConfig(0 * time.Second)
			config.Close()

			for _, observer := range scn.observers {
				_ = config.AddObserver(observer, func(old, new interface{}) {})
			}

			if check := config.HasObserver(scn.search); check != scn.expected {
				t.Errorf("returned (%v)", check)
			}
		}
	})
}

func Test_Config_AddObserver(t *testing.T) {
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

		var config *Config
		_ = config.AddObserver("path", nil)
	})

	t.Run("nil callback", func(t *testing.T) {
		search := "path"
		expectedError := "invalid nil 'callback' argument"

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		if err := config.AddObserver(search, nil); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("return the (%v) error", err)
		}
	})
}

func Test_Config_RemoveObserver(t *testing.T) {
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

		var config *Config
		config.RemoveObserver("path")
	})

	t.Run("remove a registered observer", func(t *testing.T) {
		observer1 := "node.1"
		observer2 := "node.2"
		observer3 := "node.3"

		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		_ = config.AddObserver(observer1, func(old, new interface{}) {})
		_ = config.AddObserver(observer2, func(old, new interface{}) {})
		_ = config.AddObserver(observer3, func(old, new interface{}) {})
		config.RemoveObserver(observer2)

		if config.HasObserver(observer2) {
			t.Errorf("didn't removed the observer")
		}
	})
}

func Test_Config(t *testing.T) {
	t.Run("reload on observable sources", func(t *testing.T) {
		id := "source"
		priority := 0
		node := "node"
		value := "value"
		partial := ConfigPartial{node: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(20 * time.Millisecond)
		defer config.Close()

		source := NewMockConfigObservableSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).Times(1)
		source.EXPECT().Reload().Return(false, nil).MinTimes(1)
		_ = config.AddSource(id, priority, source)

		time.Sleep(60 * time.Millisecond)
	})

	t.Run("rebuild if the observable source notify changes", func(t *testing.T) {
		id := "source"
		priority := 0
		node := "node"
		value := "value"
		partial := ConfigPartial{node: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(20 * time.Millisecond)
		defer config.Close()

		source := NewMockConfigObservableSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial).MinTimes(2)
		source.EXPECT().Reload().Return(true, nil).MinTimes(1)
		_ = config.AddSource(id, priority, source)

		time.Sleep(60 * time.Millisecond)

		if check := config.Get(node); check != value {
			t.Errorf("returned (%v)", check)
		}
	})

	t.Run("should call observer callback function on config changes", func(t *testing.T) {
		id := "source"
		priority := 0
		node := "node"
		value := "value"
		partial := ConfigPartial{node: value}
		check := false

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(20 * time.Millisecond)
		defer config.Close()

		_ = config.AddObserver(node, func(old, new interface{}) {
			check = true

			if old != nil {
				t.Errorf("callback called with (%v) as old value", old)
			}
			if new != value {
				t.Errorf("callback called with (%v) as new value", new)
			}
		})

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Close().Times(1)
		source.EXPECT().Get("").Return(partial)
		_ = config.AddSource(id, priority, source)

		if !check {
			t.Errorf("didn't actually called the callback")
		}
	})
}
