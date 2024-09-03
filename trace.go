package telegram

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Trace struct {
	mu       sync.Mutex
	Sessions map[int64]*Step // map[user id] steps
}

type Step struct {
	Endpoint []string
	Data     []string
}

func NewTrace() *Trace {
	t := &Trace{
		mu:       sync.Mutex{},
		Sessions: make(map[int64]*Step),
	}

	go func() {
		for {
			byteTrace, err := json.Marshal(t)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(byteTrace))
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()

	return t
}

func (t *Trace) SaveStep(userId int64, step string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.Sessions[userId]; !ok {
		t.Sessions[userId] = new(Step)
	}

	if step == "" {
		return
	}
	t.Sessions[userId].Endpoint = append(t.Sessions[userId].Endpoint, step)
}

func (t *Trace) SaveData(userId int64, data ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.Sessions[userId]; !ok {
		t.Sessions[userId] = new(Step)
	}

	for i := 0; i < len(data); i++ {
		if data[i] == "" {
			continue
		}
		t.Sessions[userId].Data = append(t.Sessions[userId].Data, data[i])
	}
}

func (t *Trace) ResetSession(userId int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.Sessions[userId]; !ok {
		t.Sessions[userId] = new(Step)
	}

	t.Sessions[userId].Data = make([]string, 0)
}

func (t *Trace) GetSteps(userId int64) []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if steps, ok := t.Sessions[userId]; ok {
		return steps.Endpoint
	}
	return nil
}

func (t *Trace) ResetSteps(userId int64, steps []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.Sessions[userId]; !ok {
		t.Sessions[userId] = new(Step)
	}
	t.Sessions[userId].Endpoint = steps
}
