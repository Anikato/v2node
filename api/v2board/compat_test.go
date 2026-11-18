package panel

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/wyx2685/v2node/conf"
)

func TestV2BoardNodeContract(t *testing.T) {
	const nodeETag = "\"node-current\""
	const usersETag = "\"users-current\""
	account := UserInfo{
		Id: 17, Uuid: "e0121e0a-a53b-4a91-8c92-b32b4d0d3120", SpeedLimit: 40, DeviceLimit: 2,
	}
	service := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("token") != "test-api-key" || query.Get("node_id") != "8" || query.Get("node_type") != "v2node" {
			t.Errorf("incorrect V2Board query parameters: %v", query)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("V2Board node should authenticate with the UniProxy query fields")
		}

		key := r.Method + " " + r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		switch key {
		case "GET /api/v2/server/config":
			if r.Header.Get("If-None-Match") == nodeETag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			w.Header().Set("ETag", nodeETag)
			_, _ = io.WriteString(w, "{\"protocol\":\"vless\",\"listen_ip\":\"127.0.0.1\",\"server_port\":8443,\"network\":\"tcp\","+
				"\"trusted_x_forwarded_for\":[\"X-Forwarded-For\"],\"base_config\":{\"push_interval\":30,\"pull_interval\":60}}")
		case "GET /api/v1/server/UniProxy/user":
			if r.Header.Get("If-None-Match") == usersETag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			if r.Header.Get("X-Response-Format") != "msgpack" {
				t.Error("missing UniProxy format negotiation")
			}
			w.Header().Set("ETag", usersETag)
			w.Header().Set("Content-Type", "application/x-msgpack")
			if err := msgpack.NewEncoder(w).Encode(UserListBody{Users: []UserInfo{account}}); err != nil {
				t.Error(err)
			}
		case "GET /api/v1/server/UniProxy/alivelist":
			_, _ = io.WriteString(w, "{\"alive\":{\"17\":1}}")
		case "POST /api/v1/server/UniProxy/push":
			var report map[int][]int64
			if err := json.NewDecoder(r.Body).Decode(&report); err != nil || !reflect.DeepEqual(report, map[int][]int64{17: {300, 1200}}) {
				t.Errorf("unexpected UniProxy traffic report: %v, %v", report, err)
			}
			_, _ = io.WriteString(w, "{\"data\":true}")
		case "POST /api/v1/server/UniProxy/alive":
			var report map[int][]string
			if err := json.NewDecoder(r.Body).Decode(&report); err != nil || !reflect.DeepEqual(report, map[int][]string{17: {"192.0.2.10"}}) {
				t.Errorf("unexpected UniProxy presence: %v, %v", report, err)
			}
			_, _ = io.WriteString(w, "{\"data\":true}")
		default:
			t.Errorf("unexpected request: %s", key)
			http.NotFound(w, r)
		}
	}))
	defer service.Close()

	const nodeID = 8
	retry := 0
	client, err := New(&conf.NodeConfig{APIHost: service.URL, NodeID: nodeID, Key: "test-api-key", RetryCount: &retry})
	if err != nil {
		t.Fatal(err)
	}
	ctx := t.Context()
	node, err := client.GetNodeInfo(ctx)
	if err != nil || node == nil {
		t.Fatalf("GetNodeInfo returned %v, %v", node, err)
	}
	if node.Id != nodeID || node.Type != "vless" || node.Common.ServerPort != 8443 ||
		node.Common.ListenIP != "127.0.0.1" || !reflect.DeepEqual(node.Common.TrustedXForwardedFor, []string{"X-Forwarded-For"}) ||
		node.PushInterval != 30*time.Second || node.PullInterval != time.Minute {
		t.Fatalf("unexpected node settings: %#v, %#v", node, node.Common)
	}
	if changed, err := client.GetNodeInfo(ctx); changed != nil || err != nil {
		t.Fatalf("unchanged config returned %v, %v", changed, err)
	}
	users, err := client.GetUserList(ctx)
	if err != nil || !reflect.DeepEqual(users, []UserInfo{account}) {
		t.Fatalf("GetUserList returned %v, %v", users, err)
	}
	if changed, err := client.GetUserList(ctx); changed != nil || err != nil {
		t.Fatalf("unchanged user list returned %v, %v", changed, err)
	}
	counts, err := client.GetUserAlive(ctx)
	if err != nil || !reflect.DeepEqual(counts, map[int]int{account.Id: 1}) {
		t.Fatalf("GetUserAlive returned %v, %v", counts, err)
	}
	if err := client.ReportUserTraffic(ctx, []UserTraffic{{UID: account.Id, Upload: 300, Download: 1200}}); err != nil {
		t.Fatal(err)
	}
	online := map[int][]string{account.Id: {"192.0.2.10"}}
	if err := client.ReportNodeOnlineUsers(ctx, &online); err != nil {
		t.Fatal(err)
	}
}
