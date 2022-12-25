package linarian

import (
	"bufio"
)

type Linarian struct {
	r       *bufio.Reader
	lines   []string
	current int
	latest  int
}

func New(r *bufio.Reader, capacity int) *Linarian {
	return &Linarian{
		r:       r,
		lines:   make([]string, capacity),
		current: 0,
		latest:  0,
	}
}

func (l *Linarian) ReadLine() (string, error) {
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

func (l *Linarian) JumpBack(steps int) {
	steps %= len(l.lines)

	l.current = (l.latest - steps)
	if l.current < 0 {
		l.current += len(l.lines)
	}
}
