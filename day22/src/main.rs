use std::cmp::{Ordering, Reverse};
use std::collections::BinaryHeap;
use std::fs::File;
use std::io::{prelude::*, BufReader};
use std::slice::Iter;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut lines = Vec::new();

    let mut min_x = 0;
    let mut min_y = 0;
    let mut max_x = 0;
    let mut max_y = 0;
    let mut first = true;

    for line in reader.lines() {
        let line = line?;
        let (lhs, rhs) = line.split_once('~').ok_or("Invalid line")?;
        let lhs_num_strs = lhs.split(',').collect::<Vec<_>>();
        if lhs_num_strs.len() != 3 {
            return Err("Invalid line".into());
        }
        let mut lhs_pos = Pos {
            x: lhs_num_strs[0].parse::<isize>()?,
            y: lhs_num_strs[1].parse::<isize>()?,
            z: lhs_num_strs[2].parse::<usize>()?,
        };
        let rhs_num_strs = rhs.split(',').collect::<Vec<_>>();
        if rhs_num_strs.len() != 3 {
            return Err("Invalid line".into());
        }
        let mut rhs_pos = Pos {
            x: rhs_num_strs[0].parse::<isize>()?,
            y: rhs_num_strs[1].parse::<isize>()?,
            z: rhs_num_strs[2].parse::<usize>()?,
        };
        if cmp_pos(&lhs_pos, &rhs_pos) == Ordering::Greater {
            (lhs_pos, rhs_pos) = (rhs_pos, lhs_pos);
        }
        if first {
            first = false;
            min_x = lhs_pos.x;
            min_y = lhs_pos.y;
            max_x = rhs_pos.x;
            max_y = rhs_pos.y;
        } else {
            min_x = min_x.min(lhs_pos.x);
            min_y = min_y.min(lhs_pos.y);
            max_x = max_x.max(rhs_pos.x);
            max_y = max_y.max(rhs_pos.y);
        }
        let height = rhs_pos.z - lhs_pos.z;
        lines.push(Line {
            from: lhs_pos,
            to: rhs_pos,
            height,
        })
    }

    lines.sort_unstable_by(|a, b| cmp_pos(&a.from, &b.from));

    let xw = (max_x - min_x + 1) as usize;
    let yw = (max_y - min_y + 1) as usize;

    let mut full_tower = vec![
        Line {
            from: Pos { x: 0, y: 0, z: 0 },
            to: Pos { x: 0, y: 0, z: 0 },
            height: 0
        };
        lines.len()
    ];
    let (p1, support_up, support_by) =
        place_tower(min_x, min_y, xw, yw, &lines, lines.len(), &mut full_tower);
    let mut closed_set = VecBitSet::new(lines.len());
    let mut collapse_set = VecBitSet::new(lines.len());
    let mut open_set = BinaryHeap::new();
    println!(
        "Part 1: {}\nPart 2: {}",
        p1,
        (0..lines.len()).fold(0, |acc, i| {
            acc + cascade_tower(
                &support_up,
                &support_by,
                &mut closed_set,
                &mut collapse_set,
                &mut open_set,
                i,
            )
        })
    );

    Ok(())
}

#[derive(Clone)]
struct Line {
    from: Pos,
    to: Pos,
    height: usize,
}

#[derive(Clone, Copy)]
struct Pos {
    x: isize,
    y: isize,
    z: usize,
}

#[derive(Clone)]
struct Slot {
    h: usize,
    id: usize,
}

fn cascade_tower(
    support_up: &[VecSet<usize>],
    support_by: &[VecSet<usize>],
    closed_set: &mut VecBitSet,
    collapse_set: &mut VecBitSet,
    open_set: &mut BinaryHeap<Reverse<usize>>,
    idx: usize,
) -> usize {
    closed_set.insert(idx);
    collapse_set.insert(idx);
    for &v in support_up[idx].iter() {
        if closed_set.contains(v) {
            continue;
        }
        closed_set.insert(v);
        open_set.push(Reverse(v));
    }
    while let Some(Reverse(cur)) = open_set.pop() {
        let supports = &support_by[cur];
        if supports
            .iter()
            .filter(|&&v| collapse_set.contains(v))
            .count()
            < supports.len()
        {
            continue;
        }
        collapse_set.insert(cur);
        for &v in support_up[cur].iter() {
            if closed_set.contains(v) {
                continue;
            }
            closed_set.insert(v);
            open_set.push(Reverse(v));
        }
    }
    let res = collapse_set.len() - 1;
    closed_set.zero();
    collapse_set.zero();
    res
}

