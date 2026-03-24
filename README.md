```
# 🔍 LogLens: Unified System & Docker Search

**LogLens** is a single-binary, zero-dependency web utility designed for rapid troubleshooting on Linux servers. It scans system logs and Docker container logs simultaneously, presenting them in a unified, "newest-first" web interface.

---

## 🚀 10-Second Quick Start
If you are troubleshooting a live issue, run these commands to get the UI running immediately:

```bash
# 1. Download the latest binary (Replace YOUR_USERNAME with your GitHub handle)
wget [https://github.com/YOUR_USERNAME/loglens/releases/latest/download/loglens](https://github.com/YOUR_USERNAME/loglens/releases/latest/download/loglens)

# 2. Make it executable
chmod +x loglens

# 3. Run it (sudo required for /var/log and Docker socket access)
sudo ./loglens

```

**Access the UI at:** `http://<server-ip>:8080`

* * * * *

✨ Features
----------

-   **Unified Search**: Greps through `/var/log/*` and `docker logs` in one go.

-   **Newest First**: Automatically reverses output using `tac` so the most recent events are at the top of your screen.

-   **Secondary Filtering**: Perform a broad search (e.g., "error"), then narrow down results instantly with a secondary "grep" filter (e.g., "nginx").

-   **Live Monitor**: A 60-second auto-refresh mode (every 5s) for watching logs in real-time during a "fire."

-   **Zero Dependencies**: A single Go binary. No Python, no Pip, no Node, and no environment errors.

* * * * *

🛠 Functions
------------

| **Feature** | **Description** |
| --- | --- |
| **Primary Search** | Case-insensitive recursive grep across system and container logs. |
| **Grep Filter** | A secondary "pipe" that narrows down the primary search results. |
| **Live Toggle** | Uses a meta-refresh to update the page every 5 seconds for 1 minute. |
| **IP Detection** | Automatically identifies and prints the LAN IP on startup for easy access. |

* * * * *

🔒 Security Note
----------------

This tool is intended for **temporary troubleshooting sessions only**. It does not have built-in authentication or encryption.

-   **DO NOT** leave this running on public-facing IPs.

-   **DO** use a firewall: `sudo ufw allow 8080/tcp`.

-   **DO** kill the process and `sudo ufw delete allow 8080` when finished.

* * * * *

*Built by [Built Networks, LLC](https://builtnetworks.com/)*
