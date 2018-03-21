package collectors

type Collector interface {
	Collect()
	Key() string
}

