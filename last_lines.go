package main

import (
	"bufio"
)

type LastLines struct {
	r       *bufio.Reader
	lines   []string
	current int
	latest  int
}

func NewLastLines(r *bufio.Reader, capacity int) *LastLines {
	return &LastLines{
		r:       r,
		lines:   make([]string, capacity),
		current: 0,
		latest:  0,
	}
}

func (l *LastLines) ReadLine() (string, error) {
	if l.current != l.latest {
		l.current = (1 + l.current) % len(l.lines)
		return l.lines[l.current], nil
	}

	line, err := l.r.ReadString('\n')
	if err != nil {
		return "", err
	}

	l.current = (1 + l.current) % len(l.lines)
	l.latest = l.current
	l.lines[l.latest] = line
	return line, err
}

func (l LastLines) GetFromLatest(backSteps int) (string, error) {
	backSteps %= len(l.lines)

	return l.lines[(backSteps+l.latest+len(l.lines))%len(l.lines)], nil
}

func (l LastLines) Get(backSteps int) (string, error) {
	distance := l.latest - l.current
	if distance < 0 {
		distance += len(l.lines)
	}

	return l.GetFromLatest(distance + backSteps)
}

func (l *LastLines) JumpBack(steps int) {
	steps %= len(l.lines)

	l.current = (l.latest - steps)
	if l.current < 0 {
		l.current += len(l.lines)
	}
}
