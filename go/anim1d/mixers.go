// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// mixers is all the patterns that are constructions of other patterns.

package anim1d

// Gradient does a gradient between 2 patterns.
//
// A good example is using two colors but it can also be animations.
//
// TODO(maruel): Support N colors at M positions.
type Gradient struct {
	Left  SPattern
	Right SPattern
	Curve Curve
	buf   Frame
}

func (g *Gradient) NextFrame(pixels Frame, timeMS uint32) {
	if g.Left.Pattern == nil || g.Right.Pattern == nil {
		return
	}
	l := len(pixels) - 1
	g.buf.reset(len(pixels))
	g.Left.NextFrame(pixels, timeMS)
	g.Right.NextFrame(g.buf, timeMS)
	if l == 0 {
		pixels.Mix(g.buf, g.Curve.Scale8(65535>>1))
	} else {
		max := len(pixels) - 1
		for i := range pixels {
			intensity := uint16(i * 65535 / max)
			pixels[i].Mix(g.buf[i], g.Curve.Scale8(intensity))
		}
	}
}

// Transition changes from Before to After over time. It doesn't repeat.
//
// In gets timeMS that is subtracted by OffsetMS.
type Transition struct {
	Before     SPattern // Old pattern that is disappearing
	After      SPattern // New pattern to show
	OffsetMS   uint32   // Offset at which the transiton from Before->In starts
	DurationMS uint32   // Duration of the transition while both are rendered
	Curve      Curve    // Type of transition, defaults to EaseOut if not set
	buf        Frame
}

func (t *Transition) NextFrame(pixels Frame, timeMS uint32) {
	if timeMS <= t.OffsetMS {
		// Before transition.
		if t.Before.Pattern != nil {
			t.Before.NextFrame(pixels, timeMS)
		}
		return
	}
	if t.After.Pattern != nil {
		t.After.NextFrame(pixels, timeMS-t.OffsetMS)
	}
	if timeMS >= t.OffsetMS+t.DurationMS {
		// After transition.
		t.buf = nil
		return
	}
	t.buf.reset(len(pixels))

	// TODO(maruel): Add lateral animation and others.
	if t.Before.Pattern != nil {
		t.Before.NextFrame(t.buf, timeMS)
	}
	intensity := uint16((timeMS - t.OffsetMS) * 65535 / (t.DurationMS))
	pixels.Mix(t.buf, 255.-t.Curve.Scale8(intensity))
}

// Cycle cycles between multiple patterns. It can be used as an animatable
// looping frame.
//
// TODO(maruel): Blend between frames with Curve, defaults to step.
// TODO(maruel): Merge with Loop.
type Cycle struct {
	Frames          []SPattern
	FrameDurationMS uint32
}

func (c *Cycle) NextFrame(pixels Frame, timeMS uint32) {
	if len(c.Frames) == 0 {
		return
	}
	c.Frames[int(timeMS/c.FrameDurationMS)%len(c.Frames)].NextFrame(pixels, timeMS)
}

// Loop rotates between all the animations.
//
// Display starts with one DurationShow for Patterns[0], then starts looping.
// timeMS is not modified so it's like as all animations continued animating
// behind.
// TODO(maruel): Add lateral transition and others.
type Loop struct {
	Patterns             []SPattern
	DurationShowMS       uint32 // Duration for each pattern to be shown as pure
	DurationTransitionMS uint32 // Duration of the transition between two patterns
	Curve                Curve  // Type of transition, defaults to EaseOut if not set
	buf                  Frame
}

func (l *Loop) NextFrame(pixels Frame, timeMS uint32) {
	lp := uint32(len(l.Patterns))
	if lp == 0 {
		return
	}
	cycleDuration := l.DurationShowMS + l.DurationTransitionMS
	if cycleDuration == 0 {
		l.Patterns[0].NextFrame(pixels, timeMS)
		return
	}
	l.buf.reset(len(pixels))
	cycleNumber := timeMS / cycleDuration
	a := l.Patterns[cycleNumber%lp]
	a.NextFrame(pixels, timeMS)
	cycleOffset := cycleNumber * cycleDuration
	if cycleOffset <= l.DurationShowMS {
		return
	}
	b := l.Patterns[(cycleNumber+1)%lp]
	b.NextFrame(l.buf, timeMS)
	intensity := uint16((cycleOffset - l.DurationShowMS) * 65535 / l.DurationTransitionMS)
	pixels.Mix(l.buf, l.Curve.Scale8(intensity))
}

// Rotate rotates a pattern that can also cycle either way.
//
// Use negative to go left. Can be used for 'candy bar'.
//
// Similar to PingPong{} except that it doesn't bounce.
//
// Use 5x oversampling with Scale{} to create smoother animation.
type Rotate struct {
	Child       SPattern
	MovesPerSec float32 // Expressed in number of light jumps per second.
	buf         Frame
}

func (r *Rotate) NextFrame(pixels Frame, timeMS uint32) {
	l := len(pixels)
	if l == 0 || r.Child.Pattern == nil {
		return
	}
	r.buf.reset(l)
	r.Child.NextFrame(r.buf, timeMS)
	offset := int(float32(timeMS)*0.001*r.MovesPerSec) % l
	if offset < 0 {
		offset = l + offset
	}
	copy(pixels[offset:], r.buf)
	copy(pixels[:offset], r.buf[l-offset:])
}

// Chronometer moves 3 lights to the right, each indicating second, minute, and
// hour passed since the start.
type Chronometer struct {
	Child SPattern
	buf   Frame
}

