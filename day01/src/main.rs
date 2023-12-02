use once_cell::sync::Lazy;
use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

static DIGIT_ONLY_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"\d").unwrap());
static DIGIT_REGEX: Lazy<Regex> =
    Lazy::new(|| Regex::new(r"\d|one|two|three|four|five|six|seven|eight|nine").unwrap());
static REV_DIGIT_REGEX: Lazy<Regex> =
    Lazy::new(|| Regex::new(r"\d|enin|thgie|neves|xis|evif|ruof|eerht|owt|eno").unwrap());

fn word_to_val(s: &str) -> Option<u32> {
    match s {
        "one" => Some(1),
        "two" => Some(2),
        "three" => Some(3),
        "four" => Some(4),
        "five" => Some(5),
        "six" => Some(6),
        "seven" => Some(7),
        "eight" => Some(8),
        "nine" => Some(9),
        "enin" => Some(9),
        "thgie" => Some(8),
        "neves" => Some(7),
        "xis" => Some(6),
        "evif" => Some(5),
        "ruof" => Some(4),
        "eerht" => Some(3),
        "owt" => Some(2),
        "eno" => Some(1),
        v => v.parse::<u32>().ok(),
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut sum1 = 0;
    let mut sum2 = 0;
    for line in reader.lines() {
        let line = line?;
        let rev = line.chars().rev().collect::<String>();
        {
            let first = DIGIT_ONLY_REGEX
                .find(&line)
                .and_then(|v| v.as_str().parse::<u32>().ok())
                .ok_or("Invalid digit")?;
            let last = DIGIT_ONLY_REGEX
                .find(&rev)
                .and_then(|v| v.as_str().parse::<u32>().ok())
                .ok_or("Invalid digit")?;
            sum1 += first * 10 + last;
        }
        {
            let first = DIGIT_REGEX
                .find(&line)
                .and_then(|v| word_to_val(v.as_str()))
                .ok_or("Invalid digit")?;
            let last = REV_DIGIT_REGEX
                .find(&rev)
                .and_then(|v| word_to_val(v.as_str()))
                .ok_or("Invalid digit")?;
            sum2 += first * 10 + last;
        }
    }

    println!("Part 1: {}", sum1);
    println!("Part 2: {}", sum2);
    Ok(())
}
