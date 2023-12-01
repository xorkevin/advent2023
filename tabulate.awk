BEGIN {
  kind = "unknown"
}

{ print $0; fflush() }

$1 == "prg" {
  kind = $2
}

$1 == "time" {
  totals[kind] += $2
  print "running total", kind, totals[kind] " ms"; fflush()
}

END {
  n = asorti(totals, keys)
  for (i = 1; i <= n; i++) {
      print "total", keys[i], totals[keys[i]] " ms"
  }
}
