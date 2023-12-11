use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

#[derive(Clone, Copy)]
struct Coord {
    x: usize,
    y: usize,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut empty_rows = Vec::new();
    let mut coords = Vec::new();
    let mut cols = Vec::new();
    for (y, line) in reader.lines().enumerate() {
        let line = line?;
        if cols.len() == 0 {
            cols = vec![Vec::new(); line.len()];
        }
        let mut is_empty = true;
        for (x, i) in line.bytes().enumerate() {
            if i == b'#' {
                is_empty = false;
                coords.push(Coord { x, y });
            }
            cols[x].push(i);
        }
        if is_empty {
            empty_rows.push(y);
        }
    }

    let coords = coords;
    let empty_rows = empty_rows;
    let empty_cols = cols
        .into_iter()
        .enumerate()
        .flat_map(|(x, i)| if !i.contains(&b'#') { Some(x) } else { None })
        .collect::<Vec<_>>();

    let (dist, expansion) = coords
        .iter()
        .enumerate()
        .flat_map(|(n, i)| coords[n + 1..].iter().map(move |j| (i, j)))
        .fold((0, 0), |(sd, se), (a, b)| {
            let dist = manhattan_distance(&a, &b);
            let expansion =
                calc_expansion(&empty_rows, a.y, b.y) + calc_expansion(&empty_cols, a.x, b.x);
            (sd + dist, se + expansion)
        });
    println!(
        "Part 1: {}\nPart 2: {}",
        dist + expansion,
        dist + expansion * 999999
    );

    Ok(())
}

fn calc_expansion(empty: &[usize], a: usize, b: usize) -> usize {
    let (a, b) = (a.min(b), a.max(b));
    let left = empty.iter().position(|&v| v > a).unwrap_or(empty.len());
    let right = empty
        .iter()
        .rev()
        .position(|&v| v < b)
        .unwrap_or(empty.len());
    empty.len() - right - left
}

fn dist(a: usize, b: usize) -> usize {
    a.max(b) - a.min(b)
}

fn manhattan_distance(a: &Coord, b: &Coord) -> usize {
    dist(a.x, b.x) + dist(a.y, b.y)
}
