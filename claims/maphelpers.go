package claims

// Helpers for retrieving standard claims from a map
func getIssuer(m map[string]interface{}) string {
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

func getSubject(m map[string]interface{}) string {
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

func getAudience(m map[string]interface{}) []string {
	aud, ok := m["aud"]
	if !ok {
		return []string{}
	}
	audSlc, ok := aud.([]interface{})
	if !ok {
		return []string{}
	}
	audStrs := make([]string, len(audSlc))
	for n, v := range audSlc {
		val, ok := v.(string)
		if !ok {
			return []string{}
		}
		audStrs[n] = val
	}
	return audStrs
}

func getExpiresAt(m map[string]interface{}) int64 {
	exp, ok := m["exp"]
	if !ok {
		return int64(0)
	}
	expFloat, ok := exp.(float64)
	if !ok {
		return int64(0)
	}
	return int64(expFloat)
}

func getNotBefore(m map[string]interface{}) int64 {
	nbf, ok := m["nbf"]
	if !ok {
		return int64(0)
	}
	nbfFloat, ok := nbf.(float64)
	if !ok {
		return int64(0)
	}
	return int64(nbfFloat)
}

func getIssuedAt(m map[string]interface{}) int64 {
	iat, ok := m["iat"]
	if !ok {
		return int64(0)
	}
	iatFloat, ok := iat.(float64)
	if !ok {
		return int64(0)
	}
	return int64(iatFloat)
}

func getID(m map[string]interface{}) string {
	jti, ok := m["jti"]
	if !ok {
		return ""
	}
	jtiStr, ok := jti.(string)
	if !ok {
		return ""
	}
	return jtiStr
}
