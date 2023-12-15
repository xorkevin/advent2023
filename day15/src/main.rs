use std::fs::File;
use std::io::prelude::*;

const PUZZLEINPUT: &str = "input.txt";

struct BoxValue {
    id: Vec<u8>,
    v: u8,
}

struct LightBox {
    values: Vec<BoxValue>,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut file = File::open(PUZZLEINPUT)?;

    let mut buf = Vec::new();
    file.read_to_end(&mut buf)?;
    let left = buf
        .iter()
        .position(|v| !v.is_ascii_whitespace())
        .unwrap_or(buf.len());
    let right = buf
        .iter()
        .rposition(|v| !v.is_ascii_whitespace())
        .unwrap_or(buf.len());
    let buf = &buf[left..=right];

    let mut boxes: [LightBox; 256] = std::array::from_fn(|_| LightBox { values: Vec::new() });
    let mut sum = 0;
    for i in buf.split(|&v| v == b',') {
        sum += hash_word(i) as usize;
        process_instr(&i, &mut boxes)?;
    }

    println!(
        "Part 1: {}\nPart 2: {}",
        sum,
        boxes.iter().enumerate().fold(0, |acc, (n, i)| acc
            + (n + 1)
                * i.values
                    .iter()
                    .enumerate()
                    .fold(0, |acc, (k, j)| acc + (k + 1) * (j.v as usize)))
    );

    Ok(())
}

fn process_instr(
    instr: &[u8],
    boxes: &mut [LightBox; 256],
) -> Result<(), Box<dyn std::error::Error>> {
    if let Some(label) = instr.strip_suffix(&[b'-']) {
        let h = hash_word(label);
        let values = &mut boxes[h as usize].values;
        if let Some(n) = values.iter().position(|v| v.id == label) {
            values.remove(n);
        }
        return Ok(());
    }
    if let Some(n) = instr.iter().position(|&v| v == b'=') {
        let (label, b) = instr.split_at(n);
        let num = std::str::from_utf8(&b[1..])?.parse::<u8>()?;
        let h = hash_word(label);
        let values = &mut boxes[h as usize].values;
        if let Some(n) = values.iter().position(|v| v.id == label) {
            values[n].v = num;
        } else {
            values.push(BoxValue {
                id: label.to_vec(),
                v: num,
            });
        }
        return Ok(());
    }
    Err("Invalid instr".into())
}

fn hash_word(v: &[u8]) -> u8 {
    v.iter().fold(0, |acc, &v| hash_step(acc, v))
}

fn hash_step(v: u8, b: u8) -> u8 {
    (((v as u32 + b as u32) * 17) % 256) as u8
}
