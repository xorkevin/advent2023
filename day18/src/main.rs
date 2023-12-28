use std::fs::File;
use std::io::{prelude::*, BufReader};

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut current = Pos { x: 0, y: 0 };
    let mut current2 = Pos { x: 0, y: 0 };
    let mut area = 0;
    let mut area2 = 0;
    let mut perimeter = 0;
    let mut perimeter2 = 0;

    for line in reader.lines() {
        let line = line?;
        let (p1, p2) = line.rsplit_once(' ').ok_or("Invalid line")?;
        {
            let (dir_str, num_str) = p1.split_once(' ').ok_or("Invalid line")?;
            let dir = dir_str.as_bytes()[0];
            let num = num_str.parse::<isize>()?;
            current = move_dir(current, dir, num);
            match dir {
                b'L' => area -= current.y * num,
                b'R' => area += current.y * num,
                _ => {}
            }
            perimeter += num;
        }
        {
            let p2 = &p2[2..8];
            let dir = p2.as_bytes()[5] - b'0';
            let num = isize::from_str_radix(&p2[..5], 16)?;
            current2 = move_dir2(current2, dir, num);
            match dir {
                2 => area2 -= current2.y * num,
                0 => area2 += current2.y * num,
                _ => {}
            }
            perimeter2 += num;
        }
    }

    if perimeter % 2 != 0 {
        return Err("Perimeter not aligned to grid".into());
    }
    if perimeter2 % 2 != 0 {
        return Err("Perimeter not aligned to grid".into());
    }
    area = area.abs();
    area2 = area2.abs();

    println!(
        "Part 1: {}\nPart 2: {}",
        area + perimeter / 2 + 1,
        area2 + perimeter2 / 2 + 1,
    );

    Ok(())
}

struct Pos {
    x: isize,
    y: isize,
}

fn move_dir(mut p: Pos, dir: u8, num: isize) -> Pos {
    match dir {
        b'U' => p.y -= num,
        b'D' => p.y += num,
        b'L' => p.x -= num,
        b'R' => p.x += num,
        _ => {}
    }
    p
}

fn move_dir2(mut p: Pos, dir: u8, num: isize) -> Pos {
    match dir {
        3 => p.y -= num,
        1 => p.y += num,
        2 => p.x -= num,
        0 => p.x += num,
        _ => {}
    }
    p
}
