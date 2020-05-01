package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

/// ---------------------------------------------------------------------------
/// ConfigPartial
/// ---------------------------------------------------------------------------

func Test_ConfigPartial_Has(t *testing.T) {
	t.Run("check if a valid path exists", func(t *testing.T) {
		scenarios := []struct {
			partial ConfigPartial
			search  string
		}{
			{ // test empty partial, search for everything
				partial: ConfigPartial{},
				search:  "",
			},
			{ // test single node, search for root node
				partial: ConfigPartial{"node": "value"},
				search:  "",
			},
			{ // test single node search
				partial: ConfigPartial{"node": "value"},
				search:  "node",
			},
			{ // test multiple node, search for root node
				partial: ConfigPartial{"node1": "value", "node2": "value"},
				search:  "",
			},
			{ // test multiple node search for first
				partial: ConfigPartial{"node1": "value", "node2": "value"},
				search:  "node1",
			},
			{ // test multiple node search for non-first
				partial: ConfigPartial{"node1": "value", "node2": "value"},
				search:  "node2",
			},
			{ // test tree, search for root node
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				search:  "",
			},
			{ // test tree, search for root level node
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				search:  "node1",
			},
			{ // test tree, search for sub node
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				search:  "node1.node2",
			},
		}

		for _, scn := range scenarios {
			if result := scn.partial.Has(scn.search); !result {
				t.Errorf("didn't found the (%s) path in (%v)", scn.search, scn.partial)
			}
		}
	})

	t.Run("check if a invalid path do not exists", func(t *testing.T) {
		scenarios := []struct {
			partial ConfigPartial
			search  string
		}{
			{ // test single node search (invalid)
				partial: ConfigPartial{"node": "value"},
				search:  "node2",
			},
			{ // test multiple node search for invalid node
				partial: ConfigPartial{"node1": "value", "node2": "value"},
				search:  "node3",
			},
			{ // test tree search for invalid root node
				partial: ConfigPartial{"node": ConfigPartial{"node": "value"}},
				search:  "node1",
			},
			{ // test tree search for invalid sub node
				partial: ConfigPartial{"node": ConfigPartial{"node": "value"}},
				search:  "node.node1",
			},
			{ // test tree search for invalid sub-sub-node
				partial: ConfigPartial{"node": ConfigPartial{"node": "value"}},
				search:  "node.node.node",
			},
		}

		for _, scn := range scenarios {
			if result := scn.partial.Has(scn.search); result {
				t.Errorf("founded the (%s) path in (%v)", scn.search, scn.partial)
			}
		}
	})
}

func Test_ConfigPartial_Get(t *testing.T) {
	t.Run("retrieve a value of a existent path", func(t *testing.T) {
		scenarios := []struct {
			partial  ConfigPartial
			search   string
			expected interface{}
		}{
			{ // test empty partial, search for everything
				partial:  ConfigPartial{},
				search:   "",
				expected: ConfigPartial{},
			},
			{ // test single node, search for root node
				partial:  ConfigPartial{"node": "value"},
				search:   "",
				expected: ConfigPartial{"node": "value"},
			},
			{ // test single node search
				partial:  ConfigPartial{"node": "value"},
				search:   "node",
				expected: "value",
			},
			{ // test multiple node, search for root node
				partial:  ConfigPartial{"node1": "value1", "node2": "value2"},
				search:   "",
				expected: ConfigPartial{"node1": "value1", "node2": "value2"},
			},
			{ // test multiple node search for first
				partial:  ConfigPartial{"node1": "value1", "node2": "value2"},
				search:   "node1",
				expected: "value1",
			},
			{ // test multiple node search for non-first
				partial:  ConfigPartial{"node1": "value1", "node2": "value2"},
				expected: "value2",
				search:   "node2",
			},
			{ // test tree, search for root node
				partial:  ConfigPartial{"node": ConfigPartial{"node": "value"}},
				search:   "",
				expected: ConfigPartial{"node": ConfigPartial{"node": "value"}},
			},
			{ // test tree, search for root level node
				partial:  ConfigPartial{"node": ConfigPartial{"node": "value"}},
				search:   "node",
				expected: ConfigPartial{"node": "value"},
			},
			{ // test tree, search for sub node
				partial:  ConfigPartial{"node": ConfigPartial{"node": "value"}},
				search:   "node.node",
				expected: "value",
			},
		}

		for _, scn := range scenarios {
			result := scn.partial.Get(scn.search)
			if !reflect.DeepEqual(result, scn.expected) {
				t.Errorf("returned (%v) when retrieving (%v), expected (%v)", result, scn.search, scn.expected)
			}
		}
	})

	t.Run("return nil if a path don't exists", func(t *testing.T) {
		scenarios := []struct {
			partial ConfigPartial
			search  string
		}{
			{ // test empty partial search for non-existent node
				partial: ConfigPartial{},
				search:  "node",
			},
			{ // test single node search for non-existent node
				partial: ConfigPartial{"node": "value"},
				search:  "node2",
			},
			{ // test multiple node search for non-existent node
				partial: ConfigPartial{"node1": "value1", "node2": "value2"},
				search:  "node3",
			},
			{ // test tree search for non-existent root node
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				search:  "node2",
			},
			{ // test tree search for non-existent sub node
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				search:  "node1.node1",
			},
			{ // test tree search for non-existent sub-sub-node
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				search:  "node1.node2.node3",
			},
		}

		for _, scn := range scenarios {
			if result := scn.partial.Get(scn.search); result != nil {
				t.Errorf("returned (%v) when retrieving (%v)", result, scn.search)
			}
		}
	})

	t.Run("return nil if the node actually stores nil", func(t *testing.T) {
		p := ConfigPartial{"node1": nil, "node2": "value2"}

		if value := p.Get("node1", "default_value"); value != nil {
			t.Errorf("returned the (%v) value", value)
		}
	})

	t.Run("return the default value if a path don't exists", func(t *testing.T) {
		scenarios := []struct {
			partial ConfigPartial
			search  string
		}{
			{ // test empty partial search for non-existent node
				partial: ConfigPartial{},
				search:  "node",
			},
			{ // test single node search for non-existent node
				partial: ConfigPartial{"node": "value"},
				search:  "node2",
			},
			{ // test multiple node search for non-existent node
				partial: ConfigPartial{"node1": "value1", "node2": "value2"},
				search:  "node3",
			},
			{ // test tree search for non-existent root node
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				search:  "node2",
			},
			{ // test tree search for non-existent sub node
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				search:  "node1.node1",
			},
			{ // test tree search for non-existent sub-sub-node
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				search:  "node1.node2.node3",
			},
		}

		defValue := "default_value"
		for _, scn := range scenarios {
			if result := scn.partial.Get(scn.search, defValue); result != defValue {
				t.Errorf("returned (%v) when retrieving (%v)", result, scn.search)
			}
		}
	})
}

