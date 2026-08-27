package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/utils"
)

// Service is the interface every service in a Laatoo plugin satisfies. Services are the main
// units of logic in Laatoo.
//
// A plugin service struct embeds core.Service. At load time the server reflects over the struct,
// looks for a field literally named "Service", and assigns its own implementation into it
// (laatooserver core/serviceelement.go:47-56); a struct that does not embed core.Service under
// that exact field name fails to load with Core_Error_Type_Mismatch. Every Add*/Set*/Get* method
// below is therefore the server's implementation, not the plugin's.
//
// Lifecycle, in the order the server drives it:
//
//	Describe -> config parsed against what Describe declared -> configuration values injected
//	  into struct fields -> Initialize -> injected services resolved -> Start
//	  -> (per request) Invoke -> Stop -> Unload
//
// Anything declared after Describe -- in Initialize or later -- is never matched against
// configuration.
type Service interface {
	ConfigurableObject
	// Metadata returns this service's ServiceInfo: declared configurations, request and response
	// params, description and version. The server seeds it from the object spec
	// (registry/objects/<plugin>.<GoStructName>.yml) and then lets Describe add to it. Declared
	// values are only populated once configuration has been processed, so reading it from
	// Describe yields declarations whose values are nil.
	Metadata() ServiceInfo
	// Describe is the one place to declare the service's shape -- AddParam*, Add*Configuration,
	// SetDescription, SetComponent. The server calls it once at load, BEFORE any configuration
	// is parsed and before Initialize (laatooserver core/serviceelement.go:66). Unlike
	// Module.Describe and ServiceFactory.Describe, the error returned here IS checked and aborts
	// the load.
	Describe(ServerContext) error
	// Initialize runs after configuration has been parsed against the declarations Describe made
	// and after configuration-backed struct fields have been injected. conf is the RAW service
	// config, so conf.GetString(ctx, k) sees every key in the file whether declared or not --
	// whereas GetStringConfiguration sees only declared keys. Services named by an `inject:`
	// configuration are NOT resolved yet; they are wired just before Start.
	Initialize(ctx ServerContext, conf config.Config) error
	// Start runs after every injected service has been resolved and assigned to its struct field
	// and after the middleware chain has been built. This is the earliest point at which an
	// injected DataComponent or peer service is usable.
	Start(ctx ServerContext) error
	// Stop is called first during service unload. If it returns an error the server logs it and
	// returns immediately, so Unload is NOT called (laatooserver core/servicemanager.go:548-552).
	Stop(ctx ServerContext) error
	// Unload is called immediately after a successful Stop. Release connections and goroutines
	// here.
	Unload(ctx ServerContext) error
	// RequestParameters returns the LIVE map of declared request params -- the template, not a
	// copy, so mutating it mutates the service definition. The Param values in it are always
	// nil: per-request values live on clones held by the Request.
	RequestParameters(ctx ServerContext) map[string]Param
	// AddParams declares each name/type pair as a request param with an empty description, not a
	// collection and not streaming. The bool is `required` and applies to all of them. It returns
	// on the first error; map iteration order means which one that is, is not deterministic.
	AddParams(ServerContext, map[string]datatypes.DataType, bool) error
	//	AddStringParams(ctx ServerContext, params utils.StringsMap, defaultValues []string)
	// AddStringParam declares a string request param -- as a COLLECTION. It calls
	// AddParam(..., collection=true, required=false, stream=false)
	// (laatooserver core/serviceimpl.go:132-134), so the param binds []string and a plain string
	// value fails the request with Core_Error_Bad_Arg. For a scalar string use
	// AddOptionalParamWithType(ctx, name, desc, datatypes.String). It also discards the error
	// AddParam returns.
	AddStringParam(ctx ServerContext, name string, desc string)
	// AddCustomObjectParam declares a param of datatypes.Object whose concrete type is the
	// registered object name customObjectType (for example "myplugin.MyEntity"); the ObjectLoader
	// instantiates it when an encoded body is decoded. The Param interface exposes no accessor
	// for customObjectType afterwards.
	AddCustomObjectParam(ctx ServerContext, name string, desc string, customObjectType string, collection, required, stream bool) error
	// AddParam declares a request param of a specific data type. Argument order is
	// (collection, required, stream) -- note it differs from MetaDataProvider.CreateParam, which
	// takes (collection, isStream, required). required is enforced when the request is built: an
	// absent required param fails with Core_Error_Missing_Arg before the service is invoked.
	AddParam(ctx ServerContext, name string, desc string, datatype datatypes.DataType, collection, required, stream bool) error
	// AddParamWithType declares an OPTIONAL, STREAMING param -- despite the name it does not mark
	// the param required. It calls AddParam(..., collection=false, required=false, stream=true)
	// (laatooserver core/serviceimpl.go:98-100), so the only difference from
	// AddOptionalParamWithType is the stream flag. To declare a required scalar param, call
	// AddParam directly.
	AddParamWithType(ctx ServerContext, name string, desc string, datatype datatypes.DataType) error
	// AddOptionalParamWithType declares an optional, non-collection, non-streaming param
	// (laatooserver core/serviceimpl.go:102-104).
	AddOptionalParamWithType(ctx ServerContext, name string, desc string, datatype datatypes.DataType) error
	//AddCollectionParams(ctx ServerContext, params map[string]datatypes.DataType) error
	//	SetRequestType(ctx ServerContext, datatype string, collection bool, stream bool)
	//	SetResponseType(ctx ServerContext, stream bool)
	// SetDescription overwrites the description carried in metadata and reused in generated
	// skill and MCP tool descriptors. The server seeds the description from the object spec
	// before calling Describe, so a SetDescription made inside Describe wins.
	SetDescription(ServerContext, string)
	// SetComponent marks the service as a component. NOTHING READS THE FLAG: it is stored on the
	// service's info and surfaced by ServiceInfo.IsComponent, and no code in laatooserver or
	// laatoomodules calls IsComponent. Several platform plugins call SetComponent(ctx, true) in
	// Describe and it changes nothing about registration, injection or channel exposure.
	SetComponent(ServerContext, bool)
	// GetTags returns the service's tags. The agent/skill manager reads them to fill a skill
	// descriptor's Metadata.Tags and to answer skillHasTag (laatooserver core/skillmanager.go:207
	// and :255). The server-supplied implementation always returns nil -- the backing slice is
	// never assigned anywhere in the server and this interface has no AddTag/SetTags -- so a
	// service that needs tags must implement GetTags on its own struct.
	GetTags(ServerContext) []*Tag

	// GetNamespace returns the security namespace handed to SecurityHandler.AuthorizeService on
	// every request (laatooserver core/serviceelement.go:203). The server-supplied implementation
	// returns the constant "default" for every service (core/serviceinfo.go:103-105); it is not
	// read from configuration. Override it on the service struct to place a service in another
	// namespace.
	GetNamespace(ctx ServerContext) string
	// IsStreaming reports whether the service was configured with `streaming: true` (or the
	// legacy `stream: true`) in its service YAML. It is not settable through this interface. It
	// is checked at channel-wiring time: a streaming service bound to a non-streaming
	// ResponseHandler fails startup with Core_Bad_Conf
	// (laatooserver core/channelmanager.go:310-317).
	IsStreaming() bool

	// ServerElement returns the ServerElementService proxy wrapping this service, for code that
	// needs the service's name, context or owning module.
	ServerElement() ServerElement
	//ConfigureService(ctx ServerContext, requestType string, collection bool, stream bool, params []string, config []string, description string)
	//ConfigureService(ctx ServerContext, params []string, config []string, description string)
}

