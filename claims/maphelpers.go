package claims

// GetIssuer gets "iss" from the map as a string
func GetIssuer(m map[string]interface{}) string {
	iss, ok := m["iss"]
	if !ok {
		return ""
	}
	issStr, ok := iss.(string)
	if !ok {
		return ""
	}
	return issStr
}

// GetSubject gets "sub" from the map as a string
func GetSubject(m map[string]interface{}) string {
	sub, ok := m["sub"]
	if !ok {
		return ""
	}
	subStr, ok := sub.(string)
	if !ok {
		return ""
	}
	return subStr
}

// GetAudience gets "aud" from the map as a string slice
func GetAudience(m map[string]interface{}) []string {
	aud, ok := m["aud"]
	if !ok {
		return []string{}
	}
	switch audSlc := aud.(type) {
	case []interface{}:
		audStrs := make([]string, len(audSlc))
		for n, v := range audSlc {
			val, ok := v.(string)
			if !ok {
				return []string{}
			}
			audStrs[n] = val
		}
		return audStrs
	case []string:
		return audSlc
	}
	return []string{}
}

// GetExpiresAt gets "exp" from the map as an int64 that represents epoch time
func GetExpiresAt(m map[string]interface{}) int64 {
	exp, ok := m["exp"]
	if !ok {
		return int64(0)
	}
	switch v := exp.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	}
	return int64(0)
}

// GetIssuedAt gets "iat" from the map as an int64 that represents epoch time
func GetIssuedAt(m map[string]interface{}) int64 {
	exp, ok := m["iat"]
	if !ok {
		return int64(0)
	}
	switch v := exp.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	}
	return int64(0)
}
