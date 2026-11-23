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

// TaskManager is the BACKEND contract for a task queue provider (embedded NATS/JetStream,
// GCP Cloud Tasks, beanstalkd, gaetasks, an in-process worker pool). Plugin code does not call it
// directly — it calls ctx.PushTask, which resolves the backend that owns the queue from the
// queue's registry declaration and delegates here (laatooserver/src/core/taskmanager.go:353-378).
//
// THERE IS NO SCHEDULING ANYWHERE IN THIS INTERFACE OR IN THE PLATFORM'S TASK MANAGER. A task
// fires only because something called PushTask. A queue declaration
// (src/server/registry/tasks/<name>.yml) binds a queue to a processor service; it carries no cron
// expression, interval or delay, and nothing in the server's task manager reads one. Recurring
// work must be driven by something outside the task subsystem.
//
// A queue must be DECLARED before either side works: an undeclared queue has no manager mapping,
// so PushTask fails with CORE_ERROR_BAD_CONF rather than creating one on demand
// (taskmanager.go:371-374, :628-635).
type TaskManager interface {
	// PushTask enqueues one task and returns the id it was accepted under.
	//
	// IDENTITY IS ALREADY STAMPED BY THE TIME A BACKEND SEES THE TASK. The platform's task manager
	// builds the Task with User: ctx.GetUser() and Tenant: ctx.GetTenant() and generates the Id,
	// before any backend is called (taskmanager.go:369). A backend therefore does not thread
	// identity separately — it only has to preserve the Task's User and Tenant across
	// serialisation, and the processor then runs as that user. A task whose identity cannot be
	// rebuilt on the way out is REFUSED rather than run as system, because a nil user reaching
	// CreateSystemRequest disables authorization entirely (taskmanager.go:645-650).
	//
	// The returned id and the Task.Id set by the caller are expected to be the same value; every
	// shipped backend returns t.Id.
	//
	// ACCEPTANCE SEMANTICS VARY BY BACKEND AND SOME ACCEPT SILENTLY WITHOUT ENQUEUEING. The NATS
	// backend publishes through JetStream and waits for the server's persistence acknowledgement,
	// so a nil error means the task is durably stored (taskmanager.go:942-972). The in-process
	// tunnymemtasks backend, by contrast, looks up the queue's worker pool and, finding none —
	// which is the case whenever SubsribeQueue was never called for that queue — DROPS the task
	// and returns the id with a nil error
	// (laatoomodules/tasks/.../tunnymemtasks/src/server/go/taskprocessor.go:31-40).
	PushTask(ctx core.RequestContext, task *Task) (string, error)

	// SubsribeQueue starts consuming a queue in THIS process and dispatching its tasks to the
	// queue's registered processor service. (The spelling is the SDK's; it is load-bearing.)
	//
	// It is called by the platform at Start for every queue whose declaration names a processor
	// this level hosts. It is the producer/consumer split: a pod that only pushes never calls it.
	//
	// The NATS backend binds a DURABLE consumer whose name is shared by every pod running the
	// queue, which is what load-balances the work, with MaxAckPending 1 so a worker is handed one
	// task at a time. The ack policy is the real contract: a decode failure is terminated and
	// dead-lettered (retrying undecodable bytes cannot help), a processor error is Nak'd and
	// retried up to MaxDeliver before being dead-lettered, and acceptance is acked
	// (taskmanager.go:974-1089).
	//
	// "ACCEPTED" IS NOT "FINISHED". A queue declared async: true acks as soon as the processor
	// accepts the task, because JetStream redelivers anything unacked past AckWait and a
	// human-scale wait would either redeliver mid-approval or need an AckWait long enough to
	// disable redelivery for genuine crashes. Durability for the long-running half lives in the
	// task record and the completion topic instead.
	SubsribeQueue(ctx core.ServerContext, queue string) error

	// UnsubsribeQueue stops this process consuming a queue, draining what is in flight.
	//
	// Used by module hot-unload. NOT every backend implements it: beanstalktasks returns
	// errors.NotImplemented (beanstalkmanager.go:142-144), so unloading a module whose queue is
	// served by that backend leaves the consumer running.
	UnsubsribeQueue(ctx core.ServerContext, queue string) error

	// GetTask returns the task previously pushed to queue under the id PushTask handed back.
	// Every backend answers this — either from its own addressable store (Cloud Tasks indexes
	// by task name) or from a store Laatoo provides. A backend that can do neither returns
	// errors.NotImplemented until one is wired in.
	GetTask(ctx core.RequestContext, queue string, id string) (*Task, error)
}

