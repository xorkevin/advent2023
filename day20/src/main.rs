use std::collections::{HashMap, VecDeque};
use std::fs::File;
use std::io::{prelude::*, BufReader};

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut mod_names = HashMap::new();
    let mut com_mods = Vec::new();

    let mut start_mod_id = 0;
    let mut end_mod_id = 0;

    for line in reader.lines() {
        let line = line?;
        let (lhs, rhs) = line.split_once(" -> ").ok_or("Invalid line")?;
        let (kind, lhs) = match lhs.as_bytes()[0] {
            v @ (b'%' | b'&') => (v, &lhs[1..]),
            _ => (0, lhs),
        };
        let mut dest = Vec::new();
        for i in rhs.split(", ") {
            let num_mods = mod_names.len();
            dest.push(*mod_names.entry(i.to_string()).or_insert(num_mods));
        }
        let num_mods = mod_names.len();
        let id = *mod_names.entry(lhs.to_string()).or_insert(num_mods);
        match lhs {
            "broadcaster" => start_mod_id = id,
            "dg" => end_mod_id = id,
            _ => {}
        }
        let num_mods = mod_names.len();
        if num_mods >= com_mods.len() {
            com_mods.resize(
                num_mods,
                ComMod {
                    kind: 0,
                    state: false,
                    mem: Vec::new(),
                    inp: Vec::new(),
                    dest: Vec::new(),
                },
            );
        }
        com_mods[id] = ComMod {
            kind,
            state: false,
            mem: Vec::new(),
            inp: Vec::new(),
            dest,
        };
    }

    let num_mods = com_mods.len();
    let mut inps = vec![Vec::new(); num_mods];
    for (id, i) in com_mods.iter_mut().enumerate() {
        i.mem.resize(num_mods, false);
        for &j in &i.dest {
            inps[j].push(id);
        }
    }
    for (id, i) in inps.into_iter().enumerate() {
        com_mods[id].inp = i;
    }

    let mut sum_hi = 0;
    let mut sum_lo = 0;
    let mut idx = 0;
    let mut inbox = VecDeque::new();
    let mut target_packets = Vec::new();
    let mut target_cycles = vec![
        Cycle {
            prev: 0,
            size: 0,
            rem: 0
        };
        num_mods
    ];
    let mut total_revisits = 0;
    while idx < 1000 {
        inbox.push_back(Packet {
            from: start_mod_id,
            to: start_mod_id,
            sig: false,
        });
        let (hi, lo) = run_til_end(&mut com_mods, &mut inbox, end_mod_id, &mut target_packets);
        sum_hi += hi;
        sum_lo += lo;
        idx += 1;
        while let Some(packet) = target_packets.pop() {
            let c = &mut target_cycles[packet];
            if c.prev == 0 {
                c.prev = idx;
            } else if c.size == 0 {
                let size = idx - c.prev;
                c.prev = idx;
                c.size = size;
                c.rem = idx % size;
                total_revisits += 1;
            } else {
                let size = idx - c.prev;
                if size != c.size {
                    return Err("Multiple cycle lengths".into());
                }
                c.prev = idx
            }
        }
    }
    while total_revisits < 4 {
        inbox.push_back(Packet {
            from: start_mod_id,
            to: start_mod_id,
            sig: false,
        });
        run_til_end(&mut com_mods, &mut inbox, end_mod_id, &mut target_packets);
        idx += 1;
        while let Some(packet) = target_packets.pop() {
            let c = &mut target_cycles[packet];
            if c.prev == 0 {
                c.prev = idx;
            } else if c.size == 0 {
                let size = idx - c.prev;
                c.prev = idx;
                c.size = size;
                c.rem = idx % size;
                total_revisits += 1;
            } else {
                let size = idx - c.prev;
                if size != c.size {
                    return Err("Multiple cycle lengths".into());
                }
                c.prev = idx
            }
        }
    }

    let mut a = 0;
    let mut m = 0;
    for i in target_cycles {
        if i.size == 0 {
            continue;
        }
        if m == 0 {
            a = i.rem;
            m = i.size;
            continue;
        }
        (a, m) = crt(a, m, i.rem, i.size).ok_or("Unsolvable constraints")?;
    }
    if a == 0 {
        a += m;
    }

    println!("Part 1: {}\nPart 2: {}", sum_hi * sum_lo, a);

    Ok(())
}

