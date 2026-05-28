package mr

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"

type TaskStatus int
type TaskType int

const (
	Pending TaskStatus = iota
	InProgress
	Completed
)

const (
	Map TaskType = iota
	Reduce
)

type StateMap struct {
	workerId int
	state TaskStatus
	startTime time.Time
	filename string
	taskType TaskType
}

type Coordinator struct {
	// Your definitions here.
	StateMaps map[string]StateMap
	nReduce int
}

// Your code here -- RPC handlers for the worker to call.
func (c * Coordinator) AssignTask(args *TaskArgs, reply *TaskReply) error {

	return nil
}

// RPC handler for when a worker completes it's task
func (c *Coordinator) CompleteTask(args *TaskArgs, reply *TaskReply) error {


	return nil
}

// Updates the state of a map/reduce task.
func (c *Coordinator) updateState(
	tasktype TaskType, 
	key string, 
	targetState string,	
) error {

	return nil
}

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}


//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.


	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Your code here.


	c.server()
	return &c
}
