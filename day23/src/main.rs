use std::collections::hash_map::Entry;
use std::collections::{HashMap, VecDeque};
use std::fs::File;
use std::hash::Hash;
use std::io::{prelude::*, BufReader};

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let grid = reader
        .lines()
        .map(|line| Ok::<_, Box<dyn std::error::Error>>(line?.into_bytes()))
        .collect::<Result<Vec<_>, _>>()?;

    let h = grid.len();
    let w = grid[0].len();

    let start = Pos {
        x: grid[0].iter().position(|&v| v == b'.').ok_or("No start")?,
        y: 0,
    };
    let end = Pos {
        x: grid[h - 1]
            .iter()
            .position(|&v| v == b'.')
            .ok_or("No end")?,
        y: h - 1,
    };

    let (undirected_graph, directed_graph) = contract_paths(start, end, &GridGraph { w, h, grid });

    println!(
        "Part 1: {}\nPart 2: {}",
        directed_graph.cost + search_longest_path(&directed_graph),
        undirected_graph.cost + search_longest_path(&undirected_graph)
    );

    Ok(())
}

struct ContractedGraph {
    start: usize,
    end: usize,
    graph: Vec<Vec<(usize, usize)>>,
    cost: usize,
}

#[derive(Clone, Copy)]
struct Node<T> {
    value: T,
    g: usize,
    depth: usize,
}

fn search_longest_path(graph: &ContractedGraph) -> usize {
    let mut max_g = 0;
    let mut closed_set = VecBitSet::new(graph.graph.len());
    let mut open_set = Vec::new();
    let mut cur_path = Vec::new();
    open_set.push(Node {
        value: graph.start,
        g: 0,
        depth: 0,
    });
    while let Some(cur) = open_set.pop() {
        if cur.value == graph.end {
            max_g = max_g.max(cur.g);
            continue;
        }

        while cur_path.len() > cur.depth {
            let v: Node<_> = cur_path.pop().unwrap();
            closed_set.remove(v.value);
        }
        closed_set.insert(cur.value);
        cur_path.push(cur);

        for &(i, g) in &graph.graph[cur.value] {
            if closed_set.contains(i) {
                continue;
            }
            open_set.push(Node {
                value: i,
                g: cur.g + g,
                depth: cur.depth + 1,
            });
        }
    }
    max_g
}

#[derive(Clone, Copy, PartialEq, Eq, Hash)]
struct Pos {
    x: usize,
    y: usize,
}

fn contract_paths(start: Pos, end: Pos, graph: &GridGraph) -> (ContractedGraph, ContractedGraph) {
    let mut node_names = NameIDMap::new();
    let start_id = node_names.get(start.y * graph.w + start.x);
    let end_id = node_names.get(end.y * graph.w + end.x);
    let mut contracted_graph = HashMap::new();
    let mut contracted_directed_graph = HashMap::new();

    let mut closed_set = VecBitSet::new(graph.w * graph.h);
    let mut open_set = VecDeque::new();
    let mut explored_set = VecBitSet::new(graph.w * graph.h);
    open_set.push_back(start);
    explored_set.insert(start.y * graph.w + start.x);
    let mut edges = Vec::new();
    let mut candidate_edges = Vec::new();
    while let Some(cur) = open_set.pop_front() {
        let cur_key = cur.y * graph.w + cur.x;
        let cur_id = node_names.get(cur_key);
        closed_set.insert(cur_key);
        let num_edges = get_path_edges(
            cur,
            end,
            &graph,
            &mut edges,
            &mut candidate_edges,
            &mut closed_set,
        );
        if cur != start && num_edges < 2 && (edges.len() == 0 || edges[0].pos != end) {
            // eliminate dead ends
            continue;
        }
        while let Some(o) = edges.pop() {
            let key = o.pos.y * graph.w + o.pos.x;
            let id = node_names.get(key);
            if id != end_id && !explored_set.contains(key) {
                open_set.push_back(o.pos);
                explored_set.insert(key);
            }
            match contracted_graph.entry((cur_id, id)) {
                Entry::Occupied(mut v) => {
                    let v = v.get_mut();
                    if o.cost > *v {
                        *v = o.cost;
                    }
                }
                Entry::Vacant(v) => {
                    v.insert(o.cost);
                }
            }
            match contracted_graph.entry((id, cur_id)) {
                Entry::Occupied(mut v) => {
                    let v = v.get_mut();
                    if o.cost > *v {
                        *v = o.cost;
                    }
                }
                Entry::Vacant(v) => {
                    v.insert(o.cost);
                }
            }
            if o.forward {
                match contracted_directed_graph.entry((cur_id, id)) {
                    Entry::Occupied(mut v) => {
                        let v = v.get_mut();
                        if o.cost > *v {
                            *v = o.cost;
                        }
                    }
                    Entry::Vacant(v) => {
                        v.insert(o.cost);
                    }
                }
            }
            if o.rev {
                match contracted_directed_graph.entry((id, cur_id)) {
                    Entry::Occupied(mut v) => {
                        let v = v.get_mut();
                        if o.cost > *v {
                            *v = o.cost;
                        }
                    }
                    Entry::Vacant(v) => {
                        v.insert(o.cost);
                    }
                }
            }
        }
    }

    let num_nodes = node_names.len();
    let undirected_graph = contracted_graph.into_iter().fold(
        vec![Vec::new(); num_nodes],
        |mut acc, ((from, to), v)| {
            acc[from].push((to, v));
            acc
        },
    );
    let directed_graph = contracted_directed_graph.into_iter().fold(
        vec![Vec::new(); num_nodes],
        |mut acc, ((from, to), v)| {
            acc[from].push((to, v));
            acc
        },
    );
    let (contracted_start_id, start_cost) = find_branch(start_id, &undirected_graph);
    let (contracted_end_id, end_cost) = find_branch(end_id, &undirected_graph);
    (
        ContractedGraph {
            start: contracted_start_id,
            end: contracted_end_id,
            graph: undirected_graph,
            cost: start_cost + end_cost,
        },
        ContractedGraph {
            start: contracted_start_id,
            end: contracted_end_id,
            graph: directed_graph,
            cost: start_cost + end_cost,
        },
    )
}

