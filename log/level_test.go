package log

import "testing"

func Test_Level(t *testing.T) {
	t.Run("levels have correct priorities", func(t *testing.T) {
		scenarios := []struct {
			lower      Level
			lowerName  string
			higher     Level
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

func Test_LevelMap(t *testing.T) {
	t.Run("level map have correct priorities", func(t *testing.T) {
		scenarios := []struct {
			name      string
			level     Level
			levelName string
		}{
			{
				name:      "fatal",
				level:     FATAL,
				levelName: "FATAL",
			},
			{
				name:      "error",
				level:     ERROR,
				levelName: "ERROR",
			},
			{
				name:      "warning",
				level:     WARNING,
				levelName: "WARNING",
			},
			{
				name:      "notice",
				level:     NOTICE,
				levelName: "NOTICE",
			},
			{
				name:      "info",
				level:     INFO,
				levelName: "INFO",
			},
			{
				name:      "debug",
				level:     DEBUG,
				levelName: "DEBUG",
			},
		}

		for _, scn := range scenarios {
			if scn.level != LevelMap[scn.name] {
				t.Errorf("(%s) did not correspond to (%s) level", scn.name, scn.levelName)
			}
		}
	})
}
