use std::cmp::Reverse;
use std::collections::BinaryHeap;
use std::fs::File;
use std::io::{prelude::*, BufReader};

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let grid = reader
        .lines()
        .map(|line| {
            Ok::<_, Box<dyn std::error::Error>>(
                line?
                    .into_bytes()
                    .into_iter()
                    .map(|v| v - b'0')
                    .collect::<Vec<_>>(),
            )
        })
        .collect::<Result<Vec<_>, _>>()?;

    let h = grid.len();
    let w = grid[0].len();

    let start = Pos { x: 0, y: 0 };
    let goal = Pos { x: w - 1, y: h - 1 };

    let mut graph = GridGraph {
        w,
        h,
        goal,
        grid,
        p2: false,
    };

    let (_, p1) = astar(
        Node {
            value: State {
                pos: start,
                dir: Dir::North,
                same_dir: 0,
            },
            g: 0,
            f: manhattan_distance(&start, &goal),
        },
        &graph,
        &mut ParentSetNoop {},
        &mut ClosedSetVec {
            w,
            vec: vec![0; w * h],
        },
    )
    .ok_or("No path to goal")?;

    graph.p2 = true;

    let (_, p2) = astar(
        Node {
            value: State {
                pos: start,
                dir: Dir::North,
                same_dir: 0,
            },
            g: 0,
            f: manhattan_distance(&start, &goal),
        },
        &graph,
        &mut ParentSetNoop {},
        &mut ClosedSetVec {
            w,
            vec: vec![0; w * h],
        },
    )
    .ok_or("No path to goal")?;

    println!("Part 1: {}\nPart 2: {}", p1, p2);

    Ok(())
}

#[derive(Eq, PartialEq, Clone, Copy)]
struct Pos {
    x: usize,
    y: usize,
}

#[derive(Eq, PartialEq, Clone, Copy)]
enum Dir {
    North = 0,
    East = 1,
    South = 2,
    West = 3,
}

#[derive(Eq, PartialEq, Clone)]
struct State {
    pos: Pos,
    dir: Dir,
    same_dir: u8,
}

impl State {
    fn forward(&self) -> Self {
        let mut next = self.clone();
        match self.dir {
            Dir::North => next.pos.y = next.pos.y.wrapping_sub(1),
            Dir::East => next.pos.x += 1,
            Dir::South => next.pos.y += 1,
            Dir::West => next.pos.x = next.pos.x.wrapping_sub(1),
        }
        next.same_dir += 1;
        next
    }

    fn left(&self) -> Self {
        let mut next = self.clone();
        match self.dir {
            Dir::North => {
                next.dir = Dir::West;
                next.pos.x = next.pos.x.wrapping_sub(1);
            }
            Dir::East => {
                next.dir = Dir::North;
                next.pos.y = next.pos.y.wrapping_sub(1);
            }
            Dir::South => {
                next.dir = Dir::East;
                next.pos.x += 1;
            }
            Dir::West => {
                next.dir = Dir::South;
                next.pos.y += 1;
            }
        }
        next.same_dir = 1;
        next
    }

    fn right(&self) -> Self {
        let mut next = self.clone();
        match self.dir {
            Dir::North => {
                next.dir = Dir::East;
                next.pos.x += 1;
            }
            Dir::East => {
                next.dir = Dir::South;
                next.pos.y += 1;
            }
            Dir::South => {
                next.dir = Dir::West;
                next.pos.x = next.pos.x.wrapping_sub(1);
            }
            Dir::West => {
                next.dir = Dir::North;
                next.pos.y = next.pos.y.wrapping_sub(1);
            }
        }
        next.same_dir = 1;
        next
    }
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
    goal: Pos,
    grid: Vec<Vec<u8>>,
    p2: bool,
}

impl Graph<State> for GridGraph {
    fn is_goal(&self, a: &State) -> bool {
        a.pos == self.goal
    }

