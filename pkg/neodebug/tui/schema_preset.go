package tui

import (
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"reflect"
)

// Decode independent wire values so editing never mutates the configured preset.
func (m *inputRequestCreator) loadPreset(data any) error {
	payload, err := msgp.AppendIntf(nil, data)
	if err != nil {
		return err
	}
	value, _, err := msgp.ReadIntfBytes(payload)
	if err != nil {
		return err
	}
	m.rootNode.Unselect()
	if err := loadPresetNode(m.rootNode, value); err != nil {
		return err
	}
	m.rootNode.Init()
	m.rootNode.SelectFromTop()
	return nil
}

func (v *ValueNode[T]) setPresetValue(value any) error {
	target := reflect.TypeFor[T]()
	source := reflect.ValueOf(value)
	if !source.IsValid() || !source.Type().ConvertibleTo(target) {
		return fmt.Errorf("cannot load %T into %v", value, target)
	}
	v.value = source.Convert(target).Interface().(T)
	return nil
}

func loadPresetNode(node SchemaNode, value any) error {
	switch node := node.(type) {
	case *NullableNode:
		node.null = value == nil
		if node.null {
			return nil
		}
		return loadPresetNode(node.other, value)
	case *StructNode:
		if node.tuple {
			values, ok := value.([]any)
			if !ok || len(values) != len(node.children) {
				return fmt.Errorf("invalid tuple body")
			}
			for i, field := range node.children {
				if err := loadPresetNode(field.Node, values[i]); err != nil {
					return fmt.Errorf("%s: %w", field.Name, err)
				}
			}
			return nil
		}
		values, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object, got %T", value)
		}
		for _, field := range node.children {
			if fieldValue, exists := values[field.Name]; exists {
				if err := loadPresetNode(field.Node, fieldValue); err != nil {
					return fmt.Errorf("%s: %w", field.Name, err)
				}
			}
		}
		return nil
	case *SliceNode:
		values, ok := value.([]any)
		if value == nil {
			values, ok = nil, true
		}
		if !ok {
			return fmt.Errorf("expected array, got %T", value)
		}
		node.items = make([]SchemaNode, 0, len(values))
		for _, item := range values {
			child := createNode(node.element, node.registry)
			if err := loadPresetNode(child, item); err != nil {
				return err
			}
			node.items = append(node.items, child)
		}
		return nil
	case *MapNode:
		values := reflect.ValueOf(value)
		if value == nil {
			node.items = nil
			return nil
		}
		if values.Kind() != reflect.Map {
			return fmt.Errorf("expected map, got %T", value)
		}
		node.items = make([]MapEntry, 0, values.Len())
		iter := values.MapRange()
		for iter.Next() {
			entry := MapEntry{Key: createNode(node.keyType, node.registry), Value: createNode(node.valueType, node.registry)}
			if err := loadPresetNode(entry.Key, iter.Key().Interface()); err != nil {
				return err
			}
			if err := loadPresetNode(entry.Value, iter.Value().Interface()); err != nil {
				return err
			}
			node.items = append(node.items, entry)
		}
		return nil
	case *FileNode:
		if value == nil {
			node.value = nil
			return nil
		}
		data, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("expected bytes, got %T", value)
		}
		node.value = append([]byte(nil), data...)
		return nil
	case interface{ setPresetValue(any) error }:
		return node.setPresetValue(value)
	default:
		return fmt.Errorf("unsupported preset editor %T", node)
	}
}
