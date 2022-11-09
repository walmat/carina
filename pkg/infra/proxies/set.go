package proxies

import "sync"

type proxySet struct {
	Index map[string]int `json:"index"`
	Data  []ProxyData    `json:"data"`

	mu sync.RWMutex
}

func newProxySet() *proxySet {
	return &proxySet{
		Index: make(map[string]int),
		Data:  []ProxyData{},
	}
}

func (s *proxySet) Add(data ProxyData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data = append(s.Data, data)
	s.Index[data.Id] = len(s.Data) - 1
	return nil
}

func (s *proxySet) Set(data ProxyData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := s.Index[data.Id]
	s.Data[idx] = data
	return nil
}

func (s *proxySet) Del(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	oIdx := s.Index[id]
	s.Data = append(s.Data[:oIdx], s.Data[oIdx+1:]...)
	delete(s.Index, id)
	for pId, idx := range s.Index {
		if idx < oIdx {
			continue
		}
		s.Index[pId] = idx - 1
	}
	return nil
}

func (s *proxySet) Get(id string) (*ProxyData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	idx := s.Index[id]
	return &s.Data[idx], nil
}

func (s *proxySet) Values() []ProxyData {
	return s.Data
}
