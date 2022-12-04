package eelchat

//go:generate ogen --clean --no-client --package spec --target internal/api/spec -v spec/swagger.yaml
//go:generate ogen --clean --no-server --package api --target client/internal/api -v spec/swagger.yaml
//go:generate sqlboiler mysql -o internal/db/model -p model --no-tests