func Test_ConfigPartial_Int(t *testing.T) {
	t.Run("panic on a invalid path or a non-integer value", func(t *testing.T) {
		scenarios := []struct {
			partial ConfigPartial
			path    string
		}{
			{ // test when the path doesn't exists
				partial: ConfigPartial{},
				path:    "node1",
			},
			{ // test when the path is storing a nil value
				partial: ConfigPartial{"node1": nil},
				path:    "node1",
			},
			{ // test when the path is storing a string value
				partial: ConfigPartial{"node1": "value1"},
				path:    "node1",
			},
			{ // test when the path is storing a object value
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value1"}},
				path:    "node1",
			},
		}

		for _, scn := range scenarios {
			test := func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("didn't panic")
					} else {
						switch e := r.(type) {
						case error:
							if strings.Index(e.Error(), "interface conversion") != 0 {
								t.Errorf("panic with the (%v) error", e)
							}
						default:
							t.Error("didn't panic with an error")
						}
					}
				}()
				scn.partial.Int(scn.path)
			}
			test()
		}
	})

	value := 101
	p := ConfigPartial{"node1": ConfigPartial{"node2": value}}

	t.Run("retrieve a integer value", func(t *testing.T) {
		if result := p.Int("node1.node2"); result != value {
			t.Errorf("returned the (%v) value", result)
		}
	})

	t.Run("return the given default value if invalid path", func(t *testing.T) {
		defValue := 3

		if result := p.Int("node3", defValue); result != defValue {
			t.Errorf("returned the (%v) value", result)
		}
	})
}

func Test_ConfigPartial_String(t *testing.T) {
	t.Run("panic on a invalid path or a non-string value", func(t *testing.T) {
		scenarios := []struct {
			partial ConfigPartial
			path    string
		}{
			{ // test when the path doesn't exists
				partial: ConfigPartial{},
				path:    "node1",
			},
			{ // test when the path is storing a nil value
				partial: ConfigPartial{"node1": nil},
				path:    "node1",
			},
			{ // test when the path is storing a int value
				partial: ConfigPartial{"node1": 101},
				path:    "node1",
			},
			{ // test when the path is storing an object value
				partial: ConfigPartial{"node1": ConfigPartial{"node2": "value1"}},
				path:    "node1",
			},
		}

		for _, scn := range scenarios {
			test := func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("didn't panic")
					} else {
						switch e := r.(type) {
						case error:
							if strings.Index(e.Error(), "interface conversion") != 0 {
								t.Errorf("panic with the (%v) error", e)
							}
						default:
							t.Error("didn't panic with an error")
						}
					}
				}()
				scn.partial.String(scn.path)
			}
			test()
		}
	})

	value := "value1"
	p := ConfigPartial{"node1": ConfigPartial{"node2": value}}

	t.Run("retrieve a string value", func(t *testing.T) {
		if result := p.String("node1.node2"); result != value {
			t.Errorf("returned the (%v) value", result)
		}
	})

	t.Run("return the given default value if invalid path", func(t *testing.T) {
		defValue := "value"

		if result := p.String("node3", defValue); result != defValue {
			t.Errorf("returned the (%v) value", result)
		}
	})
}

func Test_ConfigPartial_Config(t *testing.T) {
	t.Run("panic on a invalid path or a non-partial value", func(t *testing.T) {
		scenarios := []struct {
			partial ConfigPartial
			path    string
		}{
			{ // test when the path doesn't exists
				partial: ConfigPartial{},
				path:    "node1",
			},
			{ // test when the path is storing a nil value
				partial: ConfigPartial{"node1": nil},
				path:    "node1",
			},
			{ // test when the path is storing a int value
				partial: ConfigPartial{"node1": 101},
				path:    "node1",
			},
			{ // test when the path is storing a string value
				partial: ConfigPartial{"node1": "value1"},
				path:    "node1",
			},
		}

		for _, scn := range scenarios {
			test := func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("didn't panic")
					} else {
						switch e := r.(type) {
						case error:
							if strings.Index(e.Error(), "interface conversion") != 0 {
								t.Errorf("panic with the (%v) error", e)
							}
						default:
							t.Error("didn't panic with an error")
						}
					}
				}()
				scn.partial.Config(scn.path)
			}
			test()
		}
	})

	value := ConfigPartial{"node2": "value1"}
	p := ConfigPartial{"node1": value}

	t.Run("retrieve a config partial", func(t *testing.T) {
		result := p.Config("node1")
		if !reflect.DeepEqual(result, value) {
			t.Errorf("returned the (%v) value", result)
		}
	})

	t.Run("return the given default value if invalid path", func(t *testing.T) {
		defValue := ConfigPartial{"node3": 345}

		result := p.Config("node3", defValue)
		if !reflect.DeepEqual(result, defValue) {
			t.Errorf("returned the (%v) value", result)
		}
	})
}