    fn get_edges(&self, a: &State, res: &mut Vec<Edge<State>>) {
        if self.p2 {
            if a.same_dir < 10 {
                let s = a.forward();
                if in_bounds(s.pos, self.w, self.h) {
                    let cost = self.grid[s.pos.y][s.pos.x] as usize;
                    let h = manhattan_distance(&s.pos, &self.goal);
                    res.push(Edge { value: s, cost, h });
                }
            }
            if a.same_dir == 0 || a.same_dir >= 4 {
                let s = a.left();
                if in_bounds(s.pos, self.w, self.h) {
                    let cost = self.grid[s.pos.y][s.pos.x] as usize;
                    let h = manhattan_distance(&s.pos, &self.goal);
                    res.push(Edge { value: s, cost, h });
                }
                let s = a.right();
                if in_bounds(s.pos, self.w, self.h) {
                    let cost = self.grid[s.pos.y][s.pos.x] as usize;
                    let h = manhattan_distance(&s.pos, &self.goal);
                    res.push(Edge { value: s, cost, h });
                }
            }
        } else {
            if a.same_dir < 3 {
                let s = a.forward();
                if in_bounds(s.pos, self.w, self.h) {
                    let cost = self.grid[s.pos.y][s.pos.x] as usize;
                    let h = manhattan_distance(&s.pos, &self.goal);
                    res.push(Edge { value: s, cost, h });
                }
            }
            let s = a.left();
            if in_bounds(s.pos, self.w, self.h) {
                let cost = self.grid[s.pos.y][s.pos.x] as usize;
                let h = manhattan_distance(&s.pos, &self.goal);
                res.push(Edge { value: s, cost, h });
            }
            let s = a.right();
            if in_bounds(s.pos, self.w, self.h) {
                let cost = self.grid[s.pos.y][s.pos.x] as usize;
                let h = manhattan_distance(&s.pos, &self.goal);
                res.push(Edge { value: s, cost, h });
            }
        }
    }
}

struct ClosedSetVec {
    w: usize,
    vec: Vec<u64>,
}

impl ClosedSet<State> for ClosedSetVec {
    fn insert(&mut self, a: State) {
        let bit = 1 << a.same_dir * 4 + a.dir as u8;
        self.vec[self.w * a.pos.y + a.pos.x] |= bit;
    }

    fn contains(&self, a: &State) -> bool {
        let bit = 1 << a.same_dir * 4 + a.dir as u8;
        self.vec[self.w * a.pos.y + a.pos.x] & bit != 0
    }
}

struct ParentSetNoop {}

impl ParentSet<State> for ParentSetNoop {
    fn get(&self, _a: &State) -> Option<usize> {
        None
    }

    fn insert(&mut self, _a: &State, _g: usize, _p: &State) {}
}

#[derive(Eq, PartialEq)]
struct Node<T> {
    value: T,
    g: usize,
    f: usize,
}

struct Edge<T> {
    value: T,
    cost: usize,
    h: usize,
}

impl<T: Eq> Ord for Node<T> {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        self.f.cmp(&other.f).then_with(|| self.g.cmp(&other.g))
    }
}

impl<T: Eq> PartialOrd for Node<T> {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        Some(self.cmp(other))
    }
}

trait Graph<T> {
    fn is_goal(&self, a: &T) -> bool;
    fn get_edges(&self, a: &T, res: &mut Vec<Edge<T>>);
}

trait ParentSet<T> {
    fn get(&self, a: &T) -> Option<usize>;
    fn insert(&mut self, a: &T, g: usize, p: &T);
}

trait ClosedSet<T> {
    fn insert(&mut self, a: T);
    fn contains(&self, a: &T) -> bool;
}

fn astar<T: Eq>(
    start: Node<T>,
    graph: &impl Graph<T>,
    parent_set: &mut impl ParentSet<T>,
    closed_set: &mut impl ClosedSet<T>,
) -> Option<(T, usize)> {
    let mut open_set = BinaryHeap::new();
    open_set.push(Reverse(start));
    let mut edges = Vec::new();
    while let Some(Reverse(cur)) = open_set.pop() {
        if closed_set.contains(&cur.value) {
            continue;
        }
        if graph.is_goal(&cur.value) {
            return Some((cur.value, cur.g));
        }
        graph.get_edges(&cur.value, &mut edges);
        while let Some(edge) = edges.pop() {
            if edge.value == cur.value || closed_set.contains(&edge.value) {
                continue;
            }
            let g = cur.g + edge.cost;
            if match parent_set.get(&edge.value) {
                Some(v) => g < v,
                None => true,
            } {
                parent_set.insert(&edge.value, g, &cur.value);
                open_set.push(Reverse(Node {
                    value: edge.value,
                    g,
                    f: g + edge.h,
                }));
            }
        }
        closed_set.insert(cur.value);
    }
    None
}
