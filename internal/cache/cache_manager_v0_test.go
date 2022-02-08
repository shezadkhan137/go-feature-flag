package cache_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomaspoignant/go-feature-flag/internal/cache"
	"github.com/thomaspoignant/go-feature-flag/internal/flag"
	"github.com/thomaspoignant/go-feature-flag/internal/notifier"
	"github.com/thomaspoignant/go-feature-flag/testutils/testconvert"
	"testing"
)

// These tests are here to validate the backward compatibility with the version v0 of go-feature-flag
//
// We place them in a dedicated file to be able to drop them easily in the future.

func Test_FlagCache_v0(t *testing.T) {
	yamlFile := []byte(`test-flag:
  rule: key eq "random-key"
  percentage: 100
  true: true
  false: false
  default: false
  trackEvents: false
`)

	jsonFile := []byte(`{
  "test-flag": {
    "rule": "key eq \"random-key\"",
    "percentage": 100,
    "true": true,
    "false": false,
    "default": false
  }
}
`)

	tomlFile := []byte(`[test-flag]
rule = "key eq \"random-key\""
percentage = 100.0
true = true
false = false
default = false
disable = false`)

	type args struct {
		loadedFlags []byte
	}
	tests := []struct {
		name       string
		args       args
		expected   map[string]flag.FlagData
		wantErr    bool
		flagFormat string
	}{
		{
			name:       "Yaml valid",
			flagFormat: "yaml",
			args: args{
				loadedFlags: yamlFile,
			},
			expected: map[string]flag.FlagData{
				"test-flag": {
					Variations: &map[string]*interface{}{
						"Default": testconvert.Interface(false),
						"False":   testconvert.Interface(false),
						"True":    testconvert.Interface(true),
					},
					Rules: &map[string]flag.Rule{
						flag.LegacyRuleName: {
							Query:           testconvert.String("key eq \"random-key\""),
							VariationResult: nil,
							Percentages: &map[string]float64{
								"True":  100,
								"False": 0,
							},
						},
					},
					DefaultRule: &flag.Rule{
						VariationResult: testconvert.String("Default"),
					},
					TrackEvents: testconvert.Bool(false),
				},
			},
			wantErr: false,
		},
		{
			name:       "Yaml invalid file",
			flagFormat: "yaml",
			args: args{
				loadedFlags: []byte(`test-flag:
  rule: key eq "random-key"
  percentage: "toot"
  true: true
  false: false
  default: false
`),
			},
			wantErr: true,
		},
		{
			name: "JSON valid",
			args: args{
				loadedFlags: jsonFile,
			},
			flagFormat: "json",
			expected: map[string]flag.FlagData{
				"test-flag": {
					Variations: &map[string]*interface{}{
						"Default": testconvert.Interface(false),
						"False":   testconvert.Interface(false),
						"True":    testconvert.Interface(true),
					},
					Rules: &map[string]flag.Rule{
						flag.LegacyRuleName: {
							Query:           testconvert.String("key eq \"random-key\""),
							VariationResult: nil,
							Percentages: &map[string]float64{
								"True":  100,
								"False": 0,
							},
						},
					},
					DefaultRule: &flag.Rule{
						VariationResult: testconvert.String("Default"),
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "JSON invalid file",
			flagFormat: "json",
			args: args{
				loadedFlags: []byte(`{
  "test-flag": {
    "rule": "key eq \"random-key\"",
    "percentage": 100,
    "true": true,
    "false": false,
    "default": false"
  }
}`),
			},
			wantErr: true,
		},
		{
			name: "TOML valid",
			args: args{
				loadedFlags: tomlFile,
			},
			flagFormat: "toml",
			expected: map[string]flag.FlagData{
				"test-flag": {
					Variations: &map[string]*interface{}{
						"Default": testconvert.Interface(false),
						"False":   testconvert.Interface(false),
						"True":    testconvert.Interface(true),
					},
					Rules: &map[string]flag.Rule{
						flag.LegacyRuleName: {
							Query:           testconvert.String("key eq \"random-key\""),
							VariationResult: nil,
							Percentages: &map[string]float64{
								"True":  100,
								"False": 0,
							},
						},
					},
					DefaultRule: &flag.Rule{
						VariationResult: testconvert.String("Default"),
					},
					Disable: testconvert.Bool(false),
				},
			},
			wantErr: false,
		},
		{
			name: "TOML invalid file",
			args: args{
				loadedFlags: []byte(`[test-flag]
rule = "key eq \"random-key\""
percentage = 100.0
true = true
false = false
default = false"
disable = false`),
			},
			flagFormat: "toml",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fCache := cache.New(cache.NewNotificationService([]notifier.Notifier{}))
			err := fCache.UpdateCache(tt.args.loadedFlags, tt.flagFormat)
			if tt.wantErr {
				assert.Error(t, err, "UpdateCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NoError(t, err, "UpdateCache() error = %v, wantErr %v", err, tt.wantErr)
			// If no error we compare with expected
			for key, expected := range tt.expected {
				got, _ := fCache.GetFlag(key)
				assert.Equal(t, &expected, got) // nolint
			}
			fCache.Close()
		})
	}
}

func Test_AllFlags_v0(t *testing.T) {
	yamlFile := []byte(`test-flag:
  rule: key eq "random-key"
  percentage: 100
  true: true
  false: false
  default: false
  trackEvents: false
`)

	type args struct {
		loadedFlags []byte
	}
	tests := []struct {
		name       string
		args       args
		expected   map[string]flag.FlagData
		wantErr    bool
		flagFormat string
	}{
		{
			name:       "Yaml valid",
			flagFormat: "yaml",
			args: args{
				loadedFlags: yamlFile,
			},
			expected: map[string]flag.FlagData{
				"test-flag": {
					Variations: &map[string]*interface{}{
						"Default": testconvert.Interface(false),
						"False":   testconvert.Interface(false),
						"True":    testconvert.Interface(true),
					},
					Rules: &map[string]flag.Rule{
						flag.LegacyRuleName: {
							Query:           testconvert.String("key eq \"random-key\""),
							VariationResult: nil,
							Percentages: &map[string]float64{
								"True":  100,
								"False": 0,
							},
						},
					},
					DefaultRule: &flag.Rule{
						VariationResult: testconvert.String("Default"),
					},
					TrackEvents: testconvert.Bool(false),
				},
			},
			wantErr: false,
		},
		{
			name:       "Yaml multiple flags",
			flagFormat: "yaml",
			args: args{
				loadedFlags: []byte(`test-flag:
  rule: key eq "random-key"
  percentage: 100
  true: true
  false: false
  default: false
  trackEvents: false
test-flag2:
  rule: key eq "random-key"
  percentage: 0
  true: "true"
  false: "false"
  default: "false"
  trackEvents: false
`),
			},
			expected: map[string]flag.FlagData{
				"test-flag": {
					Variations: &map[string]*interface{}{
						"Default": testconvert.Interface(false),
						"False":   testconvert.Interface(false),
						"True":    testconvert.Interface(true),
					},
					Rules: &map[string]flag.Rule{
						flag.LegacyRuleName: {
							Query:           testconvert.String("key eq \"random-key\""),
							VariationResult: nil,
							Percentages: &map[string]float64{
								"True":  100,
								"False": 0,
							},
						},
					},
					DefaultRule: &flag.Rule{
						VariationResult: testconvert.String("Default"),
					},
					TrackEvents: testconvert.Bool(false),
				},
				"test-flag2": {
					Variations: &map[string]*interface{}{
						"Default": testconvert.Interface("false"),
						"False":   testconvert.Interface("false"),
						"True":    testconvert.Interface("true"),
					},
					Rules: &map[string]flag.Rule{
						flag.LegacyRuleName: {
							Query:           testconvert.String("key eq \"random-key\""),
							VariationResult: nil,
							Percentages: &map[string]float64{
								"True":  0,
								"False": 100,
							},
						},
					},
					DefaultRule: &flag.Rule{
						VariationResult: testconvert.String("Default"),
					},
					TrackEvents: testconvert.Bool(false),
				},
			},
			wantErr: false,
		}, {
			name:       "empty",
			flagFormat: "yaml",
			args: args{
				loadedFlags: []byte(``),
			},
			expected: map[string]flag.FlagData{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fCache := cache.New(cache.NewNotificationService([]notifier.Notifier{}))
			_ = fCache.UpdateCache(tt.args.loadedFlags, tt.flagFormat)

			allFlags, err := fCache.AllFlags()
			if tt.wantErr {
				assert.Error(t, err, "UpdateCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NoError(t, err)

			// If no error we compare with expected
			for key, expected := range tt.expected {
				got := allFlags[key]
				assert.Equal(t, &expected, got) //nolint: gosec
			}
			fCache.Close()
		})
	}
}
