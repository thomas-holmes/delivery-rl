#!/usr/bin/env bash

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$(pwd)/linux

exec ./delivery-rl