// Tag labels a service for discovery -- the skill manager matches agents to skills by tag.
// ParentTag allows a hierarchy; matching walks it.
type Tag struct {
	Name        string
	Description string
	ParentTag   *Tag
}

// UserInvokableService is a Service that can actually handle a request.
//
// The distinction is load-bearing and it FAILS SILENTLY. The server type-asserts each service to
// this interface and, when the assertion fails, only logs at Debug level and stores nil
// (laatooserver core/serviceelement.go:71-76); the invocation chain then skips the nil entry
// (core/serviceelement.go:319-321). A service that lacks an Invoke method -- or whose Invoke has
// a wrong signature -- therefore answers every call with HTTP 200 and an empty body, with no
// error logged above Debug and nothing failing at startup. Embedding core.Service does NOT supply
// an Invoke: the embedded interface's method set does not include one.
type UserInvokableService interface {
	Service
	// Invoke handles one request. Read parameters through ctx.Get*Param and publish the result
	// with ctx.SetResponse.
	//
	// Returning a non-nil error DISCARDS whatever response was set: the HTTP response handler
	// replaces it with 401 for Core_Error_Unauthorized and with 400 Bad Request for every other
	// error (laatooserver engine/http/responsehandler.go:24-32) -- there is no path to a 500.
	// Returning nil without setting a response yields 200 with no body.
	//
	// When the service runs behind middleware, a middleware that sets any response
	// short-circuits the chain and the service itself never runs
	// (laatooserver core/serviceelement.go:333-335).
	Invoke(RequestContext) error
}

