package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func Test_makeRequest(t *testing.T) {
	src := rand.New(rand.NewSource(time.Now().UnixMilli()))

	port := 49152 + src.Intn(65535-49152)
	listen := fmt.Sprintf("127.0.0.1:%v", port)

	url := PrepareUrl(listen)

	emptyMd5 := []byte("\xd4\x1d\x8c\xd9\x8f\x00\xb2\x04\xe9\x80\x09\x98\xec\xf8\x42\x7e")
	stringMd5 := []byte("\xb4\x5c\xff\xe0\x84\xdd\x3d\x20\xd9\x28\xbe\xe8\x5e\x7b\x0f\x21")

	tests := []struct {
		name    string
		timeout time.Duration
		url     string
		hash    []byte
		errorFn func(error) bool
	}{
		{
			"empty body",
			1 * time.Second,
			url + "/empty",
			emptyMd5,
			func(err error) bool {
				return err == nil
			},
		},
		{
			"not found",
			1 * time.Second,
			url + "/404",
			nil,
			func(err error) bool {
				return err.Error() == "404 Not Found"
			},

		},
		{
			"string",
			1 * time.Second,
			url + "/string",
			stringMd5,
			func(err error) bool {
				return err == nil
			},

		},
		{
			"redirect",
			1 * time.Second,
			url + "/redirect",
			stringMd5,
			func(err error) bool {
				return err == nil
			},
		},
		{
			"timeout",
			1 * time.Millisecond,
			url + "/timeout",
			nil,
			func(err error) bool {
				return errors.Is(err, context.DeadlineExceeded)
			},

		},
	}

	t.Run("prepare server", func(t *testing.T) {

		mux := &http.ServeMux{}

		mux.HandleFunc("/empty", func(_ http.ResponseWriter, _ *http.Request) {})

		mux.HandleFunc("/string", func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprintf(w, "string")
		})

		mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/string", http.StatusFound)
		})

		mux.HandleFunc("/timeout", func(w http.ResponseWriter, _ *http.Request) {
			<-time.After(1 * time.Second)
		})

		s := &http.Server{
			Addr:    listen,
			Handler: mux,
		}

		go func() {
			log.Fatal(s.ListenAndServe())
		}()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
				defer cancel()

				result := makeRequest(ctx, Job{Url: tt.url})
				if !reflect.DeepEqual(result.Hash, tt.hash) {
					t.Errorf("makeRequest().Hash = %v, want %v", result.Hash, tt.hash)
				}
				if !tt.errorFn(result.Error){
					t.Errorf("unexpected makeRequest().Error = %v", result.Error)
				}
			})
		}

		s.Close()
	})
}
