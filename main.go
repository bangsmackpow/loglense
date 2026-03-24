package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>LogLens - Unified Search</title>
    {{if .IsLive}}<meta http-equiv="refresh" content="5">{{end}}
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; padding: 30px; background: #1e1e1e; color: #d4d4d4; }
        .header { border-bottom: 2px solid #333; padding-bottom: 10px; margin-bottom: 20px; }
        input, button { padding: 12px; border-radius: 4px; border: 1px solid #333; background: #2d2d2d; color: white; }
        input[type="text"] { width: 250px; }
        button { background: #0078d4; cursor: pointer; border: none; font-weight: bold; }
        button:hover { background: #005a9e; }
        pre { background: #000; color: #adff2f; padding: 20px; border-radius: 8px; border: 1px solid #444; overflow-x: auto; line-height: 1.5; font-size: 13px; }
        .meta { color: #888; font-size: 0.9em; margin-top: 5px; }
        .section-header { color: #569cd6; font-weight: bold; margin-top: 20px; text-transform: uppercase; letter-spacing: 1px; }
    </style>
</head>
<body>
    <div class="header">
        <h2>🔍 LogLens <span style="font-size: 0.5em; color: #888;">Built Networks Edition</span></h2>
    </div>

    <form method="GET">
        <input type="text" name="q" placeholder="Primary Search..." value="{{.Query}}">
        <input type="text" name="f" placeholder="Filter Results (Grep)..." value="{{.Filter}}">
        <label style="margin-left: 10px;">
            <input type="checkbox" name="live" {{if .IsLive}}checked{{end}}> Live (60s)
        </label>
        {{if .IsLive}}<input type="hidden" name="start" value="{{.StartTime}}">{{end}}
        <button type="submit">SEARCH LOGS</button>
    </form>

    <div class="meta">
        {{if .IsLive}}🔴 Monitoring active... Auto-stops in {{.TimeRemaining}}s{{else}}View: Static{{end}}
    </div>

    <hr style="border: 0; border-top: 1px solid #333; margin: 20px 0;">

    {{if .Results}}
        <pre>{{.Results}}</pre>
    {{else if .Query}}
        <p>No matches found for "{{.Query}}".</p>
    {{else}}
        <p style="color: #666;">Enter a search term (e.g., "error", "nginx", "docker") to begin.</p>
    {{end}}
</body>
</html>`

type PageData struct {
	Query         string
	Filter        string
	Results       string
	IsLive        bool
	StartTime     int64
	TimeRemaining int64
}

func getLocalIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "127.0.0.1"
}

func searchLogs(query, filter string) string {
	if query == "" { return "" }
	var out strings.Builder

	// Filesystem Logic: Grep -> Filter -> Tail (Last 100) -> Tac (Reverse)
	fsBase := fmt.Sprintf("grep -r -i '%s' /var/log 2>/dev/null", query)
	if filter != "" {
		fsBase += fmt.Sprintf(" | grep -i '%s'", filter)
	}
	fsCmd := fmt.Sprintf("%s | tail -n 100 | tac", fsBase)
	
	fsOutput, _ := exec.Command("bash", "-c", fsCmd).CombinedOutput()
	if len(fsOutput) > 0 {
		out.WriteString("=== SYSTEM LOGS (Newest First) ===\n")
		out.Write(fsOutput)
		out.WriteString("\n\n")
	}

	// Docker Logic: Logs -> Grep -> Filter -> Tail -> Tac
	dockBase := fmt.Sprintf("docker ps -q | xargs -I {} docker logs --tail 500 {} 2>&1 | grep -i '%s'", query)
	if filter != "" {
		dockBase += fmt.Sprintf(" | grep -i '%s'", filter)
	}
	dockCmd := fmt.Sprintf("%s | tail -n 100 | tac", dockBase)
	
	dockOutput, _ := exec.Command("bash", "-c", dockCmd).CombinedOutput()
	if len(dockOutput) > 0 {
		out.WriteString("=== DOCKER LOGS (Newest First) ===\n")
		out.Write(dockOutput)
	}

	return out.String()
}

func main() {
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		filter := r.URL.Query().Get("f")
		liveParam := r.URL.Query().Get("live") == "on"
		startParam, _ := strconv.ParseInt(r.URL.Query().Get("start"), 10, 64)

		now := time.Now().Unix()
		if liveParam && startParam == 0 {
			startParam = now
		}

		remaining := 60 - (now - startParam)
		isLive := liveParam && remaining > 0

		results := searchLogs(query, filter)
		tmpl.Execute(w, PageData{
			Query:         query,
			Filter:        filter,
			Results:       results,
			IsLive:        isLive,
			StartTime:     startParam,
			TimeRemaining: remaining,
		})
	})

	lanIP := getLocalIP()
	
	// Cleaned up startup message to avoid "newline in string" error
	fmt.Println("\n------------------------------------")
	fmt.Println("🚀 LogLens Unified Search Active")
	fmt.Printf("Local:   http://localhost:8080\n")
	fmt.Printf("Network: http://%s:8080\n", lanIP)
	fmt.Println("------------------------------------\n")
	
	http.ListenAndServe("0.0.0.0:8080", nil)
}