// Param describes one declared request or response parameter.
//
// Two kinds of Param exist and only one carries data. Those reached through
// Service.RequestParameters or ServiceInfo.GetRequestInfo are DECLARATIONS whose GetValue is
// always nil; those reached through Request are per-request clones populated from the incoming
// call.
type Param interface {
	// GetName returns the parameter name as declared in the object spec's request.params, or as
	// passed to Service.AddParam.
	GetName() string
	// GetDescription returns the declared description. It is "" for params declared through
	// AddParams, which passes an empty description for all of them.
	GetDescription() string
	// IsCollection reports whether the param binds a list rather than a single value.
	IsCollection() bool
	// IsStream reports whether the param was declared with stream: true.
	IsStream() bool
	// IsRequired reports whether the param is mandatory. It is enforced while the request is
	// built: an absent required param fails with Core_Error_Missing_Arg before the service runs
	// (laatooserver core/request.go:185-188). A param that is present but whose value is nil
	// counts as absent.
	IsRequired() bool
	// GetDataType returns the declared data type. For datatypes.Object the concrete registered
	// object name is held internally and is not reachable through this interface.
	GetDataType() datatypes.DataType
	// GetValue returns the bound value, or nil on a declaration (see the interface comment).
	// A value is only bound for the data types the param decoder handles explicitly -- Stringmap,
	// Stringsmap, Config, Bytes, Int, Datetime, String, Files, Stringarr, Bool and Object
	// (laatooserver core/param.go:122-152). Any other declared type leaves the value unset and
	// fails the request with Core_Error_Bad_Arg.
	GetValue() interface{}
}

// ServiceFunc is a bare request handler: the shape a listener or generated route handler takes
// when there is no Service object behind it. The HTTP engine builds one per route, and the
// pub/sub module stores subscribers as ServiceFuncs.
type ServiceFunc func(ctx RequestContext) error

// Request is the parameter set for one service invocation, reached from a RequestContext. Every
// accessor takes a RequestContext that the server-side implementation ignores entirely
// (laatooserver core/request.go); pass the current one anyway.
//
// The map it reads is built by an ALLOWLIST pass over the service's DECLARED params
// (core/request.go:161-196): a value the channel supplies that the object spec does not declare
// is dropped and can never be read, and every declared param is inserted whether or not the
// caller supplied one.
//
// The consequence for callers is that the second return value means "declared", not "supplied".
// GetParam and GetParamValue return true for a declared-but-absent param, with a nil value.
type Request interface {
	// GetService ALWAYS returns nil. The server constructs the request with only its Params map
	// and never sets the service field (laatooserver core/requestcontext.go:158-160; the nil is
	// asserted by core/request_params_test.go:359), so a caller that chains off the result
	// panics. There is no supported way to reach the invoked Service through this interface.
	GetService(RequestContext) Service
	//GetBody() interface{}
	// GetParam returns the Param for name. ok is true for any DECLARED param, including one the
	// caller did not supply -- in which case Param.GetValue() is nil. It is false only for a name
	// the object spec does not declare.
	GetParam(RequestContext, string) (Param, bool)
	// GetParams returns the live per-request param map, not a copy.
	GetParams(RequestContext) map[string]Param
	// GetParamValue returns the param's bound value. It returns (nil, true) for a declared param
	// the caller did not supply, so test the value rather than the flag.
	GetParamValue(RequestContext, string) (interface{}, bool)
	// GetIntParam returns the value when it is an int. It returns -1, not 0, when the param is
	// missing or holds any other type -- so a returned -1 is indistinguishable from a literal -1
	// unless the flag is checked.
	GetIntParam(RequestContext, string) (int, bool)
	// GetStringParam returns ("", false) when the param is absent or does not hold a string.
	GetStringParam(RequestContext, string) (string, bool)
	// GetStringMapParam accepts a value stored either as utils.StringMap or as a plain
	// map[string]interface{}; anything else returns (nil, false).
	GetStringMapParam(RequestContext, string) (utils.StringMap, bool)
	// GetStringsMapParam accepts a value stored either as utils.StringsMap or as a plain
	// map[string]string. This is the accessor for a param declared `type: stringsmap`, which is
	// how a whole JSON body arrives when a channel wraps it with bodyparamname.
	GetStringsMapParam(RequestContext, string) (utils.StringsMap, bool)
}

