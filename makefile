PROTO_DIR = proto

PYTHON_VENV = audio-analyzer/venv

PYTHON_OUT = audio-analyzer/app/audioproto
GRPC_PYTHON_OUT = audio-analyzer/app/audioproto

TABGEN_GO_OUT = tab-generator/app
TABGEN_GO_GRPC_OUT = tab-generator/app

API_GO_OUT = api-gateway/app
API_GO_GRPC_OUT = api-gateway/app

py-venv:
	cd $(PYTHON_VENV)

proto-gen: proto-go-gen-tabgen-audio proto-go-gen-tabgen-tab proto-go-gen-api-tab proto-py-gen

proto-go-gen-tabgen-audio:
	protoc --go_out=$(TABGEN_GO_OUT) --go-grpc_out=$(TABGEN_GO_GRPC_OUT) -Iproto $(PROTO_DIR)/audio.proto

proto-go-gen-tabgen-tab:
	protoc --go_out=$(TABGEN_GO_OUT) --go-grpc_out=$(TABGEN_GO_GRPC_OUT) -Iproto $(PROTO_DIR)/tab.proto

proto-go-gen-api-tab:
	protoc --go_out=$(API_GO_OUT) --go-grpc_out=$(API_GO_GRPC_OUT) -Iproto $(PROTO_DIR)/tab.proto

proto-py-gen:
	$(PYTHON_VENV)/Scripts/python.exe -m grpc_tools.protoc -I proto --python_out=$(PYTHON_OUT) --pyi_out=$(PYTHON_OUT) --grpc_python_out=$(GRPC_PYTHON_OUT) $(PROTO_DIR)/audio.proto