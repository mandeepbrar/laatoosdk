package ai

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/components/data"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// MemoryItem is a logical unit of memory (DTO).
// It acts as a data transfer object between Agents and MemoryBanks.
type MemoryItem struct {
	data.StorageInfo
	data.TenantInfo
	Type       string          `json:"Type"` // e.g., "message", "artifact", "datarecord"
	Content    []byte          `json:"Content"`
	Importance float64         `json:"Importance"`
	Timestamp  string          `json:"Timestamp"` // ISO8601 string
	Tags       []string        `json:"Tags"`
	Metadata   utils.StringMap `json:"Metadata"`
	Vector     []float32       `json:"Vector"`
}

func (mi *MemoryItem) GetImportance() float64       { return mi.Importance }
func (mi *MemoryItem) GetTimestamp() string         { return mi.Timestamp }
func (mi *MemoryItem) GetTags() []string            { return mi.Tags }
func (mi *MemoryItem) GetMetadata() utils.StringMap { return mi.Metadata }
func (mi *MemoryItem) GetContent() any              { return mi.Content }

func (mi *MemoryItem) Config() *core.StorableConfig {
	return &core.StorableConfig{
		ObjectType:  "ai.MemoryItem",
		LabelField:  "Id",
		PreSave:     false,
		PostSave:    false,
		Workflow:    false,
		PostUpdate:  false,
		PostLoad:    false,
		Trackable:   false,
		Multitenant: true,
		Collection:  "MemoryItem",
		Cacheable:   true,
	}
}

func (ent *MemoryItem) ReadAll(c ctx.Context, cdc datatypes.Codec, rdr datatypes.SerializableReader) error {
	var err error

	if err = rdr.ReadMap(c, cdc, "Metadata", &ent.Metadata); err != nil {
		return err
	}

	if err = rdr.ReadString(c, cdc, "Type", &ent.Type); err != nil {
		return err
	}

	cont, err := rdr.ReadBytes(c, cdc, "Content")
	if err != nil {
		return err
	} else {
		ent.Content = cont
	}

	if err = rdr.ReadFloat64(c, cdc, "Importance", &ent.Importance); err != nil {
		return err
	}

	if err = rdr.ReadString(c, cdc, "Timestamp", &ent.Timestamp); err != nil {
		return err
	}

	if err = rdr.ReadArray(c, cdc, "Tags", &ent.Tags); err != nil {
		return err
	}

	if err = rdr.ReadArray(c, cdc, "Vector", &ent.Vector); err != nil {
		return err
	}

	err = ent.TenantInfo.ReadAll(c, cdc, rdr)
	if err != nil {
		return err
	}

	err = ent.StorageInfo.ReadAll(c, cdc, rdr)
	if err != nil {
		return err
	}
	return nil
}

func (ent *MemoryItem) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	var err error

	if err = wtr.WriteMap(c, cdc, "Metadata", &ent.Metadata); err != nil {
		return err
	}

	if err = wtr.WriteString(c, cdc, "Type", &ent.Type); err != nil {
		return err
	}

	if err = wtr.WriteBytes(c, cdc, "Content", &ent.Content); err != nil {
		return err
	}

	if err = wtr.WriteFloat64(c, cdc, "Importance", &ent.Importance); err != nil {
		return err
	}

	if err = wtr.WriteString(c, cdc, "Timestamp", &ent.Timestamp); err != nil {
		return err
	}

	if err = wtr.WriteArray(c, cdc, "Tags", &ent.Tags); err != nil {
		return err
	}

	if err = wtr.WriteArray(c, cdc, "Vector", &ent.Vector); err != nil {
		return err
	}

	err = ent.TenantInfo.WriteAll(c, cdc, wtr)
	if err != nil {
		return err
	}

	return ent.StorageInfo.WriteAll(c, cdc, wtr)
}
