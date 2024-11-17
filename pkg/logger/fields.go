package logger

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}

func Int64(key string, value int64) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
