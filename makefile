# Папка с proto-контрактами
PROTO_DIR = contracts/proto

# Python analyzer
ANALYZER_PY_VENV = analyzer/venv
ANALYZER_PY_OUT = analyzer/app
ANALYZER_GRPC_PY_OUT = analyzer/app

# Go tab-generator
TABGEN_GO_OUT = tab-generator/app
TABGEN_GO_GRPC_OUT = tab-generator/app

# Go api-gateway
API_GO_OUT = api-gateway/app
API_GO_GRPC_OUT = api-gateway/app

# Go tab-service
TAB_SERVICE_GO_OUT = tab-service/app
TAB_SERVICE_GO_GRPC_OUT = tab-service/app

# Go pb out
TAB_PB_OUT = contracts/tab
ANALYZER_PB_OUT = contracts/note_analyzer

# Audio preprocessor C++ (если нужен)
PREPROC_CPP_OUT = audio-preprocessor/app
PREPROC_CPP_GRPC_OUT = audio-preprocessor/app

# Протокол и плагины
CONAN_PROTOC = C:\Users\Lenovo\.conan2\p\proto1344852724c4b\p\bin\protoc.exe
GRPC_PLUGIN = C:\Users\Lenovo\.conan2\p\grpc2a6788fd4476e\p\bin\grpc_cpp_plugin.exe

# Общая цель генерации
proto-gen: proto-go-gen proto-py-gen

# Python
proto-py-gen: proto-py-gen-analyzer

proto-py-gen-analyzer:
	$(ANALYZER_PY_VENV)/Scripts/python.exe \
	-m grpc_tools.protoc \
	-I $(PROTO_DIR) \
	--python_out=$(ANALYZER_PY_OUT) \
	--pyi_out=$(ANALYZER_PY_OUT) \
	--grpc_python_out=$(ANALYZER_GRPC_PY_OUT) \
	$(PROTO_DIR)/note_analyzer.proto

# Go
proto-go-gen: proto-go-gen-analyzer proto-go-gen-tab-service

# Go
proto-go-gen-tab-service:
	protoc \
	  --go_out=. \
	  --go-grpc_out=. \
	  -I $(PROTO_DIR) \
	  $(PROTO_DIR)/tab.proto

proto-go-gen-analyzer:
	protoc \
	  --go_out=. \
	  --go-grpc_out=. \
	  -I $(PROTO_DIR) \
	  $(PROTO_DIR)/note_analyzer.proto