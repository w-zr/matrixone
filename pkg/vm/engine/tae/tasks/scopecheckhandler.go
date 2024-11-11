package tasks

import (
	"sync"

	"github.com/matrixorigin/matrixone/pkg/objectio"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/common"
)

type ScopeCheckTaskHandler struct {
	sync.RWMutex
	actives map[common.ID]struct{}

	handler *TypedTaskHandler
}

func (s *ScopeCheckTaskHandler) Start() {
	s.RLock()
	defer s.RUnlock()
	s.handler.Start()
}

func (s *ScopeCheckTaskHandler) Close() error {
	s.Lock()
	defer s.Unlock()
	return s.handler.Close()
}

func (s *ScopeCheckTaskHandler) Enqueue(task Task) {
	s.handler.Enqueue(task)
}

func NewScopeCheckTaskHandlerDispatcher() *ScopeCheckTaskHandler {
	return &ScopeCheckTaskHandler{
		actives: make(map[common.ID]struct{}),
		handler: NewTypedTaskHandler(),
	}
}

func (s *ScopeCheckTaskHandler) AddHandle(t TaskType, h TaskHandler) {
	s.Lock()
	defer s.Unlock()
	s.handler.AddHandle(t, h)
}

func (s *ScopeCheckTaskHandler) checkConflictLocked(scopes []common.ID) (err error) {
	for active := range s.actives {
		for _, scope := range scopes {
			if err = scopeConflictCheck(&active, &scope); err != nil {
				return
			}
		}
	}
	return
}

func (s *ScopeCheckTaskHandler) TryDispatch(task Task) (err error) {
	mscoped := task.(MScopedTask)
	scopes := mscoped.Scopes()
	if len(scopes) == 0 {
		s.Enqueue(task)
		return
	}
	s.Lock()
	if err = s.checkConflictLocked(scopes); err != nil {
		s.Unlock()
		return
	}
	for _, scope := range scopes {
		s.actives[scope] = struct{}{}
	}
	task.AddObserver(s)
	s.Unlock()
	s.Enqueue(task)
	return
}

func (s *ScopeCheckTaskHandler) OnExecDone(v any) {
	task := v.(MScopedTask)
	scopes := task.Scopes()
	s.Lock()
	for _, scope := range scopes {
		delete(s.actives, scope)
	}
	s.Unlock()
}

func scopeConflictCheck(oldScope, newScope *common.ID) (err error) {
	if oldScope.TableID != newScope.TableID {
		return
	}
	if !oldScope.SegmentID().Eq(*newScope.SegmentID()) &&
		!objectio.IsEmptySegid(oldScope.SegmentID()) &&
		!objectio.IsEmptySegid(newScope.SegmentID()) {
		return
	}
	if oldScope.BlockID != newScope.BlockID &&
		!objectio.IsEmptyBlkid(&oldScope.BlockID) &&
		!objectio.IsEmptyBlkid(&newScope.BlockID) {
		return
	}
	return ErrScheduleScopeConflict
}
