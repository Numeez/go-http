It is a custom implementation of  http protocol over tcp

# 🌐 go-http

A lightweight HTTP server written **from scratch in Go** — without relying on Go's built-in `net/http` package.  
This project demonstrates a deep understanding of the **HTTP protocol**, **TCP networking**, and **request/response lifecycles** by manually implementing core server logic.  

---

## 🚀 Features

- 🔌 **Custom TCP Listener**  
  Built the foundation using Go's `net` package to accept raw TCP connections.  

- 📜 **Manual HTTP Parsing**  
  - Extracted **Request Line** (method, path, protocol)  
  - Parsed **Headers** by hand  
  - Supported **Request Body** parsing  

- 🖋 **Response Writing**  
  - Implemented logic to construct valid HTTP responses  
  - Supported sending text and binary content  

- ⚙️ **Custom Routing & Handlers**  
  - Users can register functions based on routes  
  - Each handler processes the request and writes back the response  

- 🎬 **Static File Serving**  
  - Served a **video file** over HTTP  
  - Showcased capability to handle large binary payloads  

- ✅ **Tests Included**  
  - Wrote tests for request parsing, response writing, and routing logic  

---

## 🛠️ Tech Stack

- **Language**: Go  
- **Core Packages Used**:  
  - `net` (for TCP sockets)  
  - `bufio` (for parsing incoming data)  
  - `strings`, `bytes` (for request parsing)  

---

## ⚡ Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/your-username/go-http.git
cd go-http
```

### 2. Running the server

```bash
go run cmd/httpserver/main.go
```
