package tasks

type TypedTaskHandler struct {
	handlers map[TaskType]TaskHandler
}

func NewTypedTaskHandler() *TypedTaskHandler {
	return &TypedTaskHandler{
		handlers: make(map[TaskType]TaskHandler),
	}
}

func (d *TypedTaskHandler) Enqueue(task Task) {
	handler, ok := d.handlers[task.Type()]
	if !ok {
		panic(ErrDispatchWrongTask)
	}
	handler.Enqueue(task)
}

func (d *TypedTaskHandler) AddHandle(t TaskType, h TaskHandler) {
	d.handlers[t] = h
}

func (d *TypedTaskHandler) Start() {
	for _, h := range d.handlers {
		h.Start()
	}
}

func (d *TypedTaskHandler) Close() error {
	for _, h := range d.handlers {
		h.Close()
	}
	return nil
}
