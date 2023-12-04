use once_cell::sync::Lazy;
use regex::Regex;
use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

static DIGIT_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"\d+").unwrap());

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut sum = 0;
    let mut total_cards = 0;
    let mut bonus_cards = vec![0; 256];

    for line in reader.lines() {
        let line = line?;
        let (a, b) = line.split_once(": ").ok_or("Invalid line")?;
        let card_num = DIGIT_REGEX
            .find(a)
            .and_then(|v| v.as_str().parse::<usize>().ok())
            .ok_or("Invalid card num")?;
        let (a, b) = b.split_once(" | ").ok_or("Invalid line")?;
        let mut winning = HashSet::new();
        for i in DIGIT_REGEX.find_iter(a) {
            winning.insert(i.as_str().parse::<u32>()?);
        }
        let mut count = 0;
        let mut points = 0;
        for i in DIGIT_REGEX.find_iter(b) {
            let num = i.as_str().parse::<u32>()?;
            if winning.contains(&num) {
                count += 1;
                points = if points == 0 { 1 } else { points * 2 };
            }
        }
        sum += points;

        let current_multiplier = bonus_cards[card_num] + 1;
        total_cards += current_multiplier;
        for i in card_num + 1..card_num + 1 + count {
            bonus_cards[i] += current_multiplier;
        }
    }

    println!("Part 1: {}", sum);
    println!("Part 2: {}", total_cards);

    Ok(())
}
