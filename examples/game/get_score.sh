#!/bin/sh
# Query the score nodes in the main game world
../../fishdb console -exec "get score"
# Query the conf node in the main game world
../../fishdb console -exec "get conf"
