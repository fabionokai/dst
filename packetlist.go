// Copyright (C) 2014 Jakob Borg and Contributors (see the CONTRIBUTORS file).
// All rights reserved. Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dst

type packetList struct {
	packets []packet
	slot    int
}

// CutLessSeq cuts packets from the start of the list with sequence numbers
// lower than seq. Returns the number of packets that were cut.
func (l *packetList) CutLessSeq(seq uint32) int {
	var i int
	for i = range l.packets {
		if i == l.slot {
			break
		}
		if !l.packets[i].LessSeq(seq) {
			break
		}
	}
	if i > 0 {
		l.Cut(i)
	}
	return i
}

func (l *packetList) Cut(n int) {
	copy(l.packets, l.packets[n:])
	l.slot -= n
}

func (l *packetList) Full() bool {
	return l.slot == len(l.packets)
}

func (l *packetList) All() []packet {
	return l.packets[:l.slot]
}

func (l *packetList) Append(pkt packet) bool {
	if l.slot == len(l.packets) {
		return false
	}
	l.packets[l.slot] = pkt
	l.slot++
	return true
}

func (l *packetList) AppendAll(pkts []packet) {
	l.packets = append(l.packets[:l.slot], pkts...)
	l.slot += len(pkts)
}

func (l *packetList) Cap() int {
	return len(l.packets)
}

func (l *packetList) Len() int {
	return l.slot
}

func (l *packetList) Resize(s int) {
	if s <= cap(l.packets) {
		l.packets = l.packets[:s]
	} else {
		t := make([]packet, s)
		copy(t, l.packets)
		l.packets = t
	}
}

func (l *packetList) InsertSorted(pkt packet) {
	for i := range l.packets {
		if i >= l.slot {
			l.packets[i] = pkt
			l.slot++
			return
		}
		if pkt.hdr.sequenceNo == l.packets[i].hdr.sequenceNo {
			return
		}
		if pkt.Less(l.packets[i]) {
			copy(l.packets[i+1:], l.packets[i:])
			l.packets[i] = pkt
			if l.slot < len(l.packets) {
				l.slot++
			}
			return
		}
	}
}

func (l *packetList) LowestSeq() uint32 {
	return l.packets[0].hdr.sequenceNo
}

func (l *packetList) PopSequence() []packet {
	highSeq := l.packets[0].hdr.sequenceNo
	var i int
	for i = 1; i < l.slot; i++ {
		if l.packets[i].hdr.sequenceNo != highSeq+1 {
			break
		}
		highSeq++
	}
	pkts := make([]packet, i)
	copy(pkts, l.packets[:i])
	l.Cut(i)
	return pkts
}
