package main

import (
	"context"
	"log"
    "fmt"
    "io/ioutil"
	"net"
    "net/http"
    "bytes"
    "encoding/base64"
    "strconv"

    "math/rand"
    "os"
    "os/exec"
    "time"
    "encoding/json"

	"google.golang.org/grpc"
	pb "init/proto"

)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
    pb.UnimplementedInitDaemonServer
}

// gPRC methods
func (s *server) RunPhauls(ctx context.Context, in *pb.ServerInfo) (*pb.Result, error) {
    cmd := exec.Command("/root/phaul/s/run.sh")
    cmd.Start()

	return &pb.Result{Res: 1}, nil
}

func (s *server) RunPhaulc(ctx context.Context, in *pb.ClientInfo) (*pb.Result, error) {
    ip := in.GetAddr()

    servicename := in.GetTarget()

    err := registerdns(ip, servicename)

    f, _ := ioutil.ReadFile("pid.json")
    var data interface{}
    err = json.Unmarshal(f, &data)

    pid := data.(map[string]interface{})[in.GetTarget()]
    delete(data.(map[string]interface{}), in.GetTarget())

    new_data, err := json.Marshal(data)

    content := []byte(new_data)
    ioutil.WriteFile("pid.json", content, os.ModePerm)

    fmt.Println(strconv.Itoa(int(pid.(float64))))
    cmd := exec.Command("/root/phaul/c/run.sh", strconv.Itoa(int(pid.(float64))), in.GetAddr())
    out, err := cmd.Output()

    fmt.Println(string(out))
    fmt.Println(err)
    return &pb.Result{Res: 1}, nil
}

func serve() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterInitDaemonServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func GetIP() {
    netInterfaceAddresses, _ := net.InterfaceAddrs()
    for _, netInterfaceAddress := range netInterfaceAddresses {
        networkIp, ok := netInterfaceAddress.(*net.IPNet)
        if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
            ip := networkIp.IP.String()
            fmt.Println("Resolved Host IP: " + ip)
        }
    }
}

func registerdns(ip string , host string) error {
// TODO: It is good to use etcd client
// Cant install smooth because
// Warning: As etcd 3.5.0 was not yet released, the command above does not work. After first pre-release of 3.5.0 #12498, etcd can be referenced using:

    key := "/skydns/test/" + host + "/"
    val := "{\"host\":\"" + ip + "\",\"ttl\":60}"

    key64 := base64.StdEncoding.EncodeToString([]byte(key))
    val64:= base64.StdEncoding.EncodeToString([]byte(val))


    fmt.Println(key64)
    fmt.Println(val64)

    //jsonStr := `{"key":"L3NreWRucy90ZXN0L2hvZ2Uv","value":"eyJob3N0IjoiMS4xLjEuMiIsInR0bCI6NjB9"}`
    jsonStr := `{"key":"` + key64 + `","value":"` + val64 + `"}`

    fmt.Println(jsonStr)
    req, err := http.NewRequest(
            "POST",
            "http://etcd:2379/v3alpha/kv/put",
            bytes.NewBuffer([]byte(jsonStr)),
    )
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()


    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Request error:", err)
        return err
    }

    str_json := string(body)
    fmt.Println(str_json)

    return err

}

func exec_cmd(command string, args []string) error {
// set randome pid for conflict when process migration
// TODO: restore process in new pid namespace
    rand.Seed(time.Now().UnixNano())
    lastpid := 100 + rand.Intn(500)


    file, err := os.OpenFile("/proc/sys/kernel/ns_last_pid",os.O_WRONLY, 0666)
    //file, err := os.OpenFile("/home/vagrant/test.txt", os.O_WRONLY, 0666)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(fmt.Sprintf("%d\n", lastpid))
    if err != nil {
        return err
    }

// exe command
    cmd := exec.Command(command, args...)
    cmd.Start()

    a := map[string]interface{}{
        os.Getenv("SERVICE_NAME"): cmd.Process.Pid,
    }

    data, err := json.Marshal(a)

    content := []byte(data)
    ioutil.WriteFile("pid.json", content, os.ModePerm)


    return err
}


func main(){

    netInterfaceAddresses, _ := net.InterfaceAddrs()
    ip := ""
    for _, netInterfaceAddress := range netInterfaceAddresses {
        networkIp, ok := netInterfaceAddress.(*net.IPNet)
        if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
            ip = networkIp.IP.String()
            fmt.Println("Resolved Host IP: " + ip)
        }
    }

    servicename := os.Getenv("SERVICE_NAME")

    err := registerdns(ip, servicename)

    if err != nil {
		log.Fatalf("fail registerdns: %v", err)
        return
    }
    err = exec_cmd(os.Args[1], os.Args[2:])
    if err != nil {
		log.Fatalf("failed execcmd: %v", err)
        return
    }

    fmt.Println("start")
    serve()
}
