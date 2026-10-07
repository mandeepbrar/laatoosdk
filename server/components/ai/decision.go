package ai

import (
	"laatoo.io/sdk/server/core"
)

// Decisions are typed judgements answered by a decision model -- a model that takes one state and a
// set of named, typed questions and returns a probability distribution per question instead of
// generating text. They are made through the AgentManager (Evaluate, Classify, Score, Validate,
// Route) and answered by a DecisionProvider, never by the chat LLM path.
//
// KNOWLEDGE IS THE CALLER'S. The caller supplies its own state, carrying whatever knowledge the
// decision needs -- gathered with the tools of the vocabulary it works in -- and, optionally, a
// subject. The server adds what only it can reach, the session: memory items found for the subject
// and the session's conversation, so every provider receives the same context:
//
//	answer, err := agentMgr.Classify(ctx, &ai.DecisionRequest{
//		State:   map[string]any{"message": "I was charged twice", "order": orderFacts},
//		Subject: &ai.DecisionSubject{ID: "order 1042"},
//		MinConfidence: 0.8,
//	}, "Which team should handle this?", []ai.DecisionOption{
//		{Key: "billing", Description: "Charges and refunds"},
//		{Key: "technical", Description: "Bugs and outages"},
//	})
//	if !answer.Decided { /* escalate: an LLM, a human, a default -- the caller's choice */ }

// DecisionQuestionType is the shape of a question's answer.
type DecisionQuestionType string

const (
	// DecisionBoolean is a yes/no question, answered with the probability of yes.
	DecisionBoolean DecisionQuestionType = "boolean"
	// DecisionChoice is a question with an unordered set of options, answered with the chosen
	// option and a probability for every option.
	DecisionChoice DecisionQuestionType = "choice"
	// DecisionScore is a question over ordered levels (the server bounds how many), answered with the probability-weighted
	// level and the distribution over levels.
	DecisionScore DecisionQuestionType = "score"
)

// DecisionOption is one option of a choice question or one level of a score question. Key is what
// an answer names; Description tells the model what the option means and is required -- a bare key
// gives the model nothing to judge against.
type DecisionOption struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

// DecisionQuestion is one named question of a decision request.
//
// Options holds a choice question's options, or a score question's levels ordered lowest first. A
// boolean question may carry exactly two options keyed "true" and "false" describing what each
// answer means; otherwise it carries none.
//
// MinConfidence and MinProbability override the request's defaults for this question when set
// (greater than zero). See DecisionAnswer.Decided for how they apply.
type DecisionQuestion struct {
	Name           string               `json:"name"`
	Type           DecisionQuestionType `json:"type"`
	Instructions   string               `json:"instructions"`
	Options        []DecisionOption     `json:"options,omitempty"`
	MinConfidence  float64              `json:"minConfidence,omitempty"`
	MinProbability float64              `json:"minProbability,omitempty"`
}

// DecisionMemorySource names one memory bank to search for the subject: its type (Reference or
// Data -- the banks that search by similarity; a Session bank ignores its query) and its id. Limit
// bounds the items taken from it; zero means the server's default.
type DecisionMemorySource struct {
	Type  MemoryType `json:"type"`
	ID    string     `json:"id"`
	Limit int        `json:"limit,omitempty"`
}

// DecisionSubject is what a decision is about, for the session memory the server searches.
//
//   - ID identifies the subject; it is the query text of the memory search.
//   - Memory names the banks to search. Nil means the session's own Reference and Data banks (a bank
//     of each type whose id is the session id), where they exist.
//
// The server reads no ontology for a decision: knowledge about the subject is the caller's to gather
// and pass in DecisionRequest.State.
type DecisionSubject struct {
	ID     string                 `json:"id"`
	Memory []DecisionMemorySource `json:"memory,omitempty"`
}

// DecisionStateSection names one section of the state the model receives.
type DecisionStateSection string

const (
	// DecisionStateCaller is the caller's own state. It is always included and never trimmed.
	DecisionStateCaller DecisionStateSection = "state"
	// DecisionStateMemory is the memory items found for the subject.
	DecisionStateMemory DecisionStateSection = "memory"
	// DecisionStateConversation is the session's earlier turns, bounded by agents.conversationturns.
	DecisionStateConversation DecisionStateSection = "conversation"
)

