package ai

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

type StepContext struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	ExitCode      int32                  `protobuf:"varint,3,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"`
	LogOutput     string                 `protobuf:"bytes,4,opt,name=log_output,json=logOutput,proto3" json:"log_output,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StepContext) Reset() {
	*x = StepContext{}
	mi := &file_ai_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StepContext) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StepContext) ProtoMessage() {}

func (x *StepContext) ProtoReflect() protoreflect.Message {
	mi := &file_ai_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*StepContext) Descriptor() ([]byte, []int) {
	return file_ai_proto_rawDescGZIP(), []int{0}
}

func (x *StepContext) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *StepContext) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *StepContext) GetExitCode() int32 {
	if x != nil {
		return x.ExitCode
	}
	return 0
}

func (x *StepContext) GetLogOutput() string {
	if x != nil {
		return x.LogOutput
	}
	return ""
}

type AnalyzeFailureRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	JobName       string                 `protobuf:"bytes,1,opt,name=job_name,json=jobName,proto3" json:"job_name,omitempty"`
	JobStatus     string                 `protobuf:"bytes,2,opt,name=job_status,json=jobStatus,proto3" json:"job_status,omitempty"`
	PipelineYaml  string                 `protobuf:"bytes,3,opt,name=pipeline_yaml,json=pipelineYaml,proto3" json:"pipeline_yaml,omitempty"`
	Steps         []*StepContext         `protobuf:"bytes,4,rep,name=steps,proto3" json:"steps,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AnalyzeFailureRequest) Reset() {
	*x = AnalyzeFailureRequest{}
	mi := &file_ai_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AnalyzeFailureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeFailureRequest) ProtoMessage() {}

func (x *AnalyzeFailureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ai_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*AnalyzeFailureRequest) Descriptor() ([]byte, []int) {
	return file_ai_proto_rawDescGZIP(), []int{1}
}

func (x *AnalyzeFailureRequest) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *AnalyzeFailureRequest) GetJobStatus() string {
	if x != nil {
		return x.JobStatus
	}
	return ""
}

func (x *AnalyzeFailureRequest) GetPipelineYaml() string {
	if x != nil {
		return x.PipelineYaml
	}
	return ""
}

func (x *AnalyzeFailureRequest) GetSteps() []*StepContext {
	if x != nil {
		return x.Steps
	}
	return nil
}

type AnalyzeFailureResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Summary       string                 `protobuf:"bytes,1,opt,name=summary,proto3" json:"summary,omitempty"`
	RootCause     string                 `protobuf:"bytes,2,opt,name=root_cause,json=rootCause,proto3" json:"root_cause,omitempty"`
	Fix           string                 `protobuf:"bytes,3,opt,name=fix,proto3" json:"fix,omitempty"`
	RelevantLines []string               `protobuf:"bytes,4,rep,name=relevant_lines,json=relevantLines,proto3" json:"relevant_lines,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AnalyzeFailureResponse) Reset() {
	*x = AnalyzeFailureResponse{}
	mi := &file_ai_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AnalyzeFailureResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeFailureResponse) ProtoMessage() {}

func (x *AnalyzeFailureResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ai_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*AnalyzeFailureResponse) Descriptor() ([]byte, []int) {
	return file_ai_proto_rawDescGZIP(), []int{2}
}

func (x *AnalyzeFailureResponse) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *AnalyzeFailureResponse) GetRootCause() string {
	if x != nil {
		return x.RootCause
	}
	return ""
}

func (x *AnalyzeFailureResponse) GetFix() string {
	if x != nil {
		return x.Fix
	}
	return ""
}

func (x *AnalyzeFailureResponse) GetRelevantLines() []string {
	if x != nil {
		return x.RelevantLines
	}
	return nil
}

type GeneratePipelineRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Description   string                 `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GeneratePipelineRequest) Reset() {
	*x = GeneratePipelineRequest{}
	mi := &file_ai_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GeneratePipelineRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeneratePipelineRequest) ProtoMessage() {}

func (x *GeneratePipelineRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ai_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GeneratePipelineRequest) Descriptor() ([]byte, []int) {
	return file_ai_proto_rawDescGZIP(), []int{3}
}

