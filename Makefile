default:
	$(MAKE) deps
	$(MAKE) all
deps:
	bash -c "./scripts/deps.sh"
test:
	bash -c "./scripts/test.sh $(TEST_TYPE)"
test-legacy:
	bash -c "./scripts/test.sh legacy"
test-small:
	bash -c "./scripts/test.sh small"
test-medium:
	bash -c "./scripts/test.sh medium"
test-large:
	bash -c "./scripts/test.sh large"
test-all:
	$(MAKE) test-small
	$(MAKE) test-medium
	$(MAKE) test-large
check:
	$(MAKE) test
all:
	bash -c "./scripts/build.sh"