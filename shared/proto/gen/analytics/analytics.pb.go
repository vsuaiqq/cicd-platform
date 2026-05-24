package analytics

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)

	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DashboardRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProjectId     string                 `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Period        string                 `protobuf:"bytes,2,opt,name=period,proto3" json:"period,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DashboardRequest) Reset() {
	*x = DashboardRequest{}
	mi := &file_analytics_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DashboardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DashboardRequest) ProtoMessage() {}

func (x *DashboardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_analytics_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*DashboardRequest) Descriptor() ([]byte, []int) {
	return file_analytics_proto_rawDescGZIP(), []int{0}
}

func (x *DashboardRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *DashboardRequest) GetPeriod() string {
	if x != nil {
		return x.Period
	}
	return ""
}

type DailyPoint struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Date          string                 `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
	Total         int32                  `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	Success       int32                  `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
	Failed        int32                  `protobuf:"varint,4,opt,name=failed,proto3" json:"failed,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DailyPoint) Reset() {
	*x = DailyPoint{}
	mi := &file_analytics_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DailyPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DailyPoint) ProtoMessage() {}

func (x *DailyPoint) ProtoReflect() protoreflect.Message {
	mi := &file_analytics_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*DailyPoint) Descriptor() ([]byte, []int) {
	return file_analytics_proto_rawDescGZIP(), []int{1}
}

func (x *DailyPoint) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *DailyPoint) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *DailyPoint) GetSuccess() int32 {
	if x != nil {
		return x.Success
	}
	return 0
}

func (x *DailyPoint) GetFailed() int32 {
	if x != nil {
		return x.Failed
	}
	return 0
}

type JobStat struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	JobName        string                 `protobuf:"bytes,1,opt,name=job_name,json=jobName,proto3" json:"job_name,omitempty"`
	TotalRuns      int32                  `protobuf:"varint,2,opt,name=total_runs,json=totalRuns,proto3" json:"total_runs,omitempty"`
	FailureRate    float64                `protobuf:"fixed64,3,opt,name=failure_rate,json=failureRate,proto3" json:"failure_rate,omitempty"`
	AvgDurationSec float64                `protobuf:"fixed64,4,opt,name=avg_duration_sec,json=avgDurationSec,proto3" json:"avg_duration_sec,omitempty"`
	AvgAttempts    float64                `protobuf:"fixed64,5,opt,name=avg_attempts,json=avgAttempts,proto3" json:"avg_attempts,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *JobStat) Reset() {
	*x = JobStat{}
	mi := &file_analytics_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobStat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobStat) ProtoMessage() {}

func (x *JobStat) ProtoReflect() protoreflect.Message {
	mi := &file_analytics_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*JobStat) Descriptor() ([]byte, []int) {
	return file_analytics_proto_rawDescGZIP(), []int{2}
}

func (x *JobStat) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *JobStat) GetTotalRuns() int32 {
	if x != nil {
		return x.TotalRuns
	}
	return 0
}

func (x *JobStat) GetFailureRate() float64 {
	if x != nil {
		return x.FailureRate
	}
	return 0
}

func (x *JobStat) GetAvgDurationSec() float64 {
	if x != nil {
		return x.AvgDurationSec
	}
	return 0
}

func (x *JobStat) GetAvgAttempts() float64 {
	if x != nil {
		return x.AvgAttempts
	}
	return 0
}

type DashboardResponse struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	TotalRuns      int64                  `protobuf:"varint,1,opt,name=total_runs,json=totalRuns,proto3" json:"total_runs,omitempty"`
	SuccessCount   int64                  `protobuf:"varint,2,opt,name=success_count,json=successCount,proto3" json:"success_count,omitempty"`
	FailedCount    int64                  `protobuf:"varint,3,opt,name=failed_count,json=failedCount,proto3" json:"failed_count,omitempty"`
	CancelledCount int64                  `protobuf:"varint,4,opt,name=cancelled_count,json=cancelledCount,proto3" json:"cancelled_count,omitempty"`
	SuccessRate    float64                `protobuf:"fixed64,5,opt,name=success_rate,json=successRate,proto3" json:"success_rate,omitempty"`
	AvgDurationSec float64                `protobuf:"fixed64,6,opt,name=avg_duration_sec,json=avgDurationSec,proto3" json:"avg_duration_sec,omitempty"`
	P50DurationSec float64                `protobuf:"fixed64,7,opt,name=p50_duration_sec,json=p50DurationSec,proto3" json:"p50_duration_sec,omitempty"`
	P95DurationSec float64                `protobuf:"fixed64,8,opt,name=p95_duration_sec,json=p95DurationSec,proto3" json:"p95_duration_sec,omitempty"`

	Trend []*DailyPoint `protobuf:"bytes,9,rep,name=trend,proto3" json:"trend,omitempty"`

	TopFailingJobs []*JobStat `protobuf:"bytes,10,rep,name=top_failing_jobs,json=topFailingJobs,proto3" json:"top_failing_jobs,omitempty"`

	TopSlowJobs []*JobStat `protobuf:"bytes,11,rep,name=top_slow_jobs,json=topSlowJobs,proto3" json:"top_slow_jobs,omitempty"`

	FlakyJobs     []*JobStat `protobuf:"bytes,12,rep,name=flaky_jobs,json=flakyJobs,proto3" json:"flaky_jobs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DashboardResponse) Reset() {
	*x = DashboardResponse{}
	mi := &file_analytics_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DashboardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DashboardResponse) ProtoMessage() {}

