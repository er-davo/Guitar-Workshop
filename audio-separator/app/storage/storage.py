import asyncio
from typing import BinaryIO
from urllib.parse import urljoin
from botocore.exceptions import ClientError
import aiobotocore.session


class Storage:
    def __init__(self, endpoint: str, access_key: str, secret_key: str, region_name: str = "us-east-1"):
        self.endpoint = endpoint
        self.access_key = access_key
        self.secret_key = secret_key
        self.region_name = region_name
        self.session = aiobotocore.session.get_session()
        self.client = aiobotocore.session.ClientCreatorContext

    async def __aenter__(self):
        self.client = await self.session.create_client(
            "s3",
            endpoint_url=self.endpoint,
            aws_access_key_id=self.access_key,
            aws_secret_access_key=self.secret_key,
            region_name=self.region_name
        ).__aenter__()
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        await self.client.__aexit__(exc_type, exc_val, exc_tb)

    async def upload_file(self, bucket_name: str, object_name: str, data: BinaryIO, length: int):
        await self.client.put_object(Bucket=bucket_name, Key=object_name, Body=data, ContentLength=length)

    async def get_file(self, bucket_name: str, object_name: str):
        return await self.client.get_object(Bucket=bucket_name, Key=object_name)

    async def remove_file(self, bucket_name: str, object_name: str):
        await self.client.delete_object(Bucket=bucket_name, Key=object_name)

    async def list_files_by_prefix(self, bucket_name: str, prefix: str):
        paginator = self.client.get_paginator("list_objects_v2")
        result = []
        async for page in paginator.paginate(Bucket=bucket_name, Prefix=prefix):
            for obj in page.get("Contents", []):
                result.append(obj)
        return result

    async def presigned_get(self, bucket_name: str, object_name: str, expires_in: int = 3600):
        return await self.client.generate_presigned_url(
            "get_object",
            Params={"Bucket": bucket_name, "Key": object_name},
            ExpiresIn=expires_in
        )
