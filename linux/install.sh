#!/bin/sh

BASEDIR=$(dirname "$0")

mkdir /mineralos
cp $BASEDIR/mineralos/* /mineralos

ln -s /mineralos/bin/* /usr/bin/