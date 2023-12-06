struct Race {
    time: i64,
    dist: i64,
}

fn main() {
    let races = [
        Race {
            time: 47,
            dist: 207,
        },
        Race {
            time: 84,
            dist: 1394,
        },
        Race {
            time: 74,
            dist: 1209,
        },
        Race {
            time: 67,
            dist: 1014,
        },
    ];
    println!(
        "Part 1: {}",
        races.iter().fold(1, |acc, v| acc * simulate(v))
    );
    println!(
        "Part 2: {}",
        simulate(&Race {
            time: 47847467,
            dist: 207139412091014,
        })
    );
}

fn simulate(race: &Race) -> i64 {
    let axis = 0.5 * race.time as f64;
    let max_possible = axis.powi(2).floor() as i64;
    if max_possible <= race.dist {
        // no real valued solutions
        return 0;
    }

    // equation x * (time - x) = dist
    // x * time - x * x = dist
    // 0 = x*x - x * time + dist
    // x = (time + sqrt(time^2 - 4*dist)) / 2
    // x = (time - sqrt(time^2 - 4*dist)) / 2

    let disc = ((race.time as f64).powi(2) - (4.0 * race.dist as f64)).sqrt() * 0.5;
    let start = (axis - disc).ceil() as i64;
    let end = (axis + disc).floor() as i64;
    end - start + 1
}
