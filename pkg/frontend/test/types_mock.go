// Code generated by MockGen. DO NOT EDIT.
// Source: ../types.go

// Package mock_frontend is a generated GoMock package.
package mock_frontend

import (
	"bytes"
	context "context"
	reflect "reflect"

	"github.com/matrixorigin/matrixone/pkg/perfcounter"
	"github.com/matrixorigin/matrixone/pkg/sql/models"
	"github.com/matrixorigin/matrixone/pkg/util/trace/impl/motrace/statistic"

	gomock "github.com/golang/mock/gomock"
	batch "github.com/matrixorigin/matrixone/pkg/container/batch"
	types "github.com/matrixorigin/matrixone/pkg/container/types"
	plan "github.com/matrixorigin/matrixone/pkg/pb/plan"
	tree "github.com/matrixorigin/matrixone/pkg/sql/parsers/tree"
	util "github.com/matrixorigin/matrixone/pkg/util"
	process "github.com/matrixorigin/matrixone/pkg/vm/process"
)

// MockComputationRunner is a mock of ComputationRunner interface.
type MockComputationRunner struct {
	ctrl     *gomock.Controller
	recorder *MockComputationRunnerMockRecorder
}

// MockComputationRunnerMockRecorder is the mock recorder for MockComputationRunner.
type MockComputationRunnerMockRecorder struct {
	mock *MockComputationRunner
}

// NewMockComputationRunner creates a new mock instance.
func NewMockComputationRunner(ctrl *gomock.Controller) *MockComputationRunner {
	mock := &MockComputationRunner{ctrl: ctrl}
	mock.recorder = &MockComputationRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComputationRunner) EXPECT() *MockComputationRunnerMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockComputationRunner) Run(ts uint64) (*util.RunResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ts)
	ret0, _ := ret[0].(*util.RunResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockComputationRunnerMockRecorder) Run(ts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockComputationRunner)(nil).Run), ts)
}

// MockComputationWrapper is a mock of ComputationWrapper interface.
type MockComputationWrapper struct {
	ctrl     *gomock.Controller
	recorder *MockComputationWrapperMockRecorder
}

// MockComputationWrapperMockRecorder is the mock recorder for MockComputationWrapper.
type MockComputationWrapperMockRecorder struct {
	mock *MockComputationWrapper
}

// NewMockComputationWrapper creates a new mock instance.
func NewMockComputationWrapper(ctrl *gomock.Controller) *MockComputationWrapper {
	mock := &MockComputationWrapper{ctrl: ctrl}
	mock.recorder = &MockComputationWrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComputationWrapper) EXPECT() *MockComputationWrapperMockRecorder {
	return m.recorder
}

// BinaryExecute mocks base method.
func (m *MockComputationWrapper) BinaryExecute() (bool, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BinaryExecute")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// BinaryExecute indicates an expected call of BinaryExecute.
func (mr *MockComputationWrapperMockRecorder) BinaryExecute() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BinaryExecute", reflect.TypeOf((*MockComputationWrapper)(nil).BinaryExecute))
}

// Clear mocks base method.
func (m *MockComputationWrapper) Clear() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Clear")
}

// Clear indicates an expected call of Clear.
func (mr *MockComputationWrapperMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockComputationWrapper)(nil).Clear))
}

// Compile mocks base method.
func (m *MockComputationWrapper) Compile(any any, fill func(*batch.Batch, *perfcounter.CounterSet) error) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compile", any, fill)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Compile indicates an expected call of Compile.
func (mr *MockComputationWrapperMockRecorder) Compile(any, fill interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compile", reflect.TypeOf((*MockComputationWrapper)(nil).Compile), any, fill)
}

// Free mocks base method.
func (m *MockComputationWrapper) Free() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Free")
}

// Free indicates an expected call of Free.
func (mr *MockComputationWrapperMockRecorder) Free() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Free", reflect.TypeOf((*MockComputationWrapper)(nil).Free))
}

