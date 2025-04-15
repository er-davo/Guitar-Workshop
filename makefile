PROTO_DIR = proto

PYTHON_VENV = audio-analyzer/venv

PYTHON_OUT = audio-analyzer
GRPC_PYTHON_OUT = audio-analyzer

GO_OUT = tab-generator
GO_GRPC_OUT = tab-generator

py-venv:
	cd $(PYTHON_VENV)

proto-gen: proto-go-gen proto-py-gen

proto-go-gen:
	protoc --go_out=$(GO_OUT) --go-grpc_out=$(GO_GRPC_OUT) -Iproto $(PROTO_DIR)/audio.proto

proto-py-gen:
	$(PYTHON_VENV)/Scripts/python.exe -m grpc_tools.protoc -I proto --python_out=$(PYTHON_OUT) --pyi_out=$(PYTHON_OUT) --grpc_python_out=$(GRPC_PYTHON_OUT) $(PROTO_DIR)/audio.proto