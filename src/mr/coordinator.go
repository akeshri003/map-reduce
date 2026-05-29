package mr

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"
import "time"
import "sync"

type TaskStatus int
type TaskType int

const (
	Idle TaskStatus = iota
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
	mu sync.Mutex
	StateMaps map[string]StateMap
	nReduce int
}

// Your code here -- RPC handlers for the worker to call.
func (c * Coordinator) AssignTask(args *TaskArgs, reply *TaskReply) error {
	c.mu.Lock()
	allMapDone := c.checkAllMapComplete()
	allReduceDone := c.checkAllReduceComplete()
	taskAssigned := false

	if !allMapDone {
		for key, val := range c.StateMaps {
			if val.taskType == Map && val.state == Idle {
				val.workerId = args.Id
				val.state = InProgress
				val.startTime = time.Now()
				c.StateMaps[key] = val

				reply.Id = val.workerId
				reply.Taskname = val.filename
				reply.Tasktype = Map
				reply.NReduce = c.nReduce
				reply.IsDone = false
				taskAssigned = true
				break
			}
		}
		if !taskAssigned {
			reply.Id = args.Id
			reply.IsDone = false
			reply.Taskname = ""
			reply.Tasktype = Map
			reply.NReduce = c.nReduce
		}
	} else if !allReduceDone {
		for key, val := range c.StateMaps {
			if val.taskType == Reduce && val.state == Idle {
				val.workerId = args.Id
				val.state = InProgress
				val.startTime = time.Now()
				c.StateMaps[key] = val

				reply.Id = val.workerId
				reply.Taskname = val.filename
				reply.Tasktype = Reduce
				reply.NReduce = c.nReduce
				reply.IsDone = false
				taskAssigned = true
				break
			}
		}
		if !taskAssigned {
			reply.Id = args.Id
			reply.IsDone = false
			reply.Taskname = ""
			reply.Tasktype = Reduce
			reply.NReduce = c.nReduce
		}
	} else {
		reply.Id = args.Id
		reply.IsDone = true
	}

	c.mu.Unlock()
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

func (c *Coordinator) checkAllMapComplete() bool {
	
	for _, v := range c.StateMaps {
		if v.taskType == Map && v.state != Completed {
			return false
		}
	}

	return true
}

func (c *Coordinator) checkAllReduceComplete() bool {

	for _, v := range c.StateMaps {
		if v.taskType == Reduce && v.state != Completed {
			return false
		}
	}

	return true
}

func backgroundWorkerHealthChecker() {

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
	c.mu.Lock()

	ret = c.checkAllReduceComplete()
	
	defer c.mu.Unlock()
	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	c.nReduce = nReduce
	c.StateMaps = make(map[string]StateMap)

	for _, file := range files {
		task := StateMap{}
		task.state = Idle
		task.filename = file
		task.taskType = Map
		task.workerId = -1
		task.startTime = time.Time{}
		c.StateMaps[file] = task
	}
	

	c.server()
	return &c
}
