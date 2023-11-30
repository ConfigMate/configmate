package types

func GetTypes() []string {
	return []string{
		"bool",
		"int",
		"float",
		"string",
		"list",
		"file",
		"host",
		"port",
		"host_port",
		"custom_object",
	}
}

func GetTypeInfo(typename string) map[string]string {
	info := map[string]map[string]string{
		"bool":          tBoolMethodsDescriptions,
		"int":           tIntMethodsDescriptions,
		"float":         tFloatMethodsDescriptions,
		"string":        tStringMethodsDescriptions,
		"list":          tListMethodsDescriptions,
		"file":          tFileMethodsDescriptions,
		"host":          tHostMethodsDescriptions,
		"port":          tPortMethodsDescriptions,
		"host_port":     tHostPortMethodsDescriptions,
		"custom_object": tCustomObjectMethodsDescriptions,
	}

	if _, ok := info[typename]; !ok {
		return nil
	}

	return info[typename]
}
