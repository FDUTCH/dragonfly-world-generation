package resourcecompiler

import "testing"

func TestComplications(t *testing.T) {
	SetPaths("/home/ikarolyi/BDS/definitions", "/home/ikarolyi/BDS/behavior_packs/vanilla/structures")
	CompileAll()
}