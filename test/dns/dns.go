package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
)

var records = map[string]string{
	"kafka.": "172.17.0.4",
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			ip := records[q.Name]
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func handle(w dns.ResponseWriter, r *dns.Msg) {
	var m dns.Msg
	m.SetReply(r)
	m.Compress = false
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(&m)
	}
	w.WriteMsg(&m)
}

func main() {
	// attach request handler func
	dns.HandleFunc(".", handle)

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	// start server
	server := &dns.Server{Listener: listener}
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer server.Shutdown()
		err = server.ListenAndServe()
		wg.Done()
	}()
	time.Sleep(time.Second)
	fmt.Println(server.Listener.Addr().String())
	wg.Wait()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
