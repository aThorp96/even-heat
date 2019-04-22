#!/bin/sh
go build; ./even-heat > out; cat out animationPrototype.py > getAnimation.py; python3 getAnimation.py