func (r *Chronometer) NextFrame(pixels Frame, timeMS uint32) {
	l := uint32(len(pixels))
	if l == 0 || r.Child.Pattern == nil {
		return
	}
	r.buf.reset(4)
	r.Child.NextFrame(r.buf, timeMS)

	seconds := timeMS / 1000
	mins := seconds / 60
	hours := mins / 60

	secPos := (l*(seconds%60) + 30) / 60
	minPos := (l*(mins%60) + 30) / 60
	hourPos := hours % l

	for i := range pixels {
		switch uint32(i) {
		case secPos:
			pixels[i] = r.buf[1]
		case minPos:
			pixels[i] = r.buf[2]
		case hourPos:
			pixels[i] = r.buf[3]
		default:
			pixels[i] = r.buf[0]
		}
	}
}

// PingPong shows a 'ball' with a trail that bounces from one side to
// the other.
//
// Can be used for a ball, a water wave or K2000 (Knight Rider) style light.
// The trail can be a Frame or a dynamic pattern.
//
// To get smoothed movement, use Scale{} with a 5x factor or so.
// TODO(maruel): That's a bit inefficient, enable Interpolation here.
type PingPong struct {
	Child SPattern // [0] is the front pixel so the pixels are effectively drawn in reverse order.
	// TODO(maruel): PerMoveMS uint32
	MovesPerSec float32 // Expressed in number of light jumps per second.
	buf         Frame
}

func (p *PingPong) NextFrame(pixels Frame, timeMS uint32) {
	if len(pixels) == 0 || p.Child.Pattern == nil {
		return
	}
	p.buf.reset(len(pixels)*2 - 1)
	p.Child.NextFrame(p.buf, timeMS)
	// The last point of each extremity is only lit on one tick but every other
	// points are lit twice during a full cycle. This means the full cycle is
	// 2*(len(pixels)-1). For a 3 pixels line, the cycle is: x00, 0x0, 00x, 0x0.
	//
	// For Child being Frame "01234567":
	//   move == 0  -> "01234567"
	//   move == 2  -> "21056789"
	//   move == 5  -> "543210ab"
	//   move == 7  -> "76543210"
	//   move == 9  -> "98765012"
	//   move == 11 -> "ba901234"
	//   move == 13 -> "d0123456"
	//   move 14 -> move 0; "2*(8-1)"
	cycle := 2 * (len(pixels) - 1)
	// TODO(maruel): Smoothing with Curve, defaults to Step.
	// TODO(maruel): Convert to uint32.
	pos := int(float32(timeMS)*0.001*p.MovesPerSec) % cycle

	// Once it works the following code looks trivial but everytime it takes me
	// an absurd amount of time to rewrite it.
	if pos >= len(pixels)-1 {
		// Head runs left.
		// pos2 is the position from the right.
		pos2 := pos + 1 - len(pixels)
		// limit is the offset at which order change.
		limit := len(pixels) - pos2 - 1
		for i := range pixels {
			if i < limit {
				// Going right.
				pixels[i] = p.buf[len(pixels)-i+pos2-1]
			} else {
				// Going left.
				pixels[i] = p.buf[i-limit]
			}
		}
	} else {
		// Head runs right.
		for i := range pixels {
			if i <= pos {
				// Going right.
				pixels[i] = p.buf[pos-i]
			} else {
				// Going left.
				pixels[i] = p.buf[pos+i]
			}
		}
	}
}

// Crop draws a subset of a strip, not touching the rest.
type Crop struct {
	Child  SPattern
	Start  int // Starting pixels to skip
	Length int // Length of the pixels to affect
}

func (s *Crop) NextFrame(pixels Frame, timeMS uint32) {
	if s.Child.Pattern != nil {
		s.Child.NextFrame(pixels[s.Start:s.Length], timeMS)
	}
}

// Mixer is a generic mixer that merges the output from multiple patterns.
//
// It doesn't animate.
type Mixer struct {
	Patterns []SPattern
	Weights  []float32 // In theory Sum(Weights) should be 1 but it doesn't need to. For example, mixing a night sky will likely have all of the Weights set to 1.
	bufs     []Frame
}

func (m *Mixer) NextFrame(pixels Frame, timeMS uint32) {
	if len(m.Patterns) != len(m.Weights) {
		return
	}
	if len(m.bufs) != len(m.Patterns) {
		m.bufs = make([]Frame, len(m.Patterns))
	}
	for i := range m.bufs {
		m.bufs[i].reset(len(pixels))
	}

	// Draw each pattern.
	for i := range m.Patterns {
		m.Patterns[i].NextFrame(m.bufs[i], timeMS)
	}

	// Merge patterns.
	for i := range pixels {
		// TODO(maruel): Use uint32 calculation.
		var r, g, b float32
		for j := range m.bufs {
			c := m.bufs[j][i]
			w := m.Weights[j]
			r += float32(c.R) * w
			g += float32(c.G) * w
			b += float32(c.B) * w
		}
		pixels[i].R = FloatToUint8(r)
		pixels[i].G = FloatToUint8(g)
		pixels[i].B = FloatToUint8(b)
	}
}

// Scale adapts a larger or smaller patterns to the Strip size
//
// This is useful to create smoother animations or scale down images.
type Scale struct {
	Child         SPattern
	Interpolation Interpolation // Defaults to Linear
	Length        int           // A buffer of this length will be provided to Child and will be scaled to the actual pixels length
	Ratio         float32       // Scaling ratio to use, <1 means smaller, >1 means larger. Only one of Length or Ratio can be used
	buf           Frame
}

func (s *Scale) NextFrame(pixels Frame, timeMS uint32) {
	if s.Child.Pattern == nil {
		return
	}
	l := s.Length
	if l == 0 {
		l = int(ceil(s.Ratio * float32(len(pixels))))
	}
	s.buf.reset(l)
	s.Child.NextFrame(s.buf, timeMS)
	s.Interpolation.Scale(s.buf, pixels)
}