// DecisionRequest is a decision call: the caller's state, the questions, and what the server should
// inject.
//
// Model names the decision model, bare or <provider>::<model>; empty means agents.decisionmodel from
// the nearest namespace that sets it. Exclude switches injected sections off for this call -- a
// narrow judgement that must see only its own state excludes all three. MinConfidence and
// MinProbability are the defaults for every question that sets none of its own.
//
// For Classify, Score, Validate and Route the question comes from the method's arguments and
// Questions must be empty.
type DecisionRequest struct {
	Model          string                 `json:"model,omitempty"`
	State          any                    `json:"state,omitempty"`
	Questions      []DecisionQuestion     `json:"questions,omitempty"`
	Subject        *DecisionSubject       `json:"subject,omitempty"`
	Exclude        []DecisionStateSection `json:"exclude,omitempty"`
	MinConfidence  float64                `json:"minConfidence,omitempty"`
	MinProbability float64                `json:"minProbability,omitempty"`
}

// DecisionAnswer is the answer to one question.
//
//   - A choice answer sets Choice, Probability (of the chosen option) and Distribution (every option).
//   - A score answer sets Choice (the most probable level's key), Score (the probability-weighted
//     level, 0 for the lowest) and Distribution (every level).
//   - A boolean answer sets Probability (of yes) and Choice ("true" or "false", the likelier).
//
// Confidence is the provider's measure of how concentrated the distribution is, 0 to 1; a provider
// that reports none leaves it 0.
//
// Decided is set by the server, never by the provider. A choice or score answer is undecided when its
// Confidence is below the question's minimum confidence, or its Probability below the minimum
// probability. A boolean answer is undecided when its Probability lies strictly between
// 1-minProbability and minProbability. With no threshold set, every answer is decided. An undecided
// answer still carries its distribution; the server never acts on it.
type DecisionAnswer struct {
	Name         string               `json:"name"`
	Type         DecisionQuestionType `json:"type"`
	Choice       string               `json:"choice,omitempty"`
	Score        float64              `json:"score,omitempty"`
	Probability  float64              `json:"probability,omitempty"`
	Distribution map[string]float64   `json:"distribution,omitempty"`
	Confidence   float64              `json:"confidence,omitempty"`
	Decided      bool                 `json:"decided"`
}

// DecisionSectionReport says what one section of the state carried: how many items it held, how many
// the state budget dropped, and a note when it is empty for a reason worth reading (no session, no
// such bank, nothing about the subject).
type DecisionSectionReport struct {
	Section DecisionStateSection `json:"section"`
	Items   int                  `json:"items"`
	Dropped int                  `json:"dropped,omitempty"`
	Note    string               `json:"note,omitempty"`
}

// DecisionResponse is the result of a decision call: one answer per question, keyed by question
// name, with the model that answered, its usage, and what each state section carried.
type DecisionResponse struct {
	Answers   map[string]*DecisionAnswer `json:"answers"`
	Model     string                     `json:"model,omitempty"`
	Provider  string                     `json:"provider,omitempty"`
	Tokens    TokenUsage                 `json:"tokens"`
	Cost      Cost                       `json:"cost"`
	RequestID string                     `json:"requestId,omitempty"`
	Sections  []DecisionSectionReport    `json:"sections,omitempty"`
}

// DecisionEvaluation is what a DecisionProvider is asked: the provider's own model name (never the
// qualified reference), the assembled state as JSON text, and the questions.
type DecisionEvaluation struct {
	Model     string
	State     string
	Questions []DecisionQuestion
}

// DecisionProvider is the contract a decision-model plugin implements and registers with the
// AgentManager (RegisterDecisionProvider). It is deliberately separate from LLMProvider: a decision
// model has no completion, streaming or token counting, and a chat model has no typed questions.
//
// The provider never sees the subject or the memory; the server assembles them, with the caller's
// state, into one state, so every provider receives the same context.
type DecisionProvider interface {
	// Name is the provider's registration name, e.g. "typesafe".
	Name() string
	// ListModels names the models this provider serves. They are indexed once, at registration, and
	// each must carry ModelCapabilities.SupportsDecisions in its GetConfig: registration is refused
	// when one does not.
	ListModels(ctx core.ServerContext) ([]string, error)
	// GetConfig answers a listed model's configuration -- its context, prices and capabilities -- or
	// an error for a model the provider does not serve. The server reads SupportsDecisions from it
	// when the provider registers.
	GetConfig(ctx core.ServerContext, model string) (*ModelConfig, error)
	// StateLimit is the largest state, in characters, the model accepts. The server trims injected
	// context to fit and refuses a call whose caller state alone exceeds it. Zero means no limit.
	StateLimit(ctx core.ServerContext, model string) (int, error)
	// Evaluate answers every question against the state in one call and returns one answer per
	// question, keyed by question name, with usage. It leaves DecisionAnswer.Decided unset; the
	// server computes it.
	Evaluate(ctx core.RequestContext, req *DecisionEvaluation) (*DecisionResponse, error)
}
