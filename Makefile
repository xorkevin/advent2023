.PHONY: bench

bench:
	find . -maxdepth 1 -type d -name 'day*' \
		| sort \
		| xargs -I{} sh -c 'cd {} && pwd && make bench' 2>/dev/null \
		| stdbuf -oL grep 'Benchmark\|Time' \
		| stdbuf -oL sed \
			-e 's#Benchmark.*bin/day#prg go day#' \
			-e 's#Benchmark.*release/day#prg rs day#' \
			-e 's#.*Time.*):\s\+#time #' \
			-e 's#\s\+ms.*# ms#' \
		| awk -f tabulate.awk
