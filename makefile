PROTO_DIR = proto

PYTHON_VENV = audio-analyzer/venv

PYTHON_OUT = audio-analyzer/app
GRPC_PYTHON_OUT = audio-analyzer/app

TABGEN_GO_OUT = tab-generator/app
TABGEN_GO_GRPC_OUT = tab-generator/app

API_GO_OUT = api-gateway/app
API_GO_GRPC_OUT = api-gateway/app

FX_CPP_OUT = fx-processor/app/gen
FX_CPP_GRPC_OUT = fx-processor/app/gen

GRPC_PLUGIN = .plugins/grpc_cpp_plugin.exe

docker run --rm -v "C:\Users\Lenovo\Desktop\VScodeFiles\Guitar Workshop:/proto" znly/protoc --cpp_out=/proto/generated --grpc_out=/proto/generated --plugin=protoc-gen-grpc=/usr/bin/grpc_cpp_plugin -I /proto /proto/proto/fx_processor.proto

py-venv:
	cd $(PYTHON_VENV)

proto-gen: proto-go-gen proto-py-gen proto-cpp-gen

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
	protoc --cpp_out=$(FX_CPP_OUT) --grpc_out=$(FX_CPP_GRPC_OUT) --plugin=protoc-gen-grpc=$(GRPC_PLUGIN) $(PROTO_DIR)/fx_processor.proto

proto-py-gen:
	$(PYTHON_VENV)/Scripts/python.exe -m grpc_tools.protoc -I proto --python_out=$(PYTHON_OUT) --pyi_out=$(PYTHON_OUT) --grpc_python_out=$(GRPC_PYTHON_OUT) $(PROTO_DIR)/audio.proto