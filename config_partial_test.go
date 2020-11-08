package servlet

import (
	"reflect"
	"strings"
	"testing"
)

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
