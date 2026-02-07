//go:build !plugin
// +build !plugin

package vst2

import (
	"testing"
)

func TestPluginOpcodeString(t *testing.T) {
	tests := []struct {
		opcode PluginOpcode
		want   string
	}{
		{plugOpen, "plugOpen"},
		{plugClose, "plugClose"},
		{plugSetProgram, "plugSetProgram"},
		{plugGetProgram, "plugGetProgram"},
		{plugSetProgramName, "plugSetProgramName"},
		{plugGetProgramName, "plugGetProgramName"},
		{plugGetParamLabel, "plugGetParamLabel"},
		{plugGetParamDisplay, "plugGetParamDisplay"},
		{plugGetParamName, "plugGetParamName"},
		{plugSetSampleRate, "plugSetSampleRate"},
		{plugSetBufferSize, "plugSetBufferSize"},
		{plugStateChanged, "plugStateChanged"},
		{plugEditorGetRect, "plugEditorGetRect"},
		{plugEditorOpen, "plugEditorOpen"},
		{plugEditorClose, "plugEditorClose"},
		{plugGetChunk, "plugGetChunk"},
		{plugSetChunk, "plugSetChunk"},
		{plugProcessEvents, "plugProcessEvents"},
		{plugCanBeAutomated, "plugCanBeAutomated"},
		{plugGetProgramNameIndexed, "plugGetProgramNameIndexed"},
		{plugGetPluginProperties, "plugGetPluginProperties"},
		{plugGetVstVersion, "plugGetVstVersion"},
		{plugCanDo, "plugCanDo"},
		{plugGetTailSize, "plugGetTailSize"},
		{plugGetParameterProperties, "plugGetParameterProperties"},
		{plugGetVendorString, "plugGetVendorString"},
		{plugGetProductString, "plugGetProductString"},
		{plugGetVendorVersion, "plugGetVendorVersion"},
		{plugVendorSpecific, "plugVendorSpecific"},
		{plugGetPluginName, "plugGetPluginName"},
		{plugSetSpeakerArrangement, "plugSetSpeakerArrangement"},
		{plugStartProcess, "plugStartProcess"},
		{plugStopProcess, "plugStopProcess"},
		{plugEditorIdle, "plugEditorIdle"},
		{plugGetInputProperties, "plugGetInputProperties"},
		{plugGetOutputProperties, "plugGetOutputProperties"},
		{plugGetPluginCategory, "plugGetPluginCategory"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.opcode.String(); got != tt.want {
				t.Errorf("PluginOpcode(%d).String() = %v, want %v", tt.opcode, got, tt.want)
			}
		})
	}
}

func TestHostOpcodeString(t *testing.T) {
	tests := []struct {
		opcode HostOpcode
		want   string
	}{
		{HostAutomate, "HostAutomate"},
		{HostVersion, "HostVersion"},
		{HostCurrentID, "HostCurrentID"},
		{HostIdle, "HostIdle"},
		{HostGetTime, "HostGetTime"},
		{HostProcessEvents, "HostProcessEvents"},
		{HostIOChanged, "HostIOChanged"},
		{HostSizeWindow, "HostSizeWindow"},
		{HostGetSampleRate, "HostGetSampleRate"},
		{HostGetBufferSize, "HostGetBufferSize"},
		{HostGetCurrentProcessLevel, "HostGetCurrentProcessLevel"},
		{HostGetAutomationState, "HostGetAutomationState"},
		{HostVendorSpecific, "HostVendorSpecific"},
		{HostCanDo, "HostCanDo"},
		{HostGetLanguage, "HostGetLanguage"},
		{HostUpdateDisplay, "HostUpdateDisplay"},
		{HostBeginEdit, "HostBeginEdit"},
		{HostEndEdit, "HostEndEdit"},
		{HostOpenFileSelector, "HostOpenFileSelector"},
		{HostCloseFileSelector, "HostCloseFileSelector"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.opcode.String(); got != tt.want {
				t.Errorf("HostOpcode(%d).String() = %v, want %v", tt.opcode, got, tt.want)
			}
		})
	}
}

func TestAsciiString(t *testing.T) {
	t.Run("ascii8", func(t *testing.T) {
		var s ascii8
		copy(s[:], "test")
		if got := s.String(); got != "test" {
			t.Errorf("ascii8.String() = %v, want %v", got, "test")
		}

		// Test with null terminator
		copy(s[:], "abc\x00xyz")
		if got := s.String(); got != "abc" {
			t.Errorf("ascii8.String() = %v, want %v", got, "abc")
		}
	})

	t.Run("ascii24", func(t *testing.T) {
		var s ascii24
		copy(s[:], "test string")
		if got := s.String(); got != "test string" {
			t.Errorf("ascii24.String() = %v, want %v", got, "test string")
		}
	})

	t.Run("ascii32", func(t *testing.T) {
		var s ascii32
		copy(s[:], "another test")
		if got := s.String(); got != "another test" {
			t.Errorf("ascii32.String() = %v, want %v", got, "another test")
		}
	})

	t.Run("ascii64", func(t *testing.T) {
		var s ascii64
		copy(s[:], "long test string")
		if got := s.String(); got != "long test string" {
			t.Errorf("ascii64.String() = %v, want %v", got, "long test string")
		}
	})
}