// GetAst mocks base method.
func (m *MockComputationWrapper) GetAst() tree.Statement {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAst")
	ret0, _ := ret[0].(tree.Statement)
	return ret0
}

// GetAst indicates an expected call of GetAst.
func (mr *MockComputationWrapperMockRecorder) GetAst() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAst", reflect.TypeOf((*MockComputationWrapper)(nil).GetAst))
}

// GetColumns mocks base method.
func (m *MockComputationWrapper) GetColumns(ctx context.Context) ([]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetColumns", ctx)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetColumns indicates an expected call of GetColumns.
func (mr *MockComputationWrapperMockRecorder) GetColumns(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetColumns", reflect.TypeOf((*MockComputationWrapper)(nil).GetColumns), ctx)
}

// GetLoadTag mocks base method.
func (m *MockComputationWrapper) GetLoadTag() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoadTag")
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetLoadTag indicates an expected call of GetLoadTag.
func (mr *MockComputationWrapperMockRecorder) GetLoadTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoadTag", reflect.TypeOf((*MockComputationWrapper)(nil).GetLoadTag))
}

// GetProcess mocks base method.
func (m *MockComputationWrapper) GetProcess() *process.Process {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProcess")
	ret0, _ := ret[0].(*process.Process)
	return ret0
}

// GetProcess indicates an expected call of GetProcess.
func (mr *MockComputationWrapperMockRecorder) GetProcess() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProcess", reflect.TypeOf((*MockComputationWrapper)(nil).GetProcess))
}

// GetServerStatus mocks base method.
func (m *MockComputationWrapper) GetServerStatus() uint16 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServerStatus")
	ret0, _ := ret[0].(uint16)
	return ret0
}

// GetServerStatus indicates an expected call of GetServerStatus.
func (mr *MockComputationWrapperMockRecorder) GetServerStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServerStatus", reflect.TypeOf((*MockComputationWrapper)(nil).GetServerStatus))
}

// GetUUID mocks base method.
func (m *MockComputationWrapper) GetUUID() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUUID")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetUUID indicates an expected call of GetUUID.
func (mr *MockComputationWrapperMockRecorder) GetUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUUID", reflect.TypeOf((*MockComputationWrapper)(nil).GetUUID))
}

// ParamVals mocks base method.
func (m *MockComputationWrapper) ParamVals() []any {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParamVals")
	ret0, _ := ret[0].([]any)
	return ret0
}

// ParamVals indicates an expected call of ParamVals.
func (mr *MockComputationWrapperMockRecorder) ParamVals() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParamVals", reflect.TypeOf((*MockComputationWrapper)(nil).ParamVals))
}

// Plan mocks base method.
func (m *MockComputationWrapper) Plan() *plan.Plan {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Plan")
	ret0, _ := ret[0].(*plan.Plan)
	return ret0
}

// Plan indicates an expected call of Plan.
func (mr *MockComputationWrapperMockRecorder) Plan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Plan", reflect.TypeOf((*MockComputationWrapper)(nil).Plan))
}

// RecordCompoundStmt mocks base method.
func (m *MockComputationWrapper) RecordCompoundStmt(ctx context.Context, statsBytes statistic.StatsArray) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordCompoundStmt", ctx, statsBytes)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecordCompoundStmt indicates an expected call of RecordCompoundStmt.
func (mr *MockComputationWrapperMockRecorder) RecordCompoundStmt(ctx, statsBytes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordCompoundStmt", reflect.TypeOf((*MockComputationWrapper)(nil).RecordCompoundStmt), ctx, statsBytes)
}

// RecordExecPlan mocks base method.
func (m *MockComputationWrapper) RecordExecPlan(ctx context.Context, phyPlan *models.PhyPlan) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordExecPlan", ctx, phyPlan)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecordExecPlan indicates an expected call of RecordExecPlan.
func (mr *MockComputationWrapperMockRecorder) RecordExecPlan(ctx, phyPlan interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordExecPlan", reflect.TypeOf((*MockComputationWrapper)(nil).RecordExecPlan), ctx, phyPlan)
}

