use once_cell::sync::Lazy;
use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::ops::Index;
use std::ops::IndexMut;

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
    let mut pipe_grid = Grid::new(vec![0; grid.width * grid.height], grid.width);
    let mut bounds = vec![None; grid.height];

    let (mut na, mut nb) = get_start_neighbors(&grid, start).ok_or("Invalid neighbors")?;
    {
        let start_char = get_start_char(na.dir, nb.dir).ok_or("Invalid neighbors")?;
        pipe_grid[start] = start_char;
        bounds[start.y] = Some(Coord {
            x: start.x,
            y: start.x,
        });
        pipe_grid[na.coord] = grid[na.coord];
        pipe_grid[nb.coord] = grid[nb.coord];
        match bounds[na.coord.y] {
            Some(Coord { x, y }) => {
                bounds[na.coord.y] = Some(Coord {
                    x: x.min(na.coord.x),
                    y: y.max(na.coord.x),
                });
            }
            None => {
                bounds[na.coord.y] = Some(Coord {
                    x: na.coord.x,
                    y: na.coord.x,
                });
            }
        }
        match bounds[nb.coord.y] {
            Some(Coord { x, y }) => {
                bounds[nb.coord.y] = Some(Coord {
                    x: x.min(nb.coord.x),
                    y: y.max(nb.coord.x),
                });
            }
            None => {
                bounds[nb.coord.y] = Some(Coord {
                    x: nb.coord.x,
                    y: nb.coord.x,
                });
            }
        }
    }
    let mut steps = 1;
    while na.coord != nb.coord {
        {
            let dir = TILE_DIR_MAP
                .get(&grid[na.coord])
                .ok_or("Invalid pipe path 1")?[na.dir]
                .ok_or("Invalid pipe path 2")?;
            let coord = na.coord.step(dir);
            na = CoordDir { coord, dir };
        }
        {
            let dir = TILE_DIR_MAP
                .get(&grid[nb.coord])
                .ok_or("Invalid pipe path 3")?[nb.dir]
                .ok_or("Invalid pipe path 4")?;
            let coord = nb.coord.step(dir);
            nb = CoordDir { coord, dir };
        }
        pipe_grid[na.coord] = grid[na.coord];
        pipe_grid[nb.coord] = grid[nb.coord];
        match bounds[na.coord.y] {
            Some(Coord { x, y }) => {
                bounds[na.coord.y] = Some(Coord {
                    x: x.min(na.coord.x),
                    y: y.max(na.coord.x),
                });
            }
            None => {
                bounds[na.coord.y] = Some(Coord {
                    x: na.coord.x,
                    y: na.coord.x,
                });
            }
        }
        match bounds[nb.coord.y] {
            Some(Coord { x, y }) => {
                bounds[nb.coord.y] = Some(Coord {
                    x: x.min(nb.coord.x),
                    y: y.max(nb.coord.x),
                });
            }
            None => {
                bounds[nb.coord.y] = Some(Coord {
                    x: nb.coord.x,
                    y: nb.coord.x,
                });
            }
        }
        steps += 1;
    }
    println!("Part 1: {}", steps);

    let pipe_grid = pipe_grid;

    let mut sum = 0;
    for i in 0..pipe_grid.height {
        if let Some(b) = bounds[i] {
            let a = pipe_grid.width * i;
            sum += get_inside_area(&pipe_grid.elements[a + b.x..a + b.y]);
        }
    }
    println!("Part 2: {}", sum);

    Ok(())
}

fn get_inside_area(row: &[u8]) -> usize {
    let mut count = 0;
    let mut inside = false;
    let mut prev = 0;
    for i in row {
        match i {
            0 => {
                if inside {
                    count += 1
                }
            }
            b'|' => {
                inside = !inside;
                prev = 0;
            }
            b'L' => prev = b'L',
            b'J' => {
                if prev == b'F' {
                    inside = !inside;
                }
                prev = 0;
            }
            b'7' => {
                if prev == b'L' {
                    inside = !inside;
                }
                prev = 0;
            }
            b'F' => prev = b'F',
            _ => {}
        }
    }
    count
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

impl IndexMut<Coord> for Grid {
    fn index_mut(&mut self, pos: Coord) -> &mut Self::Output {
        &mut self.elements[self.width * pos.y + pos.x]
    }
}

#[derive(Debug, Clone, Copy)]
struct CoordDir {
    coord: Coord,
    dir: Dir,
}

fn get_start_neighbors(grid: &Grid, pos: Coord) -> Option<(CoordDir, CoordDir)> {
    let mut neighbors = Vec::with_capacity(4);
    let top = pos.step(Dir::North);
    match grid[top] {
        b'|' | b'7' | b'F' => neighbors.push(CoordDir {
            coord: top,
            dir: Dir::North,
        }),
        _ => {}
    }
    let right = pos.step(Dir::East);
    match grid[right] {
        b'-' | b'J' | b'7' => neighbors.push(CoordDir {
            coord: right,
            dir: Dir::East,
        }),
        _ => {}
    }
    let bot = pos.step(Dir::South);
    match grid[bot] {
        b'|' | b'L' | b'J' => neighbors.push(CoordDir {
            coord: bot,
            dir: Dir::South,
        }),
        _ => {}
    }
    let left = pos.step(Dir::West);
    match grid[left] {
        b'-' | b'L' | b'F' => neighbors.push(CoordDir {
            coord: left,
            dir: Dir::West,
        }),
        _ => {}
    }
    if let &[a, b] = &neighbors[..] {
        Some((a, b))
    } else {
        None
    }
}

fn get_start_char(mut a: Dir, mut b: Dir) -> Option<u8> {
    if a as isize > b as isize {
        (a, b) = (b, a)
    }
    match (a, b) {
        (Dir::North, Dir::East) => Some(b'L'),
        (Dir::North, Dir::South) => Some(b'|'),
        (Dir::North, Dir::West) => Some(b'J'),
        (Dir::East, Dir::South) => Some(b'F'),
        (Dir::East, Dir::West) => Some(b'-'),
        (Dir::South, Dir::West) => Some(b'7'),
        _ => None,
    }
}
