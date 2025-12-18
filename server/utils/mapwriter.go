package utils

import (
	"io"
	"strings"
	"time"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/errors"
	"laatoo.io/sdk/utils"
)

func CreateObjectFromMap(ctx core.ServerContext, objType string, smap utils.StringMap, transformations utils.StringMap) (datatypes.Serializable, error) {
	obj, err := ctx.CreateObject(objType)
	if err != nil {
		return nil, errors.WrapError(ctx, err)
	}
	serObj, ok := obj.(datatypes.Serializable)
	if !ok {
		return nil, errors.SerializationError(ctx, "Object type is not serializable", objType)
	}

	finalMap := smap
	if transformations != nil {
		finalMap = applyTransformations(smap, transformations)
	}

	wr := &MapSerializableWriter{finalMap}
	err = serObj.WriteAll(ctx, nil, wr)
	if err != nil {
		return nil, errors.WrapError(ctx, err)
	}
	return serObj, nil
}

func applyTransformations(data utils.StringMap, transformations utils.StringMap) utils.StringMap {
	if len(transformations) == 0 {
		return data
	}

	finalMap := make(utils.StringMap)
	for k, v := range data {
		newKey := k
		var subTransforms utils.StringMap

		if tVal, ok := transformations[k]; ok {
			if sVal, ok := tVal.(string); ok {
				newKey = sVal
			} else if mVal, ok := tVal.(map[string]interface{}); ok {
				subTransforms = utils.StringMap(mVal)
				if keyVal, ok := subTransforms.GetString("__key"); ok {
					newKey = keyVal
				}
			} else if mVal, ok := tVal.(utils.StringMap); ok {
				subTransforms = mVal
				if keyVal, ok := subTransforms.GetString("__key"); ok {
					newKey = keyVal
				}
			}
		}

		if subMap, ok := data.GetStringMap(k); ok {
			// If we have sub-transformations for this key (which is a map), apply them recursively
			transformedSubMap := applyTransformations(subMap, subTransforms)
			// Then insert the transformed map at the new key locaion
			insertDeep(finalMap, newKey, transformedSubMap)
		} else {
			insertDeep(finalMap, newKey, v)
		}
	}
	return finalMap
}

func insertDeep(data utils.StringMap, path string, value interface{}) {
	parts := strings.Split(path, ".")
	if len(parts) == 1 {
		data[parts[0]] = value
		return
	}

	current := data
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if _, ok := current[part]; !ok {
			current[part] = make(utils.StringMap)
		}

		if nextMap, ok := current[part].(utils.StringMap); ok {
			current = nextMap
		} else if nextMap, ok := current[part].(map[string]interface{}); ok {
			newMap := utils.StringMap(nextMap)
			current[part] = newMap
			current = newMap
		} else {
			// Conflict: path segment exists but is not a map. Overwrite with new map.
			newMap := make(utils.StringMap)
			current[part] = newMap
			current = newMap
		}
	}
	current[parts[len(parts)-1]] = value
}

type MapSerializableWriter struct {
	Data utils.StringMap
}

func (w *MapSerializableWriter) Start() error {
	return nil
}

func (w *MapSerializableWriter) Close() error {
	return nil
}

