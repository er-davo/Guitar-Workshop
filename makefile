PROTO_DIR = proto

ANALYZER_PY_VENV = analyzer/venv

ANALYZER_PY_OUT = analyzer/app
ANALYZER_GRPC_PY_OUT = analyzer/app

SEPARATOR_PY_VENV = audio-separator/venv

SEPARATOR_PY_OUT = audio-separator/app
SEPARATOR_GRPC_PY_OUT = audio-separator/app

TABGEN_GO_OUT = tab-generator/app
TABGEN_GO_GRPC_OUT = tab-generator/app

API_GO_OUT = api-gateway/app
API_GO_GRPC_OUT = api-gateway/app

PREPROC_CPP_OUT = audio-preprocessor/app
PREPROC_CPP_GRPC_OUT = audio-preprocessor/app
CONAN_PROTOC = C:\Users\Lenovo\.conan2\p\proto1344852724c4b\p\bin\protoc.exe

GRPC_PLUGIN = C:\Users\Lenovo\.conan2\p\grpc2a6788fd4476e\p\bin\grpc_cpp_plugin.exe

proto-gen: proto-go-gen proto-py-gen proto-cpp-gen

proto-py-gen: proto-py-gen-separator proto-py-gen-analyzer

proto-go-gen: proto-go-gen-api proto-go-gen-tabgen

proto-go-gen-api: proto-go-gen-api-tab proto-go-gen-api-separator proto-go-gen-api-processor

proto-go-gen-tabgen: proto-go-gen-tabgen-analyzer proto-go-gen-tabgen-tab

proto-go-gen-tabgen-analyzer:
	protoc --go_out=$(TABGEN_GO_OUT) \
	--go-grpc_out=$(TABGEN_GO_GRPC_OUT) \
	-Iproto $(PROTO_DIR)/note_analyzer.proto

proto-go-gen-tabgen-tab:
	protoc --go_out=$(TABGEN_GO_OUT) \
	--go-grpc_out=$(TABGEN_GO_GRPC_OUT) \
	-Iproto $(PROTO_DIR)/tab.proto

proto-go-gen-api-tab:
	protoc --go_out=$(API_GO_OUT) \
	--go-grpc_out=$(API_GO_GRPC_OUT) \
	-Iproto $(PROTO_DIR)/tab.proto

proto-go-gen-api-separator:
	protoc --go_out=$(API_GO_OUT) \
	--go-grpc_out=$(API_GO_GRPC_OUT) \
	-Iproto $(PROTO_DIR)/separator.proto

proto-go-gen-api-processor:
	protoc --go_out=$(API_GO_OUT) \
	--go-grpc_out=$(API_GO_GRPC_OUT) \
	-Iproto $(PROTO_DIR)/processor.proto

proto-cpp-gen:
	$(CONAN_PROTOC) --cpp_out=$(PREPROC_CPP_OUT) \
	--grpc_out=$(PREPROC_CPP_GRPC_OUT) \
	--plugin=protoc-gen-grpc=$(GRPC_PLUGIN) \
	$(PROTO_DIR)/processor.proto

proto-py-gen-separator:
	$(SEPARATOR_PY_VENV)/Scripts/python.exe \
	-m grpc_tools.protoc -I proto \
	--python_out=$(SEPARATOR_PY_OUT) \
	--pyi_out=$(SEPARATOR_PY_OUT) \
	--grpc_python_out=$(SEPARATOR_GRPC_PY_OUT) \
	$(PROTO_DIR)/separator.proto

proto-py-gen-analyzer:
	$(ANALYZER_PY_VENV)/Scripts/python.exe \
	-m grpc_tools.protoc -I proto \
	--python_out=$(ANALYZER_PY_OUT) \
	--pyi_out=$(ANALYZER_PY_OUT) \
	--grpc_python_out=$(ANALYZER_GRPC_PY_OUT) \
	$(PROTO_DIR)/note_analyzer.proto