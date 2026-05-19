#!/bin/bash

APP="controlDB"
OUTDIR="build"

mkdir -p "$OUTDIR"

echo "Compilando para Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o "$OUTDIR/${APP}" .
echo "✔ $OUTDIR/${APP}"

echo "Compilando para Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o "$OUTDIR/${APP}.exe" .
echo "✔ $OUTDIR/${APP}.exe"

echo ""
echo "Binarios generados en ./$OUTDIR/"