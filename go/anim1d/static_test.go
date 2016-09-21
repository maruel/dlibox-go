// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"fmt"
	"testing"

	"github.com/maruel/ut"
)

func TestColor(t *testing.T) {
	p := &Color{255, 255, 255}
	e := []expectation{{3000, Frame{{255, 255, 255}}}}
	testFrames(t, p, e)
}

func TestColorMix(t *testing.T) {
	w := Color{255, 255, 255}
	b := Color{0, 0, 0}
	c := w
	c.Mix(b, 0)
	ut.AssertEqual(t, c, w)
	c = w
	c.Mix(b, 255)
	ut.AssertEqual(t, c, b)

	// Not sure where this difference comes from.
	c = w
	c.Mix(b, 128)
	ut.AssertEqual(t, c, Color{127, 127, 127})
	c = b
	c.Mix(w, 128)
	ut.AssertEqual(t, c, Color{128, 128, 128})

	// Test for overflow.
	c = w
	c.Mix(w, 0)
	ut.AssertEqual(t, c, w)
	c.Mix(w, 128)
	ut.AssertEqual(t, c, w)
	c.Mix(w, 255)
	ut.AssertEqual(t, c, w)

	// Verify channels.
	a := Color{0x10, 0x20, 0x30}
	c = a
	c.Mix(b, 0)
	ut.AssertEqual(t, c, a)
	c = b
	c.Mix(a, 255)
	ut.AssertEqual(t, c, a)
}

