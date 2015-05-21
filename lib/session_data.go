package lib

import "encoding/json"

// SessionDataEntry ...
type SessionDataEntry struct {
	ID    SessionID
	Key   string
	Value string
}

// SessionData ...
type SessionData map[string]string

// DecodeSessionDataFromJSON ...
func DecodeSessionDataFromJSON(data []byte) (SessionData, error) {
	var sessionData SessionData
	err := json.Unmarshal(data, &sessionData)
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

// EncodeToJSON ...
func (sd SessionData) EncodeToJSON() ([]byte, error) {
	return json.Marshal(sd)
}

// Set ...
func (sd *SessionData) Set(key, value string) {
	(*sd)[key] = value
}

// Get ...
func (sd *SessionData) Get(key string) string {
	return (*sd)[key]
}
