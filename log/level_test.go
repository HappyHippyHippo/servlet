package log

import "testing"

func Test_Level(t *testing.T) {
	t.Run("defined levels should have correct priorities", func(t *testing.T) {
		action := "Comparing log levels priority"

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
				higher:     DEBUG,
				higherName: "DEBUG",
			},
		}

		for _, scn := range scenarios {
			if scn.lower > scn.higher {
				t.Errorf("%s invalidated, lower %s (%v) greater then higher %s (%v)", action, scn.lowerName, scn.lower, scn.higherName, scn.higher)
			}
		}
	})
}

func Test_LevelMap(t *testing.T) {
	t.Run("defined level map should have correct priorities", func(t *testing.T) {
		action := "Checking the level map value "

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
				name:      "debug",
				level:     DEBUG,
				levelName: "DEBUG",
			},
		}

		for _, scn := range scenarios {
			if scn.level != LevelMap[scn.name] {
				t.Errorf("%s invalidated, (%s) did not correspond to (%s) level", action, scn.name, scn.levelName)
			}
		}
	})
}
