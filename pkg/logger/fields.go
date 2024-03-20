package logger

func Error(err error) Field {
	return Field{
		"error",
		err.Error(),
	}
}

func Int64(key string, val int64) Field {
	return Field{
		key,
		val,
	}
}

func Int(key string, val int) Field {
	return Field{
		key,
		val,
	}
}
