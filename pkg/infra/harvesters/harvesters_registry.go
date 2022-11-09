package harvesters

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lithammer/shortuuid/v3"
)

var (
	ErrHarvesterNotFound = errors.New("harvester does not exist")
	Harvesters           = make(map[string]HarvesterData)
	harvestersMutex      = sync.RWMutex{}

	ErrGmailNotFound = errors.New("gmail does not exist")
	Gmails           = make(map[string]GmailData)
	GmailsMutex      = sync.RWMutex{}
)

func AddHarvester() HarvesterData {
	harvestersMutex.Lock()
	defer harvestersMutex.Unlock()

	index := len(Harvesters)

	h := HarvesterData{
		ID:      shortuuid.New(),
		Index:   index,
		Focused: false,
		Name:    fmt.Sprintf("Solver %d", index),
	}

	Harvesters[h.ID] = h

	Create(h.ID)

	return h
}

func EditHarvester(newHarvester HarvesterData) (HarvesterData, error) {
	harvestersMutex.Lock()
	defer harvestersMutex.Unlock()

	_, ok := Harvesters[newHarvester.ID]
	if !ok {
		return HarvesterData{}, ErrHarvesterNotFound
	}

	if Harvesters[newHarvester.ID].Proxy != newHarvester.Proxy {
		SetProxy(newHarvester.ID, newHarvester.Proxy)
	}

	Harvesters[newHarvester.ID] = newHarvester

	return newHarvester, nil
}

func focus(ID string) {
	for _, h := range Harvesters {
		if h.Focused {
			o := Harvesters[h.ID]
			o.Focused = false
			Harvesters[h.ID] = o
			Hide(h.ID)
		}
	}

	n := Harvesters[ID]
	n.Focused = true
	Harvesters[ID] = n
	Show(ID)
}

func FocusHarvester(ID string) (id string, Error error) {
	harvestersMutex.Lock()
	defer harvestersMutex.Unlock()

	_, ok := Harvesters[ID]
	if !ok {
		return "", ErrHarvesterNotFound
	}

	focus(ID)

	return ID, nil
}

func ShowHarvester() {
	for _, h := range Harvesters {
		if h.Focused {
			Show(h.ID)
		}
	}

}

func HideHarvester() {
	for _, h := range Harvesters {
		if h.Focused {
			Hide(h.ID)
		}
	}
}

func DeleteHarvester(id string) (oid string, Error error) {
	if !doesHarvesterExist(id) {
		return "", ErrHarvesterNotFound
	}

	harvestersMutex.RLock()
	defer harvestersMutex.RUnlock()

	delete(Harvesters, id)

	Remove(id)

	focus("default")

	return id, nil
}

func GetHarvesterIds() []string {
	harvestersMutex.RLock()
	defer harvestersMutex.RUnlock()

	var harvesterIds []string
	for id := range Harvesters {
		harvesterIds = append(harvesterIds, id)
	}
	return harvesterIds
}

func GetHarvester(id string) (HarvesterData, error) {
	harvestersMutex.RLock()
	defer harvestersMutex.RUnlock()

	if har, ok := Harvesters[id]; ok {
		return har, nil
	} else {
		return HarvesterData{}, ErrHarvesterNotFound
	}
}

func doesHarvesterExist(id string) bool {
	harvestersMutex.RLock()
	defer harvestersMutex.RUnlock()
	_, ok := Harvesters[id]
	return ok
}

func AddDefaultHarvester() {
	harvestersMutex.Lock()
	defer harvestersMutex.Unlock()

	Harvesters["default"] = HarvesterData{
		ID:      "default",
		Index:   0,
		Focused: true,
		Name:    "default",
	}
}

func AddGmail() GmailData {
	GmailsMutex.Lock()
	defer GmailsMutex.Unlock()

	id := shortuuid.New()

	g := GmailData{
		ID: id,
	}

	Gmails[g.ID] = g

	GmailsMutex.Unlock()
	gm, _ := SelectGmail(g.ID)
	GmailsMutex.Lock()

	return gm
}

func SelectGmail(id string) (gmail GmailData, err error) {
	_, ok := Gmails[id]
	if !ok {
		return
	}

	if Gmails[id].Cookies == nil {
		c, err := RequestSolve(Solver{
			ID:   id,
			Type: Gmail_Login,
			Store: &Store{
				Url: "https://accounts.google.com/signin/v2/identifier?service=youtube&uilel=3&passive=true&continue=https%3A%2F%2Fwww.youtube.com%2Fsignin%3Faction_handle_signin%3Dtrue%26app%3Ddesktop%26hl%3Den%26next%3D%252F&hl=en&ec=65620&flowName=GlifWebSignIn&flowEntry=ServiceLogin",
			},
		})

		if err != nil {
			return gmail, err
		}

		s, err := GetSolve(c)

		if err != nil {
			return gmail, err
		}

		h := Gmails[id]
		h.Email = s.Data.Email
		h.Cookies = &s.Data.Cookies
		EditGmail(h)
	}

	gmail = Gmails[id]

	for _, h := range Harvesters {
		if h.Focused {
			g := Gmails[id]
			h.Account = &g
			EditHarvester(h)
		}
	}

	return
}

func EditGmail(newGmail GmailData) (GmailData, error) {
	GmailsMutex.Lock()
	defer GmailsMutex.Unlock()

	_, ok := Gmails[newGmail.ID]
	if !ok {
		return GmailData{}, ErrGmailNotFound
	}

	Gmails[newGmail.ID] = newGmail

	return newGmail, nil
}

func DeleteGmail(id string) (oid string, Error error) {
	if !DoesGmailExist(id) {
		return "", ErrGmailNotFound
	}

	GmailsMutex.RLock()
	defer GmailsMutex.RUnlock()
	delete(Gmails, id)

	for _, h := range Harvesters {
		if h.Account != nil && h.Account.ID == id {
			h.Account = nil
			EditHarvester(h)
		}

		if h.Solver != nil && h.Solver.ID == id {
			RemoveSolve(id)
		}
	}

	return id, nil
}

func DoesGmailExist(id string) bool {
	GmailsMutex.RLock()
	defer GmailsMutex.RUnlock()
	_, ok := Gmails[id]
	return ok
}

func GetGmailIds() []string {
	GmailsMutex.RLock()
	defer GmailsMutex.RUnlock()

	var gmailIds []string
	for id := range Gmails {
		gmailIds = append(gmailIds, id)
	}
	return gmailIds
}

func GetGmail(id string) (gmail GmailData, Error error) {
	GmailsMutex.RLock()
	defer GmailsMutex.RUnlock()

	if gmail, ok := Gmails[id]; ok {
		return gmail, nil
	} else {
		return GmailData{}, ErrGmailNotFound
	}
}
