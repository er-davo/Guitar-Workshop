#!/bin/sh
# minio-init.sh
sleep 5

for bucket in audio tab; do
  mc mb http://$MINIO_ROOT_USER:$MINIO_ROOT_PASSWORD@minio:9000/$bucket || echo "Bucket $bucket already exists"
done

echo "MinIO buckets ready: audio, tab"