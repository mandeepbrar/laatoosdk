package core

import (
	"context"
)

// StreamChunk represents a single chunk in a streaming response.
type StreamChunk struct {
	Status   int
	Data     interface{}
	MetaInfo map[string]interface{}
	Error    error
	Final    bool // true when this is the last chunk
}

// ResponseStream is the producer-side interface a service uses to push data.
type ResponseStream struct {
	chunks chan *StreamChunk
	done   chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
}

func NewResponseStream(ctx context.Context, bufferSize int) *ResponseStream {
	streamCtx, cancel := context.WithCancel(ctx)
	return &ResponseStream{
		chunks: make(chan *StreamChunk, bufferSize),
		done:   make(chan struct{}),
		ctx:    streamCtx,
		cancel: cancel,
	}
}

// Send sends a chunk to the stream. Blocks if buffer is full.
func (rs *ResponseStream) Send(chunk *StreamChunk) error {
	select {
	case <-rs.ctx.Done():
		return rs.ctx.Err()
	case rs.chunks <- chunk:
		return nil
	}
}

// Close signals the end of the stream.
func (rs *ResponseStream) Close() {
	close(rs.chunks)
}

// Cancel aborts the stream from the consumer side.
func (rs *ResponseStream) Cancel() {
	rs.cancel()
}

// Chunks returns the read-only channel for consumers.
func (rs *ResponseStream) Chunks() <-chan *StreamChunk {
	return rs.chunks
}
