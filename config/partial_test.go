package config

import (
	"reflect"
	"testing"
)

func Test_Partial_Has(t *testing.T) {
	t.Run("should correctly check if a valid path exists", func(t *testing.T) {
		action := "Checking the existence of a valid path in the partial"

		scenarios := []struct {
			partial partial
			search  string
		}{
			{ // test empty partial, search for everything
				partial: partial{},
				search:  "",
			},
			{ // test single node, search for root node
				partial: partial{"node": "value"},
				search:  "",
			},
			{ // test single node search
				partial: partial{"node": "value"},
				search:  "node",
			},
			{ // test multiple node, search for root node
				partial: partial{"node1": "value", "node2": "value"},
				search:  "",
			},
			{ // test multiple node search for first
				partial: partial{"node1": "value", "node2": "value"},
				search:  "node1",
			},
			{ // test multiple node search for non-first
				partial: partial{"node1": "value", "node2": "value"},
				search:  "node2",
			},
			{ // test tree, search for root node
				partial: partial{"node1": partial{"node2": "value"}},
				search:  "",
			},
			{ // test tree, search for root level node
				partial: partial{"node1": partial{"node2": "value"}},
				search:  "node1",
			},
			{ // test tree, search for subnode
				partial: partial{"node1": partial{"node2": "value"}},
				search:  "node1.node2",
			},
		}

		for _, scn := range scenarios {
			if result := scn.partial.Has(scn.search); !result {
				t.Errorf("%s didn't found the (%s) path in (%v)", action, scn.search, scn.partial)
			}
		}
	})

	t.Run("should correctly check if a invalid path do not exists", func(t *testing.T) {
		action := "Checking the existence of a valid path in the partial"

		scenarios := []struct {
			partial partial
			search  string
		}{
			{ // test single node search (invalid)
				partial: partial{"node": "value"},
				search:  "node2",
			},
			{ // test multiple node search for invalid node
				partial: partial{"node1": "value", "node2": "value"},
				search:  "node3",
			},
			{ // test tree search for invalid root node
				partial: partial{"node": partial{"node": "value"}},
				search:  "node1",
			},
			{ // test tree search for invalid subnode
				partial: partial{"node": partial{"node": "value"}},
				search:  "node.node1",
			},
			{ // test tree search for invalid sub-sub-node
				partial: partial{"node": partial{"node": "value"}},
				search:  "node.node.node",
			},
		}

		for _, scn := range scenarios {
			if result := scn.partial.Has(scn.search); result {
				t.Errorf("%s unexpectedly found the (%s) path in (%v)", action, scn.search, scn.partial)
			}
		}
	})
}