func Test_ConfigPartial_Merge(t *testing.T) {
	t.Run("merges two partials", func(t *testing.T) {
		scenarios := []struct {
			partial1 ConfigPartial
			partial2 ConfigPartial
			expected ConfigPartial
		}{
			{ // test merging nil partial source
				partial1: ConfigPartial{},
				partial2: nil,
				expected: ConfigPartial{},
			},
			{ // test merging empty partial
				partial1: ConfigPartial{},
				partial2: ConfigPartial{},
				expected: ConfigPartial{},
			},
			{ // test merging empty partial on non empty partial
				partial1: ConfigPartial{"node1": "value1"},
				partial2: ConfigPartial{},
				expected: ConfigPartial{"node1": "value1"},
			},
			{ // test merging partial into empty partial
				partial1: ConfigPartial{},
				partial2: ConfigPartial{"node1": "value1"},
				expected: ConfigPartial{"node1": "value1"},
			},
			{ // test merging override source value
				partial1: ConfigPartial{"node1": "value1"},
				partial2: ConfigPartial{"node1": "value2"},
				expected: ConfigPartial{"node1": "value2"},
			},
			{ // test merging does not override non-present value in merged partial (create)
				partial1: ConfigPartial{"node1": "value1"},
				partial2: ConfigPartial{"node2": "value2"},
				expected: ConfigPartial{"node1": "value1", "node2": "value2"},
			},
			{ // test merging does not override non-present value in merged partial (override)
				partial1: ConfigPartial{"node1": "value1", "node2": "value2"},
				partial2: ConfigPartial{"node2": "value3"},
				expected: ConfigPartial{"node1": "value1", "node2": "value3"},
			},
			{ // test merging override source value to a subtree
				partial1: ConfigPartial{"node1": "value1"},
				partial2: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				expected: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
			},
			{ // test merging override source value in a subtree (to a value)
				partial1: ConfigPartial{"node1": ConfigPartial{"node2": "value1"}},
				partial2: ConfigPartial{"node1": ConfigPartial{"node2": "value2"}},
				expected: ConfigPartial{"node1": ConfigPartial{"node2": "value2"}},
			},
			{ // test merging override source value in a subtree (to a subtree)
				partial1: ConfigPartial{"node1": ConfigPartial{"node2": "value"}},
				partial2: ConfigPartial{"node1": ConfigPartial{"node2": ConfigPartial{"node3": "value"}}},
				expected: ConfigPartial{"node1": ConfigPartial{"node2": ConfigPartial{"node3": "value"}}},
			},
		}

		for _, scn := range scenarios {
			result := scn.partial1.merge(scn.partial2)
			if !reflect.DeepEqual(result, scn.expected) {
				t.Errorf("resulted in (%s) when merging (%v) and (%v), expecting (%v)", result, scn.partial1, scn.partial2, scn.expected)
			}
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigDecoderFactory
/// ---------------------------------------------------------------------------

func Test_NewConfigDecoderFactory(t *testing.T) {
	t.Run("create a new config decoder factory", func(t *testing.T) {
		if factory := NewConfigDecoderFactory(); factory == nil {
			t.Error("didn't returned a valid reference")
		} else if factory.strategies == nil {
			t.Errorf("didn't instantiated the strategies storing array")
		}
	})
}

func Test_ConfigDecoderFactory_Register(t *testing.T) {
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

		var factory *ConfigDecoderFactory
		_ = factory.Register(nil)
	})

	factory := NewConfigDecoderFactory()

	t.Run("nil strategy", func(t *testing.T) {
		if err := factory.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'strategy' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the decoder factory strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockConfigDecoderFactoryStrategy(ctrl)

		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.strategies[0] != strategy {
			t.Error("didn't stored the strategy")
		}
	})
}

func Test_ConfigDecoderFactory_Create(t *testing.T) {
	format := "format"

	t.Run("error if the format is unrecognized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigDecoderFactory()

		reader := NewMockReader(ctrl)
		strategy := NewMockConfigDecoderFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format, reader).Return(false).Times(1)
		_ = factory.Register(strategy)

		if result, err := factory.Create(format, reader); result != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized format type : format" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("should create the requested yaml config decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigDecoderFactory()

		reader := NewMockReader(ctrl)
		decoder := NewMockConfigDecoder(ctrl)
		strategy := NewMockConfigDecoderFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format, reader).Return(true).Times(1)
		strategy.EXPECT().Create(reader).Return(decoder, nil).Times(1)
		_ = factory.Register(strategy)

		if check, err := factory.Create(format, reader); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(check, decoder) {
			t.Error("didn't returned the created strategy")
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigYamlDecoder
/// ---------------------------------------------------------------------------

func Test_NewConfigYamlDecoder(t *testing.T) {
	t.Run("nil reader", func(t *testing.T) {
		if decoder, err := NewConfigYamlDecoder(nil); decoder != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'reader' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new yaml decoder adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl)
		reader.EXPECT().Close().Times(1)

		if decoder, err := NewConfigYamlDecoder(reader); decoder == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			defer decoder.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if decoder.reader != reader {
				t.Error("didn't store the reader reference")
			}
		}
	})
}

func Test_ConfigYamlDecoder_Close(t *testing.T) {
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

		var decoder *ConfigYamlDecoder
		decoder.Close()
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockReader(ctrl)
	reader.EXPECT().Close().Times(1)
	decoder, _ := NewConfigYamlDecoder(reader)

	t.Run("call close method on reader only once", func(t *testing.T) {
		decoder.Close()
		decoder.Close()
	})
}

func Test_ConfigYamlDecoder_Decode(t *testing.T) {
	value := ConfigPartial{"node": "value"}
	expectedError := "error"

	t.Run("return decode error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl)
		reader.EXPECT().Close().Times(1)
		decoder, _ := NewConfigYamlDecoder(reader)
		defer decoder.Close()

		underlyingDecoder := NewMockunderlyingConfigYamlDecoder(ctrl)
		underlyingDecoder.EXPECT().Decode(&ConfigPartial{}).DoAndReturn(func(p interface{}) error {
			return fmt.Errorf(expectedError)
		}).Times(1)
		decoder.decoder = underlyingDecoder

		if result, err := decoder.Decode(); result != nil {
			t.Error("returned an reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("redirect to the underlying decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl)
		reader.EXPECT().Close().Times(1)
		decoder, _ := NewConfigYamlDecoder(reader)
		defer decoder.Close()

		underlyingDecoder := NewMockunderlyingConfigYamlDecoder(ctrl)
		underlyingDecoder.EXPECT().Decode(&ConfigPartial{}).DoAndReturn(func(p interface{}) error {
			p = p.(*ConfigPartial).merge(value)
			return nil
		}).Times(1)
		decoder.decoder = underlyingDecoder

		if result, err := decoder.Decode(); result == nil {
			t.Error("returned a nil value")
		} else if !reflect.DeepEqual(result, value) {
			t.Errorf("returned (%v)", result)
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigYamlDecoderFactoryStrategy
/// ---------------------------------------------------------------------------

func Test_NewConfigYamlDecoderFactoryStrategy(t *testing.T) {
	t.Run("new strategy", func(t *testing.T) {
		if strategy := NewConfigYamlDecoderFactoryStrategy(); strategy == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_ConfigYamlDecoderFactoryStrategy_Accept(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockReader(ctrl)
	strategy := NewConfigYamlDecoderFactoryStrategy()

	t.Run("accept only yaml format", func(t *testing.T) {
		scenarios := []struct {
			format   string
			expected bool
		}{
			{ // test yaml format
				format:   ConfigDecoderFormatYAML,
				expected: true,
			},
			{ // test non-yaml format (json)
				format:   "json",
				expected: false,
			},
		}

		for _, scn := range scenarios {
			if check := strategy.Accept(scn.format, reader); check != scn.expected {
				t.Errorf("returned (%v) when checking (%s) format", check, scn.format)
			}
		}
	})

	t.Run("no extra arguments", func(t *testing.T) {
		if strategy.Accept(ConfigDecoderFormatYAML) {
			t.Error("returned true")
		}
	})

	t.Run("first extra argument is not a io.Reader interface", func(t *testing.T) {
		if strategy.Accept(ConfigDecoderFormatYAML, "string") {
			t.Error("returned true")
		}
	})
}

func Test_ConfigYamlDecoderFactoryStrategy_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockReader(ctrl)
	strategy := NewConfigYamlDecoderFactoryStrategy()

	t.Run("create the decoder", func(t *testing.T) {
		if decoder, err := strategy.Create(reader); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if decoder == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch decoder.(type) {
			case *ConfigYamlDecoder:
			default:
				t.Error("didn't returned a YAML decoder")
			}
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigSource
/// ---------------------------------------------------------------------------

func Test_ConfigBaseSource_Close(t *testing.T) {
	s := &ConfigBaseSource{&sync.Mutex{}, ConfigPartial{}}
	s.Close()
}

func Test_ConfigBaseSource_Has(t *testing.T) {
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

		var source *ConfigBaseSource
		_ = source.Has("path")
	})

	search := "path"
	expected := true
	partial := ConfigPartial{search: "value"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mutex := NewMockLocker(ctrl)
	mutex.EXPECT().Lock().Times(1)
	mutex.EXPECT().Unlock().Times(1)

	s := &ConfigBaseSource{mutex: mutex, partial: partial}

	t.Run("lock and redirect to the stored partial", func(t *testing.T) {
		if value := s.Has(search); value != expected {
			t.Errorf("returned the (%v) value", value)
		}
	})
}

func Test_ConfigBaseSource_Get(t *testing.T) {
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

		var source *ConfigBaseSource
		_ = source.Get("path")
	})

	search := "path"
	expected := "value"
	partial := ConfigPartial{search: expected}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mutex := NewMockLocker(ctrl)
	mutex.EXPECT().Lock().Times(1)
	mutex.EXPECT().Unlock().Times(1)

	s := &ConfigBaseSource{mutex: mutex, partial: partial}

	t.Run("lock and redirect to the stored partial", func(t *testing.T) {
		if value := s.Get(search); value != expected {
			t.Errorf("returned the (%v) value", value)
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigSourceFactory
/// ---------------------------------------------------------------------------

func Test_NewConfigSourceFactory(t *testing.T) {
	t.Run("new config source factory", func(t *testing.T) {
		if factory := NewConfigSourceFactory(); factory == nil {
			t.Error("didn't returned a valid reference")
		} else if factory.strategies == nil {
			t.Error("didn't instantiated the strategies storing array")
		}
	})
}

func Test_ConfigSourceFactory_Register(t *testing.T) {
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

		var factory *ConfigSourceFactory
		_ = factory.Register(nil)
	})

	t.Run("nil strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigSourceFactory()

		if err := factory.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'strategy' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the source factory strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockConfigSourceFactoryStrategy(ctrl)
		factory := NewConfigSourceFactory()

		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.strategies[0] != strategy {
			t.Error("didn't stored the strategy")
		}
	})
}

func Test_ConfigSourceFactory_Create(t *testing.T) {
	sourceType := "type"
	path := "path"
	format := "format"

	t.Run("error on unrecognized format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigSourceFactory()

		strategy := NewMockConfigSourceFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(sourceType, path, format).Return(false).Times(1)
		_ = factory.Register(strategy)

		if source, err := factory.Create(sourceType, path, format); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized source type : type" {
			t.Errorf("returned the (%v) error", err)
		} else if source != nil {
			t.Error("didn't returned the source")
		}
	})

	t.Run("create the requested config source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigSourceFactory()

		source := &ConfigBaseSource{}
		strategy := NewMockConfigSourceFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(sourceType, path, format).Return(true).Times(1)
		strategy.EXPECT().Create(path, format).Return(source, nil).Times(1)
		_ = factory.Register(strategy)

		if check, err := factory.Create(sourceType, path, format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(check, source) {
			t.Error("didn't returned the created source")
		}
	})
}

func Test_ConfigSourceFactory_CreateConfig(t *testing.T) {
	t.Run("error on unrecognized format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigSourceFactory()

		conf := ConfigPartial{}
		strategy := NewMockConfigSourceFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(false).Times(1)
		_ = factory.Register(strategy)

		if source, err := factory.CreateConfig(conf); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != fmt.Sprintf("unrecognized source config : %v", conf) {
			t.Errorf("returned the (%v) error", err)
		} else if source != nil {
			t.Error("returned a valid reference")
		}
	})

	t.Run("create the config source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigSourceFactory()

		conf := ConfigPartial{}
		source := &ConfigBaseSource{}
		strategy := NewMockConfigSourceFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(true).Times(1)
		strategy.EXPECT().CreateConfig(conf).Return(source, nil).Times(1)
		_ = factory.Register(strategy)

		if check, err := factory.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(check, source) {
			t.Error("didn't returned the created source")
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigFileSource
/// ---------------------------------------------------------------------------

func Test_NewConfigFileSource(t *testing.T) {
	t.Run("nil file system adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		decoderFactory := NewConfigDecoderFactory()

		if stream, err := NewConfigFileSource("path", ConfigDecoderFormatYAML, nil, decoderFactory); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil decoder factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		if stream, err := NewConfigFileSource("path", ConfigDecoderFormatYAML, fileSystem, nil); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'decoderFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when opening the file", func(t *testing.T) {
		path := "path"
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)
		decoderFactory := NewConfigDecoderFactory()

		if stream, err := NewConfigFileSource(path, ConfigDecoderFormatYAML, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when creating the decoder", func(t *testing.T) {
		path := "path"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()

		if stream, err := NewConfigFileSource(path, "invalid_format", fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized format type : invalid_format" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when running the decoder", func(t *testing.T) {
		path := "path"
		errorMessage := "error"
		expectedError := fmt.Sprintf("yaml: input error: %s", errorMessage)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			return 0, fmt.Errorf(errorMessage)
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())

		if stream, err := NewConfigFileSource(path, ConfigDecoderFormatYAML, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("creates the config file source", func(t *testing.T) {
		path := "path"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, "field: value")
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())

		if stream, err := NewConfigFileSource(path, ConfigDecoderFormatYAML, fileSystem, decoderFactory); stream == nil {
			t.Error("didn't returned a valid reference")
		} else {
			defer stream.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if stream.mutex == nil {
				t.Error("didn't created the access mutex")
			} else if stream.path != path {
				t.Error("didn't stored the file path")
			} else if stream.format != ConfigDecoderFormatYAML {
				t.Error("didn't stored the file content format")
			} else if stream.fileSystem != fileSystem {
				t.Error("didn't stored the file system adapter reference")
			} else if stream.decoderFactory != decoderFactory {
				t.Error("didn't stored the decoder factory reference")
			}
		}
	})

	t.Run("store the decoded partial", func(t *testing.T) {
		path := "path"
		field := "field"
		value := "value"
		expected := ConfigPartial{field: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())

		stream, _ := NewConfigFileSource(path, ConfigDecoderFormatYAML, fileSystem, decoderFactory)

		if check := stream.partial; !reflect.DeepEqual(check, expected) {
			t.Error("didn't correctly stored the decoded partial")
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigFileSourceFactoryStrategy
/// ---------------------------------------------------------------------------

func Test_NewConfigFileSourceFactoryStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	decoderFactory := NewConfigDecoderFactory()

	t.Run("nil file system adapter", func(t *testing.T) {
		if strategy, err := NewConfigFileSourceFactoryStrategy(nil, decoderFactory); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil decoder factory", func(t *testing.T) {
		if strategy, err := NewConfigFileSourceFactoryStrategy(fileSystem, nil); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'decoderFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new file source factory strategy", func(t *testing.T) {
		if strategy, err := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if strategy == nil {
			t.Error("didn't returned a valid reference")
		} else if strategy.fileSystem != fileSystem {
			t.Error("didn't stored the file system adapter reference")
		} else if strategy.decoderFactory != decoderFactory {
			t.Error("didn't stored the decoder factory reference")
		}
	})
}

func Test_ConfigFileSourceFactoryStrategy_Accept(t *testing.T) {
	sourceType := ConfigSourceTypeFile
	path := "path"
	format := ConfigDecoderFormatYAML

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	decoderFactory := NewConfigDecoderFactory()
	strategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)

	t.Run("don't accept if at least 2 extra arguments are passed", func(t *testing.T) {
		if strategy.Accept(sourceType, path) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the path is not a string", func(t *testing.T) {
		if strategy.Accept(sourceType, 1, format) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the format is not a string", func(t *testing.T) {
		if strategy.Accept(sourceType, path, 1) {
			t.Error("returned true")
		}
	})

	t.Run("accept only file type", func(t *testing.T) {
		scenarios := []struct {
			sourceType string
			expected   bool
		}{
			{ // test file type
				sourceType: ConfigSourceTypeFile,
				expected:   true,
			},
			{ // test non-file type (observable_file)
				sourceType: ConfigSourceTypeObservableFile,
				expected:   false,
			},
		}

		for _, scn := range scenarios {
			if check := strategy.Accept(scn.sourceType, path, format); check != scn.expected {
				t.Errorf("for the type (%s), returned (%v)", scn.sourceType, check)
			}
		}
	})
}

func Test_ConfigFileSourceFactoryStrategy_AcceptConfig(t *testing.T) {
	sourceType := ConfigSourceTypeFile
	path := "path"
	format := ConfigDecoderFormatYAML

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	decoderFactory := NewConfigDecoderFactory()
	strategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)

	t.Run("don't accept if type is missing", func(t *testing.T) {
		partial := ConfigPartial{}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if type is not a string", func(t *testing.T) {
		partial := ConfigPartial{"type": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if path is missing", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if path is not a string", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if format is missing", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if format is not a string", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if invalid type", func(t *testing.T) {
		partial := ConfigPartial{"type": ConfigSourceTypeObservableFile, "path": path, "format": format}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("accept config", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": format}
		if !strategy.AcceptConfig(partial) {
			t.Error("returned false")
		}
	})
}

func Test_ConfigFileSourceFactoryStrategy_Create(t *testing.T) {
	t.Run("non-string path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if source, err := strategy.Create(123, "format"); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-string format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if source, err := strategy.Create("path", 123); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the file source", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		field := "field"
		value := "value"
		expected := ConfigPartial{field: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		strategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if source, err := strategy.Create(path, format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if source == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch s := source.(type) {
			case *ConfigFileSource:
				if !reflect.DeepEqual(s.partial, expected) {
					t.Error("didn't loaded the content correctly")
				}
			default:
				t.Error("didn't returned a new file source")
			}
		}
	})
}

func Test_ConfigFileSourceFactoryStrategy_CreateConfig(t *testing.T) {
	t.Run("non-string path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)

		conf := ConfigPartial{"path": 123, "format": "format"}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-string format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)

		conf := ConfigPartial{"path": "path", "format": 123}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the file source", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		field := "field"
		value := "value"
		expected := ConfigPartial{field: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		strategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)

		conf := ConfigPartial{"path": path, "format": format}

		if source, err := strategy.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if source == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch s := source.(type) {
			case *ConfigFileSource:
				if !reflect.DeepEqual(s.partial, expected) {
					t.Error("didn't loaded the content correctly")
				}
			default:
				t.Error("didn't returned a new file source")
			}
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigObservableFileSource
/// ---------------------------------------------------------------------------

func Test_NewConfigObservableFileSource(t *testing.T) {
	path := "path"
	format := ConfigDecoderFormatYAML

	decoderFactory := NewConfigDecoderFactory()
	_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())

	t.Run("nil file system adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		if stream, err := NewConfigObservableFileSource(path, format, nil, decoderFactory); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil decoder factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		if stream, err := NewConfigObservableFileSource(path, format, fileSystem, nil); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'decoderFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when retrieving the file info", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(nil, fmt.Errorf(expectedError)).Times(1)

		if stream, err := NewConfigObservableFileSource(path, format, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when opening the file", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)

		if stream, err := NewConfigObservableFileSource(path, format, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when creating the decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		if stream, err := NewConfigObservableFileSource(path, "invalid_format", fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized format type : invalid_format" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when running the decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("{"))
			return 1, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		if stream, err := NewConfigObservableFileSource(path, format, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "yaml: line 1: did not find expected node content" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the config observable file source", func(t *testing.T) {
		field := "field"
		value := "value"
		expected := ConfigPartial{field: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		if stream, err := NewConfigObservableFileSource(path, format, fileSystem, decoderFactory); stream == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			defer stream.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if stream.mutex == nil {
				t.Error("didn't created the access mutex")
			} else if stream.path != path {
				t.Error("didn't stored the file path")
			} else if stream.format != ConfigDecoderFormatYAML {
				t.Error("didn't stored the file content format")
			} else if stream.fileSystem != fileSystem {
				t.Error("didn't stored the file system adapter reference")
			} else if stream.decoderFactory != decoderFactory {
				t.Error("didn't stored the decoder factory reference")
			} else if !reflect.DeepEqual(stream.partial, expected) {
				t.Error("didn't loaded the content correctly")
			}
		}
	})

	t.Run("store the decoded partial", func(t *testing.T) {
		field := "field"
		value := "value"
		expected := ConfigPartial{field: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		stream, _ := NewConfigObservableFileSource(path, format, fileSystem, decoderFactory)
		defer stream.Close()

		if check := stream.partial; !reflect.DeepEqual(check, expected) {
			t.Error("didn't correctly stored the decoded partial")
		}
	})
}

func Test_ConfigObservableFileSource_Reload(t *testing.T) {
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

		var source *ConfigObservableFileSource
		_, _ = source.Reload()
	})

	path := "path"
	format := ConfigDecoderFormatYAML

	decoderFactory := NewConfigDecoderFactory()
	_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())

	t.Run("error if fail to retrieving the file info", func(t *testing.T) {
		expectedError := "error"
		field := "field"
		value := "value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		gomock.InOrder(
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
			fileSystem.EXPECT().Stat(path).Return(nil, fmt.Errorf(expectedError)),
		)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		stream, _ := NewConfigObservableFileSource(path, format, fileSystem, decoderFactory)
		defer stream.Close()

		if reloaded, err := stream.Reload(); reloaded {
			t.Error("flagged that was reloaded")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error if fails to load the file content", func(t *testing.T) {
		expectedError := "error"
		field := "field"
		value := "value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		gomock.InOrder(
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)),
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 2)),
		)
		fileSystem := NewMockFs(ctrl)
		gomock.InOrder(
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
		)
		gomock.InOrder(
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil),
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)),
		)

		stream, _ := NewConfigObservableFileSource(path, format, fileSystem, decoderFactory)
		defer stream.Close()

		if reloaded, err := stream.Reload(); reloaded {
			t.Error("flagged that was reloaded")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("prevent reload of a unchanged source", func(t *testing.T) {
		field := "field"
		value := "value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(2)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(2)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		stream, _ := NewConfigObservableFileSource(path, format, fileSystem, decoderFactory)

		if reloaded, err := stream.Reload(); reloaded {
			t.Error("flagged that was reloaded")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("should reload a changed source", func(t *testing.T) {
		field := "field"
		value1 := "value1"
		value2 := "value2"
		expected := ConfigPartial{field: value2}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file1 := NewMockFile(ctrl)
		file1.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value1))
			return 13, io.EOF
		})
		file1.EXPECT().Close().Times(1)
		file2 := NewMockFile(ctrl)
		file2.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value2))
			return 13, io.EOF
		})
		file2.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		gomock.InOrder(
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)),
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 2)),
		)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(2)
		gomock.InOrder(
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file1, nil),
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file2, nil),
		)

		stream, _ := NewConfigObservableFileSource(path, format, fileSystem, decoderFactory)

		if reloaded, err := stream.Reload(); !reloaded {
			t.Error("flagged that was not reloaded")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(expected, stream.partial) {
			t.Error("didn't stored the reloaded configuration")
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigObservableFileSourceFactoryStrategy
/// ---------------------------------------------------------------------------

func Test_NewConfigObservableFileSourceFactoryStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	decoderFactory := NewConfigDecoderFactory()

	t.Run("nil file system adapter", func(t *testing.T) {
		if strategy, err := NewConfigObservableFileSourceFactoryStrategy(nil, decoderFactory); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil decoder factory", func(t *testing.T) {
		if strategy, err := NewConfigObservableFileSourceFactoryStrategy(fileSystem, nil); strategy != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'decoderFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new file source factory strategy", func(t *testing.T) {
		if strategy, err := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if strategy == nil {
			t.Error("didn't returned a valid reference")
		} else if strategy.fileSystem != fileSystem {
			t.Error("didn't stored the file system adapter reference")
		} else if strategy.decoderFactory != decoderFactory {
			t.Error("didn't stored the decoder factory reference")
		}
	})
}

func Test_ConfigObservableFileSourceFactoryStrategy_Accept(t *testing.T) {
	sourceType := ConfigSourceTypeObservableFile
	path := "path"
	format := ConfigDecoderFormatYAML

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	decoderFactory := NewConfigDecoderFactory()
	strategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

	t.Run("don't accept if at least 2 extra arguments are passed", func(t *testing.T) {
		if strategy.Accept(sourceType, path) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the path is not a string", func(t *testing.T) {
		if strategy.Accept(sourceType, 1, format) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if the format is not a string", func(t *testing.T) {
		if strategy.Accept(sourceType, path, 1) {
			t.Error("returned true")
		}
	})

	t.Run("accept only file type", func(t *testing.T) {
		scenarios := []struct {
			sourceType string
			expected   bool
		}{
			{ // test file type
				sourceType: ConfigSourceTypeObservableFile,
				expected:   true,
			},
			{ // test non-file type (file)
				sourceType: ConfigSourceTypeFile,
				expected:   false,
			},
		}

		for _, scn := range scenarios {
			if check := strategy.Accept(scn.sourceType, path, format); check != scn.expected {
				t.Errorf("for the type (%s), returned (%v)", scn.sourceType, check)
			}
		}
	})
}

func Test_ConfigObservableFileSourceFactoryStrategy_AcceptConfig(t *testing.T) {
	sourceType := ConfigSourceTypeObservableFile
	path := "path"
	format := ConfigDecoderFormatYAML

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	decoderFactory := NewConfigDecoderFactory()
	strategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

	t.Run("don't accept if type is missing", func(t *testing.T) {
		partial := ConfigPartial{}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if type is not a string", func(t *testing.T) {
		partial := ConfigPartial{"type": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if path is missing", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if path is not a string", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if format is missing", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if format is not a string", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": 123}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("don't accept if invalid type", func(t *testing.T) {
		partial := ConfigPartial{"type": ConfigSourceTypeFile, "path": path, "format": format}
		if strategy.AcceptConfig(partial) {
			t.Error("returned true")
		}
	})

	t.Run("accept config", func(t *testing.T) {
		partial := ConfigPartial{"type": sourceType, "path": path, "format": format}
		if !strategy.AcceptConfig(partial) {
			t.Error("returned false")
		}
	})
}

func Test_ConfigObservableFileSourceFactoryStrategy_Create(t *testing.T) {
	t.Run("non-string path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if source, err := strategy.Create(123, "format"); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-string format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if source, err := strategy.Create("path", 123); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the file source", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		field := "field"
		value := "value"
		expected := ConfigPartial{field: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		strategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if source, err := strategy.Create(path, format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if source == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch s := source.(type) {
			case *ConfigObservableFileSource:
				if !reflect.DeepEqual(s.partial, expected) {
					t.Error("didn't loaded the content correctly")
				}
			default:
				t.Error("didn't returned a new file source")
			}
		}
	})
}

func Test_ConfigObservableFileSourceFactoryStrategy_CreateConfig(t *testing.T) {
	t.Run("non-string path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		conf := ConfigPartial{"path": 123, "format": "format"}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("non-string format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewConfigDecoderFactory()
		strategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)

		conf := ConfigPartial{"path": "path", "format": 123}
		if source, err := strategy.CreateConfig(conf); source != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the file source", func(t *testing.T) {
		path := "path"
		format := ConfigDecoderFormatYAML
		field := "field"
		value := "value"
		expected := ConfigPartial{field: value}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, fmt.Sprintf("%s: %s", field, value))
			return 12, io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		strategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		conf := ConfigPartial{"path": path, "format": format}

		if source, err := strategy.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if source == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch s := source.(type) {
			case *ConfigObservableFileSource:
				if !reflect.DeepEqual(s.partial, expected) {
					t.Error("didn't loaded the content correctly")
				}
			default:
				t.Error("didn't returned a new file source")
			}
		}
	})
}