func (x *DashboardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_analytics_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*DashboardResponse) Descriptor() ([]byte, []int) {
	return file_analytics_proto_rawDescGZIP(), []int{3}
}

func (x *DashboardResponse) GetTotalRuns() int64 {
	if x != nil {
		return x.TotalRuns
	}
	return 0
}

func (x *DashboardResponse) GetSuccessCount() int64 {
	if x != nil {
		return x.SuccessCount
	}
	return 0
}

func (x *DashboardResponse) GetFailedCount() int64 {
	if x != nil {
		return x.FailedCount
	}
	return 0
}

func (x *DashboardResponse) GetCancelledCount() int64 {
	if x != nil {
		return x.CancelledCount
	}
	return 0
}

func (x *DashboardResponse) GetSuccessRate() float64 {
	if x != nil {
		return x.SuccessRate
	}
	return 0
}

func (x *DashboardResponse) GetAvgDurationSec() float64 {
	if x != nil {
		return x.AvgDurationSec
	}
	return 0
}

func (x *DashboardResponse) GetP50DurationSec() float64 {
	if x != nil {
		return x.P50DurationSec
	}
	return 0
}

func (x *DashboardResponse) GetP95DurationSec() float64 {
	if x != nil {
		return x.P95DurationSec
	}
	return 0
}

func (x *DashboardResponse) GetTrend() []*DailyPoint {
	if x != nil {
		return x.Trend
	}
	return nil
}

func (x *DashboardResponse) GetTopFailingJobs() []*JobStat {
	if x != nil {
		return x.TopFailingJobs
	}
	return nil
}

func (x *DashboardResponse) GetTopSlowJobs() []*JobStat {
	if x != nil {
		return x.TopSlowJobs
	}
	return nil
}

func (x *DashboardResponse) GetFlakyJobs() []*JobStat {
	if x != nil {
		return x.FlakyJobs
	}
	return nil
}

type CopilotContextRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProjectId     string                 `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CopilotContextRequest) Reset() {
	*x = CopilotContextRequest{}
	mi := &file_analytics_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CopilotContextRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CopilotContextRequest) ProtoMessage() {}

func (x *CopilotContextRequest) ProtoReflect() protoreflect.Message {
	mi := &file_analytics_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*CopilotContextRequest) Descriptor() ([]byte, []int) {
	return file_analytics_proto_rawDescGZIP(), []int{4}
}

func (x *CopilotContextRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type RecentRunSummary struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RunId         string                 `protobuf:"bytes,1,opt,name=run_id,json=runId,proto3" json:"run_id,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Branch        string                 `protobuf:"bytes,3,opt,name=branch,proto3" json:"branch,omitempty"`
	DurationSec   uint32                 `protobuf:"varint,4,opt,name=duration_sec,json=durationSec,proto3" json:"duration_sec,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RecentRunSummary) Reset() {
	*x = RecentRunSummary{}
	mi := &file_analytics_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecentRunSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecentRunSummary) ProtoMessage() {}

func (x *RecentRunSummary) ProtoReflect() protoreflect.Message {
	mi := &file_analytics_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*RecentRunSummary) Descriptor() ([]byte, []int) {
	return file_analytics_proto_rawDescGZIP(), []int{5}
}

func (x *RecentRunSummary) GetRunId() string {
	if x != nil {
		return x.RunId
	}
	return ""
}

func (x *RecentRunSummary) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *RecentRunSummary) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

func (x *RecentRunSummary) GetDurationSec() uint32 {
	if x != nil {
		return x.DurationSec
	}
	return 0
}

func (x *RecentRunSummary) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

type CopilotContextResponse struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	RecentRuns     []*RecentRunSummary    `protobuf:"bytes,1,rep,name=recent_runs,json=recentRuns,proto3" json:"recent_runs,omitempty"`
	TopFailingJobs []*JobStat             `protobuf:"bytes,2,rep,name=top_failing_jobs,json=topFailingJobs,proto3" json:"top_failing_jobs,omitempty"`
	TopSlowJobs    []*JobStat             `protobuf:"bytes,3,rep,name=top_slow_jobs,json=topSlowJobs,proto3" json:"top_slow_jobs,omitempty"`
	FlakyJobs      []*JobStat             `protobuf:"bytes,4,rep,name=flaky_jobs,json=flakyJobs,proto3" json:"flaky_jobs,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *CopilotContextResponse) Reset() {
	*x = CopilotContextResponse{}
	mi := &file_analytics_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CopilotContextResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CopilotContextResponse) ProtoMessage() {}

func (x *CopilotContextResponse) ProtoReflect() protoreflect.Message {
	mi := &file_analytics_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*CopilotContextResponse) Descriptor() ([]byte, []int) {
	return file_analytics_proto_rawDescGZIP(), []int{6}
}

func (x *CopilotContextResponse) GetRecentRuns() []*RecentRunSummary {
	if x != nil {
		return x.RecentRuns
	}
	return nil
}

func (x *CopilotContextResponse) GetTopFailingJobs() []*JobStat {
	if x != nil {
		return x.TopFailingJobs
	}
	return nil
}

func (x *CopilotContextResponse) GetTopSlowJobs() []*JobStat {
	if x != nil {
		return x.TopSlowJobs
	}
	return nil
}

func (x *CopilotContextResponse) GetFlakyJobs() []*JobStat {
	if x != nil {
		return x.FlakyJobs
	}
	return nil
}

var File_analytics_proto protoreflect.FileDescriptor

