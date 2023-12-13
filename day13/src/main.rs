use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut sum = 0;
    let mut sum2 = 0;

    let mut grid = Vec::new();

    for line in reader.lines() {
        let line = line?;
        if line != "" {
            grid.push(line.bytes().collect::<Vec<_>>());
            continue;
        }
        let (s, s2) = find_reflections(&grid).ok_or("No mirror")?;
        sum += s;
        sum2 += s2;
        grid = Vec::new();
    }

    let (s, s2) = find_reflections(&grid).ok_or("No mirror")?;
    sum += s;
    sum2 += s2;

    println!("Part 1: {}\nPart 2: {}", sum, sum2);

    Ok(())
}

fn find_reflections(grid: &Vec<Vec<u8>>) -> Option<(usize, usize)> {
    let tgrid = transpose(grid);
    let (s, eh, ev) = match find_reflection(grid, &tgrid) {
        Some(v) => v,
        None => return None,
    };
    let s2 = match find_smudge_reflection(grid, &tgrid, eh, ev) {
        Some(v) => v,
        None => return None,
    };
    Some((s, s2))
}

fn find_reflection(
    grid: &Vec<Vec<u8>>,
    transpose: &Vec<Vec<u8>>,
) -> Option<(usize, Option<usize>, Option<usize>)> {
    if let Some(v) = find_mirror(grid) {
        return Some((v * 100, Some(v), None));
    }
    if let Some(v) = find_mirror(transpose) {
        return Some((v, None, Some(v)));
    }
    None
}

fn find_smudge_reflection(
    grid: &Vec<Vec<u8>>,
    transpose: &Vec<Vec<u8>>,
    eh: Option<usize>,
    ev: Option<usize>,
) -> Option<usize> {
    if let Some(v) = find_smudge(grid, eh) {
        return Some(v * 100);
    }
    if let Some(v) = find_smudge(transpose, ev) {
        return Some(v);
    }
    None
}

fn find_smudge(grid: &Vec<Vec<u8>>, except: Option<usize>) -> Option<usize> {
    for i in 1..grid.len() {
        if let Some(v) = except {
            if i == v {
                continue;
            }
        }
        if is_almost_mirrored_at(grid, i) {
            return Some(i);
        }
    }
    None
}

fn is_almost_mirrored_at(grid: &Vec<Vec<u8>>, r: usize) -> bool {
    let height = grid.len();
    let lim = r.min(height - r);
    let mut has_diff = false;
    for i in 0..lim {
        let a = &grid[r - i - 1];
        let b = &grid[r + i];
        if a != b {
            if !is_edit_dist_1(a, b) {
                return false;
            }
            if has_diff {
                return false;
            }
            has_diff = true
        }
    }
    has_diff
}

fn is_edit_dist_1(a: &Vec<u8>, b: &Vec<u8>) -> bool {
    let mut has_diff = false;
    for (&a, &b) in a.iter().zip(b.iter()) {
        if a != b {
            if has_diff {
                return false;
            }
            has_diff = true
        }
    }
    has_diff
}

fn find_mirror(grid: &Vec<Vec<u8>>) -> Option<usize> {
    for i in 1..grid.len() {
        if is_mirrored_at(grid, i) {
            return Some(i);
        }
    }
    None
}

fn is_mirrored_at(grid: &Vec<Vec<u8>>, r: usize) -> bool {
    let height = grid.len();
    let lim = r.min(height - r);
    for i in 0..lim {
        if grid[r - i - 1] != grid[r + i] {
            return false;
        }
    }
    true
}

fn transpose(grid: &Vec<Vec<u8>>) -> Vec<Vec<u8>> {
    let height = grid.len();
    let width = grid[0].len();
    let mut res = vec![vec![0; height]; width];
    for i in 0..width {
        for j in 0..height {
            res[i][j] = grid[j][i];
        }
    }
    res
}