// ResetPlanAndStmt mocks base method.
func (m *MockComputationWrapper) ResetPlanAndStmt(stmt tree.Statement) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetPlanAndStmt", stmt)
}

// ResetPlanAndStmt indicates an expected call of ResetPlanAndStmt.
func (mr *MockComputationWrapperMockRecorder) ResetPlanAndStmt(stmt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPlanAndStmt", reflect.TypeOf((*MockComputationWrapper)(nil).ResetPlanAndStmt), stmt)
}

// Run mocks base method.
func (m *MockComputationWrapper) Run(ts uint64) (*util.RunResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ts)
	ret0, _ := ret[0].(*util.RunResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockComputationWrapperMockRecorder) Run(ts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockComputationWrapper)(nil).Run), ts)
}

// SetExplainBuffer mocks base method.
func (m *MockComputationWrapper) SetExplainBuffer(buf *bytes.Buffer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetExplainBuffer", buf)
}

// SetExplainBuffer indicates an expected call of SetExplainBuffer.
func (mr *MockComputationWrapperMockRecorder) SetExplainBuffer(buf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetExplainBuffer", reflect.TypeOf((*MockComputationWrapper)(nil).SetExplainBuffer), buf)
}

// StatsCompositeSubStmtResource mocks base method.
func (m *MockComputationWrapper) StatsCompositeSubStmtResource(ctx context.Context) statistic.StatsArray {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StatsCompositeSubStmtResource", ctx)
	ret0, _ := ret[0].(statistic.StatsArray)
	return ret0
}

// StatsCompositeSubStmtResource indicates an expected call of StatsCompositeSubStmtResource.
func (mr *MockComputationWrapperMockRecorder) StatsCompositeSubStmtResource(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StatsCompositeSubStmtResource", reflect.TypeOf((*MockComputationWrapper)(nil).StatsCompositeSubStmtResource), ctx)
}

// MockColumnInfo is a mock of ColumnInfo interface.
type MockColumnInfo struct {
	ctrl     *gomock.Controller
	recorder *MockColumnInfoMockRecorder
}

// MockColumnInfoMockRecorder is the mock recorder for MockColumnInfo.
type MockColumnInfoMockRecorder struct {
	mock *MockColumnInfo
}

// NewMockColumnInfo creates a new mock instance.
func NewMockColumnInfo(ctrl *gomock.Controller) *MockColumnInfo {
	mock := &MockColumnInfo{ctrl: ctrl}
	mock.recorder = &MockColumnInfoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockColumnInfo) EXPECT() *MockColumnInfoMockRecorder {
	return m.recorder
}

// GetName mocks base method.
func (m *MockColumnInfo) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockColumnInfoMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockColumnInfo)(nil).GetName))
}

// GetType mocks base method.
func (m *MockColumnInfo) GetType() types.T {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetType")
	ret0, _ := ret[0].(types.T)
	return ret0
}

// GetType indicates an expected call of GetType.
func (mr *MockColumnInfoMockRecorder) GetType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetType", reflect.TypeOf((*MockColumnInfo)(nil).GetType))
}

// MockTableInfo is a mock of TableInfo interface.
type MockTableInfo struct {
	ctrl     *gomock.Controller
	recorder *MockTableInfoMockRecorder
}

// MockTableInfoMockRecorder is the mock recorder for MockTableInfo.
type MockTableInfoMockRecorder struct {
	mock *MockTableInfo
}

// NewMockTableInfo creates a new mock instance.
func NewMockTableInfo(ctrl *gomock.Controller) *MockTableInfo {
	mock := &MockTableInfo{ctrl: ctrl}
	mock.recorder = &MockTableInfoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTableInfo) EXPECT() *MockTableInfoMockRecorder {
	return m.recorder
}

// GetColumns mocks base method.
func (m *MockTableInfo) GetColumns() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetColumns")
}

