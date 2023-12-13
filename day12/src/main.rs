use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut sum = 0;
    let mut sum2 = 0;
    for line in reader.lines() {
        let line = line?;
        let (a, b) = line.split_once(' ').ok_or("Invalid line")?;
        let nums = b
            .split(',')
            .map(|v| v.parse::<usize>())
            .collect::<Result<Vec<_>, _>>()?;
        sum += get_num_arrangements(
            a.as_bytes(),
            &nums,
            &mut vec![0; a.len() * nums.len()],
            a.len(),
        );
        let ai = a.bytes().collect::<Vec<_>>();
        let mut aa = ai.clone();
        let mut bb = nums.clone();
        for _ in 0..4 {
            aa.push(b'?');
            aa.append(&mut ai.clone());
            bb.append(&mut nums.clone());
        }
        sum2 += get_num_arrangements(&aa[..], &bb, &mut vec![0; aa.len() * bb.len()], aa.len());
    }

    println!("Part 1: {}\nPart 2: {}", sum, sum2);

    Ok(())
}

fn get_num_arrangements(
    b: &[u8],
    nums: &[usize],
    cache: &mut Vec<usize>,
    cache_row_width: usize,
) -> usize {
    if nums.len() == 0 {
        return if rest_no_group(b) { 1 } else { 0 };
    }
    if b.len() == 0 {
        return 0;
    }

    let key = cache_row_width * (nums.len() - 1) + b.len() - 1;
    if cache[key] > 0 {
        return cache[key] - 1;
    }

    let first = b[0];
    if first == b'.' {
        let count = get_num_arrangements(&b[1..], nums, cache, cache_row_width);
        cache[key] = count + 1;
        return count;
    }

    let first_num = nums[0];
    let (prefix_matches, is_end) = match_prefix(b, first_num);
    if is_end {
        return if nums.len() == 1 {
            cache[key] = 2;
            1
        } else {
            cache[key] = 1;
            0
        };
    }

    let mut count = 0;
    if first == b'?' {
        count += get_num_arrangements(&b[1..], nums, cache, cache_row_width);
    } else {
        if !prefix_matches {
            cache[key] = 1;
            return 0;
        }
    }

    if prefix_matches {
        count += get_num_arrangements(&b[first_num + 1..], &nums[1..], cache, cache_row_width);
    }

    cache[key] = count + 1;

    count
}

fn rest_no_group(b: &[u8]) -> bool {
    b.iter().all(|&v| v != b'#')
}

fn match_prefix(b: &[u8], num: usize) -> (bool, bool) {
    if b.len() < num {
        return (false, false);
    }
    for i in 0..num {
        if b[i] == b'.' {
            return (false, false);
        }
    }
    if b.len() == num {
        (true, true)
    } else if b[num] == b'#' {
        (false, false)
    } else {
        (true, false)
    }
}
