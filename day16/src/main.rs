use std::fs::File;
use std::io::{prelude::*, BufReader};

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let grid = reader
        .lines()
        .map(|v| Ok::<_, Box<dyn std::error::Error>>(v?.into_bytes()))
        .collect::<Result<Vec<_>, _>>()?;

    let h = grid.len();
    let w = grid[0].len();

    println!(
        "Part 1: {}\nPart 2: {}",
        simulate(
            Beam {
                x: 0,
                y: 0,
                dir: Dir::East
            },
            w,
            h,
            &grid
        ),
        find_largest(w, h, &grid)
    );

    Ok(())
}

fn find_largest(w: usize, h: usize, grid: &Vec<Vec<u8>>) -> usize {
    (0..h).fold(
        (0..w).fold(0, |acc, x| {
            acc.max(simulate(
                Beam {
                    x,
                    y: 0,
                    dir: Dir::South,
                },
                w,
                h,
                grid,
            ))
            .max(simulate(
                Beam {
                    x,
                    y: h - 1,
                    dir: Dir::North,
                },
                w,
                h,
                grid,
            ))
        }),
        |acc, y| {
            acc.max(simulate(
                Beam {
                    x: 0,
                    y,
                    dir: Dir::East,
                },
                w,
                h,
                grid,
            ))
            .max(simulate(
                Beam {
                    x: w - 1,
                    y,
                    dir: Dir::West,
                },
                w,
                h,
                grid,
            ))
        },
    )
}

#[derive(Clone, Copy)]
enum Dir {
    North = 0,
    East = 1,
    South = 2,
    West = 3,
}

struct Beam {
    x: usize,
    y: usize,
    dir: Dir,
}

fn simulate(beam: Beam, w: usize, h: usize, grid: &Vec<Vec<u8>>) -> usize {
    let mut hist = vec![false; w * h];
    let mut beam_hist = vec![false; w * h * 4];

    let mut sum = 0;

    let mut beams = Vec::new();
    beams.push(beam);
    while let Some(beam) = beams.pop() {
        sum += step_beam(w, h, beam, &mut beams, grid, &mut hist, &mut beam_hist);
    }

    sum
}

fn step_beam(
    w: usize,
    h: usize,
    beam: Beam,
    beams: &mut Vec<Beam>,
    grid: &Vec<Vec<u8>>,
    hist: &mut Vec<bool>,
    beam_hist: &mut Vec<bool>,
) -> usize {
    let hkey = beam.y * w + beam.x;
    let key = hkey * 4 + beam.dir as usize;
    if beam_hist[key] {
        return 0;
    }
    match grid[beam.y][beam.x] {
        b'/' => {
            let mut next = beam;
            match next.dir {
                Dir::North => {
                    next.x += 1;
                    next.dir = Dir::East;
                }
                Dir::East => {
                    next.y = next.y.wrapping_sub(1);
                    next.dir = Dir::North;
                }
                Dir::South => {
                    next.x = next.x.wrapping_sub(1);
                    next.dir = Dir::West;
                }
                Dir::West => {
                    next.y += 1;
                    next.dir = Dir::South;
                }
            }
            if is_in_bounds(&next, w, h) {
                beams.push(next);
            }
        }
        b'\\' => {
            let mut next = beam;
            match next.dir {
                Dir::North => {
                    next.x = next.x.wrapping_sub(1);
                    next.dir = Dir::West;
                }
                Dir::East => {
                    next.y += 1;
                    next.dir = Dir::South;
                }
                Dir::South => {
                    next.x += 1;
                    next.dir = Dir::East;
                }
                Dir::West => {
                    next.y = next.y.wrapping_sub(1);
                    next.dir = Dir::North;
                }
            }
            if is_in_bounds(&next, w, h) {
                beams.push(next);
            }
        }
        b'|' => match beam.dir {
            Dir::North => {
                let mut next = beam;
                next.y = next.y.wrapping_sub(1);
                if is_in_bounds(&next, w, h) {
                    beams.push(next);
                }
            }
            Dir::East | Dir::West => {
                {
                    let next = Beam {
                        x: beam.x,
                        y: beam.y.wrapping_sub(1),
                        dir: Dir::North,
                    };
                    if is_in_bounds(&next, w, h) {
                        beams.push(next);
                    }
                }
                {
                    let next = Beam {
                        x: beam.x,
                        y: beam.y + 1,
                        dir: Dir::South,
                    };
                    if is_in_bounds(&next, w, h) {
                        beams.push(next);
                    }
                }
            }
            Dir::South => {
                let mut next = beam;
                next.y += 1;
                if is_in_bounds(&next, w, h) {
                    beams.push(next);
                }
            }
        },
        b'-' => match beam.dir {
            Dir::North | Dir::South => {
                {
                    let next = Beam {
                        x: beam.x.wrapping_sub(1),
                        y: beam.y,
                        dir: Dir::West,
                    };
                    if is_in_bounds(&next, w, h) {
                        beams.push(next);
                    }
                }
                {
                    let next = Beam {
                        x: beam.x + 1,
                        y: beam.y,
                        dir: Dir::East,
                    };
                    if is_in_bounds(&next, w, h) {
                        beams.push(next);
                    }
                }
            }
            Dir::East => {
                let mut next = beam;
                next.x += 1;
                if is_in_bounds(&next, w, h) {
                    beams.push(next);
                }
            }
            Dir::West => {
                let mut next = beam;
                next.x = next.x.wrapping_sub(1);
                if is_in_bounds(&next, w, h) {
                    beams.push(next);
                }
            }
        },
        _ => {
            let mut next = beam;
            match next.dir {
                Dir::North => {
                    next.y = next.y.wrapping_sub(1);
                }
                Dir::East => {
                    next.x += 1;
                }
                Dir::South => {
                    next.y += 1;
                }
                Dir::West => {
                    next.x = next.x.wrapping_sub(1);
                }
            }
            if is_in_bounds(&next, w, h) {
                beams.push(next);
            }
        }
    }
    beam_hist[key] = true;
    if hist[hkey] {
        return 0;
    }
    hist[hkey] = true;
    return 1;
}

fn is_in_bounds(beam: &Beam, w: usize, h: usize) -> bool {
    beam.x < w && beam.y < h
}
