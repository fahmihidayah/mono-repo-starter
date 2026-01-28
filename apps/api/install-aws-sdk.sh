#!/bin/bash

# Install AWS SDK v1 Dependencies
# This script installs the required AWS SDK package for the storage implementation

echo "Installing AWS SDK (aws-sdk-go) dependencies..."
echo ""

# AWS SDK package
echo "Installing aws-sdk-go..."
go get github.com/aws/aws-sdk-go

echo ""
echo "Tidying up go.mod and go.sum..."
go mod tidy

echo ""
echo "✅ AWS SDK installation complete!"
echo ""
echo "You can now use S3 storage by setting STORAGE_TYPE=s3 in your .env file"
