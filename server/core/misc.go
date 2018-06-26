package core

type MessageListener func(ctx RequestContext, message interface{}, info map[string]interface{}) error
