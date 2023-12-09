use once_cell::sync::Lazy;
use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

static DIGIT_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"-?\d+").unwrap());

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut s1 = 0;
    let mut s2 = 0;
    for line in reader.lines() {
        let mut digits = DIGIT_REGEX
            .find_iter(&line?)
            .map(|v| v.as_str().parse::<i32>())
            .collect::<Result<Vec<_>, _>>()?;
        s1 += find_next_seq(&digits);
        digits.reverse();
        s2 += find_next_seq(&digits);
    }

    println!("Part 1: {}\nPart 2: {}", s1, s2);

    Ok(())
}

fn find_next_seq(nums: &[i32]) -> i32 {
    if nums.iter().all(|&v| v == 0) {
        0
    } else if nums.len() == 1 {
        *nums.first().unwrap()
    } else {
        nums.last().unwrap()
            + find_next_seq(&nums.windows(2).map(|v| v[1] - v[0]).collect::<Vec<_>>())
    }
}
