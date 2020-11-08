package servlet

import "testing"

func Test_LogLevel(t *testing.T) {
	t.Run("levels have correct priorities", func(t *testing.T) {
		scenarios := []struct {
			lower      LogLevel
			lowerName  string
			higher     LogLevel
			higherName string
		}{
			{
				lower:      FATAL,
				lowerName:  "FATAL",
				higher:     ERROR,
				higherName: "ERROR",
			},
			{
				lower:      ERROR,
				lowerName:  "ERROR",
				higher:     WARNING,
				higherName: "WARNING",
			},
			{
				lower:      WARNING,
				lowerName:  "WARNING",
				higher:     NOTICE,
				higherName: "NOTICE",
			},
			{
				lower:      NOTICE,
				lowerName:  "NOTICE",
				higher:     INFO,
				higherName: "INFO",
			},
			{
				lower:      INFO,
				lowerName:  "INFO",
				higher:     DEBUG,
				higherName: "DEBUG",
			},
		}

		for _, scn := range scenarios {
			if scn.lower > scn.higher {
				t.Errorf("lower %s greater then %s", scn.lowerName, scn.higherName)
			}
		}
	})
}

func Test_LogLevelMap(t *testing.T) {
	t.Run("level map have correct priorities", func(t *testing.T) {
		scenarios := []struct {
			name  string
			level LogLevel
		}{
			{
				name:  "fatal",
				level: FATAL,
			},
			{
				name:  "error",
				level: ERROR,
			},
			{
				name:  "warning",
				level: WARNING,
			},
			{
				name:  "notice",
				level: NOTICE,
			},
			{
				name:  "info",
				level: INFO,
			},
			{
				name:  "debug",
				level: DEBUG,
			},
		}

		for _, scn := range scenarios {
			if scn.level != LogLevelMap[scn.name] {
				t.Errorf("(%s) did not correspond to (%v) level", scn.name, scn.level)
			}
		}
	})
}

func Test_LogLevelNameMap(t *testing.T) {
	t.Run("level map have correct priorities", func(t *testing.T) {
		scenarios := []struct {
			name  string
			level LogLevel
		}{
			{
				name:  "fatal",
				level: FATAL,
			},
			{
				name:  "error",
				level: ERROR,
			},
			{
				name:  "warning",
				level: WARNING,
			},
			{
				name:  "notice",
				level: NOTICE,
			},
			{
				name:  "info",
				level: INFO,
			},
			{
				name:  "debug",
				level: DEBUG,
			},
		}

		for _, scn := range scenarios {
			if scn.name != LogLevelNameMap[scn.level] {
				t.Errorf("(%v) did not correspond to (%s) name", scn.level, scn.name)
			}
		}
	})
}