// Response is what a service publishes with ctx.SetResponse and what a ResponseHandler renders.
type Response struct {
	// Status is a core.Status* constant. A nil Response is rendered as StatusSuccess
	// (laatooserver engine/http/responsehandler.go:35-37).
	Status int
	// Data is the payload, encoded with the channel's codec. Nil Data yields an empty body.
	Data interface{}
	// MetaInfo is emitted as response headers by the HTTP handlers: each key becomes a header and
	// the key names are listed in Access-Control-Expose-Headers.
	MetaInfo map[string]interface{}
	// Error is rendered only on the bad-request branch, where it becomes the JSON body
	// {"Error": <message>} (laatooserver engine/http/responsehandler.go:136-138). On a success
	// response it is ignored.
	Error error
	// Return is not read by any code in laatooserver.
	Return bool
}

// ResponseHandler renders a service's Response onto the transport. Each channel gets one: either
// the object named by the channel's response-handler configuration, or the engine's default.
type ResponseHandler interface {
	// Initialize is called once, at channel wiring, and ONLY for a handler named in channel
	// configuration; an engine's default handler is constructed ready to use and never sees this
	// call (laatooserver core/channelmanager.go:288-306). conf is that handler's config block.
	Initialize(ctx ServerContext, conf config.Config) error
	// HandleResponse renders ctx.GetResponse() onto the transport. It is called on every
	// non-streaming path, including early exits where the service never ran (bad request, failed
	// authentication). err is whatever the pipeline produced; when it is non-nil the HTTP
	// implementation DISCARDS the service's own response and substitutes 401 for
	// Core_Error_Unauthorized and 400 Bad Request for everything else
	// (laatooserver engine/http/responsehandler.go:24-32).
	HandleResponse(ctx RequestContext, err error) error
	// IsStreaming returns true if this handler supports streaming. It is checked at
	// channel-wiring time rather than per request: binding a streaming service to a
	// non-streaming handler fails startup with Core_Bad_Conf, while the reverse is only logged
	// (laatooserver core/channelmanager.go:310-327).
	IsStreaming() bool
	// HandleStream is called after execution to drain/handle the response stream. It runs
	// INSTEAD of HandleResponse, and only when both this handler and the request are streaming --
	// a streaming handler on a request that never entered streaming mode falls back to
	// HandleResponse. The HTTP engine runs it in its own goroutine under a timeout and abandons
	// it if it does not return (laatooserver engine/http/requesthandler.go:86-99), so it must
	// terminate on its own.
	HandleStream(ctx RequestContext) error
	// GetResponseStream returns a ResponseStream for this handler and context.
	// It is called exactly ONCE, while the RequestContext is being created, and only when
	// IsStreaming reports true (laatooserver core/servercontext.go:564-568) -- not per write.
	// Returning nil leaves the context with no stream.
	GetResponseStream(ctx RequestContext) ResponseStream
	// GetName returns the name of the response handler. It is used only in startup log and error
	// messages (laatooserver core/channelmanager.go:316 and :325); nothing looks a handler up by
	// it.
	GetName() string
}
