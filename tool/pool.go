package tool

import (
	"sync"
)

type Pool struct {
	queue        chan int
	wg           *sync.WaitGroup
	total_number int
	done_number  int
}

// 创建并发控制池, 设置并发数量与任务总数量
func NewPool(cap, total int) *Pool {
	if cap < 1 {
		cap = 1
	}
	p := &Pool{
		queue: make(chan int, cap),
		wg:    new(sync.WaitGroup),
	}
	p.wg.Add(total)
	p.total_number = total
	return p
}

// 向并发队列中添加一个
func (p *Pool) AddOne() {
	p.queue <- 1
}

// 并发队列中释放一个, 并从任务总数量中减去一个
func (p *Pool) DelOne() {
	<-p.queue
	p.wg.Done()
	p.done_number += 1
}

// 获取进度
func (p *Pool) GetProgressRate() float32 {
	if p.total_number == 0 {
		return 0
	}
	req := float32(p.done_number) / float32(p.total_number)
	if req < 0 {
		return 0
	}
	return req
}

// 封装wg的wait函数
func (p *Pool) Wait() {
	p.wg.Wait()
}
