import argparse
import yaml
import urllib3
import time
import subprocess

import grpc
import initdaemon_pb2
import initdaemon_pb2_grpc

def wait_for_running(script):
    try:
        res = subprocess.check_output(script)
    except:
        print(res)
        sys.exit()



def run_end2end_test(script):
    try:
        res = subprocess.check_output(script)
    except:
        print(res)
        sys.exit()
     
def tear_down(): 
    print("AAA")

def construct(serverIP, clientIP, service):
    sip = serverIP + ":50051"
    cip = clientIP + ":50051"
    print(sip)

    with grpc.insecure_channel(sip) as channel:
        stub = initdaemon_pb2_grpc.InitDaemonStub(channel)
        stub.RunPhauls(initdaemon_pb2.ServerInfo(target=service , servicename = service))

    with grpc.insecure_channel(cip) as channel:
        stub = initdaemon_pb2_grpc.InitDaemonStub(channel)
        stub.RunPhaulc(initdaemon_pb2.ClientInfo(target=service, servicename = service, addr = serverIP))


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--composefile")
    parser.add_argument("--conffile")
    parser.add_argument("--check")
    parser.add_argument("--test")
    args = parser.parse_args()

    f = open(args.composefile, "r")
    composefile = yaml.load(f, Loader=yaml.SafeLoader)
    hosts = composefile['services'].keys() 
    f.close()


    f = open(args.conffile, "r")

    while True:
        line = f.readline()
        if line:
            conf = line.split(',')
            
            if len(conf) == 3:
                sip = conf[0]
                cip = conf[1]
                servicename = conf[2].replace('\n' , '')
                wait_for_running(args.check)
                construct(sip, cip, servicename)

            #wait_for_running(args.check)
            run_end2end_test(args.test)

        else:
            break
    f.close()
