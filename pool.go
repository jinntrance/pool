package pool

import (
	"log"
	"time"
)

const (
	defaultPoolSize       = 10
	defaultMinPoolSize    = 5
	defaultRecycleTimeout = 30 // 30 minutes
)

type Pool struct {
	shared chan interface{}
	//store the clients of the MinPoolSize
	minShared chan interface{}
	//return a new object and error
	New func() (interface{}, error)
	//close the element, usually a client
	Close       func(interface{})
	PoolSize    int
	MinPoolSize int
	//whether to wait for others to return clients when the PoolSize is exceeded
	defaultToWait  bool
	RecycleTimeout int
}

// Put element 'x' back to the pool.
func (p *Pool) Put(x interface{}, lastError error) {
	if nil == x {
		return
	} else if nil != lastError {
		//error occurs in x, so close 'x' from the pool and never push it back
		p.Close(x)
		return
	} else {
		select {
		//push back the element to the pool 'minShared' and then 'shared'
		case p.minShared <- x:
		case p.shared <- x:
		//or else remove items when item size exceeds PoolSize
		default:
			p.Close(x)
		}
	}

}

func (p *Pool) Init() *Pool {
	if 0 >= p.MinPoolSize {
		p.MinPoolSize = defaultMinPoolSize
	}
	if 0 >= p.PoolSize {
		p.PoolSize = defaultPoolSize
	}
	if p.MinPoolSize > p.PoolSize {
		p.PoolSize = p.MinPoolSize
	}
	if nil == p.shared {
		p.shared = make(chan interface{}, p.PoolSize-p.MinPoolSize)
	}
	if nil == p.minShared {
		p.minShared = make(chan interface{}, p.MinPoolSize)
	}
	if 0 >= p.RecycleTimeout {
		p.RecycleTimeout = defaultRecycleTimeout
	}
	go func() {
		for {
			//close those surplus clients
			select {
			case <-time.After(time.Duration(p.RecycleTimeout) * time.Minute):
				for {
					select {
					case c := <-p.shared:
						select {
						case p.minShared <- c:
						default:
							p.Close(c)
						}
					default:
						goto END
					}
				}
			END:
			}
		}
	}()
	return p
}

func (p *Pool) Get() (interface{}, error) {
	if nil == p.shared {
		p.Init()
	}
	if nil == p.New {
		log.Println("No 'New' method is specified for the pool")
		return nil, nil
	}
	select {
	//pop up one element from the pool
	case x := <-p.minShared:
		return x, nil
	case x := <-p.shared:
		return x, nil
	default: //or else create a new one
		return p.New()
	}
}

func (p *Pool) Cleanup() {
	if nil != p.shared {
		close(p.shared)
		close(p.minShared)
		for ele := range p.shared {
			p.Close(ele)
		}
		for ele := range p.minShared {
			p.Close(ele)
		}
	}
}
