# AWS SDK Installation Guide

## Error Message

If you see this error:
```
could not import github.com/aws/aws-sdk-go/service/s3
(no required module provides package "github.com/aws/aws-sdk-go/service/s3")
```

This means the AWS SDK packages haven't been installed yet.

---

## Quick Fix

Run the installation script:

```bash
chmod +x install-aws-sdk.sh
./install-aws-sdk.sh
```

**OR** manually install the packages (see below).

---

## Manual Installation

### Step 1: Install AWS SDK Package

Run this command in your terminal:

```bash
# Install AWS SDK for Go (v1)
go get github.com/aws/aws-sdk-go
```

### Step 2: Tidy Dependencies

```bash
go mod tidy
```

This will:
- Download all required dependencies
- Update `go.mod` with the new packages
- Update `go.sum` with checksums

---

## Verification

After installation, verify the package is installed:

```bash
# Check go.mod file
grep "aws-sdk-go" go.mod
```

You should see an entry like:
```
github.com/aws/aws-sdk-go v1.xx.x
```

### Build Test

Try building the project:

```bash
go build ./cmd/api/main.go
```

If successful, you should see no errors!

---

## What If Local Storage is Enough?

If you only want to use **local storage** and don't need S3, you have two options:

### Option 1: Install AWS SDK Anyway (Recommended)

Even if you're not using S3, installing the SDK won't hurt:
- The S3 code won't be used if `STORAGE_TYPE=local`
- No runtime overhead
- You can switch to S3 later without code changes

### Option 2: Comment Out S3 Import (Not Recommended)

If you really don't want to install AWS SDK, you can temporarily comment out the S3 storage:

**In `internal/storage/factory.go`:**
```go
func (f *Factory) CreateStorage() (Storage, error) {
    storageType := f.config.Storage.Type

    switch storageType {
    case "local":
        log.Printf("Initializing local storage...")
        return NewLocalStorage(f.config.Storage.LocalUploadDir, f.config.BaseURL), nil

    case "s3":
        // S3 storage temporarily disabled
        return nil, fmt.Errorf("S3 storage not available - AWS SDK not installed")

    default:
        return nil, fmt.Errorf("unsupported storage type: %s", storageType)
    }
}
```

**⚠️ Warning**: This approach means you can't switch to S3 without uncommenting the code and installing the SDK.

---

## Troubleshooting

### Error: "go: module github.com/aws/aws-sdk-go not found"

**Solution**: Make sure you have internet connection and try again:
```bash
go clean -modcache
go get github.com/aws/aws-sdk-go
```

### Error: "go.mod not found"

**Solution**: Make sure you're in the project root directory:
```bash
cd /path/to/blog-go
go get github.com/aws/aws-sdk-go
```

### Error: "permission denied"

**Solution**: Use sudo (Linux/Mac) or run as administrator (Windows):
```bash
# Mac/Linux
sudo chmod +x install-aws-sdk.sh
./install-aws-sdk.sh

# Or install package directly
go get github.com/aws/aws-sdk-go
```

### Still Having Issues?

1. **Check Go version**: Ensure you have Go 1.19 or higher
   ```bash
   go version
   ```

2. **Check Go environment**:
   ```bash
   go env
   ```

3. **Clear module cache**:
   ```bash
   go clean -modcache
   go mod download
   ```

4. **Verify network connectivity**:
   ```bash
   curl -I https://proxy.golang.org
   ```

---

## After Installation

Once the AWS SDK is installed, you can:

1. **Use Local Storage** (default):
   ```env
   STORAGE_TYPE=local
   ```

2. **Use S3 Storage**:
   ```env
   STORAGE_TYPE=s3
   AWS_S3_REGION=us-east-1
   AWS_S3_BUCKET=your-bucket-name
   AWS_ACCESS_KEY_ID=your-key
   AWS_SECRET_ACCESS_KEY=your-secret
   ```

3. **Switch between them** by changing `STORAGE_TYPE` in `.env`

---

## Package Versions

The storage implementation uses AWS SDK for Go (v1):

- `github.com/aws/aws-sdk-go` >= v1.44.0

The installation command will automatically install the most recent compatible version.

---

## Need More Help?

- **AWS SDK for Go Documentation**: https://docs.aws.amazon.com/sdk-for-go/
- **Go Modules Guide**: https://go.dev/ref/mod
- **Storage Implementation Guide**: See [STORAGE.md](./STORAGE.md)

---

Happy coding! 🚀
