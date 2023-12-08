use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

struct Node {
    id: String,
    left: String,
    right: String,
    is_end: bool,
}

struct NodeVisit<'a> {
    id: &'a str,
    count: usize,
    cycle: usize,
    candidate: bool,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut steps = Vec::new();
    let mut nodes = HashMap::new();
    let mut starts = Vec::new();
    for line in reader.lines() {
        let line = line?;
        if steps.len() == 0 {
            steps = line.into_bytes();
            continue;
        }
        if line == "" {
            continue;
        }
        let (lhs, rhs) = line.split_once(" = ").ok_or("Invalid line")?;
        let discard_chars: &[_] = &['(', ')'];
        let (a, b) = rhs
            .trim_matches(discard_chars)
            .split_once(", ")
            .ok_or("Invalid line")?;
        nodes.insert(
            lhs.to_owned(),
            Node {
                id: lhs.to_owned(),
                left: a.to_owned(),
                right: b.to_owned(),
                is_end: lhs.ends_with("Z"),
            },
        );
        if lhs.ends_with("A") {
            starts.push(lhs.to_owned());
        }
    }

    let mut start = nodes.get("AAA").ok_or("No start node")?;
    let end = "ZZZ";
    let mut count = 0;
    while start.id != end {
        let next = if steps[count % steps.len()] == b'L' {
            &start.left
        } else {
            &start.right
        };
        start = nodes.get(next).ok_or("Next node not found")?;
        count += 1;
    }

    println!("Part 1: {}", count);

    let mut start_nodes = starts
        .into_iter()
        .map(|v| nodes.get(&v).ok_or("Start node not found"))
        .collect::<Result<Vec<_>, _>>()?;
    let total_starts = start_nodes.len();
    let mut revisits = (0..total_starts)
        .map(|_| NodeVisit {
            id: "",
            count: 0,
            cycle: 0,
            candidate: false,
        })
        .collect::<Vec<_>>();
    let mut total_revisits = 0;
    let mut count = 0;
    while total_revisits < total_starts {
        let instr = steps[count % steps.len()];
        count += 1;
        for (i, v) in start_nodes.iter_mut().enumerate() {
            let next = if instr == b'L' { &v.left } else { &v.right };
            *v = nodes.get(next).ok_or("Next node not found")?;
            if v.is_end {
                let visit = revisits.get_mut(i).ok_or("Revisit does not exist")?;
                if visit.id == "" {
                    *visit = NodeVisit {
                        id: &v.id,
                        count,
                        cycle: 0,
                        candidate: false,
                    };
                } else if !visit.candidate {
                    if v.id != visit.id {
                        return Err("Multiple terminals for cycle".into());
                    }
                    let cycle = count - visit.count;
                    let rem = count % cycle;
                    if rem != 0 {
                        return Err("Remainder not zero".into());
                    }
                    *visit = NodeVisit {
                        id: visit.id,
                        count,
                        cycle,
                        candidate: true,
                    };
                    total_revisits += 1;
                } else {
                    if v.id != visit.id {
                        return Err("Multiple terminals for cycle".into());
                    }
                    let cycle = count - visit.count;
                    if cycle != visit.cycle {
                        return Err("Multiple cycle lengths".into());
                    }
                    visit.count = count;
                }
            }
        }
    }

    println!(
        "Part 2: {}",
        revisits.into_iter().fold(1, |acc, v| lcm(acc, v.cycle))
    );

    Ok(())
}

fn lcm(a: usize, b: usize) -> usize {
    a * b / gcd(a, b)
}

fn gcd(mut a: usize, mut b: usize) -> usize {
    if a > b {
        (a, b) = (b, a);
    }
    while a != 0 {
        (a, b) = (b % a, a);
    }
    b
}
