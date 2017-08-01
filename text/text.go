// Package text implements text styling.
package text

import (
	"gomatcha.io/matcha/comm"
	pb "gomatcha.io/matcha/pb/text"
)

// type Position struct {
// 	id   int64
// 	text *Text
// }

// // -1 if the position has been removed.
// func (p *Position) Index() int {
// 	p.text.positionMu.Lock()
// 	defer p.text.positionMu.Unlock()
// 	return p.text.positions[p.id]
// }

// type position struct {
// 	id    int64
// 	index int
// }

type Text struct {
	group comm.Relay
	bytes []byte
	// isRune        []bool
	// isGlyph       []bool
	// runeCount     int
	// glyphCount    int
	// positions     map[int64]int
	// positionMaxId int64
	// positionMu    *sync.Mutex
}

// New is a convenience function that returns a new Text that contains string b.
func New(b string) *Text {
	t := &Text{}
	t.bytes = []byte(b)
	// t.positions = map[int64]int{}
	// t.normalize()
	return t
}

func (t *Text) MarshalProtobuf() *pb.Text {
	if t == nil {
		return nil
	}
	return &pb.Text{
		Text: string(t.bytes),
	}
}

func (t *Text) UnmarshalProtobuf(pbtext *pb.Text) error {
	t.SetString(pbtext.Text)
	return nil
}

// Notify implements comm.Notify.
func (t *Text) Notify(f func()) comm.Id {
	return t.group.Notify(f)
}

// Unnotify implements comm.Notify.
func (t *Text) Unnotify(id comm.Id) {
	t.group.Unnotify(id)
}

// // Panics if idx is out of range.
// func (t *Text) ByteAt(byteIdx int) byte {
// 	return t.bytes[byteIdx]
// }

// // Panics if idx is out of range.
// func (t *Text) RuneAt(byteIdx int) rune {
// 	// Start at the position and look backwards until we find the start of the rune
// 	var runeStart int = -1
// 	for i := byteIdx; i >= 0; i -= 1 {
// 		isRune := t.isRune[i]
// 		if isRune {
// 			runeStart = i
// 			break
// 		}
// 	}

// 	if runeStart == -1 {
// 		panic("RuneAt: Couldn't find rune start")
// 	}

// 	bytes := []byte{t.bytes[runeStart]}
// 	// Add bytes until next rune
// 	for i := runeStart + 1; i < len(t.bytes); i++ {
// 		if t.isRune[i] {
// 			break
// 		}
// 		bytes = append(bytes, t.bytes[i])
// 	}
// 	return []rune(string(bytes))[0]
// }

// // Panics if idx is out of range.
// func (t *Text) GlyphAt(byteIdx int) string {
// 	// Start at the position and look backwards until we find the start of the glyph
// 	var glyphStart int = -1
// 	for i := byteIdx; i >= 0; i -= 1 {
// 		isGlyph := t.isGlyph[i]
// 		if isGlyph {
// 			glyphStart = i
// 			break
// 		}
// 	}

// 	if glyphStart == -1 {
// 		panic("GlyphAt: Couldn't find glyph start")
// 	}

// 	bytes := []byte{t.bytes[glyphStart]}
// 	// Add bytes until next glyph
// 	for i := glyphStart + 1; i < len(t.bytes); i++ {
// 		if t.isGlyph[i] {
// 			break
// 		}
// 		bytes = append(bytes, t.bytes[i])
// 	}
// 	return string(bytes)
// }

// func (t *Text) ByteIndex(byteIdx int) int {
// 	return 0
// }

// func (t *Text) RuneIndex(runeIdx int) int {
// 	return 0
// }

// func (t *Text) GlyphIndex(glyphIdx int) int {
// 	return 0
// }

// // Returns -1 if out of range.
// func (t *Text) ByteNextIndex(byteIdx int) int {
// 	idx := byteIdx + 1
// 	if idx >= len(t.bytes) {
// 		return -1
// 	}
// 	return idx
// }

