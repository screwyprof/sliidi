package main

const numOfWorkers = 3
const defaultRecordsPerRequest = 1

type result struct {
	items []*ContentItem
	err   error
}

type fetcher struct {
	ContentClients map[Provider]Client
	Config         ContentMix
}

func (f fetcher) fetchItems(userIP string, count, offset int) []*ContentItem {
	results := make(chan result, count)
	jobs := f.genJobs(count)

	for i := 0; i < numOfWorkers; i++ {
		go f.worker(jobs, results, userIP, offset)
	}

	storage := make(map[Provider][]result, count)

	items := make([]*ContentItem, 0, count)
	for i := 0; i < count; i++ {
		res := <-results

		if res.err == nil {
			// get provider name
			var source Provider
			for _, item := range res.items {
				source = Provider(item.Source)
				break
			}
			storage[source] = append(storage[source], res)
		}
	}

	for i := 0; i < count; i++ {
		cfg := f.selectCfgProviderFor(i + offset)

		if resItems, ok := f.popResultItemsFor(storage, cfg.Type); ok {
			items = append(items, resItems...)
			continue
		}

		if cfg.Fallback == nil {
			break // big failure
		}

		if resItems, ok := f.popResultItemsFor(storage, *cfg.Fallback); ok {
			items = append(items, resItems...)
			continue
		} else {
			break // big failure
		}
	}

	return items
}

func (f fetcher) genJobs(n int) <-chan int {
	jobs := make(chan int, n)
	defer close(jobs)
	for j := 0; j < n; j++ {
		jobs <- j
	}
	return jobs
}

func (f fetcher) worker(jobs <-chan int, results chan<- result, userIP string, offset int) {
	for j := range jobs {
		items, err := f.fetchItem(userIP, j+offset)
		results <- result{items: items, err: err}
	}
}

func (f fetcher) fetchItem(ip string, n int) ([]*ContentItem, error) {
	p := f.selectCfgProviderFor(n)

	items, err := f.ContentClients[p.Type].GetContent(ip, defaultRecordsPerRequest)
	if err == nil {
		return items, nil
	}

	if p.Fallback == nil {
		return nil, err
	}

	return f.ContentClients[*p.Fallback].GetContent(ip, defaultRecordsPerRequest)
}

func (f fetcher) selectCfgProviderFor(n int) ContentConfig {
	p := f.Config[n%len(f.Config)]
	return p
}

func (f fetcher) popResultItemsFor(storage map[Provider][]result, key Provider) ([]*ContentItem, bool) {
	if p, ok := storage[key]; ok {
		if len(p) > 0 {
			var r result
			r, storage[key] = p[0], p[1:] // pop element from results queue
			return r.items, true
		}
	}
	return nil, false
}