func Test_Partial_Get(t *testing.T) {
	t.Run("should correctly retrieve a value of a existent path", func(t *testing.T) {
		action := "Retrieving a valid path from the partial"

		scenarios := []struct {
			partial  partial
			search   string
			expected interface{}
		}{
			{ // test empty partial, search for everything
				partial:  partial{},
				search:   "",
				expected: partial{},
			},
			{ // test single node, search for root node
				partial:  partial{"node": "value"},
				search:   "",
				expected: partial{"node": "value"},
			},
			{ // test single node search
				partial:  partial{"node": "value"},
				search:   "node",
				expected: "value",
			},
			{ // test multiple node, search for root node
				partial:  partial{"node1": "value1", "node2": "value2"},
				search:   "",
				expected: partial{"node1": "value1", "node2": "value2"},
			},
			{ // test multiple node search for first
				partial:  partial{"node1": "value1", "node2": "value2"},
				search:   "node1",
				expected: "value1",
			},
			{ // test multiple node search for non-first
				partial:  partial{"node1": "value1", "node2": "value2"},
				expected: "value2",
				search:   "node2",
			},
			{ // test tree, search for root node
				partial:  partial{"node": partial{"node": "value"}},
				search:   "",
				expected: partial{"node": partial{"node": "value"}},
			},
			{ // test tree, search for root level node
				partial:  partial{"node": partial{"node": "value"}},
				search:   "node",
				expected: partial{"node": "value"},
			},
			{ // test tree, search for subnode
				partial:  partial{"node": partial{"node": "value"}},
				search:   "node.node",
				expected: "value",
			},
		}

		for _, scn := range scenarios {
			result := scn.partial.Get(scn.search)
			if !reflect.DeepEqual(result, scn.expected) {
				t.Errorf("%s resulted (%v) when retrieving (%v), expected (%v)", action, result, scn.search, scn.expected)
			}
		}
	})

	t.Run("should correctly return nil if a path don't exists", func(t *testing.T) {
		action := "Retrieving a non-existing path from the partial"

		scenarios := []struct {
			partial partial
			search  string
		}{
			{ // test empty partial search for non-existent node
				partial: partial{},
				search:  "node",
			},
			{ // test single node search for non-existent node
				partial: partial{"node": "value"},
				search:  "node2",
			},
			{ // test multiple node search for non-existent node
				partial: partial{"node1": "value1", "node2": "value2"},
				search:  "node3",
			},
			{ // test tree search for non-existent root node
				partial: partial{"node1": partial{"node2": "value"}},
				search:  "node2",
			},
			{ // test tree search for non-existent subnode
				partial: partial{"node1": partial{"node2": "value"}},
				search:  "node1.node1",
			},
			{ // test tree search for non-existent sub-sub-node
				partial: partial{"node1": partial{"node2": "value"}},
				search:  "node1.node2.node3",
			},
		}

		for _, scn := range scenarios {
			result := scn.partial.Get(scn.search)
			if result != nil {
				t.Errorf("%s returned (%v) when retrieving (%v), expecting nil", action, result, scn.search)
			}
		}
	})

	t.Run("should return nil if the node actually stores nil", func(t *testing.T) {
		action := "Retrieving a stored nil value"

		p := partial{"node1": nil}
		path := "node1"
		var expectedValue interface{} = nil

		value := p.Get(path, "__default_value__")
		if value != expectedValue {
			t.Errorf("%s returned the (%v) value, expected (%v)", action, value, expectedValue)
		}
	})

	t.Run("should correctly return the default value if a path don't exists", func(t *testing.T) {
		action := "Retrieving a non-existing path from the partial with a default value"

		scenarios := []struct {
			partial       partial
			search        string
			expectedValue string
		}{
			{ // test empty partial search for non-existent node
				partial:       partial{},
				search:        "node",
				expectedValue: "__default__",
			},
			{ // test single node search for non-existent node
				partial:       partial{"node": "value"},
				search:        "node2",
				expectedValue: "__default__",
			},
			{ // test multiple node search for non-existent node
				partial:       partial{"node1": "value1", "node2": "value2"},
				search:        "node3",
				expectedValue: "__default__",
			},
			{ // test tree search for non-existent root node
				partial:       partial{"node1": partial{"node2": "value"}},
				search:        "node2",
				expectedValue: "__default__",
			},
			{ // test tree search for non-existent subnode
				partial:       partial{"node1": partial{"node2": "value"}},
				search:        "node1.node1",
				expectedValue: "__default__",
			},
			{ // test tree search for non-existent sub-sub-node
				partial:       partial{"node1": partial{"node2": "value"}},
				search:        "node1.node2.node3",
				expectedValue: "__default__",
			},
		}

		for _, scn := range scenarios {
			result := scn.partial.Get(scn.search, scn.expectedValue)
			if result != scn.expectedValue {
				t.Errorf("%s returned (%v) when retrieving (%v), expecting %s", action, result, scn.search, scn.expectedValue)
			}
		}
	})
}

func Test_Partial_Int(t *testing.T) {
	t.Run("should correctly validate a invalid request", func(t *testing.T) {
		action := "Retrieving the path value as a int with an error"

		scenarios := []struct {
			partial partial
			path    string
		}{
			{ // test when the path dosen't exists
				partial: partial{},
				path:    "node1",
			},
			{ // test when the path is storing anil value
				partial: partial{"node1": nil},
				path:    "node1",
			},
			{ // test when the path is storing a string value
				partial: partial{"node1": "value1"},
				path:    "node1",
			},
			{ // test when the path is storing an object value
				partial: partial{"node1": partial{"node2": "value1"}},
				path:    "node1",
			},
		}

		for _, scn := range scenarios {
			test := func() {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("%s didn't panic", action)
					}
				}()
				scn.partial.Int(scn.path)
			}
			test()
		}
	})

	t.Run("should correctly retrieve a integer value", func(t *testing.T) {
		action := "Retrieving the path value as a int"

		p := partial{"node1": partial{"node2": 101}}
		path := "node1.node2"
		expectedValue := 101

		value := p.Int(path)
		if value != expectedValue {
			t.Errorf("%s returned the (%v) value, expected (%v)", action, value, expectedValue)
		}
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		p := partial{"node1": partial{"node2": 101}}
		path := "node3"
		expectedValue := 3

		value := p.Int(path, expectedValue)
		if value != expectedValue {
			t.Errorf("%s returned the (%v) value, expected (%v)", action, value, expectedValue)
		}
	})
}

func Test_Partial_String(t *testing.T) {
	t.Run("should correctly validate a invalid request", func(t *testing.T) {
		action := "Retrieving the path value as a string with an error"

		scenarios := []struct {
			partial partial
			path    string
		}{
			{ // test when the path dosen't exists
				partial: partial{},
				path:    "node1",
			},
			{ // test when the path is storing anil value
				partial: partial{"node1": nil},
				path:    "node1",
			},
			{ // test when the path is storing a int value
				partial: partial{"node1": 101},
				path:    "node1",
			},
			{ // test when the path is storing an object value
				partial: partial{"node1": partial{"node2": "value1"}},
				path:    "node1",
			},
		}

		for _, scn := range scenarios {
			test := func() {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("%s didn't panic", action)
					}
				}()
				scn.partial.String(scn.path)
			}
			test()
		}
	})

	t.Run("should correctly retrieve a string value", func(t *testing.T) {
		action := "Retrieving the path value as a string"

		p := partial{"node1": partial{"node2": "value1"}}
		path := "node1.node2"
		expectedValue := "value1"

		value := p.String(path)
		if value != expectedValue {
			t.Errorf("%s returned the (%v) value, expected (%v)", action, value, expectedValue)
		}
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		p := partial{"node1": partial{"node2": 101}}
		path := "node3"
		expectedValue := "value"

		value := p.String(path, expectedValue)
		if value != expectedValue {
			t.Errorf("%s returned the (%v) value, expected (%v)", action, value, expectedValue)
		}
	})
}

