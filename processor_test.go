//go:build !plugin
// +build !plugin

package vst2_test

import (
	"context"
	"testing"

	"github.com/cwbudde/vst2"
	"pipelined.dev/pipe"
	"pipelined.dev/pipe/mutable"
	"pipelined.dev/signal"
)

func TestProcessor(t *testing.T) {
	path := skipIfNoPlugin(t)
	v, err := vst2.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer v.Close()

	const (
		channels   = 2
		bufferSize = 64
		sampleRate = 44100
	)

	t.Run("allocator", func(t *testing.T) {
		t.Parallel()
		progressCalled := false
		progressFn := func(frames int) {
			progressCalled = true
			if frames != bufferSize {
				t.Errorf("expected progress with %d frames, got %d", bufferSize, frames)
			}
		}

		processor := v.Processor(vst2.Host{}, progressFn)

		// Test with init function
		initCalled := false
		initFn := func(p *vst2.Plugin) {
			initCalled = true
			p.SetParamValue(0, 0.5)
		}

		allocator := processor.Allocator(initFn)

		mctx := mutable.Context{}
		props := pipe.SignalProperties{
			Channels:   channels,
			SampleRate: sampleRate,
		}

		p, err := allocator(mctx, bufferSize, props)
		if err != nil {
			t.Fatalf("allocator failed: %v", err)
		}

		if !initCalled {
			t.Error("init function was not called")
		}

		// Test Start
		if p.StartFunc != nil {
			if err := p.StartFunc(context.Background()); err != nil {
				t.Errorf("StartFunc failed: %v", err)
			}
		}

		// Test Process
		if p.ProcessFunc != nil {
			in := signal.Allocator{
				Channels: channels,
				Length:   bufferSize,
				Capacity: bufferSize,
			}.Float64()
			out := signal.Allocator{
				Channels: channels,
				Length:   bufferSize,
				Capacity: bufferSize,
			}.Float64()

			// Fill input with test data
			for c := 0; c < channels; c++ {
				for i := 0; i < bufferSize; i++ {
					in.SetSample(in.BufferIndex(c, i), 0.5)
				}
			}

			processed, err := p.ProcessFunc(in, out)
			if err != nil {
				t.Errorf("ProcessFunc failed: %v", err)
			}
			if processed != bufferSize {
				t.Errorf("expected %d processed frames, got %d", bufferSize, processed)
			}

			if !progressCalled {
				t.Error("progress function was not called")
			}
		}

		// Test Flush
		if p.FlushFunc != nil {
			if err := p.FlushFunc(context.Background()); err != nil {
				t.Errorf("FlushFunc failed: %v", err)
			}
		}
	})

	t.Run("allocator without init", func(t *testing.T) {
		t.Parallel()
		processor := v.Processor(vst2.Host{}, nil)
		allocator := processor.Allocator(nil)

		mctx := mutable.Context{}
		props := pipe.SignalProperties{
			Channels:   channels,
			SampleRate: sampleRate,
		}

		p, err := allocator(mctx, bufferSize, props)
		if err != nil {
			t.Fatalf("allocator failed: %v", err)
		}

		// Verify properties
		if p.Channels != channels {
			t.Errorf("expected %d channels, got %d", channels, p.Channels)
		}
		if p.SampleRate != sampleRate {
			t.Errorf("expected sample rate %v, got %v", sampleRate, p.SampleRate)
		}
	})
}
