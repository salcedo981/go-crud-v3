takbo:
	go get -u ./...
	go run main.go

.PHONY: push-patch-version
push-patch-version:
	@LATEST_TAG=$$(git tag --sort=v:refname | grep -E '^[0-9]+\.[0-9]+\.[0-9]+$$' | sort -V | tail -n 1); \
	echo "Latest tag: $$LATEST_TAG"; \
	VERSION_PARTS=$$(echo $$LATEST_TAG | sed 's/^//' | tr '.' ' '); \
	MAJOR=$$(echo $$VERSION_PARTS | cut -d' ' -f1); \
	MINOR=$$(echo $$VERSION_PARTS | cut -d' ' -f2); \
	PATCH=$$(echo $$VERSION_PARTS | cut -d' ' -f3); \
	NEW_PATCH=$$((PATCH + 1)); \
	NEW_TAG="v$$MAJOR.$$MINOR.$$NEW_PATCH"; \
	echo "New tag: $$NEW_TAG"; \
	git tag $$NEW_TAG; \
	git push origin $$NEW_TAG

.PHONY: push-minor-version
push-minor-version:
	@LATEST_TAG=$$(git tag --sort=v:refname | grep -E '^[0-9]+\.[0-9]+\.[0-9]+$$' | sort -V | tail -n 1); \
	echo "Latest tag: $$LATEST_TAG"; \
	VERSION_PARTS=$$(echo $$LATEST_TAG | sed 's/^//' | tr '.' ' '); \
	MAJOR=$$(echo $$VERSION_PARTS | cut -d' ' -f1); \
	MINOR=$$(echo $$VERSION_PARTS | cut -d' ' -f2); \
	PATCH=$$(echo $$VERSION_PARTS | cut -d' ' -f3); \
	NEW_MINOR=$$((MINOR + 1)); \
	NEW_TAG="$$MAJOR.$$NEW_MINOR.$$PATCH"; \
	echo "New tag: $$NEW_TAG"; \
	git tag $$NEW_TAG; \
	git push origin $$NEW_TAG

.PHONY: push-major-version
push-major-version:
	@LATEST_TAG=$$(git tag --sort=v:refname | grep -E '^[0-9]+\.[0-9]+\.[0-9]+$$' | sort -V | tail -n 1); \
	echo "Latest tag: $$LATEST_TAG"; \
	VERSION_PARTS=$$(echo $$LATEST_TAG | sed 's/^//' | tr '.' ' '); \
	MAJOR=$$(echo $$VERSION_PARTS | cut -d' ' -f1); \
	MINOR=$$(echo $$VERSION_PARTS | cut -d' ' -f2); \
	PATCH=$$(echo $$VERSION_PARTS | cut -d' ' -f3); \
	NEW_MAJOR=$$((major + 1)); \
	NEW_TAG="$$NEW_MAJOR.$$MINOR.$$PATCH"; \
	echo "New tag: $$NEW_TAG"; \
	git tag $$NEW_TAG; \
	git push origin $$NEW_TAG