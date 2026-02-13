//go:build !plugin
// +build !plugin

package vst2_test

import (
	"testing"

	"github.com/cwbudde/vst2"
	"pipelined.dev/signal"
)

func TestDemoPlugin(t *testing.T) {
	path := skipIfNoPlugin(t)
	v, err := vst2.Open(path)
	assertEqual(t, "vst error", err, nil)
	defer v.Close()

	t.Run("plugin metadata", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		// GetVSTVersion may return 0 for simple plugins that don't handle the opcode.
		_ = p.GetVSTVersion()
	})

	t.Run("process double", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		// EditorGetRect should not panic; demoplugin has no editor so nil is expected.
		_ = p.EditorGetRect()

		// EditorOpen, EditorClose, EditorIdle should not panic even with no editor
		p.EditorOpen(nil)
		p.EditorClose()
		p.EditorIdle()
	})

	t.Run("advanced operations", func(t *testing.T) {
		t.Parallel()
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		// Test UniqueID
		_ = p.UniqueID()

		// Test InitialDelay
		_ = p.InitialDelay()

		// Test Flags
		_ = p.Flags()

		// Test CanProcessFloat32/Float64
		_ = p.CanProcessFloat32()
		_ = p.CanProcessFloat64()

		// Test Suspend
		p.Suspend()

		// Test CanDo
		_ = p.CanDo("sendVstEvents")
		_ = p.CanDo("receiveVstEvents")

		// Test GetTailSize
		_ = p.GetTailSize()
	})

	t.Run("programs", func(t *testing.T) {
		t.Parallel()
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		// Test NumPrograms
		numProgs := p.NumPrograms()
		if numProgs < 0 {
			t.Fatalf("expected non-negative NumPrograms, got %d", numProgs)
		}

		// Test Program/SetProgram
		_ = p.Program()
		p.SetProgram(0)

		// Test CurrentProgramName/SetCurrentProgramName
		_ = p.CurrentProgramName()
		p.SetCurrentProgramName("Test Program")

		// Test ProgramName (if plugin has programs)
		if numProgs > 0 {
			_ = p.ProgramName(0)
		}
	})

	t.Run("param properties", func(t *testing.T) {
		t.Parallel()
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		numParams := p.NumParams()
		if numParams > 0 {
			// Test ParamProperties
			_, _ = p.ParamProperties(0)
		}
	})

	t.Run("chunk data", func(t *testing.T) {
		t.Parallel()
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		// Test GetProgramData/SetProgramData
		data := p.GetProgramData()
		if len(data) > 0 {
			p.SetProgramData(data)
		}

		// Test GetBankData/SetBankData
		bankData := p.GetBankData()
		if len(bankData) > 0 {
			p.SetBankData(bankData)
		}
	})

	t.Run("speaker arrangement", func(t *testing.T) {
		t.Parallel()
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		// SetSpeakerArrangement should not panic
		p.SetSpeakerArrangement(nil, nil)
	})

	t.Run("channel properties", func(t *testing.T) {
		t.Parallel()
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		numInputs := p.NumInputs()
		numOutputs := p.NumOutputs()

		// Test GetInputProperties
		if numInputs > 0 {
			_, _ = p.GetInputProperties(0)
		}

		// Test GetOutputProperties
		if numOutputs > 0 {
			_, _ = p.GetOutputProperties(0)
		}
	})
}

func TestHostCallback(t *testing.T) {
	path := skipIfNoPlugin(t)
	v, err := vst2.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer v.Close()

	t.Run("host callbacks", func(t *testing.T) {
		t.Parallel()
		var sampleRate signal.Frequency = 48000
		bufferSize := 512

		host := vst2.Host{
			GetSampleRate: func() signal.Frequency {
				return sampleRate
			},
			GetBufferSize: func() int {
				return bufferSize
			},
			GetProcessLevel: func() vst2.ProcessLevel {
				return vst2.ProcessLevelRealtime
			},
			UpdateDisplay: func() bool {
				return true
			},
		}

		p := v.Plugin(host.Callback())
		defer p.Close()

		// These calls will trigger the host callbacks
		p.Start()
		p.SetSampleRate(sampleRate)
		p.SetBufferSize(bufferSize)
		p.Resume()
	})

	t.Run("send events", func(t *testing.T) {
		t.Parallel()
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		// Create MIDI events
		event1 := &vst2.MIDIEvent{
			DeltaFrames: 0,
			Data:        [3]byte{0x90, 0x3C, 0x40}, // Note On
		}

		event2 := &vst2.MIDIEvent{
			DeltaFrames: 100,
			Data:        [3]byte{0x80, 0x3C, 0x40}, // Note Off
		}

		events := vst2.Events(event1, event2)
		defer events.Free()

		// Send events to plugin
		p.SendEvents(events)
	})

	t.Run("program name by index", func(t *testing.T) {
		t.Parallel()
		p := v.Plugin(vst2.NoopHostCallback())
		defer p.Close()

		numProgs := p.NumPrograms()
		if numProgs > 0 {
			_ = p.ProgramName(0)
		}
	})
}
