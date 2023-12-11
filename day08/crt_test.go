package main

import (
	"strconv"
	"testing"
)

func TestExtGCD(t *testing.T) {
	for tcn, tc := range []struct {
		a, b    int
		g, p, q int
	}{
		{
			a: 240,
			b: 46,
			g: 2,
			p: -9,
			q: 47,
		},
		{
			a: 3,
			b: 4,
			g: 1,
			p: -1,
			q: 1,
		},
	} {
		tc := tc
		t.Run("ext gcd test case "+strconv.Itoa(tcn), func(t *testing.T) {
			g, p, q := extGCD(tc.a, tc.b)
			if g != tc.g || p != tc.p || q != tc.q {
				t.Fatalf("Invalid output g %d != %d p %d != %d q %d != %d", g, tc.g, p, tc.p, q, tc.q)
			}
			x := tc.a*p + tc.b*q
			if x != g {
				t.Fatalf("Bezout identity does not hold %d != %d", x, g)
			}
		})
	}
}

func TestCRT(t *testing.T) {
	for tcn, tc := range []struct {
		inp []struct{ a, m int }
		exp int
	}{
		{
			inp: []struct{ a, m int }{
				{a: 0, m: 3},
				{a: 3, m: 4},
			},
			exp: 3,
		},
		{
			inp: []struct{ a, m int }{
				{a: 0, m: 3},
				{a: 3, m: 4},
				{a: 4, m: 5},
			},
			exp: 39,
		},
		{
			inp: []struct{ a, m int }{
				{a: 1, m: 2},
				{a: 0, m: 3},
				{a: 3, m: 4},
				{a: 4, m: 5},
			},
			exp: 39,
		},
		{
			inp: []struct{ a, m int }{
				{a: 2, m: 3},
				{a: 3, m: 5},
				{a: 2, m: 7},
			},
			exp: 23,
		},
		{
			inp: []struct{ a, m int }{
				{a: 0, m: 17},
				{a: 11, m: 13},
				{a: 16, m: 19},
			},
			exp: 3417,
		},
		{
			inp: []struct{ a, m int }{
				{a: 0, m: 7},
				{a: 12, m: 13},
				{a: 55, m: 59},
				{a: 25, m: 31},
				{a: 12, m: 19},
			},
			exp: 1068781,
		},
	} {
		tc := tc
		t.Run("crt test case "+strconv.Itoa(tcn), func(t *testing.T) {
			if len(tc.inp) == 0 {
				t.Fatalf("Invalid crt input")
			}

			a, m := tc.inp[0].a, tc.inp[0].m
			for _, i := range tc.inp[1:] {
				a1, m1, ok := crt(a, m, i.a, i.m)
				if !ok {
					t.Fatalf("Incompatible a1=%d m1=%d a2=%d m2=%d", a, m, i.a, i.m)
				}
				a = a1
				m = m1
			}

			if a != tc.exp {
				t.Fatalf("Invalid output %d != %d", a, tc.exp)
			}
		})
	}
}
