# Storage Configuration Guide

This application supports two storage backends for media uploads: **Local File System** and **AWS S3** (including S3-compatible services).

## Table of Contents

- [Overview](#overview)
- [Local Storage](#local-storage)
- [S3 Storage](#s3-storage)
- [Switching Between Storage Types](#switching-between-storage-types)
- [Implementation Details](#implementation-details)

---

## Overview

The storage system is abstracted through a `Storage` interface, allowing you to easily switch between different storage backends without changing your application code.

### Supported Storage Types

1. **Local** - Files are stored on the local file system
2. **S3** - Files are stored on AWS S3 or S3-compatible services (MinIO, DigitalOcean Spaces, etc.)

---

## Local Storage

### Configuration

Add to your `.env` file:

```env
STORAGE_TYPE=local
LOCAL_UPLOAD_DIR=./uploads
BASE_URL=http://localhost:8080
```

### How It Works

- Files are uploaded to the directory specified in `LOCAL_UPLOAD_DIR`
- Files are accessible via `BASE_URL/uploads/filename`
- File naming: `{timestamp}_{slugified-name}.{ext}`

### Advantages

- ✅ Simple setup - no external dependencies
- ✅ No additional costs
- ✅ Fast for local development
- ✅ No network latency

### Disadvantages

- ❌ Not suitable for multi-server deployments
- ❌ No automatic backups
- ❌ Limited scalability
- ❌ Files lost if server is recreated

### Best For

- Local development
- Single-server deployments
- Small-scale applications
- Testing

---

## S3 Storage

### Prerequisites

1. AWS Account (or compatible service)
2. S3 Bucket created
3. IAM credentials with S3 access

### AWS S3 Configuration

Add to your `.env` file:

```env
STORAGE_TYPE=s3
AWS_S3_REGION=us-east-1
AWS_S3_BUCKET=your-bucket-name
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
```

### S3-Compatible Services

#### MinIO (Self-hosted S3-compatible)

```env
STORAGE_TYPE=s3
AWS_S3_REGION=us-east-1
AWS_S3_BUCKET=your-bucket-name
AWS_ACCESS_KEY_ID=minioadmin
AWS_SECRET_ACCESS_KEY=minioadmin
AWS_S3_ENDPOINT=http://localhost:9000
```

#### DigitalOcean Spaces

```env
STORAGE_TYPE=s3
AWS_S3_REGION=nyc3
AWS_S3_BUCKET=your-space-name
AWS_ACCESS_KEY_ID=your-spaces-key
AWS_SECRET_ACCESS_KEY=your-spaces-secret
AWS_S3_ENDPOINT=https://nyc3.digitaloceanspaces.com
```

#### Wasabi

```env
STORAGE_TYPE=s3
AWS_S3_REGION=us-east-1
AWS_S3_BUCKET=your-bucket-name
AWS_ACCESS_KEY_ID=your-wasabi-key
AWS_SECRET_ACCESS_KEY=your-wasabi-secret
AWS_S3_ENDPOINT=https://s3.wasabisys.com
```

### IAM Policy for AWS S3

Create an IAM user with this policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket",
        "s3:PutObjectAcl"
      ],
      "Resource": [
        "arn:aws:s3:::your-bucket-name",
        "arn:aws:s3:::your-bucket-name/*"
      ]
    }
  ]
}
```

### Bucket Configuration

1. **Create Bucket** in your AWS region
2. **Enable Public Access** for uploaded files:
   - Disable "Block all public access"
   - Or use bucket policy to allow public read

3. **Bucket Policy** (Optional - for public read):

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::your-bucket-name/*"
    }
  ]
}
```

### Advantages

- ✅ Highly scalable
- ✅ Automatic backups and durability (99.999999999%)
- ✅ CDN integration available
- ✅ Works with multi-server deployments
- ✅ Geographically distributed
- ✅ No server disk space required

### Disadvantages

- ❌ Additional cost (pay per storage and bandwidth)
- ❌ Network latency for uploads
- ❌ Requires AWS account setup
- ❌ More complex configuration

### Best For

- Production deployments
- Multi-server applications
- Scalable applications
- Applications with high storage needs

---

## Switching Between Storage Types

### Step 1: Update Environment Variable

```bash
# Switch to local
STORAGE_TYPE=local

# Switch to S3
STORAGE_TYPE=s3
```

### Step 2: Restart Application

```bash
# The application will automatically use the configured storage type
go run cmd/api/main.go
```

### Step 3: Migrate Existing Files (Optional)

If you're switching from local to S3 (or vice versa), you'll need to manually migrate existing files:

```bash
# Example: Upload local files to S3
aws s3 sync ./uploads s3://your-bucket-name/uploads/
```

---

## Implementation Details

### Architecture

```
┌─────────────────┐
│ Media Service   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Storage Interface│
└────────┬────────┘
         │
    ┌────┴─────┐
    │          │
    ▼          ▼
┌─────────┐  ┌──────────┐
│  Local  │  │    S3    │
│ Storage │  │ Storage  │
└─────────┘  └──────────┘
```

### Storage Interface

```go
type Storage interface {
    Upload(file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error)
    Delete(path string) error
    GetURL(path string) string
    Exists(path string) (bool, error)
}
```

### Factory Pattern

The application uses a factory pattern to create the appropriate storage instance:

```go
func (f *Factory) CreateStorage() (Storage, error) {
    switch f.config.Storage.Type {
    case "local":
        return NewLocalStorage(uploadDir, baseURL), nil
    case "s3":
        return NewS3Storage(s3Config)
    default:
        return nil, fmt.Errorf("unsupported storage type: %s", storageType)
    }
}
```

### File Structure

```
internal/storage/
├── storage.go          # Storage interface definition
├── local_storage.go    # Local file system implementation
├── s3_storage.go       # AWS S3 implementation
└── factory.go          # Factory for creating storage instances
```

---

## Testing

### Test Local Storage

```bash
# Upload a file
curl -X POST http://localhost:8080/api/media \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@image.jpg" \
  -F "alt=Test Image"

# File should be in ./uploads/
ls -la ./uploads/
```

### Test S3 Storage

```bash
# Upload a file
curl -X POST http://localhost:8080/api/media \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@image.jpg" \
  -F "alt=Test Image"

# Check S3 bucket
aws s3 ls s3://your-bucket-name/uploads/
```

---

## Troubleshooting

### Local Storage Issues

**Problem**: Permission denied when creating upload directory

```bash
# Solution: Set proper permissions
mkdir -p ./uploads
chmod 755 ./uploads
```

**Problem**: Files not accessible via URL

```bash
# Solution: Ensure static file serving is configured
# Check that /uploads/* route serves files from LOCAL_UPLOAD_DIR
```

### S3 Storage Issues

**Problem**: Access Denied error

```bash
# Solution: Check IAM permissions
# Ensure user has s3:PutObject, s3:GetObject, s3:DeleteObject permissions
```

**Problem**: Files uploaded but not publicly accessible

```bash
# Solution: Set ACL to public-read or configure bucket policy
# The code sets ACL: "public-read" during upload
```

**Problem**: Connection timeout

```bash
# Solution: Check AWS credentials and region
# Ensure AWS_S3_REGION matches your bucket region
```

**Problem**: Invalid endpoint

```bash
# Solution: For custom endpoints, include protocol
# Correct: AWS_S3_ENDPOINT=https://nyc3.digitaloceanspaces.com
# Wrong: AWS_S3_ENDPOINT=nyc3.digitaloceanspaces.com
```

---

## Cost Considerations

### Local Storage
- **Storage**: Cost of server disk space
- **Bandwidth**: Included in server costs
- **Typical Cost**: $5-20/month (included in VPS)

### AWS S3
- **Storage**: ~$0.023/GB/month (first 50TB)
- **Data Transfer**: ~$0.09/GB (first 10TB)
- **Requests**: ~$0.005 per 1,000 PUT requests
- **Typical Cost**: $1-100/month depending on usage

### DigitalOcean Spaces
- **Storage**: $5/month for 250GB
- **Bandwidth**: 1TB included, then $0.01/GB
- **Typical Cost**: $5-15/month

---

## Security Best Practices

1. **Never commit `.env` file** with real credentials
2. **Use IAM roles** on EC2 instead of access keys when possible
3. **Enable bucket versioning** for S3 to prevent accidental deletions
4. **Use HTTPS** for S3 endpoints
5. **Implement file type validation** before upload
6. **Limit file sizes** to prevent abuse
7. **Scan files** for malware before storing (in production)
8. **Use signed URLs** for private files

---

## Migration Guide

### Local → S3

1. Set up S3 bucket and credentials
2. Update `.env` with S3 configuration
3. Upload existing files:
   ```bash
   aws s3 sync ./uploads s3://your-bucket-name/uploads/
   ```
4. Change `STORAGE_TYPE=s3`
5. Restart application
6. Update database URLs (if storing full URLs):
   ```sql
   UPDATE media
   SET url = REPLACE(url, 'http://localhost:8080', 'https://bucket.s3.region.amazonaws.com');
   ```

### S3 → Local

1. Download existing files:
   ```bash
   aws s3 sync s3://your-bucket-name/uploads/ ./uploads/
   ```
2. Update `.env`: `STORAGE_TYPE=local`
3. Restart application
4. Update database URLs (if needed)

---

## FAQ

**Q: Can I use both local and S3 at the same time?**
A: No, the application uses one storage type at a time. You can switch by changing the environment variable.

**Q: What happens to existing files when I switch storage?**
A: Existing files remain in the old storage. You need to manually migrate them.

**Q: Does this support Azure Blob Storage or Google Cloud Storage?**
A: Not currently, but you can implement the `Storage` interface for those services.

**Q: Can I use different storage for different file types?**
A: Not out of the box, but you can modify the code to use different storage instances based on file type.

**Q: Is there a file size limit?**
A: This depends on your configuration. S3 supports files up to 5TB. Local storage is limited by disk space.

---

## Install AWS SDK Dependencies

If you haven't already, install the required AWS SDK package:

```bash
go get github.com/aws/aws-sdk-go
go mod tidy
```

---

For more information, see the [main README](../README.md).
