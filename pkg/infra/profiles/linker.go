package profiles

import (
	"errors"
	"nebula/pkg/infra/eventbus"
	"nebula/pkg/infra/tasks"
	"sync"
)

var (
	links   = make(map[string]string)
	linksMu sync.RWMutex

	NotLinkedErr = errors.New("profiles: profile not linked")
)

func init() {
	// remove links related to deleted tasks
	_ = eventbus.Subscribe(func(e *tasks.TaskRemovedEvent) {
		linksMu.Lock()
		defer linksMu.Unlock()
		for k := range links {
			if k == e.Id {
				delete(links, k)
			}
		}
	})

	// remove tasks that are linked to this profile
	_ = eventbus.Subscribe(func(e *ProfileRemovedEvent) {
		linksMu.Lock()
		defer linksMu.Unlock()
		for k, v := range links {
			if v == e.Id {
				_ = tasks.RemoveTask(k)
				delete(links, k)
			}
		}
	})
}

func AssignProfileToTask(taskId, profileId string) error {
	profilesMu.RLock()
	defer profilesMu.RUnlock()
	if !profileExists(profileId) {
		return ProfileNotExistErr
	}

	linksMu.Lock()
	defer linksMu.Unlock()
	links[taskId] = profileId
	eventbus.Publish(&ProfileAssignedToTaskEvent{
		Id:     profileId,
		TaskId: taskId,
	})
	return nil
}

func GetProfileForTask(taskId string) (*ProfileData, error) {
	linksMu.RLock()
	defer linksMu.RUnlock()
	profilesMu.RLock()
	defer profilesMu.RUnlock()
	profileId, ok := links[taskId]
	if !ok {
		return nil, NotLinkedErr
	} else if !profileExists(profileId) {
		return nil, ProfileNotExistErr
	}

	groupId := profileToGroup[profileId]
	profileData := profileGroups[groupId].Profiles[profileId]
	return &profileData, nil
}
