package components

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/server/core"
)

type Task struct {
	datatypes.Serializable
	Queue  string
	Data   []byte
	Id     string
	User   auth.User
	Tenant auth.TenantInfo
}

type TaskManager interface {
	PushTask(ctx core.RequestContext, task *Task) error
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

	return nil
}
