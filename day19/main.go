package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	sum := 0
	workflows := map[string]Workflow{}

	addWorkflows := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if addWorkflows {
			if line == "" {
				addWorkflows = false
				continue
			}

			name, rest, ok := strings.Cut(line, "{")
			if !ok {
				log.Fatalln("Invalid workflow")
			}
			rest = strings.Trim(rest, "}")
			ruleStrs := strings.Split(rest, ",")
			rules := make([]Rule, 0, len(ruleStrs))
			for _, i := range ruleStrs {
				if lhs, rhs, ok := strings.Cut(i, ":"); ok {
					if opIdx := strings.IndexAny(lhs, "<>"); opIdx >= 0 {
						if opIdx == 0 {
							log.Fatalln("No part name in rule")
						}
						imm, err := strconv.Atoi(lhs[opIdx+1:])
						if err != nil {
							log.Fatalln(err)
						}
						rules = append(rules, Rule{
							Part:   lhs[:opIdx],
							Op:     string(lhs[opIdx]),
							Imm:    imm,
							Target: rhs,
						})
						continue
					}
					log.Fatalln("Invalid rule condition")
					continue
				}
				rules = append(rules, Rule{
					Target: i,
				})
			}
			workflows[name] = Workflow{
				Name:  name,
				Rules: rules,
			}
			continue
		}

		stateStr := strings.Trim(line, "{}")
		stateMap := make(map[string]int, len(stateStr))
		rating := 0
		for _, i := range strings.Split(stateStr, ",") {
			lhs, rhs, ok := strings.Cut(i, "=")
			if !ok {
				log.Fatalln("Invalid state part assign")
			}
			num, err := strconv.Atoi(rhs)
			if err != nil {
				log.Fatalln(err)
			}
			stateMap[lhs] = num
			rating += num
		}
		if runWorkflows(workflows, stateMap) {
			sum += rating
		} else {
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum)

	stateMapRanges := map[string]Range{
		"x": {
			Left:  1,
			Right: 4001,
		},
		"m": {
			Left:  1,
			Right: 4001,
		},
		"a": {
			Left:  1,
			Right: 4001,
		},
		"s": {
			Left:  1,
			Right: 4001,
		},
	}
	fmt.Println("Part 2:", runWorkflowRanges(workflows, "in", stateMapRanges))
}

type (
	Workflow struct {
		Name  string
		Rules []Rule
	}

	Rule struct {
		Part   string
		Op     string
		Imm    int
		Target string
	}

	Range struct {
		Left  int
		Right int
	}
)

func runWorkflowRanges(workflows map[string]Workflow, current string, stateMap map[string]Range) int {
	switch current {
	case "A":
		{
			prod := 1
			for _, v := range stateMap {
				prod *= v.Right - v.Left
			}
			return prod
		}
	case "R":
		return 0
	}
	wf, ok := workflows[current]
	if !ok {
		log.Fatalln("Invalid workflow name")
	}
	sum := 0
	for _, rule := range wf.Rules {
		if rule.Op != "" {
			v, ok := stateMap[rule.Part]
			if !ok {
				log.Fatalln("Invalid rule part name")
			}
			switch rule.Op {
			case "<":
				if v.Right <= rule.Imm {
					sum += runWorkflowRanges(workflows, rule.Target, stateMap)
					return sum
				} else if v.Left >= rule.Imm {
				} else {
					childStateMap := maps.Clone(stateMap)
					childStateMap[rule.Part] = Range{
						Left:  v.Left,
						Right: rule.Imm,
					}
					sum += runWorkflowRanges(workflows, rule.Target, childStateMap)
					if v.Right == rule.Imm {
						return sum
					}
					childStateMap[rule.Part] = Range{
						Left:  rule.Imm,
						Right: v.Right,
					}
					stateMap = childStateMap
				}
			case ">":
				if v.Left > rule.Imm {
					sum += runWorkflowRanges(workflows, rule.Target, stateMap)
					return sum
				} else if v.Right <= rule.Imm+1 {
				} else {
					childStateMap := maps.Clone(stateMap)
					childStateMap[rule.Part] = Range{
						Left:  rule.Imm + 1,
						Right: v.Right,
					}
					sum += runWorkflowRanges(workflows, rule.Target, childStateMap)
					if v.Left == rule.Imm+1 {
						return sum
					}
					childStateMap[rule.Part] = Range{
						Left:  v.Left,
						Right: rule.Imm + 1,
					}
					stateMap = childStateMap
				}
			default:
				log.Fatalln("Invalid rule op")
			}
		} else {
			sum += runWorkflowRanges(workflows, rule.Target, stateMap)
			return sum
		}
	}
	return sum
}

func runWorkflows(workflows map[string]Workflow, stateMap map[string]int) bool {
	current := "in"
outer:
	for {
		wf, ok := workflows[current]
		if !ok {
			log.Fatalln("Invalid workflow name")
		}
		for _, rule := range wf.Rules {
			if rule.Op != "" {
				v, ok := stateMap[rule.Part]
				if !ok {
					log.Fatalln("Invalid rule part name")
				}
				switch rule.Op {
				case "<":
					if v < rule.Imm {
					} else {
						continue
					}
				case ">":
					if v > rule.Imm {
					} else {
						continue
					}
				default:
					log.Fatalln("Invalid rule op")
				}
			}
			switch rule.Target {
			case "A":
				return true
			case "R":
				return false
			default:
				current = rule.Target
				continue outer
			}
		}
	}
}