// GetColumns indicates an expected call of GetColumns.
func (mr *MockTableInfoMockRecorder) GetColumns() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetColumns", reflect.TypeOf((*MockTableInfo)(nil).GetColumns))
}

// MockExecResult is a mock of ExecResult interface.
type MockExecResult struct {
	ctrl     *gomock.Controller
	recorder *MockExecResultMockRecorder
}

// MockExecResultMockRecorder is the mock recorder for MockExecResult.
type MockExecResultMockRecorder struct {
	mock *MockExecResult
}

// NewMockExecResult creates a new mock instance.
func NewMockExecResult(ctrl *gomock.Controller) *MockExecResult {
	mock := &MockExecResult{ctrl: ctrl}
	mock.recorder = &MockExecResultMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExecResult) EXPECT() *MockExecResultMockRecorder {
	return m.recorder
}

// ColumnIsNull mocks base method.
func (m *MockExecResult) ColumnIsNull(ctx context.Context, rindex, cindex uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ColumnIsNull", ctx, rindex, cindex)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ColumnIsNull indicates an expected call of ColumnIsNull.
func (mr *MockExecResultMockRecorder) ColumnIsNull(ctx, rindex, cindex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ColumnIsNull", reflect.TypeOf((*MockExecResult)(nil).ColumnIsNull), ctx, rindex, cindex)
}

// GetInt64 mocks base method.
func (m *MockExecResult) GetInt64(ctx context.Context, rindex, cindex uint64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt64", ctx, rindex, cindex)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInt64 indicates an expected call of GetInt64.
func (mr *MockExecResultMockRecorder) GetInt64(ctx, rindex, cindex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt64", reflect.TypeOf((*MockExecResult)(nil).GetInt64), ctx, rindex, cindex)
}

// GetRowCount mocks base method.
func (m *MockExecResult) GetRowCount() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRowCount")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetRowCount indicates an expected call of GetRowCount.
func (mr *MockExecResultMockRecorder) GetRowCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRowCount", reflect.TypeOf((*MockExecResult)(nil).GetRowCount))
}

// GetString mocks base method.
func (m *MockExecResult) GetString(ctx context.Context, rindex, cindex uint64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", ctx, rindex, cindex)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetString indicates an expected call of GetString.
func (mr *MockExecResultMockRecorder) GetString(ctx, rindex, cindex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockExecResult)(nil).GetString), ctx, rindex, cindex)
}

// GetUint64 mocks base method.
func (m *MockExecResult) GetUint64(ctx context.Context, rindex, cindex uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUint64", ctx, rindex, cindex)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUint64 indicates an expected call of GetUint64.
func (mr *MockExecResultMockRecorder) GetUint64(ctx, rindex, cindex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUint64", reflect.TypeOf((*MockExecResult)(nil).GetUint64), ctx, rindex, cindex)
}

// MockBackgroundExec is a mock of BackgroundExec interface.
type MockBackgroundExec struct {
	ctrl     *gomock.Controller
	recorder *MockBackgroundExecMockRecorder
}

// MockBackgroundExecMockRecorder is the mock recorder for MockBackgroundExec.
type MockBackgroundExecMockRecorder struct {
	mock *MockBackgroundExec
}

// NewMockBackgroundExec creates a new mock instance.
func NewMockBackgroundExec(ctrl *gomock.Controller) *MockBackgroundExec {
	mock := &MockBackgroundExec{ctrl: ctrl}
	mock.recorder = &MockBackgroundExecMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackgroundExec) EXPECT() *MockBackgroundExecMockRecorder {
	return m.recorder
}

// Clear mocks base method.
func (m *MockBackgroundExec) Clear() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Clear")
}

// Clear indicates an expected call of Clear.
func (mr *MockBackgroundExecMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockBackgroundExec)(nil).Clear))
}

// ClearExecResultBatches mocks base method.
func (m *MockBackgroundExec) ClearExecResultBatches() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearExecResultBatches")
}