/// ---------------------------------------------------------------------------
/// Config
/// ---------------------------------------------------------------------------

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

	id := "source"
	priority := 0

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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return boolean value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return int value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return int8 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return int16 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return int32 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return int64 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return uint value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return uint8 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return int16 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return uint32 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return uint64 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return float32 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return float64 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return complex64 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return complex128 value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return rune value", func(t *testing.T) {
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

	id := "source"
	priority := 0
	search := "node"

	t.Run("return string value", func(t *testing.T) {
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

	t.Run("validate if the source is registered", func(t *testing.T) {
		if !config.HasSource(id) {
			t.Error("returned false")
		}
	})

	t.Run("invalidate if the source is not registered", func(t *testing.T) {
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

	id := "source"
	priority := 1
	p := ConfigPartial{}

	t.Run("nil source", func(t *testing.T) {
		config, _ := NewConfig(60 * time.Second)
		defer config.Close()

		if err := config.AddSource(id, priority, nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'source' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register a new source", func(t *testing.T) {
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

	node := "node"
	id1 := "source.1"
	value1 := "value1"
	p1 := ConfigPartial{node: value1}

	id2 := "source.2"
	value2 := "value2"
	p2 := ConfigPartial{node: value2}

	t.Run("override path if the insert have higher priority", func(t *testing.T) {
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

	extendedNode := "extendedNode"
	extendedValue := "extraValue"
	extendedPartial := p2.merge(ConfigPartial{extendedNode: extendedValue})

	t.Run("still be able to get not overridden paths of a inserted lower priority", func(t *testing.T) {
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

	t.Run("unregister a previously registered source", func(t *testing.T) {
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

	node := "node"
	id1 := "source.1"
	priority1 := 1
	value1 := "value1"
	partial1 := ConfigPartial{node: value1}

	id2 := "source.2"
	priority2 := 2
	value2 := "value1"
	partial2 := ConfigPartial{node: value2}

	t.Run("update the priority of the source", func(t *testing.T) {
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
	id := "source"
	priority := 0
	node := "node"
	value := "value"
	partial := ConfigPartial{node: value}

	t.Run("reload on observable sources", func(t *testing.T) {
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

	check := false

	t.Run("should call observer callback function on config changes", func(t *testing.T) {
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

/// ---------------------------------------------------------------------------
/// ConfigLoader
/// ---------------------------------------------------------------------------

func Test_NewLoader(t *testing.T) {
	config, _ := NewConfig(0 * time.Second)
	sourceFactory := NewConfigSourceFactory()

	t.Run("nil config", func(t *testing.T) {
		if loader, err := NewConfigLoader(nil, sourceFactory); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'config' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil config source factory", func(t *testing.T) {
		if loader, err := NewConfigLoader(config, nil); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'sourceFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new config loader", func(t *testing.T) {
		if loader, err := NewConfigLoader(config, sourceFactory); loader == nil {
			t.Error("didn't returned a valid reference")
		} else if err != nil {
			t.Errorf("return the (%v) error", err)
		}
	})
}

func Test_Loader_Load(t *testing.T) {
	sourceID := "base_source_id"
	sourcePath := "base_source_path"
	sourceFormat := ConfigDecoderFormatYAML

	t.Run("error getting the base source", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error storing the base source", func(t *testing.T) {
		content := "field: value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(ConfigPartial{})
		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource(sourceID, 0, source)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "duplicate source id : base_source_id" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("add base source into the config", func(t *testing.T) {
		content := "field: value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on invalid list of sources", func(t *testing.T) {
		content := "config:\n  sources: 123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "error while parsing the list of sources" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on loaded invalid id", func(t *testing.T) {
		content := `
config:
  sources:
    - id: 12`

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on loaded invalid priority", func(t *testing.T) {
		content := `
config:
  sources:
    - id: id
      priority: string`

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on loaded source factory", func(t *testing.T) {
		content := `
config:
  sources:
    - id: id
      priority: 0`

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "unrecognized source config") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on source registration", func(t *testing.T) {
		content := `
config:
  sources:
    - id: id
      priority: 0
      type: file
      path: path
      format: yaml`

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(ConfigPartial{}).AnyTimes()
		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("id", 0, source)

		file1 := NewMockFile(ctrl)
		file1.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file1.EXPECT().Close().Times(1)

		file2 := NewMockFile(ctrl)
		file2.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, "field: value")
			return 12, io.EOF
		}).Times(1)
		file2.EXPECT().Close().Times(1)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file1, nil).Times(1)
		fileSystem.EXPECT().OpenFile("path", os.O_RDONLY, os.FileMode(0644)).Return(file2, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "duplicate source id : id" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the loaded source", func(t *testing.T) {
		content := `
config:
  sources:
    - id: id
      priority: 0
      type: file
      path: path
      format: yaml`

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file1 := NewMockFile(ctrl)
		file1.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file1.EXPECT().Close().Times(1)

		file2 := NewMockFile(ctrl)
		file2.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, "field: value")
			return 12, io.EOF
		}).Times(1)
		file2.EXPECT().Close().Times(1)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file1, nil).Times(1)
		fileSystem.EXPECT().OpenFile("path", os.O_RDONLY, os.FileMode(0644)).Return(file2, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigYamlDecoderFactoryStrategy())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigParams
/// ---------------------------------------------------------------------------

func Test_NewConfigParams(t *testing.T) {
	t.Run("new parameters", func(t *testing.T) {
		parameters := NewConfigParams()

		if value := parameters.ConfigID; value != ContainerConfigID {
			t.Errorf("stored (%v) config ID", value)
		} else if value := parameters.FileSystemID; value != ContainerFileSystemID {
			t.Errorf("stored (%v) file sytem ID", value)
		} else if value := parameters.SourceFactoryID; value != ContainerConfigSourceFactoryID {
			t.Errorf("stored (%v) source factory ID", value)
		} else if value := parameters.DecoderFactoryID; value != ContainerConfigDecoderFactoryID {
			t.Errorf("stored (%v) decoder factory ID", value)
		} else if value := parameters.LoaderID; value != ContainerConfigLoaderID {
			t.Errorf("stored (%v) loader ID", value)
		} else if value := parameters.ObserveFrequency; value != ConfigObserveFrequency {
			t.Errorf("stored (%v) observe frequecy", value)
		} else if value := parameters.BaseSourceActive; value != ConfigBaseSourceActive {
			t.Errorf("stored (%v) base source active", value)
		} else if value := parameters.BaseSourceID; value != ConfigBaseSourceID {
			t.Errorf("stored (%v) base source id", value)
		} else if value := parameters.BaseSourcePath; value != ConfigBaseSourcePath {
			t.Errorf("stored (%v) base source path", value)
		} else if value := parameters.BaseSourceFormat; value != ConfigBaseSourceFormat {
			t.Errorf("stored (%v) base source format", value)
		}
	})

	t.Run("with the env config ID", func(t *testing.T) {
		configID := "config_id"
		_ = os.Setenv(EnvContainerConfigID, configID)
		defer func() { _ = os.Setenv(EnvContainerConfigID, "") }()

		parameters := NewConfigParams()
		if value := parameters.ConfigID; value != configID {
			t.Errorf("stored (%v) config ID", value)
		}
	})

	t.Run("with the env file system ID", func(t *testing.T) {
		fileSystemID := "file_system_id"
		_ = os.Setenv(EnvContainerFileSystemID, fileSystemID)
		defer func() { _ = os.Setenv(EnvContainerFileSystemID, "") }()

		parameters := NewConfigParams()
		if value := parameters.FileSystemID; value != fileSystemID {
			t.Errorf("stored (%v) file system ID", value)
		}
	})

	t.Run("with the env source factory ID", func(t *testing.T) {
		sourceFactoryID := "source_factory_id"
		_ = os.Setenv(EnvContainerConfigSourceFactoryID, sourceFactoryID)
		defer func() { _ = os.Setenv(EnvContainerConfigSourceFactoryID, "") }()

		parameters := NewConfigParams()
		if value := parameters.SourceFactoryID; value != sourceFactoryID {
			t.Errorf("stored (%v) source factory ID", value)
		}
	})

	t.Run("with the env decoder factory ID", func(t *testing.T) {
		decoderFactoryID := "decoder_factory_id"
		_ = os.Setenv(EnvContainerConfigDecoderFactoryID, decoderFactoryID)
		defer func() { _ = os.Setenv(EnvContainerConfigDecoderFactoryID, "") }()

		parameters := NewConfigParams()
		if value := parameters.DecoderFactoryID; value != decoderFactoryID {
			t.Errorf("stored (%v) decoder factory ID", value)
		}
	})

	t.Run("with the env loader ID", func(t *testing.T) {
		loaderID := "loader_id"
		_ = os.Setenv(EnvContainerConfigLoaderID, loaderID)
		defer func() { _ = os.Setenv(EnvContainerConfigLoaderID, "") }()

		parameters := NewConfigParams()
		if value := parameters.LoaderID; value != loaderID {
			t.Errorf("stored (%v) loader ID", value)
		}
	})

	t.Run("with the env observer frequency", func(t *testing.T) {
		observeFrequency := time.Second * 10
		_ = os.Setenv(EnvConfigObserveFrequency, strconv.Itoa(int(observeFrequency.Seconds())))
		defer func() { _ = os.Setenv(EnvConfigObserveFrequency, "") }()

		parameters := NewConfigParams()
		if value := parameters.ObserveFrequency; value != observeFrequency {
			t.Errorf("stored (%v) observe frequency", value)
		}
	})

	t.Run("with the env base source active", func(t *testing.T) {
		_ = os.Setenv(EnvConfigBaseSourceActive, fmt.Sprintf("%v", true))
		defer func() { _ = os.Setenv(EnvConfigBaseSourceActive, "") }()

		parameters := NewConfigParams()
		if value := parameters.BaseSourceActive; !value {
			t.Errorf("stored (%v) base source ID", value)
		}
	})

	t.Run("with the env base source ID", func(t *testing.T) {
		baseSourceID := "base_config_id"
		_ = os.Setenv(EnvConfigBaseSourceID, baseSourceID)
		defer func() { _ = os.Setenv(EnvConfigBaseSourceID, "") }()

		parameters := NewConfigParams()
		if value := parameters.BaseSourceID; value != baseSourceID {
			t.Errorf("stored (%v) base source ID", value)
		}
	})

	t.Run("with the env base source path", func(t *testing.T) {
		baseSourcePath := "base_config_path"
		_ = os.Setenv(EnvConfigBaseSourcePath, baseSourcePath)
		defer func() { _ = os.Setenv(EnvConfigBaseSourcePath, "") }()

		parameters := NewConfigParams()
		if value := parameters.BaseSourcePath; value != baseSourcePath {
			t.Errorf("stored (%v) base source path", value)
		}
	})

	t.Run("with the env base source format", func(t *testing.T) {
		baseSourceFormat := "base_config_format"
		_ = os.Setenv(EnvConfigBaseSourceFormat, baseSourceFormat)
		defer func() { _ = os.Setenv(EnvConfigBaseSourceFormat, "") }()

		parameters := NewConfigParams()
		if value := parameters.BaseSourceFormat; value != baseSourceFormat {
			t.Errorf("stored (%v) base source format", value)
		}
	})
}

/// ---------------------------------------------------------------------------
/// ConfigProvider
/// ---------------------------------------------------------------------------

func Test_NewConfigProvider(t *testing.T) {
	t.Run("without params", func(t *testing.T) {
		if provider := NewConfigProvider(nil); provider == nil {
			t.Error("didn't returned a valid reference")
		} else if !reflect.DeepEqual(NewConfigParams(), provider.params) {
			t.Errorf("stored the (%v) parameters", provider.params)
		}
	})

	t.Run("with defined params", func(t *testing.T) {
		params := NewConfigParams()
		if provider := NewConfigProvider(params); provider == nil {
			t.Error("didn't returned a valid reference")
		} else if params != provider.params {
			t.Errorf("stored the (%v) parameters", provider.params)
		}
	})
}

func Test_ConfigProvider_Register(t *testing.T) {
	t.Run("nil container", func(t *testing.T) {
		provider := NewConfigProvider(nil)
		if err := provider.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'container' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register components", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewConfigProvider(nil)

		if err := provider.Register(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !container.Has(ContainerConfigDecoderFactoryID) {
			t.Error("didnt registered the config decoder factory", err)
		} else if !container.Has(ContainerConfigSourceFactoryID) {
			t.Error("didnt registered the config source factory", err)
		} else if !container.Has(ContainerConfigID) {
			t.Error("didnt registered the config", err)
		} else if !container.Has(ContainerConfigLoaderID) {
			t.Error("didnt registered the config loader", err)
		}
	})

	t.Run("retrieving config decoder factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		if decoderFactory, err := container.Get(ContainerConfigDecoderFactoryID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if decoderFactory == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch decoderFactory.(type) {
			case *ConfigDecoderFactory:
			default:
				t.Error("didn't returned a decoder factory reference")
			}
		}
	})

	t.Run("error retrieving file system", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if loader, err := container.Get(ContainerConfigSourceFactoryID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid file system", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if loader, err := container.Get(ContainerConfigSourceFactoryID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving decoder factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigDecoderFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if loader, err := container.Get(ContainerConfigSourceFactoryID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid decoder factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigDecoderFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if loader, err := container.Get(ContainerConfigSourceFactoryID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("retrieving config source factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		if sourceFactory, err := container.Get(ContainerConfigSourceFactoryID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if sourceFactory == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch sourceFactory.(type) {
			case *ConfigSourceFactory:
			default:
				t.Error("didn't returned a source factory reference")
			}
		}
	})

	t.Run("retrieving config", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		if config, err := container.Get(ContainerConfigID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if config == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch config.(type) {
			case *Config:
			default:
				t.Error("didn't returned a config reference")
			}
		}
	})

	t.Run("error retrieving config", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if loader, err := container.Get(ContainerConfigLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if loader, err := container.Get(ContainerConfigLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving source factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigSourceFactoryID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if loader, err := container.Get(ContainerConfigLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid source factory", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		_ = container.Add(ContainerConfigSourceFactoryID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if loader, err := container.Get(ContainerConfigLoaderID); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("retrieving config loader", func(t *testing.T) {
		container := NewAppContainer()
		_ = NewFileSystemProvider(nil).Register(container)
		_ = NewConfigProvider(nil).Register(container)

		if loader, err := container.Get(ContainerConfigLoaderID); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if loader == nil {
			t.Error("didn't returned a valid reference")
		} else {
			switch loader.(type) {
			case *ConfigLoader:
			default:
				t.Error("didn't returned a loader reference")
			}
		}
	})
}

func Test_ConfigProvider_Boot(t *testing.T) {
	t.Run("no active flag", func(t *testing.T) {
		container := NewAppContainer()

		params := NewConfigParams()
		params.BaseSourceActive = false
		provider := NewConfigProvider(params)

		if err := provider.Boot(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("flag active but loader not registered", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewConfigProvider(nil)

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "entry 'servlet.config.loader' not registered in the container" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error retrieving loader", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = provider.Register(container)
		_ = container.Add(ContainerConfigLoaderID, func(*AppContainer) (interface{}, error) {
			return nil, fmt.Errorf("error")
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if err.Error() != "error" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("invalid loader", func(t *testing.T) {
		container := NewAppContainer()
		provider := NewConfigProvider(nil)
		_ = provider.Register(container)
		_ = container.Add(ContainerConfigLoaderID, func(*AppContainer) (interface{}, error) {
			return "string", nil
		})

		if err := provider.Boot(container); err == nil {
			t.Error("didn't returned the expected error", err)
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("add base source into the config", func(t *testing.T) {
		content := "field: value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(ConfigBaseSourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		container := NewAppContainer()
		_ = container.Add(ContainerFileSystemID, func(*AppContainer) (interface{}, error) {
			return fileSystem, nil
		})

		provider := NewConfigProvider(nil)
		_ = provider.Register(container)

		if err := provider.Boot(container); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}