func (x *GeneratePipelineRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type GeneratePipelineResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Yaml          string                 `protobuf:"bytes,1,opt,name=yaml,proto3" json:"yaml,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GeneratePipelineResponse) Reset() {
	*x = GeneratePipelineResponse{}
	mi := &file_ai_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GeneratePipelineResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeneratePipelineResponse) ProtoMessage() {}

func (x *GeneratePipelineResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ai_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GeneratePipelineResponse) Descriptor() ([]byte, []int) {
	return file_ai_proto_rawDescGZIP(), []int{4}
}

func (x *GeneratePipelineResponse) GetYaml() string {
	if x != nil {
		return x.Yaml
	}
	return ""
}

type ValidationIssue struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Line          int32                  `protobuf:"varint,1,opt,name=line,proto3" json:"line,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ValidationIssue) Reset() {
	*x = ValidationIssue{}
	mi := &file_ai_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidationIssue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidationIssue) ProtoMessage() {}

func (x *ValidationIssue) ProtoReflect() protoreflect.Message {
	mi := &file_ai_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ValidationIssue) Descriptor() ([]byte, []int) {
	return file_ai_proto_rawDescGZIP(), []int{5}
}

func (x *ValidationIssue) GetLine() int32 {
	if x != nil {
		return x.Line
	}
	return 0
}

func (x *ValidationIssue) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type PipelineCopilotRequest struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	ProjectId         string                 `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Yaml              string                 `protobuf:"bytes,2,opt,name=yaml,proto3" json:"yaml,omitempty"`
	Action            string                 `protobuf:"bytes,3,opt,name=action,proto3" json:"action,omitempty"`
	CustomPrompt      string                 `protobuf:"bytes,4,opt,name=custom_prompt,json=customPrompt,proto3" json:"custom_prompt,omitempty"`
	ValidationErrors  []*ValidationIssue     `protobuf:"bytes,5,rep,name=validation_errors,json=validationErrors,proto3" json:"validation_errors,omitempty"`
	Lang              string                 `protobuf:"bytes,6,opt,name=lang,proto3" json:"lang,omitempty"`
	RunHistoryContext string                 `protobuf:"bytes,7,opt,name=run_history_context,json=runHistoryContext,proto3" json:"run_history_context,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *PipelineCopilotRequest) Reset() {
	*x = PipelineCopilotRequest{}
	mi := &file_ai_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PipelineCopilotRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PipelineCopilotRequest) ProtoMessage() {}

func (x *PipelineCopilotRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ai_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*PipelineCopilotRequest) Descriptor() ([]byte, []int) {
	return file_ai_proto_rawDescGZIP(), []int{6}
}

func (x *PipelineCopilotRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *PipelineCopilotRequest) GetYaml() string {
	if x != nil {
		return x.Yaml
	}
	return ""
}

func (x *PipelineCopilotRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *PipelineCopilotRequest) GetCustomPrompt() string {
	if x != nil {
		return x.CustomPrompt
	}
	return ""
}

func (x *PipelineCopilotRequest) GetValidationErrors() []*ValidationIssue {
	if x != nil {
		return x.ValidationErrors
	}
	return nil
}

func (x *PipelineCopilotRequest) GetLang() string {
	if x != nil {
		return x.Lang
	}
	return ""
}

func (x *PipelineCopilotRequest) GetRunHistoryContext() string {
	if x != nil {
		return x.RunHistoryContext
	}
	return ""
}

type PipelineCopilotResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Yaml          string                 `protobuf:"bytes,1,opt,name=yaml,proto3" json:"yaml,omitempty"`
	Explanation   string                 `protobuf:"bytes,2,opt,name=explanation,proto3" json:"explanation,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PipelineCopilotResponse) Reset() {
	*x = PipelineCopilotResponse{}
	mi := &file_ai_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PipelineCopilotResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PipelineCopilotResponse) ProtoMessage() {}

func (x *PipelineCopilotResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ai_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*PipelineCopilotResponse) Descriptor() ([]byte, []int) {
	return file_ai_proto_rawDescGZIP(), []int{7}
}

func (x *PipelineCopilotResponse) GetYaml() string {
	if x != nil {
		return x.Yaml
	}
	return ""
}

func (x *PipelineCopilotResponse) GetExplanation() string {
	if x != nil {
		return x.Explanation
	}
	return ""
}

var File_ai_proto protoreflect.FileDescriptor

const file_ai_proto_rawDesc = "" +
	"\n" +
	"\bai.proto\x12\x02ai\"u\n" +
	"\vStepContext\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x16\n" +
	"\x06status\x18\x02 \x01(\tR\x06status\x12\x1b\n" +
	"\texit_code\x18\x03 \x01(\x05R\bexitCode\x12\x1d\n" +
	"\n" +
	"log_output\x18\x04 \x01(\tR\tlogOutput\"\x9d\x01\n" +
	"\x15AnalyzeFailureRequest\x12\x19\n" +
	"\bjob_name\x18\x01 \x01(\tR\ajobName\x12\x1d\n" +
	"\n" +
	"job_status\x18\x02 \x01(\tR\tjobStatus\x12#\n" +
	"\rpipeline_yaml\x18\x03 \x01(\tR\fpipelineYaml\x12%\n" +
	"\x05steps\x18\x04 \x03(\v2\x0f.ai.StepContextR\x05steps\"\x8a\x01\n" +
	"\x16AnalyzeFailureResponse\x12\x18\n" +
	"\asummary\x18\x01 \x01(\tR\asummary\x12\x1d\n" +
	"\n" +
	"root_cause\x18\x02 \x01(\tR\trootCause\x12\x10\n" +
	"\x03fix\x18\x03 \x01(\tR\x03fix\x12%\n" +
	"\x0erelevant_lines\x18\x04 \x03(\tR\rrelevantLines\";\n" +
	"\x17GeneratePipelineRequest\x12 \n" +
	"\vdescription\x18\x01 \x01(\tR\vdescription\".\n" +
	"\x18GeneratePipelineResponse\x12\x12\n" +
	"\x04yaml\x18\x01 \x01(\tR\x04yaml\"?\n" +
	"\x0fValidationIssue\x12\x12\n" +
	"\x04line\x18\x01 \x01(\x05R\x04line\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\"\x8e\x02\n" +
	"\x16PipelineCopilotRequest\x12\x1d\n" +
	"\n" +
	"project_id\x18\x01 \x01(\tR\tprojectId\x12\x12\n" +
	"\x04yaml\x18\x02 \x01(\tR\x04yaml\x12\x16\n" +
	"\x06action\x18\x03 \x01(\tR\x06action\x12#\n" +
	"\rcustom_prompt\x18\x04 \x01(\tR\fcustomPrompt\x12@\n" +
	"\x11validation_errors\x18\x05 \x03(\v2\x13.ai.ValidationIssueR\x10validationErrors\x12\x12\n" +
	"\x04lang\x18\x06 \x01(\tR\x04lang\x12.\n" +
	"\x13run_history_context\x18\a \x01(\tR\x11runHistoryContext\"O\n" +
	"\x17PipelineCopilotResponse\x12\x12\n" +
	"\x04yaml\x18\x01 \x01(\tR\x04yaml\x12 \n" +
	"\vexplanation\x18\x02 \x01(\tR\vexplanation2\xef\x01\n" +
	"\tAIService\x12G\n" +
	"\x0eAnalyzeFailure\x12\x19.ai.AnalyzeFailureRequest\x1a\x1a.ai.AnalyzeFailureResponse\x12M\n" +
	"\x10GeneratePipeline\x12\x1b.ai.GeneratePipelineRequest\x1a\x1c.ai.GeneratePipelineResponse\x12J\n" +
	"\x0fPipelineCopilot\x12\x1a.ai.PipelineCopilotRequest\x1a\x1b.ai.PipelineCopilotResponseB-Z+github.com/vsuaiqq/cicd/shared/proto/gen/aib\x06proto3"

var (
	file_ai_proto_rawDescOnce sync.Once
	file_ai_proto_rawDescData []byte
)

func file_ai_proto_rawDescGZIP() []byte {
	file_ai_proto_rawDescOnce.Do(func() {
		file_ai_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_ai_proto_rawDesc), len(file_ai_proto_rawDesc)))
	})
	return file_ai_proto_rawDescData
}

var file_ai_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_ai_proto_goTypes = []any{
	(*StepContext)(nil),
	(*AnalyzeFailureRequest)(nil),
	(*AnalyzeFailureResponse)(nil),
	(*GeneratePipelineRequest)(nil),
	(*GeneratePipelineResponse)(nil),
	(*ValidationIssue)(nil),
	(*PipelineCopilotRequest)(nil),
	(*PipelineCopilotResponse)(nil),
}
var file_ai_proto_depIdxs = []int32{
	0,
	5,
	1,
	3,
	6,
	2,
	4,
	7,
	5,
	2,
	2,
	2,
	0,
}

func init() { file_ai_proto_init() }
func file_ai_proto_init() {
	if File_ai_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_ai_proto_rawDesc), len(file_ai_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ai_proto_goTypes,
		DependencyIndexes: file_ai_proto_depIdxs,
		MessageInfos:      file_ai_proto_msgTypes,
	}.Build()
	File_ai_proto = out.File
	file_ai_proto_goTypes = nil
	file_ai_proto_depIdxs = nil
}
