# Execution-controller attempt recovery

The attempt observer reads a seed preflight result independently of the aggregate pod-log lookup. A failed pod-list/log request does not suppress a later successful named-container read. Its findings and refusal reason still feed classification and the task status report; a main-log lookup failure remains a soft failure.
