commit:
	@git add . && git commit -m "update" || true

push:
	@git push origin main
