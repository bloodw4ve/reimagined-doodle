package rss

import "testing"

func TestParseFeed(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Case 1 - habr.com golang RSS-feed",
			args: args{
				url: "https://habr.com/ru/rss/hub/go/all/?fl=ru",
			},
			wantErr: false,
		}, {
			name: "Case 2 - habr.com best daily RSS-feed",
			args: args{
				url: "https://habr.com/ru/rss/best/daily/?fl=ru",
			},
			wantErr: false,
		}, {
			name: "Case 3 - Golang Weekly RSS-feed",
			args: args{
				url: "https://cprss.s3.amazonaws.com/golangweekly.com.xml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFeed(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) == 0 {
				t.Errorf("Data was unparsed or RSS-feed is empty")
			}
			t.Logf("Aquired %v posts\n%v", len(got), got)
		})
	}
}
