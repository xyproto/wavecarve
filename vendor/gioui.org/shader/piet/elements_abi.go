// Code generated by gioui.org/cpu/cmd/compile DO NOT EDIT.

//go:build linux && (arm64 || arm || amd64)
// +build linux
// +build arm64 arm amd64

package piet

import "gioui.org/cpu"
import "unsafe"

/*
#cgo LDFLAGS: -lm

#include <stdint.h>
#include <stdlib.h>
#include "abi.h"
#include "runtime.h"
#include "elements_abi.h"
*/
import "C"

var ElementsProgramInfo = (*cpu.ProgramInfo)(unsafe.Pointer(&C.elements_program_info))

type ElementsDescriptorSetLayout = C.struct_elements_descriptor_set_layout

const ElementsHash = "0f18de15866045b36217068789c9c8715a63e0f9f120c53ea2d4d76f53e443c3"

func (l *ElementsDescriptorSetLayout) Binding0() *cpu.BufferDescriptor {
	return (*cpu.BufferDescriptor)(unsafe.Pointer(&l.binding0))
}

func (l *ElementsDescriptorSetLayout) Binding1() *cpu.BufferDescriptor {
	return (*cpu.BufferDescriptor)(unsafe.Pointer(&l.binding1))
}

func (l *ElementsDescriptorSetLayout) Binding2() *cpu.BufferDescriptor {
	return (*cpu.BufferDescriptor)(unsafe.Pointer(&l.binding2))
}

func (l *ElementsDescriptorSetLayout) Binding3() *cpu.BufferDescriptor {
	return (*cpu.BufferDescriptor)(unsafe.Pointer(&l.binding3))
}
