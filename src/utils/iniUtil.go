package utils

import (
	"gopkg.in/ini.v1"
)

type IniParser struct {
	confReader *ini.File // config reader
}

type IniParserError struct {
	errorInfo string
}

func (e *IniParserError) Error() string { return e.errorInfo }

func (i *IniParser) Load(configFilename string) error {
	conf, err := ini.Load(configFilename)
	if err != nil {
		i.confReader = nil
		return err
	}
	i.confReader = conf
	return nil
}

func (i *IniParser) GetString(section string, key string) string {
	if i.confReader == nil {
		return ""
	}

	s := i.confReader.Section(section)
	if s == nil {
		return ""
	}

	return s.Key(key).String()
}

func (i *IniParser) GetInt32(section string, key string) int32 {
	if i.confReader == nil {
		return 0
	}

	s := i.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Int()

	return int32(valueInt)
}

func (i *IniParser) GetUint32(section string, key string) uint32 {
	if i.confReader == nil {
		return 0
	}

	s := i.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Uint()

	return uint32(valueInt)
}

func (i *IniParser) GetInt64(section string, key string) int64 {
	if i.confReader == nil {
		return 0
	}

	s := i.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Int64()
	return valueInt
}

func (i *IniParser) GetUint64(section string, key string) uint64 {
	if i.confReader == nil {
		return 0
	}

	s := i.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Uint64()
	return valueInt
}

func (i *IniParser) GetFloat32(section string, key string) float32 {
	if i.confReader == nil {
		return 0
	}

	s := i.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueFloat, _ := s.Key(key).Float64()
	return float32(valueFloat)
}

func (i *IniParser) GetFloat64(section string, key string) float64 {
	if i.confReader == nil {
		return 0
	}

	s := i.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueFloat, _ := s.Key(key).Float64()
	return valueFloat
}

