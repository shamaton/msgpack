package common

import (
	"reflect"
	"strings"
)

// Common is used encoding/decoding
type Common struct{}

// FieldInfo holds information about a struct field including its path for embedded structs
type FieldInfo struct {
	Path      []int   // path to reach this field (indices for embedded structs)
	Name      string  // field name or tag
	Omit      bool    // omitempty flag
	Tagged    bool    // tag name explicitly set
	OmitPaths [][]int // paths to embedded fields with omitempty
}

// CollectFields collects all fields from a struct, expanding embedded structs
// following the same rules as encoding/json
func (c *Common) CollectFields(t reflect.Type, path []int) []FieldInfo {
	return c.collectFields(t, path, nil)
}

func (c *Common) collectFields(t reflect.Type, path []int, omitPaths [][]int) []FieldInfo {
	var fields []FieldInfo
	var embedded []FieldInfo // embedded fields to process later (lower priority)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check field visibility and get omitempty flag
		public, omit, name := c.CheckField(field)
		if !public {
			continue
		}

		// Get tag to check if embedded
		tag := field.Tag.Get("msgpack")
		// Extract just the name part (before comma if any)
		tagName := tag
		for j, ch := range tag {
			if ch == ',' {
				tagName = tag[:j]
				break
			}
		}

		// Check if this is an embedded struct
		isEmbedded := field.Anonymous && (tag == "" || tagName == "")
		tagged := tagName != ""

		if isEmbedded {
			// Get the actual type (dereference pointer if needed)
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			// If it's a struct, expand its fields
			if fieldType.Kind() == reflect.Struct {
				newPath := append(append([]int{}, path...), i)
				nextOmitPaths := omitPaths
				if omit {
					nextOmitPaths = appendOmitPath(omitPaths, newPath)
				}
				embeddedFields := c.collectFields(fieldType, newPath, nextOmitPaths)
				embedded = append(embedded, embeddedFields...)
				continue
			}
		}

		// Regular field or embedded non-struct
		newPath := append(append([]int{}, path...), i)
		fields = append(fields, FieldInfo{
			Path:      newPath,
			Name:      name,
			Omit:      omit,
			Tagged:    tagged,
			OmitPaths: omitPaths,
		})
	}

	// Add embedded fields after regular fields (they have lower priority)
	fields = append(fields, embedded...)

	// Remove duplicates and handle ambiguous fields
	return c.deduplicateFields(fields)
}

func appendOmitPath(paths [][]int, path []int) [][]int {
	if len(paths) == 0 {
		return [][]int{path}
	}
	newPaths := make([][]int, len(paths)+1)
	copy(newPaths, paths)
	newPaths[len(paths)] = path
	return newPaths
}

// deduplicateFields removes duplicate fields and handles ambiguous fields
// following encoding/json behavior
func (c *Common) deduplicateFields(fields []FieldInfo) []FieldInfo {
	// Group fields by name and depth, preserving order
	type fieldAtDepth struct {
		field FieldInfo
		depth int
	}
	fieldsByName := make(map[string][]fieldAtDepth)
	var seenNames []string // To preserve order

	for _, f := range fields {
		if _, seen := fieldsByName[f.Name]; !seen {
			seenNames = append(seenNames, f.Name)
		}
		fieldsByName[f.Name] = append(fieldsByName[f.Name], fieldAtDepth{
			field: f,
			depth: len(f.Path),
		})
	}

	var result []FieldInfo
	for _, name := range seenNames {
		fieldsWithDepth := fieldsByName[name]

		// Find minimum depth
		minDepth := fieldsWithDepth[0].depth
		for _, fd := range fieldsWithDepth {
			if fd.depth < minDepth {
				minDepth = fd.depth
			}
		}

		// Count fields at minimum depth
		var fieldsAtMinDepth []FieldInfo
		for _, fd := range fieldsWithDepth {
			if fd.depth == minDepth {
				fieldsAtMinDepth = append(fieldsAtMinDepth, fd.field)
			}
		}

		// If there's exactly one field at minimum depth, use it
		if len(fieldsAtMinDepth) == 1 {
			result = append(result, fieldsAtMinDepth[0])
			continue
		}

		// Prefer the tagged field if exactly one is tagged at minimum depth
		var taggedFields []FieldInfo
		for _, f := range fieldsAtMinDepth {
			if f.Tagged {
				taggedFields = append(taggedFields, f)
			}
		}
		if len(taggedFields) == 1 {
			result = append(result, taggedFields[0])
		}
		// else: ambiguous field, skip it (following encoding/json behavior)
	}

	return result
}

// CheckField returns flag whether should encode/decode or not and field name
func (c *Common) CheckField(field reflect.StructField) (public, omit bool, name string) {
	// A to Z
	if !c.isPublic(field.Name) {
		return false, false, ""
	}
	tag := field.Tag.Get("msgpack")
	if tag == "" {
		return true, false, field.Name
	}

	parts := strings.Split(tag, ",")
	// check ignore
	if parts[0] == "-" {
		return false, false, ""
	}
	// check omitempty
	for _, part := range parts[1:] {
		if part == "omitempty" {
			omit = true
		}
	}
	// check name
	name = field.Name
	if parts[0] != "" {
		name = parts[0]
	}
	return true, omit, name
}

func (c *Common) isPublic(name string) bool {
	return 0x41 <= name[0] && name[0] <= 0x5a
}
