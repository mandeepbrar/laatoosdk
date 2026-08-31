package data

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
)

// QueryComponent is one form in which a query may be WRITTEN — OData, openCypher, the
// declarative Filters form, or anything a deployment adds.
//
// Every form lowers into the same Query AST. The form is a writing choice and not a different
// execution path — nothing downstream can tell which one produced the query, which is what lets a
// provider compile one representation rather than one per surface syntax.
//
// It lives here, and not in the server, so that a surface syntax is something a PLUGIN supplies.
// Before this interface every form shipped inside the server binary: OData and openCypher were
// ~1500 lines across 65 functions, plus a third-party parser, carried by every deployment whether
// or not one dataset used them. A form is now opt-in, and adding one needs no server change.
//
// Registration is through elements.DataManager.RegisterQueryComponent, mirroring
// RegisterDataComponent: a live element the plugin already reaches through its context, not
// package-level state that would have to survive plugin.Open.
//
// A QUERY IS WRITTEN IN TWO PLACES, AND ONE INTERFACE COVERS BOTH. A dataset declares its query
// as CONFIGURATION, resolved once at load: Claims and Lower serve that. A Go caller holds it as
// TEXT at request time with no config in hand: ParseQuery serves that. A form implements all
// three, because a form that served only datasets would leave every programmatic caller back
// where DataComponent.CreateODataQuery was — declared and refusing.
type QueryComponent interface {
	// Name identifies the form. A dataset selects it explicitly under the `frontend` key, and a
	// programmatic caller names it the same way; the name is what both are matched against.
	Name() string

	// Claims reports whether this component recognises a dataset's declaration, and is what
	// selects one when the dataset does not name a form explicitly. Components are offered the
	// declaration in registration order, so an earlier registration wins one that two would claim.
	//
	// Returning false unconditionally is legitimate and sometimes REQUIRED: a form sharing a
	// config key with another form cannot claim by content without guessing which language a
	// string is written in, and guessing wrong reads the dataset as a different query rather than
	// refusing it. Such a form is selectable only by the explicit key.
	Claims(ctx core.ServerContext, conf config.Config) bool

	// Lower converts a dataset's declaration into the query AST, writing onto the Query it is
	// given. It is called at most once per dataset, at load.
	//
	// The returned slice is the projection the declaration named — OData's $select, and whatever
	// the equivalent is in another form. Return nil when the form carries no projection of its
	// own, which leaves the dataset's declared Properties in force.
	Lower(ctx core.ServerContext, conf config.Config, query *Query) ([]string, error)

	// ParseQuery converts query TEXT into the AST, writing onto the Query it is given, and
	// returns the projection the text named — nil when it named none.
	//
	// It takes the same Query the dataset path writes to, deliberately: a caller may already have
	// set parts of it, and a form that overwrote rather than composed would silently discard them.
	//
	// A form with no text spelling returns errors.NotImplemented, and that is a correct answer
	// rather than a stub. The declarative Filters form is the case: it is structured
	// configuration, not a language, and its text spelling is OData — which is a DIFFERENT
	// registered component. Refusing here sends the caller to the form that can actually read
	// what they hold.
	//
	// This is one method that may refuse, deliberately NOT a second optional interface asserted
	// for at the call site. That shape belongs to ExpandingComponent, where it exists to avoid
	// widening a DataComponent that seven providers already implement, and where the failed
	// assertion selects a real fallback — the hop executor. Neither holds here: this interface is
	// new so there is nothing to break by widening it, and a form that cannot read text has no
	// fallback to select. An error at the call site is the whole of the behaviour, so an error is
	// the whole of the contract.
	ParseQuery(ctx core.ServerContext, queryText string, query *Query) ([]string, error)
}
