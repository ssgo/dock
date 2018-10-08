package tests

import (
	"fmt"
	"strings"
	"time"
)

type RunInfo struct{
	Name string
	Image string
}
var sequences = map[string]int{}
var runs = map[string]map[string]*RunInfo{}

func TestShell(timeout time.Duration, nodeName string, args ...string) (string, int, error) {
	if runs[nodeName] == nil {
		runs[nodeName] = make(map[string]*RunInfo)
	}

	if args[0] == "run" {
		//if dock.GetStats().Nodes[nodeName] == nil {
		//	return ""
		//}
		sequences[nodeName]++
		id := fmt.Sprintf("%s<%.2d>", nodeName, sequences[nodeName])
		runs[nodeName][id] = &RunInfo{Image:args[len(args)-1], Name:args[2]}
		if runs[nodeName][id].Image == "xxx" {
			runs[nodeName][id].Image = args[len(args)-2]
		}
		return id, 100, nil
	}
	if args[0] == "ps" {
		lines := make([]string, 0)
		for id, run := range runs[nodeName] {
			lines = append(lines, id+","+run.Name+","+run.Image+",Up 1 minutes")
		}
		return strings.Join(lines, "\n"), 100, nil
	}
	if args[0] == "stop" {
		delete(runs[nodeName], args[1])
		return args[1], 100, nil
	}
	if args[0] == "rm" {
		return args[1], 100, nil
	}
	if args[0] == "exec" {
		return args[1], 100, nil
	}
	return "", 100, nil
}
