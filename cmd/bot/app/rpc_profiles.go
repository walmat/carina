package app

import (
	"nebula/pkg/infra/profiles"
)

func AddProfiles(groupId string, data []profiles.ProfileData) ([]string, error) {
	var profileIds []string
	for _, profileData := range data {
		profileId, err := profiles.AddProfile(groupId, profileData)
		if err != nil {
			return nil, err
		}
		profileIds = append(profileIds, profileId)
	}
	return profileIds, nil
}

func EditProfiles(data []profiles.ProfileData) error {
	for _, profile := range data {
		if err := profiles.EditProfile(profile.Id, profile); err != nil {
			return err
		}
	}
	return nil
}

func MoveProfiles(profileIds []string, newGroupId string) error {
	for _, profileId := range profileIds {
		if err := profiles.MoveProfile(profileId, newGroupId); err != nil {
			return err
		}
	}
	return nil
}

func RemoveProfiles(profileIds []string) error {
	for _, profileId := range profileIds {
		if err := profiles.RemoveProfile(profileId); err != nil {
			return err
		}
	}
	return nil
}
