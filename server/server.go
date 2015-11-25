package main

import (
	// "fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cubar.com/model"
	"github.com/gorilla/mux"
)

type Server struct {
	httpAddr   string
	httpServer *http.Server
	router     *mux.Router
	listener   *net.TCPListener
}

func MainServer(host string, port string) (s *Server) {

	if host == "" {
		host = "127.0.0.1"
	}

	if port == "" {
		port = "25000"
	}

	addr := host + ":" + port
	s = &Server{
		httpAddr: addr,
		router:   NewRouter(),
	}

	s.httpServer = &http.Server{
		Addr:    s.HTTPAddr(),
		Handler: s.router,
	}

	return
}

func (s *Server) HTTPAddr() string { return s.httpAddr }

func (s *Server) Start() {

	ppid := os.Getpid()
	log.Printf("pid = %d Start http listener on %v", ppid, s.httpServer.Addr)

	laddr, err := net.ResolveTCPAddr("tcp", s.httpServer.Addr)
	if nil != err {
		log.Fatalln(err)
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {
		log.Fatalln(err)
	}
	s.listener = listener
	log.Println("listening on", listener.Addr())
	s.httpServer.Serve(listener)
	log.Println("Server Listener closed")
}

func (s *Server) Stop() {

	s.listener.Close()
	log.Println("Server Listener closing....")
	log.Println("Server stop")
}

func main() {

	// 开始model
	model.Start()
	defer model.Stop()

	server := MainServer("", "")
	go server.Start()
	defer server.Stop()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR2)

	for {
		sig := <-sigCh
		log.Println("signal ", sig)
		switch sig {
		case syscall.SIGTERM:
			log.Println("get stop singal")
			// server.Stop()
			return
		case syscall.SIGUSR2:
			log.Println("get restart singal")
			// server.Restart()
			return
		}
	}
}
