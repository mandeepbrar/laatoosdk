package components

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Task struct {
	datatypes.Serializable
	Queue    string          `json:"queue"`
	Data     []byte          `json:"data"`
	Id       string          `json:"id"`
	User     auth.User       `json:"user"`
	Tenant   auth.TenantInfo `json:"tenant"`
	Metadata utils.StringMap `json:"metadata,omitempty"`
}

type TaskManager interface {
	PushTask(ctx core.RequestContext, task *Task) (string, error)
	SubsribeQueue(ctx core.ServerContext, queue string) error
	UnsubsribeQueue(ctx core.ServerContext, queue string) error
}

func (ent *Task) ReadAll(c ctx.Context, cdc datatypes.Codec, rdr datatypes.SerializableReader) error {
	var err error

	if err = rdr.ReadString(c, cdc, "Id", &ent.Id); err != nil {
		return err
	}

	if err = rdr.ReadString(c, cdc, "Queue", &ent.Queue); err != nil {
		return err
	}

	if err = rdr.ReadArray(c, cdc, "Data", &ent.Data); err != nil {
		return err
	}

	err = ent.User.ReadAll(c, cdc, rdr)
	if err != nil {
		return err
	}
	err = ent.Tenant.ReadAll(c, cdc, rdr)
	if err != nil {
		return err
	}
	if err = rdr.ReadObject(c, cdc, "Metadata", &ent.Metadata); err != nil {
		return err
	}

	return nil
}

func (ent *Task) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	var err error

	if err = wtr.WriteString(c, cdc, "Id", &ent.Id); err != nil {
		return err
	}

	if err = wtr.WriteString(c, cdc, "Queue", &ent.Queue); err != nil {
		return err
	}

	/*	if err = wtr.WriteObject(c, cdc, "User", &ent.User); err != nil {
			return err
		}

		if err = wtr.WriteObject(c, cdc, "Tenant", &ent.Tenant); err != nil {
			return err
		}
	*/
	if err = wtr.WriteArray(c, cdc, "Data", &ent.Data); err != nil {
		return err
	}

	if ent.User != nil {
		err = ent.User.WriteAll(c, cdc, wtr)
		if err != nil {
			return err
		}

	}
	if ent.Tenant != nil {
		err = ent.Tenant.WriteAll(c, cdc, wtr)
		if err != nil {
			return err
		}
	}
	if ent.Metadata != nil {
		if err = wtr.WriteObject(c, cdc, "Metadata", &ent.Metadata); err != nil {
			return err
		}
	}

	return nil
}
type TaskCompletionMessage struct {
	InvocationId string          `json:"invocation_id"`
	Queue        string          `json:"queue"`
	Result       interface{}     `json:"result"`
	Metadata     utils.StringMap `json:"metadata,omitempty"`
	Error        string          `json:"error,omitempty"`
}
