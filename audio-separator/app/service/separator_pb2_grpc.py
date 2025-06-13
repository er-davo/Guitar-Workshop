# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

import separator_pb2 as separator__pb2

GRPC_GENERATED_VERSION = '1.73.0'
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in separator_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
    )


class AudioSeparatorStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.SeparateAudio = channel.unary_unary(
                '/separator.AudioSeparator/SeparateAudio',
                request_serializer=separator__pb2.SeparateRequest.SerializeToString,
                response_deserializer=separator__pb2.SeparateResponse.FromString,
                _registered_method=True)


class AudioSeparatorServicer(object):
    """Missing associated documentation comment in .proto file."""

    def SeparateAudio(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_AudioSeparatorServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'SeparateAudio': grpc.unary_unary_rpc_method_handler(
                    servicer.SeparateAudio,
                    request_deserializer=separator__pb2.SeparateRequest.FromString,
                    response_serializer=separator__pb2.SeparateResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'separator.AudioSeparator', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('separator.AudioSeparator', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class AudioSeparator(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def SeparateAudio(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/separator.AudioSeparator/SeparateAudio',
            separator__pb2.SeparateRequest.SerializeToString,
            separator__pb2.SeparateResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
