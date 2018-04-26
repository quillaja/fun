package hex

import (
	"math"
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

const epsilon = 1e-6

func TestHexGrid_ToWorld_FlatTop(t *testing.T) {

	type args struct {
		c float64
		r float64
	}
	tests := []struct {
		name  string
		grid  *HexGrid
		args  args
		want  float64
		want1 float64
	}{
		{
			name:  "0,0",
			grid:  NewHexGrid(1, FlatTop),
			args:  args{0, 0},
			want:  0,
			want1: 0,
		},
		{
			name:  "0,1",
			grid:  NewHexGrid(1, FlatTop),
			args:  args{0, 1},
			want:  0,
			want1: math.Sqrt(3),
		},
		{
			name:  "1,0",
			grid:  NewHexGrid(1, FlatTop),
			args:  args{1, 0},
			want:  1.5,
			want1: math.Sqrt(3),
		},
		{
			name:  "1,1",
			grid:  NewHexGrid(1, FlatTop),
			args:  args{1, 1},
			want:  1.5,
			want1: 3 * math.Sqrt(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := tt.grid
			got, got1 := grid.ToWorld(tt.args.c, tt.args.r)
			if math.Abs(got-tt.want) > epsilon {
				t.Errorf("HexGrid.ToWorld() got = %v, want %v", got, tt.want)
			}
			if math.Abs(got1-tt.want1) > epsilon {
				t.Errorf("HexGrid.ToWorld() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHexGrid_ToWorld_PointyTop(t *testing.T) {

	type args struct {
		c float64
		r float64
	}
	tests := []struct {
		name  string
		grid  *HexGrid
		args  args
		want  float64
		want1 float64
	}{
		{
			name:  "0,0",
			grid:  NewHexGrid(1, PointyTop),
			args:  args{0, 0},
			want:  0,
			want1: 0,
		},
		{
			name:  "0,1",
			grid:  NewHexGrid(1, PointyTop),
			args:  args{0, 1},
			want:  math.Sqrt(3) / 2,
			want1: 1.5,
		},
		{
			name:  "1,0",
			grid:  NewHexGrid(1, PointyTop),
			args:  args{1, 0},
			want:  math.Sqrt(3),
			want1: 0,
		},
		{
			name:  "1,1",
			grid:  NewHexGrid(1, PointyTop),
			args:  args{1, 1},
			want:  1.5 * math.Sqrt(3),
			want1: 1.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := tt.grid
			got, got1 := grid.ToWorld(tt.args.c, tt.args.r)
			if math.Abs(got-tt.want) > epsilon {
				t.Errorf("HexGrid.ToWorld() got = %v, want %v", got, tt.want)
			}
			if math.Abs(got1-tt.want1) > epsilon {
				t.Errorf("HexGrid.ToWorld() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHexGrid_ToGrid_FlatTop(t *testing.T) {

	type args struct {
		c float64
		r float64
	}
	tests := []struct {
		name  string
		grid  *HexGrid
		args  args
		want  float64
		want1 float64
	}{
		{
			name:  "0,0",
			grid:  NewHexGrid(1, FlatTop),
			args:  args{0, 0},
			want:  0,
			want1: 0,
		},
		{
			name:  "0,1",
			grid:  NewHexGrid(1, FlatTop),
			args:  args{0, math.Sqrt(3)},
			want:  0,
			want1: 1,
		},
		{
			name:  "1,0",
			grid:  NewHexGrid(1, FlatTop),
			args:  args{1.5, math.Sqrt(3) / 2},
			want:  1,
			want1: 0,
		},
		{
			name:  "1,1",
			grid:  NewHexGrid(1, FlatTop),
			args:  args{1.5, 1.5 * math.Sqrt(3)},
			want:  1,
			want1: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := tt.grid
			got, got1 := grid.ToGrid(tt.args.c, tt.args.r)
			if math.Abs(got-tt.want) > epsilon {
				t.Errorf("HexGrid.ToGrid() got = %v, want %v", got, tt.want)
			}
			if math.Abs(got1-tt.want1) > epsilon {
				t.Errorf("HexGrid.ToGrid() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHexGrid_Vertices(t *testing.T) {
	t.SkipNow()
	type args struct {
		c int
		r int
	}
	tests := []struct {
		name      string
		grid      *HexGrid
		args      args
		wantVerts []mgl64.Vec2
	}{
		{
			name:      "0,0 - FlatTop",
			grid:      NewHexGrid(1, FlatTop),
			args:      args{0, 0},
			wantVerts: []mgl64.Vec2{},
		},
		{
			name:      "1,1 - FlatTop",
			grid:      NewHexGrid(1, FlatTop),
			args:      args{1, 1},
			wantVerts: []mgl64.Vec2{},
		},
		{
			name:      "0,0 - PointyTop",
			grid:      NewHexGrid(1, PointyTop),
			args:      args{0, 0},
			wantVerts: []mgl64.Vec2{},
		},
		{
			name:      "1,1 - PointyTop",
			grid:      NewHexGrid(1, PointyTop),
			args:      args{1, 1},
			wantVerts: []mgl64.Vec2{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := tt.grid
			if gotVerts := grid.Vertices(tt.args.c, tt.args.r); !reflect.DeepEqual(gotVerts, tt.wantVerts) {
				t.Errorf("HexGrid.Vertices() = %v, want %v", gotVerts, tt.wantVerts)
			}
		})
	}
}
