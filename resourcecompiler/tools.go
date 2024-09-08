package resourcecompiler

func JsonBlock(raw interface{}) uint32{
	bStr, ok := raw.(string)
	if ok{
		return findBlock(bStr)
	}

	bMap, ok := raw.(map[string]interface{})
	if ok{
		name := bMap["name"].(string)
		return findBlock(name)
		// delete(bMap, "name")

		// world.BlockByName(name, bMap)
	}

	return 0
}