package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"yoyoSystem"
)


const LOG_DIR string = "/tmp/fdir"
const LOG_FILE string = "yoyo.log"
const READ_DIR string = "/tmp/fdir"
const READ_FILE string = "yoyo.read"
const YOYO_VER string = "0.7"

func main() {
	preProcess()
	//yoyoHttp()
	yoyoGraceHttp()
}

func preProcess(){
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func yoyoHttp() {
	fmt.Println("+++ START YOYO HTTP +++")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Version : "+YOYO_VER)
		printPid(&w)
		printHostname(&w)
		printRunCore(&w)
		printEnv(&w,"AGE")
		printNetworkInfo(&w)
		printReadFileStatus(&w)
		printWriteLogStatus(&w)
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func yoyoGraceHttp(){
	fmt.Println("+++ START YOYO GRACE HTTP +++")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Version : "+YOYO_VER)
		printPid(&w)
		printHostname(&w)
		printRunCore(&w)
		printEnv(&w,"AGE")
		printNetworkInfo(&w)
		printReadFileStatus(&w)
		printWriteLogStatus(&w)
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()
	log.Println("server started")

	stopC := make(chan os.Signal,1)
	signal.Notify(stopC,syscall.SIGTERM,syscall.SIGINT,syscall.SIGKILL)

	sgn := <- stopC
	fmt.Println("signal: ", sgn)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)


	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown Error : ",err.Error())
	}
	// 5초의 타임아웃으로 ctx.Done()을 캐치합니다.
	select {
	case <-ctx.Done():
		fmt.Println("timeout of 5 seconds.")
	}

	fmt.Println("Server exiting")
}


func printHostname(w *http.ResponseWriter) {
	pHostName,_err := yoyoSystem.GetHostname()
	if _err != nil {
		pHostName = _err.Error()
	}
	fmt.Fprintln(*w, "HostName : "+pHostName)
}

func printRunCore(w *http.ResponseWriter) {
	fmt.Fprintln(*w, "Active Core :",runtime.GOMAXPROCS(0))

}

func printPid(w *http.ResponseWriter) {
	pid := yoyoSystem.GetPid()
	fmt.Fprintln(*w, "PID : "+ fmt.Sprintf("%d",pid))
}

func printEnv(w *http.ResponseWriter, enVar string){
	fmt.Fprintln(*w, "")
	fmt.Fprintln(*w, "========== ENV READ TEST ==========")
	fmt.Fprintln(*w, enVar+" : "+os.Getenv("AGE"))
}

func printNetworkInfo(w *http.ResponseWriter) {
	fmt.Fprintln(*w, "")
	fmt.Fprintln(*w, "========== NETWORK INFORMATION ==========")
	pNetInfo,_err := yoyoSystem.GetNetworkInfo()
	if _err != nil {
		fmt.Fprintln(*w,"ERR")
	}
	for _, eNetInfo := range pNetInfo {
		fmt.Fprintln(*w, eNetInfo.InfName+" - "+eNetInfo.InfIp)
	}
}

func printReadFileStatus(w *http.ResponseWriter) {
	fmt.Fprintln(*w, "")
	fmt.Fprintln(*w, "========== FILE READ TEST ==========")

	path := READ_DIR + "/" + READ_FILE
	flagStr := ""

	_err := yoyoSystem.ReadFileTest(path)

	if _err != nil {
		flagStr = "[FAIL]"
	}else{
		flagStr = "[SUCCESS]"
	}

	fmt.Fprintln(*w, flagStr +" "+ path + " ("+ _err.Error()+")")
}

func printWriteLogStatus(w *http.ResponseWriter) {

	fmt.Fprintln(*w, "")
	fmt.Fprintln(*w, "========== LOG FILE WRITE ==========")

	path := LOG_DIR + "/" + LOG_FILE
	flagStr := ""

	_err := putYoyoLog(path,"printWriteLogStatus")

	// merge Format
	if _err != nil {
		flagStr = "[FAIL]"
	}else{
		flagStr = "[SUCCESS]"
	}
	fmt.Fprintln(*w, flagStr +" "+ path + " (" + _err.Error()+")")
}


func putYoyoLog( path, contents string) (rErr error){
	/*  Make Log Format */
	// time
	_location, _ := time.LoadLocation("Asia/Seoul")
	_now := time.Now()
	_t := _now.In(_location)
	Kst := _t.String()

	// hostname
	hostName, _ := os.Hostname()

	// merge Format
	logFormat := fmt.Sprintf(Kst + " | " + hostName + " | "+ contents+ "\n")
	rErr = yoyoSystem.WriteFile(path,logFormat)

	return
}