func (w *MapSerializableWriter) Write(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (w *MapSerializableWriter) Bytes() []byte {
	return nil
}

func (w *MapSerializableWriter) WriteBytes(ctx ctx.Context, cdc datatypes.Codec, prop string, val *[]byte) error {
	if v, ok := w.Data[prop]; ok {
		if b, ok := v.([]byte); ok {
			*val = b
		} else if s, ok := v.(string); ok {
			*val = []byte(s)
		} else {
			return errors.SerializationError(ctx, "Value is not a byte array", prop)
		}
	}
	return nil
}

func (w *MapSerializableWriter) WriteInt(ctx ctx.Context, cdc datatypes.Codec, prop string, val *int) error {
	if v, ok := w.Data.GetInt(prop); ok {
		*val = v
	} else if v, ok := w.Data[prop]; ok {
		// Handle float64 which is common in JSON maps
		if f, ok := v.(float64); ok {
			*val = int(f)
		} else {
			return errors.SerializationError(ctx, "Value is not an int", prop)
		}
	}
	return nil
}

func (w *MapSerializableWriter) WriteInt32(ctx ctx.Context, cdc datatypes.Codec, prop string, val *int32) error {
	if v, ok := w.Data[prop]; ok {
		if i, ok := v.(int32); ok {
			*val = i
		} else if i, ok := v.(int); ok {
			*val = int32(i)
		} else if f, ok := v.(float64); ok {
			*val = int32(f)
		} else {
			return errors.SerializationError(ctx, "Value is not an int32", prop)
		}
	}
	return nil
}

func (w *MapSerializableWriter) WriteInt64(ctx ctx.Context, cdc datatypes.Codec, prop string, val *int64) error {
	if v, ok := w.Data[prop]; ok {
		if i, ok := v.(int64); ok {
			*val = i
		} else if i, ok := v.(int); ok {
			*val = int64(i)
		} else if f, ok := v.(float64); ok {
			*val = int64(f)
		} else {
			return errors.SerializationError(ctx, "Value is not an int64", prop)
		}
	}
	return nil
}

func (w *MapSerializableWriter) WriteString(ctx ctx.Context, cdc datatypes.Codec, prop string, val *string) error {
	if v, ok := w.Data.GetString(prop); ok {
		*val = v
	} else if _, ok := w.Data[prop]; !ok {
		// Ignore missing values
		return nil
	} else {
		return errors.SerializationError(ctx, "Value is not a string", prop)
	}
	return nil
}

func (w *MapSerializableWriter) WriteFloat32(ctx ctx.Context, cdc datatypes.Codec, prop string, val *float32) error {
	if v, ok := w.Data[prop]; ok {
		if f, ok := v.(float32); ok {
			*val = f
		} else if f, ok := v.(float64); ok {
			*val = float32(f)
		} else {
			return errors.SerializationError(ctx, "Value is not a float32", prop)
		}
	}
	return nil
}

func (w *MapSerializableWriter) WriteFloat64(ctx ctx.Context, cdc datatypes.Codec, prop string, val *float64) error {
	if v, ok := w.Data[prop]; ok {
		if f, ok := v.(float64); ok {
			*val = f
		} else if f, ok := v.(float32); ok {
			*val = float64(f)
		} else {
			return errors.SerializationError(ctx, "Value is not a float64", prop)
		}
	}
	return nil
}

func (w *MapSerializableWriter) WriteBool(ctx ctx.Context, cdc datatypes.Codec, prop string, val *bool) error {
	v, ok := w.Data.GetBool(prop)
	if ok {
		*val = v
	} else if _, ok := w.Data[prop]; !ok {
		return nil
	} else {
		return errors.SerializationError(ctx, "Field is not boolean", prop)
	}
	return nil
}

func (w *MapSerializableWriter) WriteObject(ctx ctx.Context, cdc datatypes.Codec, prop string, val interface{}) error {
	if subMap, ok := w.Data.GetStringMap(prop); ok {
		if ser, ok := val.(datatypes.Serializable); ok {
			subWriter := &MapSerializableWriter{Data: subMap}
			return ser.WriteAll(ctx, cdc, subWriter)
		}
	} else if _, ok := w.Data[prop]; !ok {
		return nil
	} else {
		return errors.SerializationError(ctx, "Field is not string map", prop)
	}
	return nil
}

func (w *MapSerializableWriter) WriteMap(ctx ctx.Context, cdc datatypes.Codec, prop string, val interface{}) error {
	// Not implemented for this use case as hydrating a map from a map is tricky with typed target
	return nil
}

func (w *MapSerializableWriter) WriteArray(ctx ctx.Context, cdc datatypes.Codec, prop string, val interface{}) error {
	// Handling arrays would require reflection or specific array Writer interfaces
	// For now leaving as no-op or basic handling if expected
	return nil
}

func (w *MapSerializableWriter) WriteTime(ctx ctx.Context, cdc datatypes.Codec, prop string, val *time.Time) error {
	if v, ok := w.Data[prop]; ok {
		if t, ok := v.(time.Time); ok {
			*val = t
		} else if s, ok := v.(string); ok {
			// Try parsing standard formats if needed, or leave strictly type matched
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				*val = t
			}
		} else {
			return errors.SerializationError(ctx, "Field is not time", prop)
		}
	}
	return nil
}
