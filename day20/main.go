package main

import (
	"bufio"
	"fmt"
	"log"
	"math/bits"
	"os"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	comMods := map[string]*ComMod{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lhs, rhs, ok := strings.Cut(line, " -> ")
		if !ok {
			log.Fatalln("Invalid line")
		}
		kind := lhs[0]
		if kind == '%' || kind == '&' {
			lhs = lhs[1:]
		} else {
			kind = 0
		}
		dest := strings.Split(rhs, ", ")
		comMods[lhs] = &ComMod{
			Name:  lhs,
			Kind:  kind,
			State: false,
			Mem:   map[string]bool{},
			Dest:  dest,
			Inbox: NewRing[Packet](),
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	for _, cm := range comMods {
		for _, i := range cm.Dest {
			if _, ok := comMods[i]; !ok {
				comMods[i] = &ComMod{
					Name:  i,
					Kind:  0,
					State: false,
					Mem:   map[string]bool{},
					Dest:  nil,
					Inbox: NewRing[Packet](),
				}
			}
			comMods[i].Mem[cm.Name] = false
		}
	}

	sumHi := 0
	sumLo := 0
	idx := 0
	targetCycles := map[string]Cycle{}
	totalRevisits := 0
	for idx < 1000 {
		comMods["broadcaster"].Inbox.Write(Packet{
			From: "",
			Sig:  false,
		})
		hi, lo, targetPackets := runTilEnd(comMods, "dg", true)
		sumHi += hi
		sumLo += lo
		idx++
		for _, i := range targetPackets {
			if c, ok := targetCycles[i.From]; ok {
				if c.Size == 0 {
					cycle := idx - c.Prev
					targetCycles[i.From] = Cycle{
						Prev: idx,
						Size: cycle,
						Rem:  idx % cycle,
					}
					totalRevisits++
				} else {
					if cycle := idx - c.Prev; cycle != c.Size {
						log.Fatalln("Multiple cycle lengths")
					}
					c.Prev = idx
					targetCycles[i.From] = c
				}
			} else {
				targetCycles[i.From] = Cycle{
					Prev: idx,
				}
			}
		}
	}

	fmt.Println("Part 1:", sumHi*sumLo)

	for totalRevisits < 4 {
		comMods["broadcaster"].Inbox.Write(Packet{
			From: "",
			Sig:  false,
		})
		_, _, targetPackets := runTilEnd(comMods, "dg", true)
		idx++
		for _, i := range targetPackets {
			if c, ok := targetCycles[i.From]; ok {
				if c.Size == 0 {
					cycle := idx - c.Prev
					targetCycles[i.From] = Cycle{
						Prev: idx,
						Size: cycle,
						Rem:  idx % cycle,
					}
					totalRevisits++
				} else {
					if cycle := idx - c.Prev; cycle != c.Size {
						log.Fatalln("Multiple cycle lengths")
					}
					c.Prev = idx
					targetCycles[i.From] = c
				}
			} else {
				targetCycles[i.From] = Cycle{
					Prev: idx,
				}
			}
		}
	}

	a, m := 0, 0
	for _, v := range targetCycles {
		if m == 0 {
			a, m = v.Rem, v.Size
			continue
		}
		var ok bool
		a, m, ok = crt(a, m, v.Rem, v.Size)
		if !ok {
			log.Fatalln("Unsolvable constraints")
		}
	}
	if a <= 0 {
		a += m
	}
	fmt.Println("Part 2:", a)
}

type (
	ComMod struct {
		Name  string
		Kind  byte
		State bool
		Mem   map[string]bool
		Dest  []string
		Inbox *Ring[Packet]
	}

	Packet struct {
		From string
		Sig  bool
	}

	Cycle struct {
		Prev int
		Size int
		Rem  int
	}
)

func runTilEnd(comMods map[string]*ComMod, target string, targetSig bool) (int, int, []Packet) {
	var packetsToTarget []Packet
	hi, lo := 0, 0
	for {
		hasSig := false
		for _, v := range comMods {
			for {
				packet, ok := v.Inbox.Read()
				if !ok {
					break
				}
				if v.Name == target {
					if packet.Sig == targetSig {
						packetsToTarget = append(packetsToTarget, packet)
					}
				}
				if packet.Sig {
					hi++
				} else {
					lo++
				}
				hasSig = true
				pulse(comMods, v.Name, packet)
			}
		}
		if !hasSig {
			return hi, lo, packetsToTarget
		}
	}
}

func pulse(comMods map[string]*ComMod, name string, packet Packet) {
	cm, ok := comMods[name]
	if !ok {
		log.Fatalln("Invalid mod name")
	}
	switch cm.Kind {
	case '%':
		{
			if packet.Sig {
				return
			}
			cm.State = !cm.State
			destSig := false
			if cm.State {
				destSig = true
			}
			for _, i := range cm.Dest {
				if _, ok := comMods[i]; !ok {
					log.Fatalln("Invalid dest")
				}
				comMods[i].Inbox.Write(Packet{
					From: name,
					Sig:  destSig,
				})
			}
		}
	case '&':
		{
			cm.Mem[packet.From] = packet.Sig
			destSig := false
			for _, v := range cm.Mem {
				if !v {
					destSig = true
					break
				}
			}
			cm.State = destSig
			for _, i := range cm.Dest {
				if _, ok := comMods[i]; !ok {
					log.Fatalln("Invalid dest")
				}
				comMods[i].Inbox.Write(Packet{
					From: name,
					Sig:  destSig,
				})
			}
		}
	case 0:
		{
			cm.State = packet.Sig
			for _, i := range cm.Dest {
				if _, ok := comMods[i]; !ok {
					log.Fatalln("Invalid dest")
				}
				comMods[i].Inbox.Write(Packet{
					From: name,
					Sig:  packet.Sig,
				})
			}
			return
		}
	default:
		log.Fatalln("Invalid mod kind")
		return
	}
}

func crt(a1, m1, a2, m2 int) (int, int, bool) {
	g, p, q := extGCD(m1, m2)
	if a1%g != a2%g {
		return 0, 0, false
	}
	m1g := m1 / g
	m2g := m2 / g
	lcm := m1g * m2
	// a1 * m2/g * q + a2 * m1/g * p (mod lcm)
	x := (mulmod(mulmod(a1, m2g, lcm), q, lcm) + mulmod(mulmod(a2, m1g, lcm), p, lcm)) % lcm
	if x < 0 {
		x += lcm
	}
	return x, lcm, true
}

func extGCD(a, b int) (int, int, int) {
	x2 := 1
	x1 := 0
	y2 := 0
	y1 := 1
	// a should be larger than b
	flip := false
	if a < b {
		a, b = b, a
		flip = true
	}
	for b > 0 {
		q := a / b
		a, b = b, a%b
		x2, x1 = x1, x2-q*x1
		y2, y1 = y1, y2-q*y1
	}
	if flip {
		x2, y2 = y2, x2
	}
	return a, x2, y2
}

func mulmod(a, b, m int) int {
	sign := 1
	if a < 0 {
		a = -a
		sign *= -1
	}
	if b < 0 {
		b = -b
		sign *= -1
	}
	a = a % m
	b = b % m
	hi, lo := bits.Mul(uint(a), uint(b))
	return sign * int(bits.Rem(hi, lo, uint(m)))
}

type (
	Ring[T any] struct {
		buf []T
		r   int
		w   int
	}
)

func NewRing[T any]() *Ring[T] {
	return &Ring[T]{
		buf: make([]T, 2),
		r:   0,
		w:   0,
	}
}

func (b *Ring[T]) resize() {
	next := make([]T, len(b.buf)*2)
	if b.r == b.w {
		b.w = 0
	} else if b.r < b.w {
		b.w = copy(next, b.buf[b.r:b.w])
	} else {
		p := copy(next, b.buf[b.r:])
		q := 0
		if b.w > 0 {
			q = copy(next[p:], b.buf[:b.w])
		}
		b.w = p + q
	}
	b.buf = next
	b.r = 0
}

func (b *Ring[T]) Write(m T) {
	next := (b.w + 1) % len(b.buf)
	if next == b.r {
		b.resize()
		b.Write(m)
		return
	}
	b.buf[b.w] = m
	b.w = next
}

func (b *Ring[T]) Read() (T, bool) {
	if b.r == b.w {
		var v T
		return v, false
	}
	next := (b.r + 1) % len(b.buf)
	m := b.buf[b.r]
	b.r = next
	return m, true
}

func (b *Ring[T]) Peek() (T, bool) {
	if b.r == b.w {
		var v T
		return v, false
	}
	m := b.buf[b.r]
	return m, true
}