func TestWaveLength2RGB(t *testing.T) {
	data := []struct {
		input    int
		expected Color
	}{
		{379, Color{0x00, 0x00, 0x00}},
		{380, Color{0x1A, 0x00, 0x1A}},
		{381, Color{0x1D, 0x00, 0x1F}},
		{382, Color{0x20, 0x00, 0x25}},
		{383, Color{0x23, 0x00, 0x2B}},
		{384, Color{0x26, 0x00, 0x30}},
		{385, Color{0x29, 0x00, 0x36}},
		{386, Color{0x2B, 0x00, 0x3C}},
		{387, Color{0x2E, 0x00, 0x42}},
		{388, Color{0x31, 0x00, 0x47}},
		{389, Color{0x34, 0x00, 0x4D}},
		{390, Color{0x37, 0x00, 0x53}},
		{391, Color{0x3A, 0x00, 0x58}},
		{392, Color{0x3C, 0x00, 0x5E}},
		{393, Color{0x3F, 0x00, 0x64}},
		{394, Color{0x42, 0x00, 0x6A}},
		{395, Color{0x45, 0x00, 0x6F}},
		{396, Color{0x48, 0x00, 0x75}},
		{397, Color{0x4B, 0x00, 0x7B}},
		{398, Color{0x4D, 0x00, 0x81}},
		{399, Color{0x50, 0x00, 0x86}},
		{400, Color{0x53, 0x00, 0x8C}},
		{401, Color{0x56, 0x00, 0x92}},
		{402, Color{0x59, 0x00, 0x97}},
		{403, Color{0x5C, 0x00, 0x9D}},
		{404, Color{0x5E, 0x00, 0xA3}},
		{405, Color{0x61, 0x00, 0xA9}},
		{406, Color{0x64, 0x00, 0xAE}},
		{407, Color{0x67, 0x00, 0xB4}},
		{408, Color{0x6A, 0x00, 0xBA}},
		{409, Color{0x6D, 0x00, 0xC0}},
		{410, Color{0x6F, 0x00, 0xC5}},
		{411, Color{0x72, 0x00, 0xCB}},
		{412, Color{0x75, 0x00, 0xD1}},
		{413, Color{0x78, 0x00, 0xD6}},
		{414, Color{0x7B, 0x00, 0xDC}},
		{415, Color{0x7E, 0x00, 0xE2}},
		{416, Color{0x80, 0x00, 0xE8}},
		{417, Color{0x83, 0x00, 0xED}},
		{418, Color{0x86, 0x00, 0xF3}},
		{419, Color{0x89, 0x00, 0xF9}},
		{420, Color{0x89, 0x00, 0xFF}},
		{421, Color{0x82, 0x00, 0xFF}},
		{422, Color{0x7B, 0x00, 0xFF}},
		{423, Color{0x74, 0x00, 0xFF}},
		{424, Color{0x6D, 0x00, 0xFF}},
		{425, Color{0x66, 0x00, 0xFF}},
		{426, Color{0x5F, 0x00, 0xFF}},
		{427, Color{0x59, 0x00, 0xFF}},
		{428, Color{0x52, 0x00, 0xFF}},
		{429, Color{0x4B, 0x00, 0xFF}},
		{430, Color{0x44, 0x00, 0xFF}},
		{431, Color{0x3D, 0x00, 0xFF}},
		{432, Color{0x36, 0x00, 0xFF}},
		{433, Color{0x2F, 0x00, 0xFF}},
		{434, Color{0x29, 0x00, 0xFF}},
		{435, Color{0x22, 0x00, 0xFF}},
		{436, Color{0x1B, 0x00, 0xFF}},
		{437, Color{0x14, 0x00, 0xFF}},
		{438, Color{0x0D, 0x00, 0xFF}},
		{439, Color{0x06, 0x00, 0xFF}},
		{440, Color{0x00, 0x00, 0xFF}},
		{441, Color{0x00, 0x05, 0xFF}},
		{442, Color{0x00, 0x0A, 0xFF}},
		{443, Color{0x00, 0x0F, 0xFF}},
		{444, Color{0x00, 0x14, 0xFF}},
		{445, Color{0x00, 0x19, 0xFF}},
		{446, Color{0x00, 0x1E, 0xFF}},
		{447, Color{0x00, 0x23, 0xFF}},
		{448, Color{0x00, 0x28, 0xFF}},
		{449, Color{0x00, 0x2D, 0xFF}},
		{450, Color{0x00, 0x33, 0xFF}},
		{451, Color{0x00, 0x38, 0xFF}},
		{452, Color{0x00, 0x3D, 0xFF}},
		{453, Color{0x00, 0x42, 0xFF}},
		{454, Color{0x00, 0x47, 0xFF}},
		{455, Color{0x00, 0x4C, 0xFF}},
		{456, Color{0x00, 0x51, 0xFF}},
		{457, Color{0x00, 0x56, 0xFF}},
		{458, Color{0x00, 0x5B, 0xFF}},
		{459, Color{0x00, 0x60, 0xFF}},
		{460, Color{0x00, 0x66, 0xFF}},
		{461, Color{0x00, 0x6B, 0xFF}},
		{462, Color{0x00, 0x70, 0xFF}},
		{463, Color{0x00, 0x75, 0xFF}},
		{464, Color{0x00, 0x7A, 0xFF}},
		{465, Color{0x00, 0x7F, 0xFF}},
		{466, Color{0x00, 0x84, 0xFF}},
		{467, Color{0x00, 0x89, 0xFF}},
		{468, Color{0x00, 0x8E, 0xFF}},
		{469, Color{0x00, 0x93, 0xFF}},
		{470, Color{0x00, 0x99, 0xFF}},
		{471, Color{0x00, 0x9E, 0xFF}},
		{472, Color{0x00, 0xA3, 0xFF}},
		{473, Color{0x00, 0xA8, 0xFF}},
		{474, Color{0x00, 0xAD, 0xFF}},
		{475, Color{0x00, 0xB2, 0xFF}},
		{476, Color{0x00, 0xB7, 0xFF}},
		{477, Color{0x00, 0xBC, 0xFF}},
		{478, Color{0x00, 0xC1, 0xFF}},
		{479, Color{0x00, 0xC6, 0xFF}},
		{480, Color{0x00, 0xCC, 0xFF}},
		{481, Color{0x00, 0xD1, 0xFF}},
		{482, Color{0x00, 0xD6, 0xFF}},
		{483, Color{0x00, 0xDB, 0xFF}},
		{484, Color{0x00, 0xE0, 0xFF}},
		{485, Color{0x00, 0xE5, 0xFF}},
		{486, Color{0x00, 0xEA, 0xFF}},
		{487, Color{0x00, 0xEF, 0xFF}},
		{488, Color{0x00, 0xF4, 0xFF}},
		{489, Color{0x00, 0xF9, 0xFF}},
		{490, Color{0x00, 0xFF, 0xFF}},
		{491, Color{0x00, 0xFF, 0xF2}},
		{492, Color{0x00, 0xFF, 0xE5}},
		{493, Color{0x00, 0xFF, 0xD8}},
		{494, Color{0x00, 0xFF, 0xCC}},
		{495, Color{0x00, 0xFF, 0xBF}},
		{496, Color{0x00, 0xFF, 0xB2}},
		{497, Color{0x00, 0xFF, 0xA5}},
		{498, Color{0x00, 0xFF, 0x99}},
		{499, Color{0x00, 0xFF, 0x8C}},
		{500, Color{0x00, 0xFF, 0x7F}},
		{501, Color{0x00, 0xFF, 0x72}},
		{502, Color{0x00, 0xFF, 0x66}},
		{503, Color{0x00, 0xFF, 0x59}},
		{504, Color{0x00, 0xFF, 0x4C}},
		{505, Color{0x00, 0xFF, 0x3F}},
		{506, Color{0x00, 0xFF, 0x33}},
		{507, Color{0x00, 0xFF, 0x26}},
		{508, Color{0x00, 0xFF, 0x19}},
		{509, Color{0x00, 0xFF, 0x0C}},
		{510, Color{0x00, 0xFF, 0x00}},
		{511, Color{0x03, 0xFF, 0x00}},
		{512, Color{0x07, 0xFF, 0x00}},
		{513, Color{0x0A, 0xFF, 0x00}},
		{514, Color{0x0E, 0xFF, 0x00}},
		{515, Color{0x12, 0xFF, 0x00}},
		{516, Color{0x15, 0xFF, 0x00}},
		{517, Color{0x19, 0xFF, 0x00}},
		{518, Color{0x1D, 0xFF, 0x00}},
		{519, Color{0x20, 0xFF, 0x00}},
		{520, Color{0x24, 0xFF, 0x00}},
		{521, Color{0x28, 0xFF, 0x00}},
		{522, Color{0x2B, 0xFF, 0x00}},
		{523, Color{0x2F, 0xFF, 0x00}},
		{524, Color{0x33, 0xFF, 0x00}},
		{525, Color{0x36, 0xFF, 0x00}},
		{526, Color{0x3A, 0xFF, 0x00}},
		{527, Color{0x3D, 0xFF, 0x00}},
		{528, Color{0x41, 0xFF, 0x00}},
		{529, Color{0x45, 0xFF, 0x00}},
		{530, Color{0x48, 0xFF, 0x00}},
		{531, Color{0x4C, 0xFF, 0x00}},
		{532, Color{0x50, 0xFF, 0x00}},
		{533, Color{0x53, 0xFF, 0x00}},
		{534, Color{0x57, 0xFF, 0x00}},
		{535, Color{0x5B, 0xFF, 0x00}},
		{536, Color{0x5E, 0xFF, 0x00}},
		{537, Color{0x62, 0xFF, 0x00}},
		{538, Color{0x66, 0xFF, 0x00}},
		{539, Color{0x69, 0xFF, 0x00}},
		{540, Color{0x6D, 0xFF, 0x00}},
		{541, Color{0x70, 0xFF, 0x00}},
		{542, Color{0x74, 0xFF, 0x00}},
		{543, Color{0x78, 0xFF, 0x00}},
		{544, Color{0x7B, 0xFF, 0x00}},
		{545, Color{0x7F, 0xFF, 0x00}},
		{546, Color{0x83, 0xFF, 0x00}},
		{547, Color{0x86, 0xFF, 0x00}},
		{548, Color{0x8A, 0xFF, 0x00}},
		{549, Color{0x8E, 0xFF, 0x00}},
		{550, Color{0x91, 0xFF, 0x00}},
		{551, Color{0x95, 0xFF, 0x00}},
		{552, Color{0x99, 0xFF, 0x00}},
		{553, Color{0x9C, 0xFF, 0x00}},
		{554, Color{0xA0, 0xFF, 0x00}},
		{555, Color{0xA3, 0xFF, 0x00}},
		{556, Color{0xA7, 0xFF, 0x00}},
		{557, Color{0xAB, 0xFF, 0x00}},
		{558, Color{0xAE, 0xFF, 0x00}},
		{559, Color{0xB2, 0xFF, 0x00}},
		{560, Color{0xB6, 0xFF, 0x00}},
		{561, Color{0xB9, 0xFF, 0x00}},
		{562, Color{0xBD, 0xFF, 0x00}},
		{563, Color{0xC1, 0xFF, 0x00}},
		{564, Color{0xC4, 0xFF, 0x00}},
		{565, Color{0xC8, 0xFF, 0x00}},
		{566, Color{0xCC, 0xFF, 0x00}},
		{567, Color{0xCF, 0xFF, 0x00}},
		{568, Color{0xD3, 0xFF, 0x00}},
		{569, Color{0xD6, 0xFF, 0x00}},
		{570, Color{0xDA, 0xFF, 0x00}},
		{571, Color{0xDE, 0xFF, 0x00}},
		{572, Color{0xE1, 0xFF, 0x00}},
		{573, Color{0xE5, 0xFF, 0x00}},
		{574, Color{0xE9, 0xFF, 0x00}},
		{575, Color{0xEC, 0xFF, 0x00}},
		{576, Color{0xF0, 0xFF, 0x00}},
		{577, Color{0xF4, 0xFF, 0x00}},
		{578, Color{0xF7, 0xFF, 0x00}},
		{579, Color{0xFB, 0xFF, 0x00}},
		{580, Color{0xFF, 0xFF, 0x00}},
		{581, Color{0xFF, 0xFB, 0x00}},
		{582, Color{0xFF, 0xF7, 0x00}},
		{583, Color{0xFF, 0xF3, 0x00}},
		{584, Color{0xFF, 0xEF, 0x00}},
		{585, Color{0xFF, 0xEB, 0x00}},
		{586, Color{0xFF, 0xE7, 0x00}},
		{587, Color{0xFF, 0xE3, 0x00}},
		{588, Color{0xFF, 0xDF, 0x00}},
		{589, Color{0xFF, 0xDB, 0x00}},
		{590, Color{0xFF, 0xD7, 0x00}},
		{591, Color{0xFF, 0xD3, 0x00}},
		{592, Color{0xFF, 0xCF, 0x00}},
		{593, Color{0xFF, 0xCC, 0x00}},
		{594, Color{0xFF, 0xC8, 0x00}},
		{595, Color{0xFF, 0xC4, 0x00}},
		{596, Color{0xFF, 0xC0, 0x00}},
		{597, Color{0xFF, 0xBC, 0x00}},
		{598, Color{0xFF, 0xB8, 0x00}},
		{599, Color{0xFF, 0xB4, 0x00}},
		{600, Color{0xFF, 0xB0, 0x00}},
		{601, Color{0xFF, 0xAC, 0x00}},
		{602, Color{0xFF, 0xA8, 0x00}},
		{603, Color{0xFF, 0xA4, 0x00}},
		{604, Color{0xFF, 0xA0, 0x00}},
		{605, Color{0xFF, 0x9C, 0x00}},
		{606, Color{0xFF, 0x99, 0x00}},
		{607, Color{0xFF, 0x95, 0x00}},
		{608, Color{0xFF, 0x91, 0x00}},
		{609, Color{0xFF, 0x8D, 0x00}},
		{610, Color{0xFF, 0x89, 0x00}},
		{611, Color{0xFF, 0x85, 0x00}},
		{612, Color{0xFF, 0x81, 0x00}},
		{613, Color{0xFF, 0x7D, 0x00}},
		{614, Color{0xFF, 0x79, 0x00}},
		{615, Color{0xFF, 0x75, 0x00}},
		{616, Color{0xFF, 0x71, 0x00}},
		{617, Color{0xFF, 0x6D, 0x00}},
		{618, Color{0xFF, 0x69, 0x00}},
		{619, Color{0xFF, 0x66, 0x00}},
		{620, Color{0xFF, 0x62, 0x00}},
		{621, Color{0xFF, 0x5E, 0x00}},
		{622, Color{0xFF, 0x5A, 0x00}},
		{623, Color{0xFF, 0x56, 0x00}},
		{624, Color{0xFF, 0x52, 0x00}},
		{625, Color{0xFF, 0x4E, 0x00}},
		{626, Color{0xFF, 0x4A, 0x00}},
		{627, Color{0xFF, 0x46, 0x00}},
		{628, Color{0xFF, 0x42, 0x00}},
		{629, Color{0xFF, 0x3E, 0x00}},
		{630, Color{0xFF, 0x3A, 0x00}},
		{631, Color{0xFF, 0x36, 0x00}},
		{632, Color{0xFF, 0x33, 0x00}},
		{633, Color{0xFF, 0x2F, 0x00}},
		{634, Color{0xFF, 0x2B, 0x00}},
		{635, Color{0xFF, 0x27, 0x00}},
		{636, Color{0xFF, 0x23, 0x00}},
		{637, Color{0xFF, 0x1F, 0x00}},
		{638, Color{0xFF, 0x1B, 0x00}},
		{639, Color{0xFF, 0x17, 0x00}},
		{640, Color{0xFF, 0x13, 0x00}},
		{641, Color{0xFF, 0x0F, 0x00}},
		{642, Color{0xFF, 0x0B, 0x00}},
		{643, Color{0xFF, 0x07, 0x00}},
		{644, Color{0xFF, 0x03, 0x00}},
		{645, Color{0xFF, 0x00, 0x00}},
		{646, Color{0xFF, 0x00, 0x00}},
		{647, Color{0xFF, 0x00, 0x00}},
		{648, Color{0xFF, 0x00, 0x00}},
		{649, Color{0xFF, 0x00, 0x00}},
		{650, Color{0xFF, 0x00, 0x00}},
		{651, Color{0xFF, 0x00, 0x00}},
		{652, Color{0xFF, 0x00, 0x00}},
		{653, Color{0xFF, 0x00, 0x00}},
		{654, Color{0xFF, 0x00, 0x00}},
		{655, Color{0xFF, 0x00, 0x00}},
		{656, Color{0xFF, 0x00, 0x00}},
		{657, Color{0xFF, 0x00, 0x00}},
		{658, Color{0xFF, 0x00, 0x00}},
		{659, Color{0xFF, 0x00, 0x00}},
		{660, Color{0xFF, 0x00, 0x00}},
		{661, Color{0xFF, 0x00, 0x00}},
		{662, Color{0xFF, 0x00, 0x00}},
		{663, Color{0xFF, 0x00, 0x00}},
		{664, Color{0xFF, 0x00, 0x00}},
		{665, Color{0xFF, 0x00, 0x00}},
		{666, Color{0xFF, 0x00, 0x00}},
		{667, Color{0xFF, 0x00, 0x00}},
		{668, Color{0xFF, 0x00, 0x00}},
		{669, Color{0xFF, 0x00, 0x00}},
		{670, Color{0xFF, 0x00, 0x00}},
		{671, Color{0xFF, 0x00, 0x00}},
		{672, Color{0xFF, 0x00, 0x00}},
		{673, Color{0xFF, 0x00, 0x00}},
		{674, Color{0xFF, 0x00, 0x00}},
		{675, Color{0xFF, 0x00, 0x00}},
		{676, Color{0xFF, 0x00, 0x00}},
		{677, Color{0xFF, 0x00, 0x00}},
		{678, Color{0xFF, 0x00, 0x00}},
		{679, Color{0xFF, 0x00, 0x00}},
		{680, Color{0xFF, 0x00, 0x00}},
		{681, Color{0xFF, 0x00, 0x00}},
		{682, Color{0xFF, 0x00, 0x00}},
		{683, Color{0xFF, 0x00, 0x00}},
		{684, Color{0xFF, 0x00, 0x00}},
		{685, Color{0xFF, 0x00, 0x00}},
		{686, Color{0xFF, 0x00, 0x00}},
		{687, Color{0xFF, 0x00, 0x00}},
		{688, Color{0xFF, 0x00, 0x00}},
		{689, Color{0xFF, 0x00, 0x00}},
		{690, Color{0xFF, 0x00, 0x00}},
		{691, Color{0xFF, 0x00, 0x00}},
		{692, Color{0xFF, 0x00, 0x00}},
		{693, Color{0xFF, 0x00, 0x00}},
		{694, Color{0xFF, 0x00, 0x00}},
		{695, Color{0xFF, 0x00, 0x00}},
		{696, Color{0xFF, 0x00, 0x00}},
		{697, Color{0xFF, 0x00, 0x00}},
		{698, Color{0xFF, 0x00, 0x00}},
		{699, Color{0xFF, 0x00, 0x00}},
		{700, Color{0xFF, 0x00, 0x00}},
		{701, Color{0xFC, 0x00, 0x00}},
		{702, Color{0xF9, 0x00, 0x00}},
		{703, Color{0xF6, 0x00, 0x00}},
		{704, Color{0xF3, 0x00, 0x00}},
		{705, Color{0xF0, 0x00, 0x00}},
		{706, Color{0xED, 0x00, 0x00}},
		{707, Color{0xEA, 0x00, 0x00}},
		{708, Color{0xE8, 0x00, 0x00}},
		{709, Color{0xE5, 0x00, 0x00}},
		{710, Color{0xE2, 0x00, 0x00}},
		{711, Color{0xDF, 0x00, 0x00}},
		{712, Color{0xDC, 0x00, 0x00}},
		{713, Color{0xD9, 0x00, 0x00}},
		{714, Color{0xD6, 0x00, 0x00}},
		{715, Color{0xD4, 0x00, 0x00}},
		{716, Color{0xD1, 0x00, 0x00}},
		{717, Color{0xCE, 0x00, 0x00}},
		{718, Color{0xCB, 0x00, 0x00}},
		{719, Color{0xC8, 0x00, 0x00}},
		{720, Color{0xC5, 0x00, 0x00}},
		{721, Color{0xC2, 0x00, 0x00}},
		{722, Color{0xC0, 0x00, 0x00}},
		{723, Color{0xBD, 0x00, 0x00}},
		{724, Color{0xBA, 0x00, 0x00}},
		{725, Color{0xB7, 0x00, 0x00}},
		{726, Color{0xB4, 0x00, 0x00}},
		{727, Color{0xB1, 0x00, 0x00}},
		{728, Color{0xAE, 0x00, 0x00}},
		{729, Color{0xAB, 0x00, 0x00}},
		{730, Color{0xA9, 0x00, 0x00}},
		{731, Color{0xA6, 0x00, 0x00}},
		{732, Color{0xA3, 0x00, 0x00}},
		{733, Color{0xA0, 0x00, 0x00}},
		{734, Color{0x9D, 0x00, 0x00}},
		{735, Color{0x9A, 0x00, 0x00}},
		{736, Color{0x97, 0x00, 0x00}},
		{737, Color{0x95, 0x00, 0x00}},
		{738, Color{0x92, 0x00, 0x00}},
		{739, Color{0x8F, 0x00, 0x00}},
		{740, Color{0x8C, 0x00, 0x00}},
		{741, Color{0x89, 0x00, 0x00}},
		{742, Color{0x86, 0x00, 0x00}},
		{743, Color{0x83, 0x00, 0x00}},
		{744, Color{0x81, 0x00, 0x00}},
		{745, Color{0x7E, 0x00, 0x00}},
		{746, Color{0x7B, 0x00, 0x00}},
		{747, Color{0x78, 0x00, 0x00}},
		{748, Color{0x75, 0x00, 0x00}},
		{749, Color{0x72, 0x00, 0x00}},
		{750, Color{0x6F, 0x00, 0x00}},
		{751, Color{0x6D, 0x00, 0x00}},
		{752, Color{0x6A, 0x00, 0x00}},
		{753, Color{0x67, 0x00, 0x00}},
		{754, Color{0x64, 0x00, 0x00}},
		{755, Color{0x61, 0x00, 0x00}},
		{756, Color{0x5E, 0x00, 0x00}},
		{757, Color{0x5B, 0x00, 0x00}},
		{758, Color{0x58, 0x00, 0x00}},
		{759, Color{0x56, 0x00, 0x00}},
		{760, Color{0x53, 0x00, 0x00}},
		{761, Color{0x50, 0x00, 0x00}},
		{762, Color{0x4D, 0x00, 0x00}},
		{763, Color{0x4A, 0x00, 0x00}},
		{764, Color{0x47, 0x00, 0x00}},
		{765, Color{0x44, 0x00, 0x00}},
		{766, Color{0x42, 0x00, 0x00}},
		{767, Color{0x3F, 0x00, 0x00}},
		{768, Color{0x3C, 0x00, 0x00}},
		{769, Color{0x39, 0x00, 0x00}},
		{770, Color{0x36, 0x00, 0x00}},
		{771, Color{0x33, 0x00, 0x00}},
		{772, Color{0x30, 0x00, 0x00}},
		{773, Color{0x2E, 0x00, 0x00}},
		{774, Color{0x2B, 0x00, 0x00}},
		{775, Color{0x28, 0x00, 0x00}},
		{776, Color{0x25, 0x00, 0x00}},
		{777, Color{0x22, 0x00, 0x00}},
		{778, Color{0x1F, 0x00, 0x00}},
		{779, Color{0x1C, 0x00, 0x00}},
		{780, Color{0x1A, 0x00, 0x00}},
		{781, Color{0x00, 0x00, 0x00}},
	}
	if false {
		for _, line := range data {
			c := waveLength2RGB(line.input)
			fmt.Printf("{%d, Color{0x%02X, 0x%02X, 0x%02X}},\n", int(line.input), c.R, c.G, c.B)
		}
	}
	for i, line := range data {
		ut.AssertEqualIndex(t, i, line.expected, waveLength2RGB(line.input))
		ut.AssertEqual(t, i+379, int(line.input))
	}
}

