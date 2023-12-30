use std::cmp::Ordering;
use std::fs::File;
use std::io::{prelude::*, BufReader};

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut lines = Vec::new();

    let mut min_x = 0;
    let mut min_y = 0;
    let mut max_x = 0;
    let mut max_y = 0;
    let mut first = true;

    for line in reader.lines() {
        let line = line?;
        let (lhs, rhs) = line.split_once('~').ok_or("Invalid line")?;
        let lhs_num_strs = lhs.split(',').collect::<Vec<_>>();
        if lhs_num_strs.len() != 3 {
            return Err("Invalid line".into());
        }
        let mut lhs_pos = Pos {
            x: lhs_num_strs[0].parse::<isize>()?,
            y: lhs_num_strs[1].parse::<isize>()?,
            z: lhs_num_strs[2].parse::<usize>()?,
        };
        let rhs_num_strs = rhs.split(',').collect::<Vec<_>>();
        if rhs_num_strs.len() != 3 {
            return Err("Invalid line".into());
        }
        let mut rhs_pos = Pos {
            x: rhs_num_strs[0].parse::<isize>()?,
            y: rhs_num_strs[1].parse::<isize>()?,
            z: rhs_num_strs[2].parse::<usize>()?,
        };
        if cmp_pos(&lhs_pos, &rhs_pos) == Ordering::Greater {
            (lhs_pos, rhs_pos) = (rhs_pos, lhs_pos);
        }
        if first {
            first = false;
            min_x = lhs_pos.x;
            min_y = lhs_pos.y;
            max_x = rhs_pos.x;
            max_y = rhs_pos.y;
        } else {
            min_x = min_x.min(lhs_pos.x);
            min_y = min_y.min(lhs_pos.y);
            max_x = max_x.max(rhs_pos.x);
            max_y = max_y.max(rhs_pos.y);
        }
        let height = rhs_pos.z - lhs_pos.z;
        lines.push(Line {
            from: lhs_pos,
            to: rhs_pos,
            height,
        })
    }

    lines.sort_by(|a, b| cmp_pos(&a.from, &b.from));

    let xw = (max_x - min_x + 1) as usize;
    let yw = (max_y - min_y + 1) as usize;

    let mut count = 0;
    let mut sum = 0;

    let mut height_map = vec![0; xw * yw];
    let mut full_tower = vec![
        Line {
            from: Pos { x: 0, y: 0, z: 0 },
            to: Pos { x: 0, y: 0, z: 0 },
            height: 0
        };
        lines.len()
    ];
    place_tower(
        min_x,
        min_y,
        xw,
        &mut height_map,
        &lines,
        lines.len(),
        &mut full_tower,
    );
    clear_height_map(&mut height_map);
    let mut candidate = vec![
        Line {
            from: Pos { x: 0, y: 0, z: 0 },
            to: Pos { x: 0, y: 0, z: 0 },
            height: 0
        };
        lines.len()
    ];
    for i in 0..lines.len() {
        place_tower(min_x, min_y, xw, &mut height_map, &lines, i, &mut candidate);
        clear_height_map(&mut height_map);
        let delta = get_tower_delta(&full_tower, &candidate, i);
        if delta == 0 {
            count += 1;
        } else {
            sum += delta;
        }
    }

    println!("Part 1: {}\nPart 2: {}", count, sum);

    Ok(())
}

#[derive(Clone)]
struct Line {
    from: Pos,
    to: Pos,
    height: usize,
}

#[derive(Clone, Copy)]
struct Pos {
    x: isize,
    y: isize,
    z: usize,
}

fn get_tower_delta(a: &[Line], b: &[Line], omit: usize) -> usize {
    let mut count = 0;
    for (n, (i, j)) in a.iter().zip(b.iter()).enumerate() {
        if n == omit {
            continue;
        }
        if i.from.z != j.from.z {
            count += 1;
        }
    }
    count
}

fn place_tower(
    min_x: isize,
    min_y: isize,
    xw: usize,
    height_map: &mut Vec<usize>,
    lines: &[Line],
    omit: usize,
    res: &mut Vec<Line>,
) {
    for (n, i) in lines.iter().enumerate() {
        if n == omit {
            continue;
        }
        let Line {
            mut from,
            mut to,
            height,
        } = i;
        let mut highest = 0;
        for y in from.y..=to.y {
            for x in from.x..=to.x {
                let key = pos_key(x, y, min_x, min_y, xw);
                highest = highest.max(height_map[key]);
            }
        }
        from.z = highest + 1;
        to.z = from.z + height;
        for y in from.y..=to.y {
            for x in from.x..=to.x {
                let key = pos_key(x, y, min_x, min_y, xw);
                height_map[key] = to.z;
            }
        }
        res[n] = Line {
            from,
            to,
            height: *height,
        };
    }
}

fn clear_height_map(height_map: &mut Vec<usize>) {
    for i in height_map.iter_mut() {
        *i = 0;
    }
}

fn pos_key(x: isize, y: isize, min_x: isize, min_y: isize, xw: usize) -> usize {
    (y - min_y) as usize * xw + (x - min_x) as usize
}

fn cmp_pos(a: &Pos, b: &Pos) -> Ordering {
    a.z.cmp(&b.z)
        .then_with(|| a.y.cmp(&b.y).then_with(|| a.x.cmp(&b.x)))
}
