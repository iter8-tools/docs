FROM node:19-buster
WORKDIR /usr/app

# Get source (including requirements)
COPY . ./

# Generate gRPC libraries from abn.proto
#COPY samples/abn-sample/frontend/abn.proto .
#RUN npm install -g grpc-tools
#RUN grpc_tools_node_protoc --js_out=import_style=commonjs,binary:. --grpc_out=. --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` abn.proto

# Install requirements
RUN npm install

CMD ["node", "main.js"]
