tw:
	@npx tailwindcss -i ./views/input.css -o ./static/assets/css/styles.css --watch

dev:
	@templ generate -watch -proxy="http://localhost:4000" -open-browser=false -cmd="go run main.go"