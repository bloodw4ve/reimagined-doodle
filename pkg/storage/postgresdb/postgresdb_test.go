package postgresdb

import (
	"fmt"
	"math/rand"
	"mod36a41/pkg/storage"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	pgPass := os.Getenv("pgPass")
	if pgPass == "" {
		t.Fatal(fmt.Errorf("No env data"))
	}
	pgConn := "postgres://postgres:" + pgPass + "@localhost:5432/news"

	s, err := New(pgConn)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
}

func TestStorage_StoreNews(t *testing.T) {
	pgPass := os.Getenv("pgPass")
	if pgPass == "" {
		t.Fatal(fmt.Errorf("No env data"))
	}
	pgConn := "postgres://postgres:" + pgPass + "@localhost:5432/news"

	s, err := New(pgConn)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	r := rand.New(rand.NewSource(time.Now().Unix()))
	randString := strconv.Itoa(r.Intn(1_000_000_000))
	news := []storage.Post{
		{
			Title:   "Test post title " + randString,
			Link:    "https://test.com/news/testpost/" + randString,
			PubTime: time.Now().Unix(),
			Content: "Test post content" + randString,
		},
	}
	if err := s.StoreNews(news); err != nil {
		t.Fatalf("Storage.StoreNews() error = %v", err)
	}
}

func TestStorage_LastNews(t *testing.T) {
	type args struct {
		n int
	}
	pgPass := os.Getenv("pgPass")
	if pgPass == "" {
		t.Fatal(fmt.Errorf("No env data"))
	}
	pgConn := "postgres://postgres:" + pgPass + "@localhost:5432/news"

	s, err := New(pgConn)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Case 1 (n = 0) - default number of news",
			args: args{
				n: 0,
			},
			want:    10,
			wantErr: false,
		}, {
			name: "Case 2 - (n = 15)",
			args: args{
				n: 15,
			},
			want:    15,
			wantErr: false,
		}, {
			name: "Case 3 - (n = 5)",
			args: args{
				n: 5,
			},
			want:    5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.LastNews(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.LastNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("Storage.LastNews() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
