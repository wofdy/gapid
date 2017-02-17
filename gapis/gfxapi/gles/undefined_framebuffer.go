// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gles

import (
	"github.com/google/gapid/core/log"
	"github.com/google/gapid/core/os/device"
	"github.com/google/gapid/gapis/atom"
	"github.com/google/gapid/gapis/atom/transform"
	"github.com/google/gapid/gapis/gfxapi"
	"github.com/google/gapid/gapis/memory"
)

// undefinedFramebuffer adds a transform that will render a pattern into the
// color buffer at the end of each frame.
func undefinedFramebuffer(ctx log.Context, device *device.Instance) transform.Transformer {
	seenSurfaces := make(map[EGLSurface]bool)
	return transform.Transform("DirtyFramebuffer", func(ctx log.Context, i atom.ID, a atom.Atom, out transform.Writer) {
		out.MutateAndWrite(ctx, i, a)
		s := out.State()
		c := GetState(s).getContext()
		if c == nil {
			return // We can't do anything without a context.
		}
		if eglMakeCurrent, ok := a.(*EglMakeCurrent); ok && !seenSurfaces[eglMakeCurrent.Draw] {
			// Render the undefined pattern for new contexts.
			drawUndefinedFramebuffer(ctx, a, device, s, c, out)
			seenSurfaces[eglMakeCurrent.Draw] = true
		}
		if a.AtomFlags().IsEndOfFrame() {
			if c != nil && !c.Info.PreserveBuffersOnSwap {
				drawUndefinedFramebuffer(ctx, a, device, s, c, out)
			}
		}
	})
}

func drawUndefinedFramebuffer(ctx log.Context, a atom.Atom, device *device.Instance, s *gfxapi.State, c *Context, out transform.Writer) error {
	const (
		aScreenCoordsLocation AttributeLocation = 0

		vertexShaderSource string = `
					precision highp float;
					attribute vec2 aScreenCoords;
					varying vec2 uv;

					void main() {
						uv = aScreenCoords;
						gl_Position = vec4(aScreenCoords.xy, 0., 1.);
					}`
		fragmentShaderSource string = `
					precision highp float;
					varying vec2 uv;

					float F(float a) { return smoothstep(0.0, 0.1, a) * smoothstep(0.4, 0.3, a); }

					void main() {
						vec2 v = uv * 5.0;
						gl_FragColor = vec4(0.8, 0.9, 0.6, 1.0) * F(fract(v.x + v.y));
					}`
	)

	version, _ := ParseVersion(device.Configuration.Drivers.OpenGL.Version)

	var (
		origProgramID     = c.BoundProgram
		origArrayBufferID = c.BoundBuffers.ArrayBuffer
	)

	// Generate new unused object IDs.
	programID := ProgramId(newUnusedID('P', func(x uint32) bool { _, ok := c.Instances.Programs[ProgramId(x)]; return ok }))
	vertexShaderID := ShaderId(newUnusedID('S', func(x uint32) bool { _, ok := c.Instances.Shaders[ShaderId(x)]; return ok }))
	fragmentShaderID := ShaderId(newUnusedID('S', func(x uint32) bool {
		_, ok := c.Instances.Shaders[ShaderId(x)]
		return ok || ShaderId(x) == vertexShaderID
	}))
	bufferID := BufferId(newUnusedID('B', func(x uint32) bool { _, ok := c.Instances.Buffers[BufferId(x)]; return ok }))
	arrayID := VertexArrayId(newUnusedID('V', func(x uint32) bool { _, ok := c.Instances.VertexArrays[VertexArrayId(x)]; return ok }))

	// 2D vertices positions for a full screen 2D triangle strip.
	positions := []float32{-1., -1., 1., -1., -1., 1., 1., 1.}

	t := tweaker{ctx: ctx, out: out, c: GetContext(s)}

	// Temporarily change rasterizing/blending state and enable VAP 0.
	t.glDisable(GLenum_GL_BLEND)
	t.glDisable(GLenum_GL_CULL_FACE)
	t.glDisable(GLenum_GL_DEPTH_TEST)
	t.glDisable(GLenum_GL_SCISSOR_TEST)
	t.glDisable(GLenum_GL_STENCIL_TEST)
	t.bindOrSaveVertexArray(version, arrayID, aScreenCoordsLocation)
	out.MutateAndWrite(ctx, atom.NoID, NewGlEnableVertexAttribArray(aScreenCoordsLocation))

	// Create the shader program
	for _, a := range BuildProgram(ctx, s, vertexShaderID, fragmentShaderID, programID, vertexShaderSource, fragmentShaderSource) {
		out.MutateAndWrite(ctx, atom.NoID, a)
	}

	out.MutateAndWrite(ctx, atom.NoID, NewGlBindAttribLocation(programID, aScreenCoordsLocation, "aScreenCoords"))
	out.MutateAndWrite(ctx, atom.NoID, NewGlLinkProgram(programID))
	out.MutateAndWrite(ctx, atom.NoID, NewGlUseProgram(programID))

	tmp := atom.Must(atom.AllocData(ctx, s, bufferID))
	out.MutateAndWrite(ctx, atom.NoID, NewGlGenBuffers(1, tmp.Ptr()).AddWrite(tmp.Data()))

	out.MutateAndWrite(ctx, atom.NoID, NewGlBindBuffer(GLenum_GL_ARRAY_BUFFER, bufferID))

	tmp = atom.Must(atom.AllocData(ctx, s, positions))
	out.MutateAndWrite(ctx, atom.NoID, NewGlBufferData(GLenum_GL_ARRAY_BUFFER, GLsizeiptr(4*len(positions)), tmp.Ptr(), GLenum_GL_STATIC_DRAW).
		AddRead(tmp.Data()))

	out.MutateAndWrite(ctx, atom.NoID, NewGlVertexAttribPointer(aScreenCoordsLocation, 2, GLenum_GL_FLOAT, GLboolean(0), 0, memory.Nullptr))
	out.MutateAndWrite(ctx, atom.NoID, NewGlDrawArrays(GLenum_GL_TRIANGLE_STRIP, 0, 4))

	t.revert()

	// Restore buffer state.
	out.MutateAndWrite(ctx, atom.NoID, NewGlBindBuffer(GLenum_GL_ARRAY_BUFFER, origArrayBufferID))
	tmp = atom.Must(atom.AllocData(ctx, s, bufferID))
	out.MutateAndWrite(ctx, atom.NoID, NewGlDeleteBuffers(1, tmp.Ptr()).AddRead(tmp.Data()))

	// Restore program state.
	out.MutateAndWrite(ctx, atom.NoID, NewGlUseProgram(origProgramID))
	out.MutateAndWrite(ctx, atom.NoID, NewGlDeleteProgram(programID))
	out.MutateAndWrite(ctx, atom.NoID, NewGlDeleteShader(vertexShaderID))
	out.MutateAndWrite(ctx, atom.NoID, NewGlDeleteShader(fragmentShaderID))

	return nil
}
