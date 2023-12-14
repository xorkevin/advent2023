use sha2::{Digest, Sha256};
use std::collections::hash_map::Entry;
use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut grid = reader
        .lines()
        .map(|v| Ok::<_, Box<dyn std::error::Error>>(v?.into_bytes()))
        .collect::<Result<Vec<_>, _>>()?;

    let height = grid.len();
    let width = grid[0].len();
    let mut other = vec![vec![0u8; height]; width];

    drop_rocks(&mut grid, &mut other);

    println!("Part 1: {}", score_rocks(&grid));

    let mut remaining = 0;
    const P2_ITERATIONS: usize = 1000000000;
    let mut cache = HashMap::new();
    for i in 0..P2_ITERATIONS {
        cycle(&mut grid, &mut other);
        let s = get_state(&grid);
        match cache.entry(s) {
            Entry::Occupied(v) => {
                let v = v.get();
                let p = i - v;
                remaining = (P2_ITERATIONS - i - 1) % p;
                break;
            }
            Entry::Vacant(v) => v.insert(i),
        };
    }
    for _ in 0..remaining {
        cycle(&mut grid, &mut other);
    }

    println!("Part 2: {}", score_rocks(&grid));

    Ok(())
}

fn get_state(grid: &Vec<Vec<u8>>) -> [u8; 32] {
    let mut hasher = Sha256::new();
    for i in grid {
        hasher.update(i);
    }
    hasher.finalize().into()
}

fn score_rocks(grid: &Vec<Vec<u8>>) -> usize {
    let height = grid.len();
    let mut sum = 0;
    for (r, i) in grid.iter().enumerate() {
        for &j in i {
            if j == b'O' {
                sum += height - r;
            }
        }
    }
    sum
}

fn cycle(grid: &mut Vec<Vec<u8>>, other: &mut Vec<Vec<u8>>) {
    // west
    drop_rocks(grid, other);
    // south
    drop_rocks(other, grid);
    // east
    drop_rocks(grid, other);
    // north
    drop_rocks(other, grid);
}

fn drop_rocks(grid: &mut Vec<Vec<u8>>, other: &mut Vec<Vec<u8>>) {
    let height = grid.len();
    for r in 0..height {
        for c in 0..grid[0].len() {
            let b = grid[r][c];
            if b == b'O' {
                drop_rock(grid, other, r, c);
            } else {
                other[c][height - r - 1] = b;
            }
        }
    }
}

fn drop_rock(grid: &mut Vec<Vec<u8>>, other: &mut Vec<Vec<u8>>, r: usize, c: usize) {
    let height = grid.len();
    let mut rest = r;
    while rest >= 1 {
        if grid[rest - 1][c] != b'.' {
            break;
        }
        rest -= 1;
    }
    grid[r][c] = b'.';
    grid[rest][c] = b'O';
    other[c][height - r - 1] = b'.';
    other[c][height - rest - 1] = b'O';
}
