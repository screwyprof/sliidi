package main

const defaultRecordsPerRequest = 1

type fetcher struct {
	ContentClients map[Provider]Client
	Config         ContentMix
}

func (a fetcher) fetchItems(userIP string, count, offset int) []*ContentItem {
	resp := make([]*ContentItem, 0, count)
	for i := 0; i < count; i++ {
		items, err := a.fetchItem(userIP, i+offset)
		if err != nil {
			break
		}
		resp = append(resp, items...)
	}
	return resp
}

func (a fetcher) fetchItem(ip string, n int) ([]*ContentItem, error) {
	p := a.selectProviderFor(n)

	items, err := a.ContentClients[p.Type].GetContent(ip, defaultRecordsPerRequest)
	if err == nil {
		return items, nil
	}

	if p.Fallback == nil {
		return nil, err
	}

	return a.ContentClients[*p.Fallback].GetContent(ip, defaultRecordsPerRequest)
}

func (a fetcher) selectProviderFor(n int) ContentConfig {
	p := a.Config[n%len(a.Config)]
	return p
}