func Test_Partial_Config(t *testing.T) {
	t.Run("should correctly validate a invalid request", func(t *testing.T) {
		action := "Retrieving the path value as a config with an error"

		scenarios := []struct {
			partial partial
			path    string
		}{
			{ // test when the path dosen't exists
				partial: partial{},
				path:    "node1",
			},
			{ // test when the path is storing anil value
				partial: partial{"node1": nil},
				path:    "node1",
			},
			{ // test when the path is storing a int value
				partial: partial{"node1": 101},
				path:    "node1",
			},
			{ // test when the path is storing a string value
				partial: partial{"node1": "value1"},
				path:    "node1",
			},
		}

		for _, scn := range scenarios {
			test := func() {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("%s didn't panic", action)
					}
				}()
				scn.partial.Config(scn.path)
			}
			test()
		}
	})

	t.Run("should correctly retrieve a config object", func(t *testing.T) {
		action := "Retrieving the path value as a config object"

		p := partial{"node1": partial{"node2": "value1"}}
		path := "node1"
		expectedValue := partial{"node2": "value1"}

		value := p.Config(path)
		if !reflect.DeepEqual(value, expectedValue) {
			t.Errorf("%s returned the (%v) value, expected (%v)", action, value, expectedValue)
		}
	})

	t.Run("should return the given default value if path is not found", func(t *testing.T) {
		action := "Retrieving a non-existing path"

		p := partial{"node1": partial{"node2": 101}}
		path := "node3"
		expectedValue := partial{"node3": 345}

		value := p.Config(path, expectedValue)
		if !reflect.DeepEqual(value, expectedValue) {
			t.Errorf("%s returned the (%v) value, expected (%v)", action, value, expectedValue)
		}
	})
}

func Test_Partial_Merge(t *testing.T) {
	t.Run("should correctly merges two partial", func(t *testing.T) {
		action := "Merging two partials"

		scenarios := []struct {
			partial1 partial
			partial2 partial
			expected partial
		}{
			{ // test merging nil partial source
				partial1: partial{},
				partial2: nil,
				expected: partial{},
			},
			{ // test merging empty partial
				partial1: partial{},
				partial2: partial{},
				expected: partial{},
			},
			{ // test merging empty partial on non empty partial
				partial1: partial{"node1": "value1"},
				partial2: partial{},
				expected: partial{"node1": "value1"},
			},
			{ // test merging partial into empty partial
				partial1: partial{},
				partial2: partial{"node1": "value1"},
				expected: partial{"node1": "value1"},
			},
			{ // test merging override source value
				partial1: partial{"node1": "value1"},
				partial2: partial{"node1": "value2"},
				expected: partial{"node1": "value2"},
			},
			{ // test merging does not override non-present value in merged partial (create)
				partial1: partial{"node1": "value1"},
				partial2: partial{"node2": "value2"},
				expected: partial{"node1": "value1", "node2": "value2"},
			},
			{ // test merging does not override non-present value in merged partial (override)
				partial1: partial{"node1": "value1", "node2": "value2"},
				partial2: partial{"node2": "value3"},
				expected: partial{"node1": "value1", "node2": "value3"},
			},
			{ // test merging override source value to a subtree
				partial1: partial{"node1": "value1"},
				partial2: partial{"node1": partial{"node2": "value"}},
				expected: partial{"node1": partial{"node2": "value"}},
			},
			{ // test merging override source value in a subtree (to a value)
				partial1: partial{"node1": partial{"node2": "value1"}},
				partial2: partial{"node1": partial{"node2": "value2"}},
				expected: partial{"node1": partial{"node2": "value2"}},
			},
			{ // test merging override source value in a subtree (to a subtree)
				partial1: partial{"node1": partial{"node2": "value"}},
				partial2: partial{"node1": partial{"node2": partial{"node3": "value"}}},
				expected: partial{"node1": partial{"node2": partial{"node3": "value"}}},
			},
		}

		for _, scn := range scenarios {
			result := scn.partial1.merge(scn.partial2)
			if !reflect.DeepEqual(result, scn.expected) {
				t.Errorf("%s resulted in (%s) when merging (%v) and (%v), expecting (%v)", action, result, scn.partial1, scn.partial2, scn.expected)
			}
		}
	})
}
