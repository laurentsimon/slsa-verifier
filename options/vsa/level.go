package vsa

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	errVsaDuplicateTrack = errors.New("duplicate Vsa track")
	errVsaMismatchTrack  = errors.New("mismatch Vsa track")
	errVsaInvalidTrack   = errors.New("invalid Vsa track")
)

type Level interface {
	ToString() string
	IsTrack(string) bool
	Track() string
	LowerThan(Level) bool
	GreaterThan(Level) bool
	EqualTo(Level) bool
	ToInt() uint
}
type BuildLevel int

const (
	BuildLevel0 BuildLevel = iota
	BuildLevel1
	BuildLevel2
	BuildLevel3
)

func LevelFromString(s string) (Level, error) {
	if strings.HasPrefix(s, "SLSA_BUILD_LEVEL_") {
		return BuildLevelFromString(s)
	}
	if strings.HasPrefix(s, "SLSA_SOURCE_LEVEL_") {
		return SourceLevelFromString(s)
	}
	return nil, errVsaInvalidTrack
}

func BuildLevelFromString(s string) (*BuildLevel, error) {
	bl, err := parseLevel[BuildLevel](s, "SLSA_BUILD_LEVEL_")
	if err != nil {
		return nil, err
	}
	return bl, nil
}

func (l BuildLevel) New() *BuildLevel {
	return &l
}

func (l BuildLevel) ToInt() uint {
	return uint(l)
}

func (l *BuildLevel) Track() string {
	return "build"
}

func (l *BuildLevel) ToString() string {
	return fmt.Sprintf("SLSA_BUILD_LEVEL_%d", *l)
}

func (l *BuildLevel) IsTrack(s string) bool {
	return strings.ToLower(s) == "build"
}

func (l *BuildLevel) LowerThan(o Level) bool {
	ol, _ := o.(*BuildLevel)
	return *l < *ol
}

func (l *BuildLevel) GreaterThan(o Level) bool {
	ol, _ := o.(*BuildLevel)
	return *l > *ol
}

func (l *BuildLevel) EqualTo(o Level) bool {
	ol, _ := o.(*BuildLevel)
	return *l == *ol
}

func SourceLevelFromString(s string) (Level, error) {
	sl, err := parseLevel[SourceLevel](s, "SLSA_SOURCE_LEVEL_")
	if err != nil {
		return nil, err
	}
	return sl, nil
}

func (l SourceLevel) New() *SourceLevel {
	return &l
}

func (l SourceLevel) ToInt() uint {
	return uint(l)
}

func (l *SourceLevel) ToString() string {
	return fmt.Sprintf("SLSA_SOURCE_LEVEL_%d", *l)
}

func (l *SourceLevel) Track() string {
	return "source"
}

func (l *SourceLevel) IsTrack(s string) bool {
	return strings.ToLower(s) == "source"
}

func (l *SourceLevel) LowerThan(o Level) bool {
	ol, _ := o.(*SourceLevel)
	return *l < *ol
}

func (l *SourceLevel) GreaterThan(o Level) bool {
	ol, _ := o.(*SourceLevel)
	return *l > *ol
}

func (l *SourceLevel) EqualTo(o Level) bool {
	ol, _ := o.(*SourceLevel)
	return *l == *ol
}

type SourceLevel int

const (
	SourceLevel0 SourceLevel = iota
	SourceLevel1
	SourceLevel2
	SourceLevel3
)

func LevelsFromArray(values []string) ([]Level, error) {
	levels := make([]Level, len(values))
	for i := range values {
		levelStr := strings.TrimSpace(values[i])
		bl, berr := BuildLevelFromString(levelStr)
		if berr != nil && !errors.Is(berr, errVsaMismatchTrack) {
			return nil, berr
		}
		if berr == nil {
			levels[i] = bl
			continue
		}

		// Source level.
		sl, serr := SourceLevelFromString(levelStr)
		if serr != nil && !errors.Is(serr, errVsaMismatchTrack) {
			return nil, serr
		}
		if serr == nil {
			levels[i] = sl
			continue
		}

		return nil, fmt.Errorf("parse: [build track:%w, source track:%w]", berr, serr)
	}

	tracks := map[string]bool{"build": false, "source": false}
	for i := range levels {
		level := levels[i]
		trackName := level.Track()
		if tracks[trackName] {
			return nil, fmt.Errorf("%w: %s", errVsaDuplicateTrack, trackName)
		}
		tracks[trackName] = true
	}

	return levels, nil
}

//nolinnt:deadcode
type levelType interface {
	BuildLevel | SourceLevel
}

func parseLevel[T levelType](levelStr, prefix string) (*T, error) {
	if !strings.HasPrefix(levelStr, prefix) ||
		len(levelStr) != len(prefix)+1 {
		return nil, fmt.Errorf("%w: expected %s. Got %s", errVsaMismatchTrack, prefix, levelStr)
	}
	level, err := strconv.Atoi(levelStr[len(levelStr)-1:])
	if err != nil {
		return nil, fmt.Errorf("%w: not a number: %q", err, levelStr)
	}
	if level < 0 || level > 4 {
		return nil, fmt.Errorf("invalid level: %q", levelStr)
	}
	ret := T(level)
	return &ret, nil
}
