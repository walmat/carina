package rates

import (
	"errors"
	"nebula/pkg/infra/profiles"
	"nebula/pkg/infra/tasks"
	"sync"
)

var (
	fetchers      = make(map[string]RateFetcher)
	fetchersMutex = sync.RWMutex{}

	shippingRates      = make(map[string]map[string]map[string][]ShippingRate)
	shippingRatesMutex = sync.RWMutex{}

	UnsupportedStoreErr = errors.New("unsupported stores")
)

type RateFetcher interface {
	FetchRates(string, *profiles.ProfileData) (*RateInformation, error)
}

type RateInformation struct {
	Discriminator string         `json:"discriminator"`
	Rates         []ShippingRate `json:"rates"`
}

type ShippingRate struct {
	Price string `json:"price"`
	Code  string `json:"code"`
}

func Initialize() error {
	taskTypes := tasks.GetTypes()
	for _, typ := range taskTypes {
		m := tasks.GetModes(typ)
		staticCtx, err := tasks.GetStaticContextForTaskType(typ, m[0])
		if err != nil {
			return err
		}

		if rateFetcher, ok := staticCtx.(RateFetcher); ok {
			fetchersMutex.Lock()
			fetchers[typ] = rateFetcher
			fetchersMutex.Unlock()
		}
	}
	return nil
}

func fetchRates(taskType, productUrl string, profile *profiles.ProfileData) (*RateInformation, error) {
	fetchersMutex.RLock()
	rateFetcher, ok := fetchers[taskType]
	fetchersMutex.RUnlock()
	if !ok {
		return nil, UnsupportedStoreErr
	}

	return rateFetcher.FetchRates(productUrl, profile)
}

func FetchRates(taskType, productUrl, profileId string) (*RateInformation, error) {
	profile, err := profiles.GetProfile(profileId)
	if err != nil {
		return nil, err
	}

	rateInfo, err := fetchRates(taskType, productUrl, profile)
	if err != nil {
		return nil, err
	}

	shippingRatesMutex.Lock()
	if discrimMap, ok := shippingRates[taskType]; ok {
		if ratesMap, ok := discrimMap[rateInfo.Discriminator]; ok {
			ratesMap[profile.Shipping.PostCode] = rateInfo.Rates
		} else {
			discrimMap[rateInfo.Discriminator] = make(map[string][]ShippingRate)
			shippingRates[taskType][rateInfo.Discriminator][profile.Shipping.PostCode] = rateInfo.Rates
		}
	} else {
		shippingRates[taskType] = make(map[string]map[string][]ShippingRate)
		shippingRates[taskType][rateInfo.Discriminator] = make(map[string][]ShippingRate)
		shippingRates[taskType][rateInfo.Discriminator][profile.Shipping.PostCode] = rateInfo.Rates
	}
	shippingRatesMutex.Unlock()

	return rateInfo, nil
}
