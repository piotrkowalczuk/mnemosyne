package shared

import "time"

// Context implements sklog.Contexter interface.
func (gr *GetRequest) Context() []interface{} {
	return []interface{}{"id", gr.Id}
}

// Context implements sklog.Contexter interface.
func (lr *ListRequest) Context() []interface{} {
	return []interface{}{
		"offset", lr.Offset,
		"limit", lr.Limit,
		"expire_at_from", lr.ExpireAtFrom,
		"expire_at_to", lr.ExpireAtTo,
	}
}

// Context implements sklog.Contexter interface.
func (er *ExistsRequest) Context() []interface{} {
	return []interface{}{"id", er.Id}
}

// Context implements sklog.Contexter interface.
func (er *CreateRequest) Context() (ctx []interface{}) {
	for key, value := range er.Data {
		ctx = append(ctx, "data_"+key, value)
	}

	return
}

// Context implements sklog.Contexter interface.
func (ar *AbandonRequest) Context() []interface{} {
	return []interface{}{
		"id", ar.Id,
	}
}

// Context implements sklog.Contexter interface.
func (sdr *SetDataRequest) Context() []interface{} {
	return []interface{}{
		"id", sdr.Id,
		"key", sdr.Key,
		"value", sdr.Value,
	}
}

// Context implements sklog.Contexter interface.
func (dr *DeleteRequest) Context() []interface{} {
	return []interface{}{
		"id", dr.Id,
		"expire_at_from", dr.ExpireAtFrom,
		"expire_at_to", dr.ExpireAtTo,
	}
}

// SetValue ...
func (s *Session) SetValue(key, value string) {
	if s.Data == nil {
		s.Data = make(map[string]string)
	}

	s.Data[key] = value
}

// GetDataForKey...
func (s *Session) Value(key string) string {
	if s.Data == nil {
		s.Data = make(map[string]string)
	}

	return s.Data[key]
}

// ParseTime ...
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
