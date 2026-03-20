.PHONY: tag-current tag-patch tag-push

# Last semver tag (falls back to v0.0.0 if none exists)
CURRENT_TAG := $(shell git tag --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' | head -1)
CURRENT_TAG := $(if $(CURRENT_TAG),$(CURRENT_TAG),v0.0.0)

# Split into major.minor.patch
TAG_MAJOR   := $(shell echo $(CURRENT_TAG) | sed 's/^v//' | cut -d. -f1)
TAG_MINOR   := $(shell echo $(CURRENT_TAG) | sed 's/^v//' | cut -d. -f2)
TAG_PATCH   := $(shell echo $(CURRENT_TAG) | sed 's/^v//' | cut -d. -f3)
NEXT_PATCH  := $(shell echo $$(($(TAG_PATCH) + 1)))
NEXT_TAG    := v$(TAG_MAJOR).$(TAG_MINOR).$(NEXT_PATCH)

## Show the current latest tag
tag-current:
	@echo $(CURRENT_TAG)

## Create a new patch tag (e.g. v0.2.1 -> v0.2.2)
tag-patch:
	@echo "Current tag: $(CURRENT_TAG)"
	@echo "New tag:     $(NEXT_TAG)"
	git tag $(NEXT_TAG)

## Push the latest tag to the remote
tag-push:
	@echo "Pushing $(NEXT_TAG) to remote..."
	git push origin $(NEXT_TAG)
