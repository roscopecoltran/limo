package gpool

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Task func()

type GPool interface {
	AddTask(task Task)
	Excute()
	Stop()
	SetTimeout(timeout time.Duration) error
}

type gpool struct {
	n        int
	taskCh   chan Task
	tasks    []Task
	cancel   context.CancelFunc
	timeout  time.Duration
	executed bool
}

func NewPool(n int) GPool {
	return &gpool{
		n:        n,
		taskCh:   make(chan Task, n),
		tasks:    make([]Task, 0),
		executed: false,
	}
}

//TODO channel满了会阻塞
func (this *gpool) AddTask(task Task) {
	this.taskCh <- task
}

func (this *gpool) SetTimeout(timeout time.Duration) error {
	if this.executed {
		return errors.New("Groutine pool is running")
	}
	this.timeout = timeout
	return nil
}

func (this *gpool) Stop() {
	if !this.executed {
		return
	}
	this.cancel()
}

func (this *gpool) Excute() {
	var ctx context.Context
	ancCtx := context.Background()
	if this.timeout != 0 {
		ancCtx, _ = context.WithTimeout(context.Background(), this.timeout)
	}
	ctx, this.cancel = context.WithCancel(ancCtx)
	this.bg(ctx)

	for _, task := range this.tasks {
		this.taskCh <- task
	}
	this.executed = true
}

func (this *gpool) bg(ctx context.Context) {
	for i := 0; i < this.n; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case task := <-this.taskCh:
					task()
				case <-ctx.Done():
					fmt.Println("Goroutine done:", ctx.Err())
				}
			}
		}(ctx)
	}
}