// The keys a task's user and tenant are written under, each as an object of its own. A user's own fields
// cannot share a key with the task's this way: written flattened beside the task's, the user's "Id" and
// the task's "Id" were one key, and a reader kept only one of them -- the JSON reader the first, so every
// task ran as a user whose id was the task's; the binary codec the last, so the task took the user's id.
const (
	taskUserKey   = "User"
	taskTenantKey = "Tenant"
)

// ReadAll reads a task: its own fields, then its user and tenant. The user and tenant are read from their
// flattened fields first, which is all a task queued before they had objects of their own carries, and
// then from their own objects, which a task queued since carries and which override the flattened read.
// A task with no principal or tenant to fill (no configured principal type, a system caller) reads its
// own fields only.
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

	if ent.User != nil {
		// the flattened layout: a task queued before the user had an object of its own
		if err = ent.User.ReadAll(c, cdc, rdr); err != nil {
			return err
		}
		// the user's own object, when the task carries one; absent, nothing is read
		if err = rdr.ReadObject(c, cdc, taskUserKey, ent.User); err != nil {
			return err
		}
	}
	if ent.Tenant != nil {
		// the flattened layout, then the tenant's own object, as for the user
		if err = ent.Tenant.ReadAll(c, cdc, rdr); err != nil {
			return err
		}
		if err = rdr.ReadObject(c, cdc, taskTenantKey, ent.Tenant); err != nil {
			return err
		}
	}
	if err = rdr.ReadObject(c, cdc, "Metadata", &ent.Metadata); err != nil {
		return err
	}

	return nil
}

// WriteAll writes a task: its own fields, then its user and tenant each as an object of its own. The user's
// and tenant's fields are also written flattened, as they were before, so a server that reads only that
// layout (an older one in a cluster during an upgrade) still reads the task's user and roles; the task's
// own "Id" is written again after them, so a writer keeping a repeated key's last write keeps the task's
// id, as a reader keeping the first does.
func (ent *Task) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	var err error

	if err = wtr.WriteString(c, cdc, "Id", &ent.Id); err != nil {
		return err
	}

	if err = wtr.WriteString(c, cdc, "Queue", &ent.Queue); err != nil {
		return err
	}

	if err = wtr.WriteArray(c, cdc, "Data", &ent.Data); err != nil {
		return err
	}

	// the flattened layout, for a reader that knows no other
	if ent.User != nil {
		if err = ent.User.WriteAll(c, cdc, wtr); err != nil {
			return err
		}
	}
	if ent.Tenant != nil {
		if err = ent.Tenant.WriteAll(c, cdc, wtr); err != nil {
			return err
		}
	}
	if ent.User != nil || ent.Tenant != nil {
		// the task's own id last as well as first: the flattened user wrote an "Id" of its own
		if err = wtr.WriteString(c, cdc, "Id", &ent.Id); err != nil {
			return err
		}
	}

	// the user and the tenant each as an object of its own, which is what a reader trusts
	if ent.User != nil {
		// the user itself, not a pointer to the interface: only the user serialises itself
		if err = wtr.WriteObject(c, cdc, taskUserKey, ent.User); err != nil {
			return err
		}
	}
	if ent.Tenant != nil {
		if err = wtr.WriteObject(c, cdc, taskTenantKey, ent.Tenant); err != nil {
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
