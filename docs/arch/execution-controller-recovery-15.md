# Execution-controller attempt recovery

Missing-Job recovery keeps an outcome already explained by recorded facts (a Job deadline, OOM kill, eviction, stop request, or structured failure). It adds job_missing only when the attempt otherwise classifies as node_lost. In particular, an observed deadline remains timeout and does not become retryable infrastructure failure.
