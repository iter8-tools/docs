
FROM python:3.9-alpine

# # Install protoc
# apk update && apk --no-cache add protobuf

# Get source
WORKDIR /
COPY samples/abn-sample/frontend/python .
RUN pip install -r requirements.txt

# # generate gRPC libraries
# COPY samples/abn-sample/frontend/abn.proto .
# RUN python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. ./abn.proto 

ENV FLASK_APP main.py
ENTRYPOINT ["flask", "run", "--port", "8090"]