package main

import (
	"fmt"
	"os"
	"os/exec"
//	"strconv"
	"strings"
	"syscall"
	rpc2 "net/rpc"
	"net/http"

	"github.com/checkpoint-restore/go-criu"
	"github.com/checkpoint-restore/go-criu/phaul"
	"github.com/checkpoint-restore/go-criu/rpc"
	"github.com/golang/protobuf/proto"
)

const(
    address="0.0.0.0"
	rpc_port = ":1234"
	pageserver_port = 5678
)

type testLocal struct {
	criu.NoNotify
	r *testRemote
}

type testRemote struct {
	srv *phaul.Server
}

type Srvapi struct {
	srv *phaul.Server
}

type Args struct {
}


func (srvapi *Srvapi) StartIter(arg *Args, reply *int) error {
	*reply = 2

	err := srvapi.srv.StartIter()
	if err != nil{
		return err
	}
	fmt.Println("done startiter")

	return nil
}

func (srvapi *Srvapi) StopIter(arg *Args, reply *int) error {
	*reply = 2

	err := srvapi.srv.StopIter()
	if err != nil{
		return err
	}
	fmt.Println("done stopiter")

	return nil
}

/* Dir where test will put dump images */
const imagesDir = "image"

func prepareImages() error {

	err := os.Mkdir(imagesDir, 0777)
	if err != nil {
		return err
	}
	/* Work dir for PhaulServer */
	err = os.Mkdir(imagesDir+"/remote", 0777)
	if err != nil {
		return err
	}

	/* Work dir for scp*/
	err = os.Mkdir("/tmp/livemig", 0777)
	if err != nil {
		return err
	}

	return nil
}

func mergeImages(dumpDir, lastPreDumpDir string) error {
	idir, err := os.Open(dumpDir)
	if err != nil {
		return err
	}

	defer idir.Close()

	imgs, err := idir.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, fname := range imgs {
		if !strings.HasSuffix(fname, ".img") {
			continue
		}

		fmt.Printf("\t%s -> %s/\n", fname, lastPreDumpDir)
		err = syscall.Link(dumpDir+"/"+fname, lastPreDumpDir+"/"+fname)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Srvapi) DoRestore(arg *Args, reply *int) error {
	lastSrvImagesDir := r.srv.LastImagesDir()
	/*
	 * In imagesDir we have images from dump, in the
	 * lastSrvImagesDir -- where server-side images
	 * (from page server, with pages and pagemaps) are.
	 * Need to put former into latter and restore from
	 * them.
	 */

	fmt.Printf("start restore\n")
	fmt.Println(lastSrvImagesDir)

	//err := mergeImages("/root/image", lastSrvImagesDir)
    scpCmd := "mv /root/image/*.img " + lastSrvImagesDir
	cmd := exec.Command("/bin/sh", "-c", scpCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println(string(output))
		return err
	}
	//if err != nil {
	//	return err
	//}


    ///XXX Change PID///
    ///Remove later
//    scpCmd = "/criu/crit/crit decode -i " + lastSrvImagesDir + "/pstree.img | sed -r \"s/\\\"pid\\\": [0-9]+/\\\"pid\\\": 1000/\" | /criu/crit/crit encode > pstree.img.new"
//	cmd = exec.Command("/bin/sh", "-c", scpCmd)
//	output, err = cmd.CombinedOutput()
//	if err != nil {
//		fmt.Println(cmd)
//		fmt.Println(string(output))
//		return err
//	}
//    scpCmd = "mv " + lastSrvImagesDir + "/pstree.img.new " + lastSrvImagesDir + "/pstree.img"
//	cmd = exec.Command("/bin/sh", "-c", scpCmd)
//	output, err = cmd.CombinedOutput()
//	if err != nil {
//		fmt.Println(cmd)
//		fmt.Println(string(output))
//		return err
//	}


	imgDir, err := os.Open(lastSrvImagesDir)
	if err != nil {
		return err
	}
	defer imgDir.Close()
	fmt.Println(imgDir.Fd())

	opts := rpc.CriuOpts{
		LogLevel:    proto.Int32(4),
		LogFile:     proto.String("restore.log"),
		ImagesDirFd: proto.Int32(int32(imgDir.Fd())),
		ShellJob:    proto.Bool(true),
	}

	cr := r.srv.GetCriu()
	fmt.Printf("Do restore\n")
	return cr.Restore(opts, nil)
}


func main() {
//	pid, _ := strconv.Atoi(os.Args[1])

err := prepareImages()
	if err != nil {
		fmt.Printf("Can't prepare dirs for images: %v\n", err)
		os.Exit(1)
		return
	}

	srv, err := phaul.MakePhaulServer(phaul.Config{
        Addr: address,
		Port: pageserver_port,
		Wdir:  imagesDir + "/remote"})
	if err != nil {
		fmt.Printf("Unable to run a server: %v", err)
		os.Exit(1)
		return
	}

	r := &Srvapi{srv}

	rpc2.Register(r)
	rpc2.HandleHTTP()

	err = http.ListenAndServe(rpc_port, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
