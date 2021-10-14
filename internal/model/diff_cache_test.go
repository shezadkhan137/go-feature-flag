package model_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomaspoignant/go-feature-flag/internal/flag"
	"github.com/thomaspoignant/go-feature-flag/internal/flagv2"
	"github.com/thomaspoignant/go-feature-flag/internal/model"
	"github.com/thomaspoignant/go-feature-flag/testutils/testconvert"
	"testing"
)

func TestDiffCache_HasDiff(t *testing.T) {
	type fields struct {
		Deleted map[string]flag.Flag
		Added   map[string]flag.Flag
		Updated map[string]model.DiffUpdated
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "null fields",
			fields: fields{},
			want:   false,
		},
		{
			name: "empty fields",
			fields: fields{
				Deleted: map[string]flag.Flag{},
				Added:   map[string]flag.Flag{},
				Updated: map[string]model.DiffUpdated{},
			},
			want: false,
		},
		{
			name: "only Deleted",
			fields: fields{
				Deleted: map[string]flag.Flag{
					"flag": &flagv2.FlagData{
						Variations: &map[string]*interface{}{
							"Default": testconvert.Interface(true),
							"False":   testconvert.Interface(true),
							"True":    testconvert.Interface(true),
						},
						DefaultRule: &flagv2.Rule{
							Percentages: &map[string]float64{"True": 100, "False": 0},
						},
					},
				},
				Added:   map[string]flag.Flag{},
				Updated: map[string]model.DiffUpdated{},
			},
			want: true,
		},
		{
			name: "only Added",
			fields: fields{
				Added: map[string]flag.Flag{
					"flag": &flagv2.FlagData{
						Variations: &map[string]*interface{}{
							"Default": testconvert.Interface(true),
							"False":   testconvert.Interface(true),
							"True":    testconvert.Interface(true),
						},
						DefaultRule: &flagv2.Rule{
							Percentages: &map[string]float64{"True": 100, "False": 0},
						},
					},
				},
				Deleted: map[string]flag.Flag{},
				Updated: map[string]model.DiffUpdated{},
			},
			want: true,
		},
		{
			name: "only Updated",
			fields: fields{
				Added:   map[string]flag.Flag{},
				Deleted: map[string]flag.Flag{},
				Updated: map[string]model.DiffUpdated{
					"flag": {
						Before: &flagv2.FlagData{
							Variations: &map[string]*interface{}{
								"Default": testconvert.Interface(true),
								"False":   testconvert.Interface(true),
								"True":    testconvert.Interface(true),
							},
							DefaultRule: &flagv2.Rule{
								Percentages: &map[string]float64{"True": 100, "False": 0},
							},
						},
						After: &flagv2.FlagData{
							Variations: &map[string]*interface{}{
								"Default": testconvert.Interface(false),
								"False":   testconvert.Interface(true),
								"True":    testconvert.Interface(true),
							},
							DefaultRule: &flagv2.Rule{
								Percentages: &map[string]float64{"True": 100, "False": 0},
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "all fields",
			fields: fields{
				Added: map[string]flag.Flag{
					"flag": &flagv2.FlagData{
						Variations: &map[string]*interface{}{
							"Default": testconvert.Interface(true),
							"False":   testconvert.Interface(true),
							"True":    testconvert.Interface(true),
						},
						DefaultRule: &flagv2.Rule{
							Percentages: &map[string]float64{"True": 100, "False": 0},
						},
					},
				},
				Deleted: map[string]flag.Flag{
					"flag": &flagv2.FlagData{
						Variations: &map[string]*interface{}{
							"Default": testconvert.Interface(true),
							"False":   testconvert.Interface(true),
							"True":    testconvert.Interface(true),
						},
						DefaultRule: &flagv2.Rule{
							Percentages: &map[string]float64{"True": 100, "False": 0},
						},
					},
				},
				Updated: map[string]model.DiffUpdated{
					"flag": {
						Before: &flagv2.FlagData{
							Variations: &map[string]*interface{}{
								"Default": testconvert.Interface(true),
								"False":   testconvert.Interface(true),
								"True":    testconvert.Interface(true),
							},
							DefaultRule: &flagv2.Rule{
								Percentages: &map[string]float64{"True": 100, "False": 0},
							},
						},
						After: &flagv2.FlagData{
							Variations: &map[string]*interface{}{
								"Default": testconvert.Interface(false),
								"False":   testconvert.Interface(true),
								"True":    testconvert.Interface(true),
							},
							DefaultRule: &flagv2.Rule{
								Percentages: &map[string]float64{"True": 100, "False": 0},
							},
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := model.DiffCache{
				Deleted: tt.fields.Deleted,
				Added:   tt.fields.Added,
				Updated: tt.fields.Updated,
			}
			assert.Equal(t, tt.want, d.HasDiff())
		})
	}
}
