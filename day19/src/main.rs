use std::collections::HashMap;
use std::fs::File;
use std::io::{prelude::*, BufReader};

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut sum = 0;
    let mut workflow_names = HashMap::new();
    let mut workflows = Vec::new();
    let mut add_workflows = true;
    let mut in_wf_id = 0;

    for line in reader.lines() {
        let line = line?;
        if add_workflows {
            if line == "" {
                add_workflows = false;
                continue;
            }

            let (name, rest) = line.split_once('{').ok_or("Invalid line")?;
            let rest = rest.trim_end_matches('}');
            let mut rules = Vec::new();
            for i in rest.split(',') {
                if let Some((lhs, rhs)) = i.split_once(':') {
                    let op_chars = &['<', '>'];
                    let op_idx = lhs.find(op_chars).ok_or("Invalid rule condition")?;
                    if op_idx == 0 {
                        return Err("Invalid part name".into());
                    }
                    let imm = lhs[op_idx + 1..].parse::<usize>()?;
                    let part = match lhs[..op_idx].as_bytes()[0] {
                        b'x' => 0,
                        b'm' => 1,
                        b'a' => 2,
                        b's' => 3,
                        _ => {
                            return Err("Invalid part name".into());
                        }
                    };
                    let num_wf = workflow_names.len();
                    let (accept, reject, target) = match rhs {
                        "A" => (true, false, 0),
                        "R" => (false, true, 0),
                        v => (
                            false,
                            false,
                            *workflow_names.entry(v.to_string()).or_insert(num_wf),
                        ),
                    };
                    rules.push(Rule {
                        part,
                        op: lhs.as_bytes()[op_idx],
                        imm,
                        accept,
                        reject,
                        target,
                    });
                    continue;
                }
                let num_wf = workflow_names.len();
                let (accept, reject, target) = match i {
                    "A" => (true, false, 0),
                    "R" => (false, true, 0),
                    v => (
                        false,
                        false,
                        *workflow_names.entry(v.to_string()).or_insert(num_wf),
                    ),
                };
                rules.push(Rule {
                    part: 0,
                    op: 0,
                    imm: 0,
                    accept,
                    reject,
                    target,
                })
            }
            let num_wf = workflow_names.len();
            let id = *workflow_names.entry(name.to_string()).or_insert(num_wf);
            if name == "in" {
                in_wf_id = id;
            }
            if id >= workflows.len() {
                workflows.resize(id + 1, Vec::new());
            }
            workflows[id] = rules;
            continue;
        }

        let line_chars: &[_] = &['{', '}'];
        let line = line.trim_matches(line_chars);
        let mut state = [0; 4];
        let mut rating = 0;
        for i in line.split(',') {
            let (lhs, rhs) = i.split_once('=').ok_or("Invalid state part assignment")?;
            let num = rhs.parse::<usize>()?;
            let part = match lhs.as_bytes()[0] {
                b'x' => 0,
                b'm' => 1,
                b'a' => 2,
                b's' => 3,
                _ => return Err("Invalid part name".into()),
            };
            state[part] = num;
            rating += num;
        }
        if run_workflow(&workflows, in_wf_id, state)? {
            sum += rating;
        }
    }

    println!(
        "Part 1: {}\nPart 2: {}",
        sum,
        run_workflow_ranges(
            &workflows,
            false,
            false,
            in_wf_id,
            [Range {
                left: 1,
                right: 4001
            }; 4]
        )?,
    );

    Ok(())
}

#[derive(Clone, Copy)]
struct Range {
    left: usize,
    right: usize,
}

fn run_workflow_ranges(
    workflows: &[Vec<Rule>],
    accept: bool,
    reject: bool,
    current: usize,
    mut state: [Range; 4],
) -> Result<usize, Box<dyn std::error::Error>> {
    match (accept, reject) {
        (true, false) => {
            return Ok(state
                .iter()
                .fold(1, |acc, v| acc * (v.right - v.left) as usize))
        }
        (false, true) => return Ok(0),
        _ => {}
    }
    let wf = &workflows[current];
    let mut sum = 0;
    for rule in wf {
        if rule.op != 0 {
            let v = &state[rule.part as usize];
            match rule.op {
                b'<' => {
                    if v.right <= rule.imm {
                        return Ok(sum
                            + run_workflow_ranges(
                                workflows,
                                rule.accept,
                                rule.reject,
                                rule.target,
                                state,
                            )?);
                    } else if v.left >= rule.imm {
                    } else {
                        let mut child_state = state.clone();
                        child_state[rule.part as usize] = Range {
                            left: v.left,
                            right: rule.imm,
                        };
                        sum += run_workflow_ranges(
                            workflows,
                            rule.accept,
                            rule.reject,
                            rule.target,
                            child_state,
                        )?;
                        if v.right == rule.imm {
                            return Ok(sum);
                        }
                        state[rule.part as usize] = Range {
                            left: rule.imm,
                            right: v.right,
                        };
                    }
                }
                b'>' => {
                    if v.left > rule.imm {
                        return Ok(sum
                            + run_workflow_ranges(
                                workflows,
                                rule.accept,
                                rule.reject,
                                rule.target,
                                state,
                            )?);
                    } else if v.right <= rule.imm + 1 {
                    } else {
                        let mut child_state = state.clone();
                        child_state[rule.part as usize] = Range {
                            left: rule.imm + 1,
                            right: v.right,
                        };
                        sum += run_workflow_ranges(
                            workflows,
                            rule.accept,
                            rule.reject,
                            rule.target,
                            child_state,
                        )?;
                        if v.left == rule.imm + 1 {
                            return Ok(sum);
                        }
                        state[rule.part as usize] = Range {
                            left: v.left,
                            right: rule.imm + 1,
                        };
                    }
                }
                _ => return Err("Invalid rule op".into()),
            }
        } else {
            return Ok(
                sum + run_workflow_ranges(workflows, rule.accept, rule.reject, rule.target, state)?
            );
        }
    }
    Err("Workflow has no default rule".into())
}

#[derive(Clone)]
struct Rule {
    part: u8,
    op: u8,
    imm: usize,
    accept: bool,
    reject: bool,
    target: usize,
}

fn run_workflow(
    workflows: &[Vec<Rule>],
    current: usize,
    state: [usize; 4],
) -> Result<bool, Box<dyn std::error::Error>> {
    let wf = &workflows[current];
    for rule in wf {
        let rule = rule;
        if rule.op != 0 {
            let v = state[rule.part as usize];
            match rule.op {
                b'<' => {
                    if v < rule.imm {
                    } else {
                        continue;
                    }
                }
                b'>' => {
                    if v > rule.imm {
                    } else {
                        continue;
                    }
                }
                _ => return Err("Invalid rule op".into()),
            }
        }
        return match (rule.accept, rule.reject) {
            (true, false) => Ok(true),
            (false, true) => Ok(false),
            _ => run_workflow(workflows, rule.target, state),
        };
    }
    Err("Workflow does not have default rule".into())
}
