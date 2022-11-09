package profiles

import (
	"errors"
	"github.com/lithammer/shortuuid/v3"
	"nebula/pkg/infra/eventbus"
	"sync"
)

var (
	profileGroups  = make(map[string]ProfileGroup)
	profileToGroup = make(map[string]string)
	profilesMu     sync.RWMutex

	ProfileNotExistErr = errors.New("profiles: profile does not exist")
	GroupNotExistErr   = errors.New("profiles: group does not exist")
)

func init() {
	profileGroups["default"] = ProfileGroup{Name: "Default", Profiles: make(map[string]ProfileData)}
}

func AddProfile(groupId string, data ProfileData) (string, error) {
	profilesMu.Lock()
	defer profilesMu.Unlock()
	if !groupExists(groupId) {
		return "", GroupNotExistErr
	}

	profileId := shortuuid.New()
	data.Id = profileId

	profileGroups[groupId].Profiles[profileId] = data
	profileToGroup[profileId] = groupId
	eventbus.Publish(&ProfileAddedEvent{
		Id:      profileId,
		GroupId: groupId,
		Data:    data,
	})
	return profileId, nil
}

func EditProfile(profileId string, data ProfileData) error {
	profilesMu.Lock()
	defer profilesMu.Unlock()
	if !profileExists(profileId) {
		return ProfileNotExistErr
	}

	groupId := profileToGroup[profileId]
	profileGroups[groupId].Profiles[profileId] = data

	eventbus.Publish(&ProfileModifiedEvent{
		Id:   profileId,
		Data: data,
	})
	return nil
}

func RemoveProfile(profileId string) error {
	profilesMu.Lock()
	defer profilesMu.Unlock()
	if !profileExists(profileId) {
		return GroupNotExistErr
	}

	groupId := profileToGroup[profileId]
	delete(profileGroups[groupId].Profiles, profileId)
	delete(profileToGroup, profileId)
	eventbus.Publish(&ProfileRemovedEvent{
		Id:      profileId,
		GroupId: groupId,
	})
	return nil
}

func GetProfile(profileId string) (*ProfileData, error) {
	profilesMu.RLock()
	defer profilesMu.RUnlock()

	groupId := profileToGroup[profileId]
	if !groupExists(groupId) {
		return nil, GroupNotExistErr
	}

	if profile, ok := profileGroups[groupId].Profiles[profileId]; ok {
		return &profile, nil
	} else {
		return nil, ProfileNotExistErr
	}
}

func AddGroup(name string) (string, error) {
	profilesMu.Lock()
	defer profilesMu.Unlock()
	groupId := shortuuid.New()
	profileGroups[groupId] = ProfileGroup{
		Name:     name,
		Profiles: make(map[string]ProfileData),
	}
	eventbus.Publish(&GroupAddedEvent{
		Id: groupId,
	})
	return groupId, nil
}

func GetGroup(groupId string) (*ProfileGroup, error) {
	profilesMu.RLock()
	defer profilesMu.RUnlock()
	if !groupExists(groupId) {
		return nil, GroupNotExistErr
	}

	group := profileGroups[groupId]
	return &group, nil
}

func GetGroups() []string {
	var groups []string
	for groupId := range profileGroups {
		groups = append(groups, groupId)
	}
	return groups
}

func RemoveGroup(groupId string) error {
	profilesMu.Lock()
	defer profilesMu.Unlock()
	if !groupExists(groupId) {
		return GroupNotExistErr
	}

	for profileId := range profileGroups[groupId].Profiles {
		delete(profileToGroup, profileId)
		eventbus.Publish(&ProfileRemovedEvent{
			Id:      profileId,
			GroupId: groupId,
		})
	}

	delete(profileGroups, groupId)
	eventbus.Publish(&GroupRemovedEvent{
		Id: groupId,
	})
	return nil
}

func GetProfileGroup(profileId string) (string, error) {
	profilesMu.RLock()
	defer profilesMu.RUnlock()
	if groupId, ok := profileToGroup[profileId]; ok {
		return groupId, nil
	} else {
		return "", ProfileNotExistErr
	}
}

func ProfileExists(profileId string) bool {
	profilesMu.RLock()
	defer profilesMu.RUnlock()
	return profileExists(profileId)
}

func profileExists(profileId string) bool {
	_, ok := profileToGroup[profileId]
	return ok
}

func GroupExists(groupId string) bool {
	profilesMu.RLock()
	defer profilesMu.RUnlock()
	return groupExists(groupId)
}

func groupExists(groupId string) bool {
	_, ok := profileGroups[groupId]
	return ok
}

func MoveProfile(profileId, newGroupId string) error {
	profilesMu.Lock()
	defer profilesMu.Unlock()
	if !groupExists(newGroupId) {
		return GroupNotExistErr
	} else if !profileExists(profileId) {
		return ProfileNotExistErr
	}

	originalGroupId := profileToGroup[profileId]
	profileData := profileGroups[originalGroupId].Profiles[profileId]
	delete(profileGroups[originalGroupId].Profiles, profileId)
	profileGroups[newGroupId].Profiles[profileId] = profileData
	profileToGroup[profileId] = newGroupId
	return nil
}
