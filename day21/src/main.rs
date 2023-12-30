use std::collections::VecDeque;
use std::fs::File;
use std::io::{prelude::*, BufReader};

const PUZZLEINPUT: &str = "input.txt";

const TARGET: usize = 26501365;
const TARGET_IS_EVEN: bool = TARGET % 2 == 0;
const P1_TARGET: usize = 64;
const P1_TARGET_IS_EVEN: bool = P1_TARGET % 2 == 0;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut start = Pos { x: 0, y: 0 };
    let mut grid = Vec::new();
    for line in reader.lines() {
        let line = line?.into_bytes();
        if let Some(x) = line.iter().position(|&v| v == b'S') {
            start = Pos { x, y: grid.len() };
        }
        grid.push(line);
    }

    let h = grid.len();
    let w = grid[0].len();
    let multiple = TARGET / h;
    let rem = TARGET % h;

    let graph = GridGraph { w, h, grid };
    let mut closed_set = ClosedSetVec {
        w,
        vec: vec![false; w * h],
        rem,
        start,
        sum: 0,
        inner_even: 0,
        inner_odd: 0,
        corner_even: 0,
        corner_odd: 0,
    };
    bfs(start, &graph, &mut closed_set);
    println!("Part 1: {}", closed_set.sum);

    if h != w {
        return Err("Grid is not square".into());
    }
    if h % 2 != 1 || start.y != (h - 1) / 2 || start.x != start.y {
        return Err("Start is not centered".into());
    }

    let multiple_is_even = multiple % 2 == 0;
    let multiple1 = multiple + 1;
    let outer_multiple = multiple1 * multiple1;
    let inner_multiple = multiple * multiple;

    let (outer_diamond, inner_diamond, outer_corner, inner_corner) =
        if TARGET_IS_EVEN == multiple_is_even {
            (
                closed_set.inner_even,
                closed_set.inner_odd,
                closed_set.corner_even,
                closed_set.corner_odd,
            )
        } else {
            (
                closed_set.inner_odd,
                closed_set.inner_even,
                closed_set.corner_odd,
                closed_set.corner_even,
            )
        };

    println!(
        "Part 2: {}",
        outer_multiple * outer_diamond
            + inner_multiple * inner_diamond
            + (outer_multiple - multiple1) * outer_corner
            + (inner_multiple + multiple) * inner_corner
    );

    Ok(())
}

#[derive(Eq, PartialEq, Clone, Copy)]
struct Pos {
    x: usize,
    y: usize,
}

fn in_bounds(a: Pos, w: usize, h: usize) -> bool {
    a.x < w && a.y < h
}

fn manhattan_distance(a: &Pos, b: &Pos) -> usize {
    a.x.abs_diff(b.x) + a.y.abs_diff(b.y)
}

struct GridGraph {
    w: usize,
    h: usize,
    grid: Vec<Vec<u8>>,
}

impl Graph<Pos> for GridGraph {
    fn is_goal(&self, _a: &Pos) -> bool {
        false
    }

    fn get_edges(&self, &a: &Pos, res: &mut Vec<Pos>) {
        let mut s = a;
        s.y = s.y.wrapping_sub(1);
        if in_bounds(s, self.w, self.h) && self.grid[s.y][s.x] != b'#' {
            res.push(s);
        }
        let mut s = a;
        s.x += 1;
        if in_bounds(s, self.w, self.h) && self.grid[s.y][s.x] != b'#' {
            res.push(s);
        }
        let mut s = a;
        s.y += 1;
        if in_bounds(s, self.w, self.h) && self.grid[s.y][s.x] != b'#' {
            res.push(s);
        }
        let mut s = a;
        s.x = s.x.wrapping_sub(1);
        if in_bounds(s, self.w, self.h) && self.grid[s.y][s.x] != b'#' {
            res.push(s);
        }
    }
}

struct ClosedSetVec {
    w: usize,
    vec: Vec<bool>,
    rem: usize,
    start: Pos,
    sum: usize,
    inner_even: usize,
    inner_odd: usize,
    corner_even: usize,
    corner_odd: usize,
}

impl ClosedSet<Pos> for ClosedSetVec {
    fn insert(&mut self, a: &Pos, g: usize) {
        self.vec[self.w * a.y + a.x] = true;
        let cur_is_even = g % 2 == 0;
        if g <= P1_TARGET && cur_is_even == P1_TARGET_IS_EVEN {
            self.sum += 1;
        }
        if manhattan_distance(&self.start, a) > self.rem {
            if cur_is_even {
                self.corner_even += 1;
            } else {
                self.corner_odd += 1;
            }
        } else {
            if cur_is_even {
                self.inner_even += 1;
            } else {
                self.inner_odd += 1;
            }
        }
    }

    fn insert_parent(&mut self, _a: &Pos, _g: usize, _p: &Pos) {}

    fn contains(&self, a: &Pos) -> bool {
        self.vec[self.w * a.y + a.x]
    }
}

#[derive(Eq, PartialEq)]
struct Node<T> {
    value: T,
    g: usize,
}

trait Graph<T> {
    fn is_goal(&self, a: &T) -> bool;
    fn get_edges(&self, a: &T, res: &mut Vec<T>);
}

trait ClosedSet<T> {
    fn insert(&mut self, a: &T, g: usize);
    fn insert_parent(&mut self, a: &T, g: usize, p: &T);
    fn contains(&self, a: &T) -> bool;
}

fn bfs<T: Eq>(
    start: T,
    graph: &impl Graph<T>,
    closed_set: &mut impl ClosedSet<T>,
) -> Option<(T, usize)> {
    let mut open_set = VecDeque::new();
    closed_set.insert(&start, 0);
    open_set.push_back(Node { value: start, g: 0 });
    let mut edges = Vec::new();
    loop {
        let current = match open_set.pop_front() {
            Some(v) => v,
            None => return None,
        };
        if graph.is_goal(&current.value) {
            return Some((current.value, current.g));
        }
        graph.get_edges(&current.value, &mut edges);
        while let Some(value) = edges.pop() {
            if closed_set.contains(&value) {
                continue;
            }
            let g = current.g + 1;
            closed_set.insert(&value, g);
            closed_set.insert_parent(&value, g, &current.value);
            open_set.push_back(Node { value, g });
        }
    }
}
