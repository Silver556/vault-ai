secrets_path="/vault-ai/secret"
echo "DEBUG: secrets_path = ${secrets_path}"
test ! -d $secrets_path && echo "ERR: /vault-ai/secret dir missing!" && return 1


export GO111MODULE=on
export GOBIN="$PWD/bin"
export GOPATH="$HOME/go"
export PATH="$PATH:$PWD/bin:$PWD/tools/protoc-3.6.1/bin"
export DOCKER_BUILDKIT=1
export OPENAI_API_KEY="$(cat ${secrets_path}/openai_api_key)"
export PINECONE_API_KEY="$(cat ${secrets_path}/pinecone_api_key)"
export PINECONE_API_ENDPOINT="$(cat ${secrets_path}/pinecone_api_endpoint)"

echo "=> Environment Variables Loaded"
