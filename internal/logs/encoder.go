package logs

func encodeFieldsToJSON(fields ...Field) string {
	log := "{"

	for i := 0; i < len(fields)-1; i++ {
		log += fields[i].toJSON()
		log += ","
	}

	log += fields[len(fields)-1].toJSON()

	log += "}"

	return log
}
