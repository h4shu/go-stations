package basic05_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

// 正常系
// 対象の API のみ Basic 認証がかかっているか、どうか。
func Test1(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()
	if srv == nil {
		t.Error("サーバーの作成に失敗しました。")
		return
	}

	testcases := map[string]struct {
		Path          string
		WantBasicAuth bool
	}{
		"BasicAuth /healthz": {
			Path:          "/healthz",
			WantBasicAuth: true,
		},
		"Not BasicAuth /todos": {
			Path:          "/todos",
			WantBasicAuth: false,
		},
		"Not BasicAuth /do-panic": {
			Path:          "/do-panic",
			WantBasicAuth: false,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			req := createRequest(t, http.MethodGet, srv.URL+tc.Path, nil)
			if req == nil {
				t.Error("リクエストの作成に失敗しました。")
				return
			}
			resp := sendRequest(t, req)
			if resp == nil {
				t.Error("リクエストの送信に失敗しました。")
				return
			}

			authHeader := resp.Header.Get("WWW-Authenticate")
			result := (authHeader == `Basic realm="SECRET AREA"`)
			if result != tc.WantBasicAuth {
				t.Errorf(`WWW-Authenticate: Basic realm="SECRET AREA" が含まれているか否か, got = %t, want = %t`, result, tc.WantBasicAuth)
				return
			}
		})
	}
}

// 正しい User ID, Password で Basic 認証をクリアしアクセスできるかどうか。
func Test2(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()
	if srv == nil {
		t.Error("サーバーの作成に失敗しました。")
		return
	}

	req := createRequest(t, http.MethodGet, srv.URL+"/healthz", nil)
	if req == nil {
		t.Error("リクエストの作成に失敗しました。")
		return
	}
	req.SetBasicAuth(os.Getenv("BASIC_AUTH_USER_ID"), os.Getenv("BASIC_AUTH_PASSWORD"))

	resp := sendRequest(t, req)
	if resp == nil {
		t.Error("リクエストの送信に失敗しました。")
		return
	}

	if resp.StatusCode == http.StatusUnauthorized {
		t.Error("Basic 認証に失敗しました。")
		return
	}
}

// 異常系
// 間違った User ID, Password を送信した場合、 Basic 認証が失敗しHTTP Status Code が 401 で返却されているかどうか。
func Test3(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()
	if srv == nil {
		t.Error("サーバーの作成に失敗しました。")
		return
	}

	testcases := map[string]struct {
		UserID   string
		Password string
	}{
		"Wrong User ID": {
			UserID:   os.Getenv("BASIC_AUTH_USER_ID") + "A",
			Password: os.Getenv("BASIC_AUTH_PASSWORD"),
		},
		"Wrong Password": {
			UserID:   os.Getenv("BASIC_AUTH_USER_ID"),
			Password: os.Getenv("BASIC_AUTH_PASSWORD") + "A",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			req := createRequest(t, http.MethodGet, srv.URL+"/healthz", nil)
			if req == nil {
				t.Error("リクエストの作成に失敗しました。")
				return
			}
			req.SetBasicAuth(tc.UserID, tc.Password)

			resp := sendRequest(t, req)
			if resp == nil {
				t.Error("リクエストの送信に失敗しました。")
				return
			}

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("期待していない HTTP status code です, got = %d, want = %d", resp.StatusCode, http.StatusUnauthorized)
				return
			}
		})
	}
}

// 空の User ID, Password を送信した場合、 Basic 認証が失敗し HTTP Status Code が 401 で返却されているかどうか。
func Test4(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()
	if srv == nil {
		t.Error("サーバーの作成に失敗しました。")
		return
	}

	req := createRequest(t, http.MethodGet, srv.URL+"/healthz", nil)
	if req == nil {
		t.Error("リクエストの作成に失敗しました。")
		return
	}
	req.SetBasicAuth("", "")

	resp := sendRequest(t, req)
	if resp == nil {
		t.Error("リクエストの送信に失敗しました。")
		return
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("期待していない HTTP status code です, got = %d, want = %d", resp.StatusCode, http.StatusUnauthorized)
		return
	}
}

// アクセス時に User ID, Password を送信しなかった場合、Basic 認証が失敗し HTTP Status Code が 401 で返却されているかどうか。
func Test5(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()
	if srv == nil {
		t.Error("サーバーの作成に失敗しました。")
		return
	}

	req := createRequest(t, http.MethodGet, srv.URL+"/healthz", nil)
	if req == nil {
		t.Error("リクエストの作成に失敗しました。")
		return
	}

	resp := sendRequest(t, req)
	if resp == nil {
		t.Error("リクエストの送信に失敗しました。")
		return
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("期待していない HTTP status code です, got = %d, want = %d", resp.StatusCode, http.StatusUnauthorized)
		return
	}
}

// 以下、ヘルパー関数
func createServer(t *testing.T) *httptest.Server {
	t.Helper()
	dbPath := "./temp_test.db"
	if err := os.Setenv("DB_PATH", dbPath); err != nil {
		t.Error("dbPathのセットに失敗しました。", err)
		return nil
	}

	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		t.Error("DBの作成に失敗しました。", err)
		return nil
	}

	t.Cleanup(func() {
		if err := todoDB.Close(); err != nil {
			t.Errorf("DBのクローズに失敗しました: %v", err)
			return
		}
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("テスト用のDBファイルの削除に失敗しました: %v", err)
			return
		}
	})

	r := router.NewRouter(todoDB)
	return httptest.NewServer(r)
}

func createRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Errorf("リクエストの作成に失敗しました: %v", err)
		return nil
	}
	return req
}

func sendRequest(t *testing.T, req *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("リクエストの送信に失敗しました: %v", err)
		return nil
	}
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("レスポンスのクローズに失敗しました: %v", err)
			return
		}
	})
	return resp
}
