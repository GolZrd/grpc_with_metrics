
generate:
	make generate-note-api

generate-note-api:
	mkdir pkg\note_v1
	protoc --proto_path api/note_v1 \
	--go_out=pkg/note_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/note_v1 --go-grpc_opt=paths=source_relative \
	api\note_v1\note.proto

grpc-load-test:
	ghz \
		--proto api/note_v1/note.proto \
		--call note_v1.NoteV1.Get \
		--data '{"id": 1}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051

grpc-error-load-test:
	ghz \
		--proto api/note_v1/note.proto \
		--call note_v1.NoteV1.Get \
		--data '{"id": 0}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051