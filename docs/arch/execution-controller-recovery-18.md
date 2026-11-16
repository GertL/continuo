# Execution-controller attempt recovery

Job observations copy the termination reason of the first failed init container when the main container never terminated. An init OOMKilled is classified oom_killed/resource, while a plain init exit 137 remains exit_code/task. A terminated main container keeps its own termination facts.
