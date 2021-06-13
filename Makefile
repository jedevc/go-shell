.PHONY: all
all:

.PHONY: checkout uncheckout
checkout: chapter0-checkout chapter1-checkout chapter2-checkout chapter3-checkout chapter4-checkout chapter5-checkout chapter6-checkout chapter7-checkout
uncheckout: chapter0-uncheckout chapter1-uncheckout chapter2-uncheckout chapter3-uncheckout chapter4-uncheckout chapter5-uncheckout chapter6-uncheckout chapter7-uncheckout

.PHONY: %-checkout %-uncheckout
%-checkout:
	git worktree add "./$*" "$*"
%-uncheckout:
	git worktree remove "./$*"
