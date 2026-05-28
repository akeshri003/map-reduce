package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.

type TaskArgs struct {
	Id int
}

type TaskReply struct {
	Id int
	Taskname string  // filename - works for both map and reduce.
	Tasktype TaskType	// task type
	NReduce int	
	IsDone bool // if False => wait, else mapreduce complete.
}

type CompleteTaskArgs struct {
	Id int
}

type CompleteTaskReply struct {
	Ack bool
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
