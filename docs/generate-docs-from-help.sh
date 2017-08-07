#!/bin/bash

if [ ! -z ${DEBUG} ]; then
    set -x
fi

CMD_DOC_FILENAME=help.md

echo "# Java Buildpack Memory Calculator Help
\`\`\`" > $CMD_DOC_FILENAME

go install ..

java-buildpack-memory-calculator -help 2>&1 | grep -v exit\ status\ 2 1>>$CMD_DOC_FILENAME

echo "
\`\`\`" >> $CMD_DOC_FILENAME

echo "Print contents of $CMD_DOC_FILENAME"
echo "==================================="
cat $CMD_DOC_FILENAME