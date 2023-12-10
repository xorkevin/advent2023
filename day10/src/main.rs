use once_cell::sync::Lazy;
use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::ops::Index;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut start = None;
    let grid = {
        let mut elements = Vec::new();
        let mut width = 0;
        for (y, line) in reader.lines().enumerate() {
            let k = line?;
            width = k.len();
            if let Some(x) = k.bytes().position(|v| v == b'S') {
                start = Some(Coord { x, y });
            }
            elements.extend(k.bytes());
        }
        Grid::new(elements, width)
    };
    let start = start.ok_or("Missing start symbol")?;
    let (mut cur_pos, mut cur_dir) =
        get_start_neighbor(&grid, start).ok_or("Missing start neighbor")?;
    let mut steps = 1;
    let mut area = match cur_dir {
        Dir::North => 0,
        Dir::East => cur_pos.y as isize,
        Dir::South => 0,
        Dir::West => -(cur_pos.y as isize),
    };
    while cur_pos != start {
        cur_dir = TILE_DIR_MAP
            .get(&grid[cur_pos])
            .ok_or("Invalid pipe path")?[cur_dir]
            .ok_or("Invalid pipe connection")?;
        cur_pos = cur_pos.step(cur_dir);
        steps += 1;
        area += match cur_dir {
            Dir::North => 0,
            Dir::East => cur_pos.y as isize,
            Dir::South => 0,
            Dir::West => -(cur_pos.y as isize),
        }
    }
    if steps % 2 != 0 {
        return Err("Pipe path not aligned to grid".into());
    }
    let half_steps = steps / 2;
    println!("Part 1: {}", half_steps);
    println!("Part 2: {}", area.abs() - half_steps + 1);
    Ok(())
}

struct TileDirMap {
    north: Option<Dir>,
    east: Option<Dir>,
    south: Option<Dir>,
    west: Option<Dir>,
}

impl Index<Dir> for TileDirMap {
    type Output = Option<Dir>;

    fn index(&self, dir: Dir) -> &Self::Output {
        match dir {
            Dir::North => &self.north,
            Dir::East => &self.east,
            Dir::South => &self.south,
            Dir::West => &self.west,
        }
    }
}

static TILE_DIR_MAP: Lazy<HashMap<u8, TileDirMap>> = Lazy::new(|| {
    HashMap::from([
        (
            b'|',
            TileDirMap {
                north: Some(Dir::North),
                east: None,
                south: Some(Dir::South),
                west: None,
            },
        ),
        (
            b'-',
            TileDirMap {
                north: None,
                east: Some(Dir::East),
                south: None,
                west: Some(Dir::West),
            },
        ),
        (
            b'L',
            TileDirMap {
                north: None,
                east: None,
                south: Some(Dir::East),
                west: Some(Dir::North),
            },
        ),
        (
            b'J',
            TileDirMap {
                north: None,
                east: Some(Dir::North),
                south: Some(Dir::West),
                west: None,
            },
        ),
        (
            b'7',
            TileDirMap {
                north: Some(Dir::West),
                east: Some(Dir::South),
                south: None,
                west: None,
            },
        ),
        (
            b'F',
            TileDirMap {
                north: Some(Dir::East),
                east: None,
                south: None,
                west: Some(Dir::South),
            },
        ),
    ])
});

#[derive(Debug, Clone, Copy)]
enum Dir {
    North,
    East,
    South,
    West,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
struct Coord {
    x: usize,
    y: usize,
}

impl Coord {
    fn step(&self, dir: Dir) -> Self {
        let &Self { x, y } = self;
        match dir {
            Dir::North => Self {
                x,
                y: y.wrapping_sub(1),
            },
            Dir::East => Self { x: self.x + 1, y },
            Dir::South => Self { x, y: self.y + 1 },
            Dir::West => Self {
                x: x.wrapping_sub(1),
                y,
            },
        }
    }
}

struct Grid {
    width: usize,
    height: usize,
    elements: Vec<u8>,
}

impl Grid {
    fn new(elements: Vec<u8>, width: usize) -> Self {
        Self {
            width,
            height: elements.len() / width,
            elements,
        }
    }
}

impl Index<Coord> for Grid {
    type Output = u8;

    fn index(&self, pos: Coord) -> &Self::Output {
        if pos.x >= self.width || pos.y >= self.height {
            &0
        } else {
            &self.elements[self.width * pos.y + pos.x]
        }
    }
}

fn get_start_neighbor(grid: &Grid, pos: Coord) -> Option<(Coord, Dir)> {
    let top = pos.step(Dir::North);
    match grid[top] {
        b'|' | b'7' | b'F' => return Some((top, Dir::North)),
        _ => {}
    }
    let right = pos.step(Dir::East);
    match grid[right] {
        b'-' | b'J' | b'7' => return Some((right, Dir::East)),
        _ => {}
    }
    let bot = pos.step(Dir::South);
    match grid[bot] {
        b'|' | b'L' | b'J' => return Some((bot, Dir::South)),
        _ => {}
    }
    let left = pos.step(Dir::West);
    match grid[left] {
        b'-' | b'L' | b'F' => return Some((left, Dir::West)),
        _ => {}
    }
    None
}