fn place_tower(
    min_x: isize,
    min_y: isize,
    xw: usize,
    yw: usize,
    lines: &[Line],
    omit: usize,
    res: &mut Vec<Line>,
) -> (usize, Vec<VecSet<usize>>, Vec<VecSet<usize>>) {
    let mut critical_supports = VecBitSet::new(lines.len());
    let mut support_up = vec![VecSet::new(); lines.len()];
    let mut support_by = Vec::with_capacity(lines.len());
    let mut height_map = vec![Slot { h: 0, id: 0 }; xw * yw];
    for (n, i) in lines.iter().enumerate() {
        if n == omit {
            continue;
        }
        let Line {
            mut from,
            mut to,
            height,
        } = i;
        let mut highest = 0;
        let mut supports = VecSet::new();
        for y in from.y..=to.y {
            for x in from.x..=to.x {
                let key = pos_key(x, y, min_x, min_y, xw);
                let v = &height_map[key];
                if v.h > highest {
                    highest = v.h;
                    supports.clear();
                    supports.insert(v.id);
                } else if highest != 0 && v.h == highest {
                    supports.insert(v.id);
                }
            }
        }
        if supports.len() == 1 {
            for &v in supports.iter() {
                critical_supports.insert(v);
            }
        }
        for &v in supports.iter() {
            support_up[v].insert(n);
        }
        support_by.push(supports);
        from.z = highest + 1;
        to.z = from.z + height;
        for y in from.y..=to.y {
            for x in from.x..=to.x {
                let key = pos_key(x, y, min_x, min_y, xw);
                height_map[key] = Slot { h: to.z, id: n };
            }
        }
        res[n] = Line {
            from,
            to,
            height: *height,
        };
    }
    (
        lines.len() - critical_supports.len(),
        support_up,
        support_by,
    )
}

fn pos_key(x: isize, y: isize, min_x: isize, min_y: isize, xw: usize) -> usize {
    (y - min_y) as usize * xw + (x - min_x) as usize
}

fn cmp_pos(a: &Pos, b: &Pos) -> Ordering {
    a.z.cmp(&b.z)
        .then_with(|| a.y.cmp(&b.y).then_with(|| a.x.cmp(&b.x)))
}

struct VecBitSet {
    bits: Vec<usize>,
    size: usize,
}

#[allow(dead_code)]
impl VecBitSet {
    fn new(size: usize) -> Self {
        Self {
            bits: vec![0; (size + (usize::BITS as usize - 1)) / usize::BITS as usize],
            size: 0,
        }
    }

    fn len(&self) -> usize {
        self.size
    }

    fn zero(&mut self) {
        for i in self.bits.iter_mut() {
            *i = 0;
        }
        self.size = 0;
    }

    fn contains(&self, i: usize) -> bool {
        let a = i / usize::BITS as usize;
        let b = i % usize::BITS as usize;
        let mask = 1 << b;
        self.bits[a] & mask != 0
    }

    fn insert(&mut self, i: usize) {
        let a = i / usize::BITS as usize;
        let b = i % usize::BITS as usize;
        let mask = 1 << b;
        if self.bits[a] & mask == 0 {
            self.bits[a] |= mask;
            self.size += 1;
        }
    }

    fn remove(&mut self, i: usize) {
        let a = i / usize::BITS as usize;
        let b = i % usize::BITS as usize;
        let mask = 1 << b;
        if self.bits[a] & mask != 0 {
            self.bits[a] &= !mask;
            self.size -= 1;
        }
    }
}

#[derive(Clone)]
struct VecSet<T> {
    elements: Vec<T>,
}

impl<T> VecSet<T> {
    fn new() -> Self {
        Self {
            elements: Vec::new(),
        }
    }

    fn clear(&mut self) {
        self.elements.clear();
    }

    fn len(&self) -> usize {
        self.elements.len()
    }

    fn iter(&self) -> Iter<T> {
        self.elements.iter()
    }
}

#[allow(dead_code)]
impl<T: Eq> VecSet<T> {
    fn contains(&self, i: &T) -> bool {
        self.elements.contains(i)
    }

    fn insert(&mut self, i: T) {
        if !self.elements.contains(&i) {
            self.elements.push(i);
        }
    }
}

#[allow(dead_code)]
impl<T: Eq + Ord> VecSet<T> {
    fn sort_unstable(&mut self) {
        self.elements.sort_unstable();
    }

    fn contains_sorted(&self, i: &T) -> bool {
        self.elements.binary_search(i).is_ok()
    }
}
