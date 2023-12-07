use std::cmp::Ordering;
use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

struct CardHand {
    kind: u8,
    kind_j: u8,
    score: u64,
    score_j: u64,
    bid: usize,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut hands = reader
        .lines()
        .map(|line| {
            let line = line?;
            let (a, b) = line.split_once(" ").ok_or("Invalid line")?;
            if a.len() != 5 {
                return Err("Invalid hand".into());
            }
            let bid = b.parse::<usize>()?;
            let a = a.as_bytes();
            Ok::<_, Box<dyn std::error::Error>>(CardHand {
                kind: hand_kind(a, false),
                kind_j: hand_kind(a, true),
                score: score_hand(a, false),
                score_j: score_hand(a, true),
                bid,
            })
        })
        .collect::<Result<Vec<_>, _>>()?;

    hands.sort_unstable_by(|a, b| match a.kind.cmp(&b.kind) {
        Ordering::Equal => a.score.cmp(&b.score),
        v => v,
    });
    println!(
        "Part 1: {}",
        hands
            .iter()
            .enumerate()
            .fold(0, |acc, (i, v)| acc + (i + 1) * v.bid)
    );

    hands.sort_unstable_by(|a, b| match a.kind_j.cmp(&b.kind_j) {
        Ordering::Equal => a.score_j.cmp(&b.score_j),
        v => v,
    });
    println!(
        "Part 2: {}",
        hands
            .iter()
            .enumerate()
            .fold(0, |acc, (i, v)| acc + (i + 1) * v.bid)
    );

    Ok(())
}

fn score_card(b: u8, with_joker: bool) -> u64 {
    match b {
        b'T' => 9,
        b'J' => {
            if with_joker {
                0
            } else {
                10
            }
        }
        b'Q' => 11,
        b'K' => 12,
        b'A' => 13,
        b if b >= b'2' && b <= b'9' => (b - b'2' + 1) as u64,
        _ => 0,
    }
}

fn hand_kind(b: &[u8], with_joker: bool) -> u8 {
    let (card_counts, jokers) = b
        .iter()
        .fold((HashMap::new(), 0u8), |(mut acc, jokers), &v| {
            if with_joker && v == b'J' {
                return (acc, jokers + 1);
            }
            let e = acc.entry(v).or_insert(0u8);
            *e += 1;
            (acc, jokers)
        });
    if card_counts.len() == 0 {
        return 6;
    }
    let mut card_counts = card_counts.into_iter().map(|v| v.1).collect::<Vec<_>>();
    card_counts.sort_unstable_by(|a, b| b.cmp(a));
    card_counts[0] += jokers;

    if card_counts[0] == 5 {
        // 5 of kind
        6
    } else if card_counts[0] == 4 {
        // 4 of kind
        5
    } else if card_counts[0] == 3 {
        if card_counts[1] == 2 {
            // full house
            4
        } else {
            // 3 of kind
            3
        }
    } else if card_counts[0] == 2 {
        if card_counts[1] == 2 {
            // two pair
            2
        } else {
            1
        }
    } else {
        0
    }
}

fn score_hand(b: &[u8], with_joker: bool) -> u64 {
    b.iter()
        .fold(0, |acc, &v| acc * 14 + score_card(v, with_joker))
}
