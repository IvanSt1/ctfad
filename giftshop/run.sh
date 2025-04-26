#!/bin/sh
socat tcp-listen:7490,reuseaddr,fork exec:"./giftshop"
