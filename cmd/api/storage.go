// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    storageEvent, err := UnmarshalStorageEvent(bytes)
//    bytes, err = storageEvent.Marshal()

package main

import "encoding/json"

func UnmarshalStorageEvent(data []byte) (StorageEvent, error) {
	var r StorageEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *StorageEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type StorageEvent struct {
	Kind                    string `json:"kind"`
	ID                      string `json:"id"`
	SelfLink                string `json:"selfLink"`
	Name                    string `json:"name"`
	Bucket                  string `json:"bucket"`
	Generation              string `json:"generation"`
	Metageneration          string `json:"metageneration"`
	ContentType             string `json:"contentType"`
	TimeCreated             string `json:"timeCreated"`
	Updated                 string `json:"updated"`
	StorageClass            string `json:"storageClass"`
	TimeStorageClassUpdated string `json:"timeStorageClassUpdated"`
	Size                    string `json:"size"`
	Md5Hash                 string `json:"md5Hash"`
	MediaLink               string `json:"mediaLink"`
	ContentEncoding         string `json:"contentEncoding"`
	Crc32C                  string `json:"crc32c"`
	Etag                    string `json:"etag"`
}
