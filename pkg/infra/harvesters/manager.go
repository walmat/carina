package harvesters

import (
	"errors"
	"nebula/pkg/integrations"
	"nebula/pkg/integrations/autosolve"
	"sync"
)

var (
	queue      []Solver
	queueMutex sync.Mutex

	solveChan = make(chan SolverResponse)

	taskChans      = make(map[string]chan SolverResponse)
	taskChansMutex sync.RWMutex

	ErrSolverNotFound = errors.New("solver not found")
)

func init() {
	go handleSolves()
}

func handleSolves() {
	for {
		c := <-solveChan

		taskChansMutex.Lock()
		if ch, ok := taskChans[c.ID]; ok {
			ch <- c
			delete(taskChans, c.ID)
		}
		taskChansMutex.Unlock()
	}
}

func RemoveSolve(taskId string) {
	removeFromQueue(taskId)

	for _, h := range Harvesters {
		if h.Solver != nil && h.Solver.ID == taskId {
			Remove(h.ID)

			h.Solver = nil
			EditHarvester(h)

			Create(h.ID)

			HandleQueue(h)
			return
		}
	}
}

//RequestSolve requests a solve from collective or autosolve solvers.
func RequestSolve(request Solver) (SolveChans, error) {
	in := integrations.Get()
	if in.Aycd.Active && request.ReCaptcha != nil {
		a, err := handleAycd(request)

		respChans := SolveChans{
			Autosolve: a,
		}

		if err == nil {
			return respChans, err
		}
	}

	c, err := RequestCollectiveSolve(request)

	respChans := SolveChans{
		Collective: c,
	}

	if err != nil {
		return SolveChans{}, err
	}

	return respChans, err
}

//GetSolve waits for a solve from the collective or autosolve.
func GetSolve(solveChans SolveChans) (resp SolverResponse, err error) {
	if solveChans.Autosolve != nil {
		a := <-solveChans.Autosolve

		resp = SolverResponse{
			ID:        a.TaskId,
			Success:   true,
			CreatedAt: a.CreatedAt,
			Data: ResponseData{
				Token: a.Token,
			},
		}

		return
	}

	if solveChans.Collective != nil {
		resp = <-solveChans.Collective
		return
	}

	return resp, ErrSolverNotFound
}

func RequestCollectiveSolve(request Solver) (respChan chan SolverResponse, err error) {
	respChan = make(chan SolverResponse)
	taskChansMutex.Lock()
	taskChans[request.ID] = respChan
	taskChansMutex.Unlock()

	var idle []HarvesterData
	for _, h := range Harvesters {
		if request.Type == Gmail_Login && h.Focused {

			if Harvesters[h.ID].Solver != nil {
				addToQueue(*Harvesters[h.ID].Solver)
			}

			harvester := Harvesters[h.ID]
			harvester.Solver = &request
			EditHarvester(harvester)

			err = Load(harvester.ID)
			return
		}

		if h.Solver == nil {
			idle = append(idle, h)
		}
	}

	if len(idle) == 0 {
		addToQueue(request)
		return
	}

	var match *HarvesterData

	for _, i := range idle {
		if filter(i, request) {
			match = &i

			if i.Focused {
				break
			}
		}
	}

	if match != nil {
		harvester := Harvesters[match.ID]
		harvester.Solver = &request
		EditHarvester(harvester)

		err = Load(harvester.ID)

		return
	}

	addToQueue(request)

	return
}

func ReceiveSolve(solve SolverResponse) {
	for _, h := range Harvesters {
		if h.Solver != nil && h.Solver.ID == solve.ID {
			Remove(h.ID)

			h.Solver = nil
			EditHarvester(h)

			Create(h.ID)

			HandleQueue(h)
		}
	}
	solveChan <- solve
}

func HandleQueue(harvester HarvesterData) Solver {
	var request Solver
	if len(queue) == 0 {
		return request
	}

	for _, r := range queue {
		if filter(harvester, r) {
			removeFromQueue(r.ID)
			return r
		}
	}

	return request
}

func addToQueue(request Solver) {
	queueMutex.Lock()
	queue = append(queue, request)
	queueMutex.Unlock()
}

func removeFromQueue(TaskId string) {
	idx := getSolveQueueIndex(TaskId)
	if idx == -1 {
		return
	}

	queueMutex.Lock()
	queue = append(queue[:idx], queue[idx+1:]...)
	queueMutex.Unlock()
}

func getSolveQueueIndex(TaskId string) int {
	queueMutex.Lock()
	for k, v := range queue {
		if TaskId == v.ID {
			queueMutex.Unlock()
			return k
		}
	}
	queueMutex.Unlock()
	return -1
}

func filter(harvester HarvesterData, request Solver) bool {
	if harvester.Store == nil {
		return true
	}

	if harvester.Store.Platform == "" && harvester.Store.Url == "" {
		return true
	}

	if harvester.Store.Url == request.Store.Url {
		return true
	}

	if harvester.Store.Platform == request.Store.Platform {
		return harvester.Store.Url == ""
	}

	return false
}

// convert the typing and then return the chan.
func handleAycd(request Solver) (respChan chan autosolve.CaptchaToken, err error) {
	var proxy string
	proxyRequired := false
	if request.Proxy != nil {
		proxyRequired = true
		proxy = *request.Proxy
	}

	cookieStr := ""
	if len(request.Cookies) > 0 {
		for _, c := range request.Cookies {
			s := c.String()
			cookieStr = cookieStr + s
		}
	}

	var captchaType int
	switch request.Type {
	case ReCaptcha_V2:
		captchaType = 0
	case ReCaptcha_V2_Invisible:
		captchaType = 1
	case ReCaptcha_V3:
		captchaType = 2
	}

	respChan, err = autosolve.RequestCaptchaToken(autosolve.CaptchaRequest{
		TaskId:           request.ID,
		Url:              request.Store.Url,
		SiteKey:          request.ReCaptcha.SiteKey,
		Version:          captchaType,
		Action:           request.ReCaptcha.Action,
		Proxy:            proxy,
		ProxyRequired:    proxyRequired,
		Cookies:          cookieStr,
		RenderParameters: request.ReCaptcha.Parameters,
		CreatedAt:        request.RequestedAt,
	})

	return
}