const file_analytics_proto_rawDesc = "" +
	"\n" +
	"\x0fanalytics.proto\x12\tanalytics\"I\n" +
	"\x10DashboardRequest\x12\x1d\n" +
	"\n" +
	"project_id\x18\x01 \x01(\tR\tprojectId\x12\x16\n" +
	"\x06period\x18\x02 \x01(\tR\x06period\"h\n" +
	"\n" +
	"DailyPoint\x12\x12\n" +
	"\x04date\x18\x01 \x01(\tR\x04date\x12\x14\n" +
	"\x05total\x18\x02 \x01(\x05R\x05total\x12\x18\n" +
	"\asuccess\x18\x03 \x01(\x05R\asuccess\x12\x16\n" +
	"\x06failed\x18\x04 \x01(\x05R\x06failed\"\xb3\x01\n" +
	"\aJobStat\x12\x19\n" +
	"\bjob_name\x18\x01 \x01(\tR\ajobName\x12\x1d\n" +
	"\n" +
	"total_runs\x18\x02 \x01(\x05R\ttotalRuns\x12!\n" +
	"\ffailure_rate\x18\x03 \x01(\x01R\vfailureRate\x12(\n" +
	"\x10avg_duration_sec\x18\x04 \x01(\x01R\x0eavgDurationSec\x12!\n" +
	"\favg_attempts\x18\x05 \x01(\x01R\vavgAttempts\"\x9a\x04\n" +
	"\x11DashboardResponse\x12\x1d\n" +
	"\n" +
	"total_runs\x18\x01 \x01(\x03R\ttotalRuns\x12#\n" +
	"\rsuccess_count\x18\x02 \x01(\x03R\fsuccessCount\x12!\n" +
	"\ffailed_count\x18\x03 \x01(\x03R\vfailedCount\x12'\n" +
	"\x0fcancelled_count\x18\x04 \x01(\x03R\x0ecancelledCount\x12!\n" +
	"\fsuccess_rate\x18\x05 \x01(\x01R\vsuccessRate\x12(\n" +
	"\x10avg_duration_sec\x18\x06 \x01(\x01R\x0eavgDurationSec\x12(\n" +
	"\x10p50_duration_sec\x18\a \x01(\x01R\x0ep50DurationSec\x12(\n" +
	"\x10p95_duration_sec\x18\b \x01(\x01R\x0ep95DurationSec\x12+\n" +
	"\x05trend\x18\t \x03(\v2\x15.analytics.DailyPointR\x05trend\x12<\n" +
	"\x10top_failing_jobs\x18\n" +
	" \x03(\v2\x12.analytics.JobStatR\x0etopFailingJobs\x126\n" +
	"\rtop_slow_jobs\x18\v \x03(\v2\x12.analytics.JobStatR\vtopSlowJobs\x121\n" +
	"\n" +
	"flaky_jobs\x18\f \x03(\v2\x12.analytics.JobStatR\tflakyJobs\"6\n" +
	"\x15CopilotContextRequest\x12\x1d\n" +
	"\n" +
	"project_id\x18\x01 \x01(\tR\tprojectId\"\x9b\x01\n" +
	"\x10RecentRunSummary\x12\x15\n" +
	"\x06run_id\x18\x01 \x01(\tR\x05runId\x12\x16\n" +
	"\x06status\x18\x02 \x01(\tR\x06status\x12\x16\n" +
	"\x06branch\x18\x03 \x01(\tR\x06branch\x12!\n" +
	"\fduration_sec\x18\x04 \x01(\rR\vdurationSec\x12\x1d\n" +
	"\n" +
	"created_at\x18\x05 \x01(\tR\tcreatedAt\"\xff\x01\n" +
	"\x16CopilotContextResponse\x12<\n" +
	"\vrecent_runs\x18\x01 \x03(\v2\x1b.analytics.RecentRunSummaryR\n" +
	"recentRuns\x12<\n" +
	"\x10top_failing_jobs\x18\x02 \x03(\v2\x12.analytics.JobStatR\x0etopFailingJobs\x126\n" +
	"\rtop_slow_jobs\x18\x03 \x03(\v2\x12.analytics.JobStatR\vtopSlowJobs\x121\n" +
	"\n" +
	"flaky_jobs\x18\x04 \x03(\v2\x12.analytics.JobStatR\tflakyJobs2\xb7\x01\n" +
	"\x10AnalyticsService\x12I\n" +
	"\fGetDashboard\x12\x1b.analytics.DashboardRequest\x1a\x1c.analytics.DashboardResponse\x12X\n" +
	"\x11GetCopilotContext\x12 .analytics.CopilotContextRequest\x1a!.analytics.CopilotContextResponseB4Z2github.com/vsuaiqq/cicd/shared/proto/gen/analyticsb\x06proto3"

var (
	file_analytics_proto_rawDescOnce sync.Once
	file_analytics_proto_rawDescData []byte
)

func file_analytics_proto_rawDescGZIP() []byte {
	file_analytics_proto_rawDescOnce.Do(func() {
		file_analytics_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_analytics_proto_rawDesc), len(file_analytics_proto_rawDesc)))
	})
	return file_analytics_proto_rawDescData
}

var file_analytics_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_analytics_proto_goTypes = []any{
	(*DashboardRequest)(nil),
	(*DailyPoint)(nil),
	(*JobStat)(nil),
	(*DashboardResponse)(nil),
	(*CopilotContextRequest)(nil),
	(*RecentRunSummary)(nil),
	(*CopilotContextResponse)(nil),
}
var file_analytics_proto_depIdxs = []int32{
	1,
	2,
	2,
	2,
	5,
	2,
	2,
	2,
	0,
	4,
	3,
	6,
	10,
	8,
	8,
	8,
	0,
}

func init() { file_analytics_proto_init() }
func file_analytics_proto_init() {
	if File_analytics_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_analytics_proto_rawDesc), len(file_analytics_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_analytics_proto_goTypes,
		DependencyIndexes: file_analytics_proto_depIdxs,
		MessageInfos:      file_analytics_proto_msgTypes,
	}.Build()
	File_analytics_proto = out.File
	file_analytics_proto_goTypes = nil
	file_analytics_proto_depIdxs = nil
}
