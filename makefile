PROTO_DIR = proto

AUDIO_PY_VENV = audio-analyzer/venv

AUDIO_PY_OUT = audio-analyzer/app
AUDIO_GRPC_PY_OUT = audio-analyzer/app

RIFF_PY_VENV = riff-ai-gen/venv

RIFF_PY_OUT = riff-ai-gen/app
RIFF_GRPC_PY_OUT = riff-ai-gen/app

TABGEN_GO_OUT = tab-generator/app
TABGEN_GO_GRPC_OUT = tab-generator/app

API_GO_OUT = api-gateway/app
API_GO_GRPC_OUT = api-gateway/app

FX_CPP_OUT = fx-processor/app/gen
FX_CPP_GRPC_OUT = fx-processor/app/gen

PREPROC_CPP_OUT = audio-preprocessor/app
PREPROC_CPP_GRPC_OUT = audio-preprocessor/app
CONAN_PROTOC = C:\Users\Lenovo\.conan2\p\proto1344852724c4b\p\bin\protoc.exe

GRPC_PLUGIN = C:\Users\Lenovo\.conan2\p\grpc2a6788fd4476e\p\bin\grpc_cpp_plugin.exe

proto-cppgen:
	docker run --rm -v "C:\Users\Lenovo\Desktop\VScodeFiles\Guitar Workshop:/proto" znly/protoc --cpp_out=/proto/generated --grpc_out=/proto/generated --plugin=protoc-gen-grpc=/usr/bin/grpc_cpp_plugin -I /proto /proto/proto/fx_processor.proto

proto-gen: proto-go-gen proto-py-gen proto-cpp-gen

proto-py-gen: proto-py-gen-audio proto-py-gen-riff

proto-go-gen: proto-go-gen-tabgen-audio proto-go-gen-tabgen-tab proto-go-gen-api-tab proto-go-gen-api-fx

proto-go-gen-tabgen-audio:
	protoc --go_out=$(TABGEN_GO_OUT) --go-grpc_out=$(TABGEN_GO_GRPC_OUT) -Iproto $(PROTO_DIR)/audio.proto

proto-go-gen-tabgen-tab:
	protoc --go_out=$(TABGEN_GO_OUT) --go-grpc_out=$(TABGEN_GO_GRPC_OUT) -Iproto $(PROTO_DIR)/tab.proto

proto-go-gen-api-tab:
	protoc --go_out=$(API_GO_OUT) --go-grpc_out=$(API_GO_GRPC_OUT) -Iproto $(PROTO_DIR)/tab.proto

proto-go-gen-api-fx:
	protoc --go_out=$(API_GO_OUT) --go-grpc_out=$(API_GO_GRPC_OUT) -Iproto $(PROTO_DIR)/fx_processor.proto

proto-cpp-gen:
	$(CONAN_PROTOC) --cpp_out=$(PREPROC_CPP_OUT) --grpc_out=$(PREPROC_CPP_GRPC_OUT) --plugin=protoc-gen-grpc=$(GRPC_PLUGIN) $(PROTO_DIR)/audio_processor.proto

proto-py-gen-audio:
	$(AUDIO_PY_VENV)/Scripts/python.exe -m grpc_tools.protoc -I proto --python_out=$(AUDIO_PY_OUT) --pyi_out=$(AUDIO_PY_OUT) --grpc_python_out=$(AUDIO_GRPC_PY_OUT) $(PROTO_DIR)/audio.proto

proto-py-gen-riff:
	$(RIFF_PY_VENV)/Scripts/python.exe -m grpc_tools.protoc -I proto --python_out=$(RIFF_PY_OUT) --pyi_out=$(RIFF_PY_OUT) --grpc_python_out=$(RIFF_GRPC_PY_OUT) $(PROTO_DIR)/riff.proto