package tool

import "sync"

type Pool struct {
	queue chan int
	WG    *sync.WaitGroup
}

// 创建并发控制池, 设置并发数量与任务总数量
func NewPool(cap, total int) *Pool {
	if cap < 1 {
		cap = 1
	}
	p := &Pool{
		queue: make(chan int, cap),
		WG:    new(sync.WaitGroup),
	}
	p.WG.Add(total)
	return p
}

// 向并发队列中添加一个
func (p *Pool) AddOne() {
	p.queue <- 1
}

// 并发队列中释放一个, 并从任务总数量中减去一个
func (p *Pool) DelOne() {
	<-p.queue
	p.WG.Done()
}
