package handlers

import "testing"

func Test_calculateAverageClosingPrice(t *testing.T) {
	type args struct {
		stockTickerEntries []StockTickerEntry
		dayCount           int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "No Values",
			args: args{
				stockTickerEntries: []StockTickerEntry{},
				dayCount:           100,
			},
			want: 0,
		},
		{
			name: "Correct Average Value",
			args: args{
				stockTickerEntries: []StockTickerEntry{
					{
						Close: 1,
					},
					{
						Close: 2,
					},
				},
				dayCount: 2,
			},
			want: 1.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateAverageClosingPrice(tt.args.stockTickerEntries, tt.args.dayCount); got != tt.want {
				t.Errorf("calculateAverageClosingPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