// // Returns -1 if out of range.
// func (t *Text) RuneNextIndex(byteIdx int) int {
// 	for i := byteIdx + 1; i < len(t.bytes); i += 1 {
// 		if t.isRune[i] {
// 			return i
// 		}
// 	}
// 	return -1
// }

// // Returns -1 if out of range.
// func (t *Text) GlyphNextIndex(byteIdx int) int {
// 	for i := byteIdx + 1; i < len(t.bytes); i += 1 {
// 		if t.isGlyph[i] {
// 			return i
// 		}
// 	}
// 	return -1
// }

// // Returns -1 if out of range.
// func (t *Text) BytePrevIndex(byteIdx int) int {
// 	idx := byteIdx - 1
// 	if idx < 0 {
// 		return -1
// 	}
// 	return idx
// }

// // Returns -1 if out of range.
// func (t *Text) RunePrevIndex(byteIdx int) int {
// 	for i := byteIdx - 1; i >= 0; i -= 1 {
// 		if t.isRune[i] {
// 			return i
// 		}
// 	}
// 	return -1
// }

// // Returns -1 if out of range.
// func (t *Text) GlyphPrevIndex(byteIdx int) int {
// 	for i := byteIdx - 1; i >= 0; i -= 1 {
// 		if t.isGlyph[i] {
// 			return i
// 		}
// 	}
// 	return -1
// }

// func (t *Text) ByteCount() int {
// 	return len(t.bytes)
// }

// func (t *Text) RuneCount() int {
// 	return t.runeCount
// }

// func (t *Text) GlyphCount() int {
// 	return t.glyphCount
// }

// // panics if minByteIdx or maxByteIdx is out of range.
// func (t *Text) Replace(minByteIdx, maxByteIdx int, new string) {
// 	if maxByteIdx > minByteIdx || minByteIdx < 0 || maxByteIdx > len(t.bytes) {
// 		panic("Index out of range")
// 	}
// }

// func (t *Text) Position(byteIdx int) *Position {
// 	t.positionMu.Lock()
// 	defer t.positionMu.Unlock()

// 	t.positionMaxId += 1
// 	t.positions[t.positionMaxId] = byteIdx

// 	p := &Position{
// 		id:   t.positionMaxId,
// 		text: t,
// 	}
// 	runtime.SetFinalizer(p, func(final *Position) {
// 		text := final.text
// 		text.positionMu.Lock()
// 		defer text.positionMu.Unlock()
// 		delete(text.positions, final.id)
// 	})
// 	return p
// }

// func (t *Text) normalize() {
// 	runeCount := 0
// 	glyphCount := 0
// 	isRune := make([]bool, 0, len(t.bytes))
// 	isGlyph := make([]bool, 0, len(t.bytes))
// 	bytes := make([]byte, 0, len(t.bytes))

// 	var iter norm.Iter
// 	iter.InitString(norm.NFD, string(t.bytes))
// 	for !iter.Done() {
// 		glyph := iter.Next()
// 		bytes = append(bytes, glyph...)

// 		for i := range glyph {
// 			isGlyph = append(isGlyph, i == 0)
// 		}
// 		glyphCount += 1

// 		isRuneSub := make([]bool, len(glyph))
// 		for i := range string(glyph) {
// 			isRuneSub[i] = true
// 			runeCount += 1
// 		}
// 		isRune = append(isRune, isRuneSub...)
// 	}
// 	t.glyphCount = glyphCount
// 	t.runeCount = runeCount
// 	t.isGlyph = isGlyph
// 	t.isRune = isRune
// 	t.bytes = bytes
// }

func (t *Text) SetString(str string) {
	t.bytes = []byte(str)
	// t.normalize()
	t.group.Signal()
}

func (t *Text) String() string {
	if t == nil {
		return "nil"
	}
	return string(t.bytes)
}

// Value implements comm.StringNotifier
func (t *Text) Value() string {
	return string(t.bytes)
}
