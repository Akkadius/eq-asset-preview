#----------------------
# Parse makefile arguments
#----------------------
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

#----------------------
# Silence GNU Make
#----------------------
ifndef VERBOSE
MAKEFLAGS += --no-print-directory
endif

#----------------------
# Terminal
#----------------------

GREEN  := $(shell tput -Txterm setaf 2)
WHITE  := $(shell tput -Txterm setaf 7)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

#------------------------------------------------------------------
# - Add the following 'help' target to your Makefile
# - Add help text after each target name starting with '\#\#'
# - A category can be added with @category
#------------------------------------------------------------------

HELP_FUN = \
	%help; \
	while(<>) { \
		push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^([a-zA-Z\-]+)\s*:.*\#\#(?:@([a-zA-Z\-]+))?\s(.*)$$/ }; \
		print "-----------------------------------------\n"; \
		print "| Makefile menu\n"; \
		print "-----------------------------------------\n"; \
		print "| usage: make [command]\n"; \
		print "-----------------------------------------\n\n"; \
		for (sort keys %help) { \
			print "${WHITE}$$_:${RESET \
		}\n"; \
		for (@{$$help{$$_}}) { \
			$$sep = " " x (32 - length $$_->[0]); \
			print "  ${YELLOW}$$_->[0]${RESET}$$sep${GREEN}$$_->[1]${RESET}\n"; \
		}; \
		print "\n"; \
	}

help: ##@other Show this help.
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

#----------------------
# build
#----------------------

build-maps: ##@build Builds JSON file maps
	tree ./assets/objects -f -J --sort=name | jq -c > ./maps/objects-map.json
	tree ./assets/npc_models -f -J --sort=name | jq -c > ./maps/npc-models-map.json
	tree ./assets/monograms -f -J --sort=name | jq -c > ./maps/monograms-map.json
	tree ./assets/item_icons -f -J --sort=name | jq -c > ./maps/item-icons-map.json
	tree ./assets/spell_icons -f -J --sort=name | jq -c > ./maps/spell-icons-map.json
	tree ./assets/spell_animations -f -J --sort=name | jq -c > ./maps/spell-animations-map.json
	tree ./assets/expansion-icons-small -f -J --sort=name | jq -c > ./maps/expansion-icons-small-map.json
	@echo "Built maps!"