func TestRainbow(t *testing.T) {
	f := make(Frame, 16)
	r := Rainbow{}
	r.NextFrame(f, 0)
	var expected Frame
	if err := expected.UnmarshalJSON([]byte("\"L4d00818300ed1b00ff0056ff00c1ff00ff7f2bff0083ff00deff00ffc000ff5600ff0000ff0000dc00007e0000000000\"")); err != nil {
		t.Fatal(err)
	}
	if !frameEqual(f, expected) {
		t.Fatalf("%s != %s", Marshal(f), Marshal(expected))
	}
}

func TestRepeated(t *testing.T) {
	a := Color{0x10, 0x10, 0x10}
	b := Color{0x20, 0x20, 0x20}
	c := Color{0x30, 0x30, 0x30}
	p := &Repeated{Frame{a, b, c}}
	e := []expectation{
		{0, Frame{a, b, c, a, b}},
		{0, Frame{a}},
	}
	testFrames(t, p, e)
}

//

type expectation struct {
	offsetMS uint32
	colors   Frame
}

func testFrames(t *testing.T, p Pattern, expectations []expectation) {
	var pixels Frame
	for frame, e := range expectations {
		pixels.reset(len(e.colors))
		p.NextFrame(pixels, e.offsetMS)
		if !frameEqual(e.colors, pixels) {
			x := Marshal(e.colors)
			t.Fatalf("frame=%d bad expectation:\n%s\n%s", frame, x, Marshal(pixels))
		}
	}
}

func testFrame(t *testing.T, p Pattern, e expectation) {
	pixels := make(Frame, len(e.colors))
	p.NextFrame(pixels, e.offsetMS)
	if !frameEqual(e.colors, pixels) {
		t.Fatalf("%s != %s", Marshal(e.colors), Marshal(pixels))
	}
}

func frameEqual(lhs, rhs Frame) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i, a := range lhs {
		if a != rhs[i] {
			return false
		}
	}
	return true
}

func frameSimilar(lhs, rhs Frame) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i, a := range lhs {
		b := rhs[i]
		dR := int(a.R) - int(b.R)
		dG := int(a.G) - int(b.G)
		dB := int(a.B) - int(b.B)
		if dR > 1 || dR < -1 || dG > 1 || dG < -1 || dB > 1 || dB < -1 {
			return false
		}
	}
	return true
}
