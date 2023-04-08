package vsa

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_LevelsFromArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		array  []string
		levels []Level
		err    error
	}{
		{
			name:   "single build level 2",
			array:  []string{"SLSA_BUILD_LEVEL_2"},
			levels: []Level{BuildLevel2.New()},
		},
		{
			name:   "single source level 2",
			array:  []string{"SLSA_SOURCE_LEVEL_2"},
			levels: []Level{SourceLevel2.New()},
		},
		{
			name:   "source and build level 2",
			array:  []string{"SLSA_SOURCE_LEVEL_2", "SLSA_BUILD_LEVEL_2"},
			levels: []Level{SourceLevel2.New(), BuildLevel2.New()},
		},
		{
			name:   "source level 1 and build level 2",
			array:  []string{"SLSA_SOURCE_LEVEL_1", "SLSA_BUILD_LEVEL_2"},
			levels: []Level{SourceLevel1.New(), BuildLevel2.New()},
		},
		{
			name:  "duplicate source level 1 and 2",
			array: []string{"SLSA_SOURCE_LEVEL_1", "SLSA_SOURCE_LEVEL_2"},
			err:   errVsaDuplicateTrack,
		},
		{
			name:  "duplicate source level 1 and 1",
			array: []string{"SLSA_SOURCE_LEVEL_1", "SLSA_SOURCE_LEVEL_1"},
			err:   errVsaDuplicateTrack,
		},
		{
			name:  "duplicate build level 1 and 2",
			array: []string{"SLSA_BUILD_LEVEL_1", "SLSA_BUILD_LEVEL_2"},
			err:   errVsaDuplicateTrack,
		},
		{
			name:  "duplicate build level 1 and 1",
			array: []string{"SLSA_BUILD_LEVEL_1", "SLSA_BUILD_LEVEL_1"},
			err:   errVsaDuplicateTrack,
		},
	}
	for _, tt := range tests {
		tt := tt // Re-initializing variable so it is not changed while executing the closure below
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			levels, err := LevelsFromArray(tt.array)
			if !cmp.Equal(err, tt.err, cmpopts.EquateErrors()) {
				t.Errorf(cmp.Diff(err, tt.err))
			}

			if err != nil {
				return
			}

			if !cmp.Equal(levels, tt.levels) {
				t.Errorf(cmp.Diff(levels, tt.levels))
			}
		})
	}
}

// TODO: tests for lower, greater, etc
