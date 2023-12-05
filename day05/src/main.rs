use once_cell::sync::Lazy;
use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

static DIGIT_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"\d+").unwrap());

struct Range {
    start: i64,
    end: i64,
}

struct Range2 {
    dest: Range,
    src: Range,
}

impl Range2 {
    fn start_in_src(&self, n: i64) -> bool {
        n >= self.src.start && n < self.src.end
    }

    fn end_in_src(&self, n: i64) -> bool {
        n > self.src.start && n <= self.src.end
    }

    fn to_dest(&self, n: i64) -> i64 {
        n - self.src.start + self.dest.start
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut seeds = Vec::new();
    let mut seeds2 = Vec::new();
    let mut range_maps = Vec::new();
    let mut last_range_map = Vec::new();

    for line in reader.lines() {
        let line = line?;
        if line.starts_with("seeds:") {
            seeds = DIGIT_REGEX
                .find_iter(&line)
                .map(|v| v.as_str().parse::<i64>())
                .collect::<Result<Vec<_>, _>>()?;
            seeds2 = seeds
                .chunks(2)
                .map(|v| match v {
                    &[a, b] => Ok(Range {
                        start: a,
                        end: a + b,
                    }),
                    _ => Err("odd number of seeds"),
                })
                .collect::<Result<Vec<_>, _>>()?;
            continue;
        } else if line.ends_with("map:") {
            if last_range_map.len() != 0 {
                range_maps.push(last_range_map);
                last_range_map = Vec::new();
            }
            continue;
        } else if line == "" {
            continue;
        }
        last_range_map.push(
            if let [a, b, c] = DIGIT_REGEX
                .find_iter(&line)
                .map(|v| v.as_str().parse::<i64>())
                .collect::<Result<Vec<_>, _>>()?[..]
            {
                Range2 {
                    dest: Range {
                        start: a,
                        end: a + c,
                    },
                    src: Range {
                        start: b,
                        end: b + c,
                    },
                }
            } else {
                Err("Invalid range")?
            },
        );
    }
    if last_range_map.len() != 0 {
        range_maps.push(last_range_map);
    }

    println!(
        "Part 1: {}",
        range_maps
            .iter()
            .fold(seeds, |acc, v| run_range(&acc, v))
            .iter()
            .min()
            .ok_or("No seeds")?
    );

    println!(
        "Part 2: {}",
        range_maps
            .iter()
            .fold(seeds2, |acc, v| run_range2(&acc, v))
            .iter()
            .min_by_key(|v| v.start)
            .ok_or("No seeds")?
            .start
    );

    Ok(())
}

fn run_range(seeds: &[i64], range_maps: &[Range2]) -> Vec<i64> {
    let mut res = Vec::with_capacity(seeds.len());
    for &i in seeds {
        let mut next = i;
        for j in range_maps {
            if j.start_in_src(i) {
                next = j.to_dest(i);
                break;
            }
        }
        res.push(next);
    }
    res
}

fn run_range2(seeds: &[Range], range_maps: &[Range2]) -> Vec<Range> {
    let mut res = Vec::with_capacity(seeds.len());
    let mut other_ranges = Vec::with_capacity(seeds.len());
    for i in seeds {
        let mut next_start = i.start;
        let mut next_end = i.end;
        for j in range_maps {
            if j.start_in_src(i.start) {
                next_start = j.to_dest(i.start);
                next_end = if j.end_in_src(i.end) {
                    j.to_dest(i.end)
                } else {
                    other_ranges.push(Range {
                        start: j.src.end,
                        end: i.end,
                    });
                    j.dest.end
                };
                break;
            }
            if j.end_in_src(i.end) {
                next_start = j.dest.start;
                next_end = j.to_dest(i.end);
                other_ranges.push(Range {
                    start: i.start,
                    end: j.src.start,
                });
                break;
            }
            if i.end > j.src.end && i.start < j.src.start {
                next_start = j.dest.start;
                next_end = j.dest.end;
                other_ranges.push(Range {
                    start: i.start,
                    end: j.src.start,
                });
                other_ranges.push(Range {
                    start: j.src.end,
                    end: i.end,
                });
                break;
            }
        }
        res.push(Range {
            start: next_start,
            end: next_end,
        });
    }
    if other_ranges.len() != 0 {
        res.append(&mut run_range2(&other_ranges, range_maps));
    }
    res
}
