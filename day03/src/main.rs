use once_cell::sync::Lazy;
use regex::Regex;
use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

static DIGIT_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"\d+").unwrap());

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let grid = {
        let mut grid = Vec::new();
        for line in reader.lines() {
            grid.push(line?);
        }
        grid
    };

    let mut sum = 0;

    let mut gears = HashMap::new();

    let lim = (grid[0].len(), grid.len());

    let mut buf = vec![(0, 0); lim.1 * 2 + 2];

    for (y, i) in grid.iter().enumerate() {
        for j in DIGIT_REGEX.find_iter(i) {
            let left = (j.start(), y);
            let right = (j.end() - 1, y);
            let num = j.as_str().parse::<usize>()?;
            let mut has_sym = false;
            let n = get_neighbors(left, right, lim, &mut buf);
            for &(x, y) in buf.iter().take(n) {
                let sym = grid[y].as_bytes()[x];
                if is_symbol(sym) {
                    has_sym = true;
                    if sym == b'*' {
                        let id = y * lim.0 + x;
                        let entry = gears.entry(id).or_insert(Vec::new());
                        entry.push(num);
                    }
                }
            }
            if has_sym {
                sum += num;
            }
        }
    }
    println!("Part 1: {}", sum);

    let sum2 = gears
        .values()
        .flat_map(|v| match &v[..] {
            [a, b] => Some(a * b),
            _ => None,
        })
        .sum::<usize>();
    println!("Part 2: {}", sum2);
    Ok(())
}

fn get_neighbors(
    p1: (usize, usize),
    p2: (usize, usize),
    lim: (usize, usize),
    buf: &mut [(usize, usize)],
) -> usize {
    let mut idx = 0;
    if p1.1 > 0 {
        let y = p1.1 - 1;
        for i in p1.0..=p2.0 {
            buf[idx] = (i, y);
            idx += 1;
        }
    }
    {
        let y = p2.1 + 1;
        if y < lim.1 {
            for i in p1.0..=p2.0 {
                buf[idx] = (i, y);
                idx += 1;
            }
        }
    }
    if p1.0 > 0 {
        let x = p1.0 - 1;
        for i in p1.1..=p2.1 {
            buf[idx] = (x, i);
            idx += 1;
        }
        if p1.1 > 0 {
            buf[idx] = (x, p1.1 - 1);
            idx += 1;
        }
        let y = p2.1 + 1;
        if y < lim.1 {
            buf[idx] = (x, y);
            idx += 1;
        }
    }
    {
        let x = p2.0 + 1;
        if x < lim.0 {
            for i in p1.1..=p2.1 {
                buf[idx] = (x, i);
                idx += 1;
            }
            if p1.1 > 0 {
                buf[idx] = (x, p1.1 - 1);
                idx += 1;
            }
            let y = p2.1 + 1;
            if y < lim.1 {
                buf[idx] = (x, y);
                idx += 1;
            }
        }
    }
    idx
}

fn is_symbol(b: u8) -> bool {
    return b != b'.' && (b < b'0' || b > b'9');
}