#[derive(Clone)]
struct Cycle {
    prev: usize,
    size: usize,
    rem: usize,
}

#[derive(Clone, Debug)]
struct ComMod {
    kind: u8,
    state: bool,
    mem: Vec<bool>,
    inp: Vec<usize>,
    dest: Vec<usize>,
}

#[derive(Clone)]
struct Packet {
    from: usize,
    to: usize,
    sig: bool,
}

fn run_til_end(
    com_mods: &mut Vec<ComMod>,
    inbox: &mut VecDeque<Packet>,
    target_mod_id: usize,
    target_packets: &mut Vec<usize>,
) -> (usize, usize) {
    let mut hi = 0;
    let mut lo = 0;
    while let Some(packet) = inbox.pop_front() {
        if packet.to == target_mod_id && packet.sig {
            target_packets.push(packet.from);
        }
        if packet.sig {
            hi += 1;
        } else {
            lo += 1;
        }
        pulse(&mut com_mods[packet.to], packet, inbox);
    }
    return (hi, lo);
}

fn pulse(com_mod: &mut ComMod, packet: Packet, inbox: &mut VecDeque<Packet>) {
    match com_mod.kind {
        b'%' => {
            if packet.sig {
                return;
            }
            com_mod.state = !com_mod.state;
            let sig = com_mod.state;
            for &i in &com_mod.dest {
                inbox.push_back(Packet {
                    from: packet.to,
                    to: i,
                    sig,
                });
            }
        }
        b'&' => {
            com_mod.mem[packet.from] = packet.sig;
            let sig = !packet.sig || com_mod.inp.iter().any(|&v| !com_mod.mem[v]);
            for &i in &com_mod.dest {
                inbox.push_back(Packet {
                    from: packet.to,
                    to: i,
                    sig,
                });
            }
        }
        _ => {
            for &i in &com_mod.dest {
                inbox.push_back(Packet {
                    from: packet.to,
                    to: i,
                    sig: packet.sig,
                });
            }
        }
    }
}

fn crt(a1: usize, m1: usize, a2: usize, m2: usize) -> Option<(usize, usize)> {
    let (g, p, q) = ext_gcd(m1, m2);
    if a1 % g != a2 % g {
        return None;
    }
    let m1g = m1 / g;
    let m2g = m2 / g;
    let lcm = m1g * m2;
    let p = num_mod(p, lcm);
    let q = num_mod(q, lcm);
    let x = mul_mod(mul_mod(a1, m2g, lcm), q, lcm) + mul_mod(mul_mod(a2, m2g, lcm), p, lcm) % lcm;
    Some((x, lcm))
}

fn ext_gcd(mut a: usize, mut b: usize) -> (usize, isize, isize) {
    let mut x2 = 1;
    let mut x1 = 0;
    let mut y2 = 0;
    let mut y1 = 1;
    let mut flip = false;
    // let a be larger than b
    if a < b {
        (a, b) = (b, a);
        flip = true;
    }
    while b > 0 {
        let q = a / b;
        (a, b) = (b, a % b);
        (x2, x1) = (x1, x2 - (q as isize) * x1);
        (y2, y1) = (y1, y2 - (q as isize) * y1);
    }
    if flip {
        (x2, y2) = (y2, x2);
    }
    (a, x2, y2)
}

fn mul_mod(a: usize, b: usize, m: usize) -> usize {
    ((a as u128 * b as u128) % m as u128) as usize
}

fn num_mod(a: isize, m: usize) -> usize {
    a.rem_euclid(m as isize) as usize
}
