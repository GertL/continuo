# Execution-controller attempt recovery

The attempt observer applies the dbt Nothing to do log marker only outside the Python node family. Successful Python nodes are classified by their exit/result rather than arbitrary diagnostic text. Production and release-leg commands both supply their node type; all Python kinds share the generated runtime predicate.
