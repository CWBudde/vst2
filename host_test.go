//go:build !plugin
// +build !plugin

package vst2_test

import (
	"testing"

	"pipelined.dev/audio/vst2"
	"pipelined.dev/signal"
)

func TestDemoPlugin(t *testing.T) {
	path := skipIfNoPlugin(t)
	v, err := vst2.Open(path)
	assertEqual(t, "vst error", err, nil)
	defer v.Close()

	t.Run("plugin metadata", func(t *testing.T) {
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		name := p.GetPluginName()
		if name == "" {
			t.Fatal("expected non-empty plugin name")
		}
		vendor := p.GetVendorString()
		if vendor == "" {
			t.Fatal("expected non-empty vendor string")
		}
		product := p.GetProductString()
		if product == "" {
			t.Fatal("expected non-empty product string")
		}
		cat := p.GetCategory()
		if cat == 0 {
			t.Fatal("expected non-zero category")
		}
		// GetVendorVersion may return 0 for simple plugins.
		_ = p.GetVendorVersion()
	})

	t.Run("channel info", func(t *testing.T) {
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		inputs := p.NumInputs()
		if inputs <= 0 {
			t.Fatalf("expected positive NumInputs, got %d", inputs)
		}
		outputs := p.NumOutputs()
		if outputs <= 0 {
			t.Fatalf("expected positive NumOutputs, got %d", outputs)
		}
	})

	t.Run("get params", func(t *testing.T) {
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		numParams := p.NumParams()
		assertEqual(t, "num params", numParams > 0, true)
		for i := 0; i < numParams; i++ {
			p.ParamName(i)
			p.ParamUnitName(i)
			p.ParamValueName(i)
		}
	})

	t.Run("set param", func(t *testing.T) {
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		original := p.ParamValue(0)
		p.SetParamValue(0, 0.75)
		updated := p.ParamValue(0)
		if updated == original {
			t.Fatalf("param value did not change: got %v", updated)
		}
	})

	t.Run("vst version", func(t *testing.T) {
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		// GetVSTVersion may return 0 for simple plugins that don't handle the opcode.
		_ = p.GetVSTVersion()
	})

	t.Run("process double", func(t *testing.T) {
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		const (
			channels = 2
			frames   = 64
		)
		p.Start()
		p.SetSampleRate(44100)
		p.SetBufferSize(frames)
		p.Resume()

		in := vst2.NewDoubleBuffer(channels, frames)
		defer in.Free()
		out := vst2.NewDoubleBuffer(channels, frames)
		defer out.Free()

		// Fill input with a non-zero signal.
		inputData := make([][]float64, channels)
		for c := range inputData {
			inputData[c] = make([]float64, frames)
			for i := range inputData[c] {
				inputData[c][i] = 0.5
			}
		}
		buf := signal.Allocator{
			Channels: channels,
			Length:   frames,
			Capacity: frames,
		}.Float64()
		signal.WriteStripedFloat64(inputData, buf)
		in.Write(buf)

		p.ProcessDouble(in, out)
		out.Read(buf)

		// Verify output has non-zero samples.
		nonZero := false
		for i := 0; i < frames; i++ {
			if buf.Sample(buf.BufferIndex(0, i)) != 0 {
				nonZero = true
				break
			}
		}
		if !nonZero {
			t.Fatal("expected non-zero output from ProcessDouble")
		}
	})

	t.Run("process float", func(t *testing.T) {
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		const (
			channels = 2
			frames   = 64
		)
		p.Start()
		p.SetSampleRate(44100)
		p.SetBufferSize(frames)
		p.Resume()

		in := vst2.NewFloatBuffer(channels, frames)
		defer in.Free()
		out := vst2.NewFloatBuffer(channels, frames)
		defer out.Free()

		// Fill input with a non-zero signal.
		inputData := make([][]float64, channels)
		for c := range inputData {
			inputData[c] = make([]float64, frames)
			for i := range inputData[c] {
				inputData[c][i] = 0.5
			}
		}
		buf := signal.Allocator{
			Channels: channels,
			Length:   frames,
			Capacity: frames,
		}.Float64()
		signal.WriteStripedFloat64(inputData, buf)
		in.Write(buf)

		p.ProcessFloat(in, out)
		out.Read(buf)

		// Verify output has non-zero samples.
		nonZero := false
		for i := 0; i < frames; i++ {
			if buf.Sample(buf.BufferIndex(0, i)) != 0 {
				nonZero = true
				break
			}
		}
		if !nonZero {
			t.Fatal("expected non-zero output from ProcessFloat")
		}
	})

	t.Run("editor", func(t *testing.T) {
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		// EditorGetRect should not panic; demoplugin has no editor so nil is expected.
		_ = p.EditorGetRect()
	})
}
