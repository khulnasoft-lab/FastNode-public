#!/usr/bin/env bash

function fastnode_wrap {
	echo "[[FASTNODE[[$1]]FASTNODE]]"
}

function fastnode_show_region_delimiter {
	fastnode_wrap 'SHOW {"region": "'$FOO'", "type": "region"}'
}

function fastnode_line {
	fastnode_wrap "LINE $1"
}
