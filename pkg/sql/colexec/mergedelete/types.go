// Copyright 2022 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mergedelete

import (
	"github.com/matrixorigin/matrixone/pkg/common/reuse"
	"github.com/matrixorigin/matrixone/pkg/container/batch"
	"github.com/matrixorigin/matrixone/pkg/pb/plan"
	"github.com/matrixorigin/matrixone/pkg/vm"
	"github.com/matrixorigin/matrixone/pkg/vm/engine"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

var _ vm.Operator = new(MergeDelete)

type container struct {
	delSource    engine.Relation
	affectedRows uint64
	bat          *batch.Batch
}
type MergeDelete struct {
	ctr             container
	AddAffectedRows bool
	Ref             *plan.ObjectRef
	Engine          engine.Engine

	vm.OperatorBase
}

func (mergeDelete *MergeDelete) GetOperatorBase() *vm.OperatorBase {
	return &mergeDelete.OperatorBase
}

func init() {
	reuse.CreatePool[MergeDelete](
		func() *MergeDelete {
			return &MergeDelete{}
		},
		func(a *MergeDelete) {
			*a = MergeDelete{}
		},
		reuse.DefaultOptions[MergeDelete]().
			WithEnableChecker(),
	)
}

func (mergeDelete MergeDelete) TypeName() string {
	return opName
}

func NewArgument() *MergeDelete {
	return reuse.Alloc[MergeDelete](nil)
}

func (mergeDelete *MergeDelete) WithObjectRef(ref *plan.ObjectRef) *MergeDelete {
	mergeDelete.Ref = ref
	return mergeDelete
}

func (mergeDelete *MergeDelete) WithEngine(eng engine.Engine) *MergeDelete {
	mergeDelete.Engine = eng
	return mergeDelete
}

func (mergeDelete *MergeDelete) WithAddAffectedRows(addAffectedRows bool) *MergeDelete {
	mergeDelete.AddAffectedRows = addAffectedRows
	return mergeDelete
}

func (mergeDelete *MergeDelete) Release() {
	if mergeDelete != nil {
		reuse.Free[MergeDelete](mergeDelete, nil)
	}
}

func (mergeDelete *MergeDelete) Reset(proc *process.Process, pipelineFailed bool, err error) {
	//can not reset affectRows because MO need get affectRows after reset
	if mergeDelete.ctr.bat != nil {
		mergeDelete.ctr.bat.CleanOnlyData()
	}
}

func (mergeDelete *MergeDelete) Free(proc *process.Process, pipelineFailed bool, err error) {
	if mergeDelete.ctr.bat != nil {
		mergeDelete.ctr.bat.Clean(proc.Mp())
		mergeDelete.ctr.bat = nil
	}
}

func (mergeDelete *MergeDelete) ExecProjection(proc *process.Process, input *batch.Batch) (*batch.Batch, error) {
	return input, nil
}

func (mergeDelete *MergeDelete) GetAffectedRows() uint64 {
	return mergeDelete.ctr.affectedRows
}
