package main

type fetcher struct {
	ContentClients map[Provider]Client
	Config         ContentMix
}

func (a fetcher) fetchItems(userIP string, count, limit int) []*ContentItem {
	resp := make([]*ContentItem, 0, count)
	for i := 0; i < count; i++ {
		items, err := a.fetchItem(i, userIP, limit)
		if err != nil {
			break
		}
		resp = append(resp, items...)
	}
	return resp
}

func (a fetcher) fetchItem(n int, ip string, limit int) ([]*ContentItem, error) {
	p := a.selectProviderFor(n)

	items, err := a.ContentClients[p.Type].GetContent(ip, limit)
	if err == nil {
		return items, nil
	}

	if p.Fallback == nil {
		return nil, err
	}

	return a.ContentClients[*p.Fallback].GetContent(ip, limit)
}

func (a fetcher) selectProviderFor(n int) ContentConfig {
	p := a.Config[n%len(a.Config)]
	return p
}
