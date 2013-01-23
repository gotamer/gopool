package gopool

type pool struct {
	max int
	resources chan interface{}
	create func() (interface{})
	destroy func(interface{})
}

/*
 * Creates a new resource pool
 */
func Initialize(max int, create func() (interface{}), destroy func(interface{}) ()) (*pool) {
	p := new(pool)
	p.max = max
	p.resources = make(chan interface{}, max)
	for i := 0; i<max; i++ {
		resource := create()
		p.resources <- resource
	}
	p.create = create
	p.destroy = destroy

	return p
}

/*
 * Obtain a resource from the pool
 */
func (p *pool) Acquire() (interface {}) {
	return <-p.resources
}

/*
 * Returns a resource back in to the pool
 */
func (p *pool) Release(resource interface{}) {
	p.resources <- resource
}

/*
 * Remove a resource from the pool.  This is helpful if the resource
 * has gone bad.  A new resource will be created in it's place.
 */
func (p *pool) Destroy(resource interface{}) {
	p.destroy(resource)
	p.resources <- p.create()
}

/*
 * Remove all resources from the pool and call the destroy method on each of
 * them.
 */
func (p *pool) Drain() {
	for {
		select {
		case r := <-p.resources:
			p.destroy(r)
		default:
			return
		}
	}
	close(p.resources)
}