// ClearExecResultBatches indicates an expected call of ClearExecResultBatches.
func (mr *MockBackgroundExecMockRecorder) ClearExecResultBatches() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearExecResultBatches", reflect.TypeOf((*MockBackgroundExec)(nil).ClearExecResultBatches))
}

// ClearExecResultSet mocks base method.
func (m *MockBackgroundExec) ClearExecResultSet() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearExecResultSet")
}

// ClearExecResultSet indicates an expected call of ClearExecResultSet.
func (mr *MockBackgroundExecMockRecorder) ClearExecResultSet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearExecResultSet", reflect.TypeOf((*MockBackgroundExec)(nil).ClearExecResultSet))
}

// Close mocks base method.
func (m *MockBackgroundExec) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockBackgroundExecMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBackgroundExec)(nil).Close))
}

// Exec mocks base method.
func (m *MockBackgroundExec) Exec(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Exec indicates an expected call of Exec.
func (mr *MockBackgroundExecMockRecorder) Exec(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockBackgroundExec)(nil).Exec), arg0, arg1)
}

// ExecRestore mocks base method.
func (m *MockBackgroundExec) ExecRestore(arg0 context.Context, arg1 string, arg2, arg3 uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecRestore", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecRestore indicates an expected call of ExecRestore.
func (mr *MockBackgroundExecMockRecorder) ExecRestore(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecRestore", reflect.TypeOf((*MockBackgroundExec)(nil).ExecRestore), arg0, arg1, arg2, arg3)
}

// ExecStmt mocks base method.
func (m *MockBackgroundExec) ExecStmt(arg0 context.Context, arg1 tree.Statement) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecStmt", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecStmt indicates an expected call of ExecStmt.
func (mr *MockBackgroundExecMockRecorder) ExecStmt(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecStmt", reflect.TypeOf((*MockBackgroundExec)(nil).ExecStmt), arg0, arg1)
}

// GetExecResultBatches mocks base method.
func (m *MockBackgroundExec) GetExecResultBatches() []*batch.Batch {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExecResultBatches")
	ret0, _ := ret[0].([]*batch.Batch)
	return ret0
}

// GetExecResultBatches indicates an expected call of GetExecResultBatches.
func (mr *MockBackgroundExecMockRecorder) GetExecResultBatches() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExecResultBatches", reflect.TypeOf((*MockBackgroundExec)(nil).GetExecResultBatches))
}

// GetExecResultSet mocks base method.
func (m *MockBackgroundExec) GetExecResultSet() []interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExecResultSet")
	ret0, _ := ret[0].([]interface{})
	return ret0
}

// GetExecResultSet indicates an expected call of GetExecResultSet.
func (mr *MockBackgroundExecMockRecorder) GetExecResultSet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExecResultSet", reflect.TypeOf((*MockBackgroundExec)(nil).GetExecResultSet))
}

// GetExecStatsArray mocks base method.
func (m *MockBackgroundExec) GetExecStatsArray() statistic.StatsArray {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExecStatsArray")
	ret0, _ := ret[0].(statistic.StatsArray)
	return ret0
}

// GetExecStatsArray indicates an expected call of GetExecStatsArray.
func (mr *MockBackgroundExecMockRecorder) GetExecStatsArray() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExecStatsArray", reflect.TypeOf((*MockBackgroundExec)(nil).GetExecStatsArray))
}

// Service mocks base method.
func (m *MockBackgroundExec) Service() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Service")
	ret0, _ := ret[0].(string)
	return ret0
}

// Service indicates an expected call of Service.
func (mr *MockBackgroundExecMockRecorder) Service() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Service", reflect.TypeOf((*MockBackgroundExec)(nil).Service))
}

// SetRestore mocks base method.
func (m *MockBackgroundExec) SetRestore(b bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRestore", b)
}

// SetRestore indicates an expected call of SetRestore.
func (mr *MockBackgroundExecMockRecorder) SetRestore(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRestore", reflect.TypeOf((*MockBackgroundExec)(nil).SetRestore), b)
}
