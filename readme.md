# Install
Clone the repo then 
`git submodule update --init --recursive`

# Build
`go build .`

# Run
First make sure you have the `FULLNODE_API_INFO` variable in your environment, in multiaddr form, for example:
`FULLNODE_API_INFO="eyJXXX:/ip4/127.0.0.1/tcp/1000/http"`

Then run the command

`filecoin-retriever <cid> <sp id>`

