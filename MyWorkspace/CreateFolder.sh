#!/bin/bash

BASE_PATH="/Users/bernatdsitumeang/Desktop/workspace-Bernatdev"

echo "========================================"
echo "   TICKET FOLDER CREATOR"
echo "========================================"
echo ""
read -p "Masukkan nama folder tiket: " FOLDER_NAME

if [ -z "$FOLDER_NAME" ]; then
    echo "Error: Nama folder tidak boleh kosong!"
    exit 1
fi

TICKET_PATH="$BASE_PATH/$FOLDER_NAME"

# Buat semua subfolder di dalamnya
mkdir -p "$TICKET_PATH/FSD"
mkdir -p "$TICKET_PATH/Analysis"
mkdir -p "$TICKET_PATH/Document"
mkdir -p "$TICKET_PATH/Document DONE"
mkdir -p "$TICKET_PATH/Revisi"
mkdir -p "$TICKET_PATH/Comparison"
mkdir -p "$TICKET_PATH/BAST"

echo ""
echo "========================================"
echo "✅ SUCCESS! Folder structure created:"
echo "========================================"
echo "Location: $TICKET_PATH"
echo ""
echo "Folder structure:"
ls -la "$TICKET_PATH"
echo ""
echo "========================================"

open "$TICKET_PATH"