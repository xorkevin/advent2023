use once_cell::sync::Lazy;
use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

static DIGIT_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"\d+").unwrap());

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut sum1 = 0;
    let mut sum2 = 0;
    for line in reader.lines() {
        let mut possible = true;
        let mut max_red = 0;
        let mut max_green = 0;
        let mut max_blue = 0;
        let line = line?;
        let (a, b) = line.split_once(": ").ok_or("Invalid line")?;
        let game_num = DIGIT_REGEX
            .find(a)
            .and_then(|v| v.as_str().parse::<u32>().ok())
            .ok_or("Invalid game num")?;
        for v in b.split("; ").flat_map(|v| v.split(", ")) {
            let (a, b) = v.split_once(" ").ok_or("Invalid cubes")?;
            let count = a.parse::<u32>()?;
            match b {
                "red" => {
                    if count > 12 {
                        possible = false;
                    }
                    max_red = max_red.max(count);
                }
                "green" => {
                    if count > 13 {
                        possible = false;
                    }
                    max_green = max_green.max(count);
                }
                "blue" => {
                    if count > 14 {
                        possible = false;
                    }
                    max_blue = max_blue.max(count);
                }
                _ => Err("Invalid color")?,
            }
        }
        if possible {
            sum1 += game_num;
        }
        sum2 += max_red * max_green * max_blue;
    }

    println!("Part 1: {}", sum1);
    println!("Part 2: {}", sum2);
    Ok(())
}