fn find_branch(a: usize, graph: &Vec<Vec<(usize, usize)>>) -> (usize, usize) {
    let edges = &graph[a];
    if edges.len() != 1 {
        return (a, 0);
    }
    *edges.first().unwrap()
}

fn get_path_edges(
    start: Pos,
    end: Pos,
    graph: &GridGraph,
    edges: &mut Vec<GridPath>,
    candidate_edges: &mut Vec<GridPath>,
    closed_set: &mut VecBitSet,
) -> usize {
    let num_edges = graph.get_edges(start, closed_set, edges);
    for i in 0..edges.len() {
        let mut cur = edges[i];
        loop {
            if cur.pos == end {
                edges[i] = cur;
                break;
            }
            if graph.get_edges(cur.pos, closed_set, candidate_edges) != 2 {
                candidate_edges.clear();
                edges[i] = cur;
                break;
            }
            let key = cur.pos.y * graph.w + cur.pos.x;
            closed_set.insert(key);
            // at least one edge must be present since nodes with two edges are always traversed
            let first = candidate_edges.pop().unwrap();
            cur = GridPath {
                pos: first.pos,
                cost: cur.cost + first.cost,
                forward: cur.forward && first.forward,
                rev: cur.rev && first.rev,
            };
        }
    }
    num_edges
}

struct GridGraph {
    w: usize,
    h: usize,
    grid: Vec<Vec<u8>>,
}

#[derive(Clone, Copy)]
struct GridPath {
    pos: Pos,
    cost: usize,
    forward: bool,
    rev: bool,
}

impl GridGraph {
    fn get_edges(&self, a: Pos, closed_set: &VecBitSet, res: &mut Vec<GridPath>) -> usize {
        let cur = self.grid[a.y][a.x];
        let mut num_edges = 0;
        let mut pos = a;
        pos.y = pos.y.wrapping_sub(1);
        if in_bounds(pos, self.w, self.h) {
            let b = self.grid[pos.y][pos.x];
            if b != b'#' {
                num_edges += 1;
                let key = pos.y * self.w + pos.x;
                if !closed_set.contains(key) {
                    res.push(GridPath {
                        pos,
                        cost: 1,
                        forward: cur == b'.' || cur == b'^',
                        rev: b == b'.' || b == b'v',
                    });
                }
            }
        }
        let mut pos = a;
        pos.x += 1;
        if in_bounds(pos, self.w, self.h) {
            let b = self.grid[pos.y][pos.x];
            if b != b'#' {
                num_edges += 1;
                let key = pos.y * self.w + pos.x;
                if !closed_set.contains(key) {
                    res.push(GridPath {
                        pos,
                        cost: 1,
                        forward: cur == b'.' || cur == b'>',
                        rev: b == b'.' || b == b'<',
                    });
                }
            }
        }
        let mut pos = a;
        pos.y += 1;
        if in_bounds(pos, self.w, self.h) {
            let b = self.grid[pos.y][pos.x];
            if b != b'#' {
                num_edges += 1;
                let key = pos.y * self.w + pos.x;
                if !closed_set.contains(key) {
                    res.push(GridPath {
                        pos,
                        cost: 1,
                        forward: cur == b'.' || cur == b'v',
                        rev: b == b'.' || b == b'^',
                    });
                }
            }
        }
        let mut pos = a;
        pos.x = pos.x.wrapping_sub(1);
        if in_bounds(pos, self.w, self.h) {
            let b = self.grid[pos.y][pos.x];
            if b != b'#' {
                num_edges += 1;
                let key = pos.y * self.w + pos.x;
                if !closed_set.contains(key) {
                    res.push(GridPath {
                        pos,
                        cost: 1,
                        forward: cur == b'.' || cur == b'<',
                        rev: b == b'.' || b == b'>',
                    });
                }
            }
        }
        num_edges
    }
}

fn in_bounds(a: Pos, w: usize, h: usize) -> bool {
    a.x < w && a.y < h
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

struct NameIDMap<T> {
    names: HashMap<T, usize>,
}

impl<T: Eq + Hash> NameIDMap<T> {
    fn new() -> Self {
        Self {
            names: HashMap::new(),
        }
    }

    fn len(&self) -> usize {
        self.names.len()
    }

    fn get(&mut self, a: T) -> usize {
        let next = self.names.len();
        *self.names.entry(a).or_insert(next)
    }
}
