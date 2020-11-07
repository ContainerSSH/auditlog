package s3

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (q *uploadQueue) initializeMultiPartUpload(s3Connection *s3.S3, name string, metadata queueEntryMetadata) (*string, error) {
	q.logger.Debugf("initializing multipart upload for audit log %s...", name)
	multipartUpload, err := s3Connection.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		ACL:         q.acl,
		Bucket:      aws.String(q.bucket),
		ContentType: aws.String("application/octet-stream"),
		Key:         aws.String(name),
		Metadata:    metadata.ToMap(q.metadataUsername, q.metadataIP),
	})
	if err != nil {
		q.logger.Warningf("failed to initialize audit log file upload %s (%w)", name, err)
		return nil, err
	}
	return multipartUpload.UploadId, nil
}

func (q *uploadQueue) processMultiPartUploadPart(
	s3Connection *s3.S3,
	name string,
	uploadID string,
	partNumber int64,
	handle *os.File,
	startingByte int64,
	endingByte int64,
) (int64, string, error) {
	q.logger.Debugf("uploading part %d of audit log %s (part size %d bytes)...", partNumber, name, endingByte-startingByte)
	contentLength := endingByte - startingByte
	response, err := s3Connection.UploadPart(&s3.UploadPartInput{
		Body:          io.NewSectionReader(handle, startingByte, contentLength),
		Bucket:        aws.String(q.bucket),
		ContentLength: aws.Int64(contentLength),
		Key:           aws.String(name),
		PartNumber:    aws.Int64(partNumber),
		UploadId:      aws.String(uploadID),
	})
	etag := ""
	if err != nil {
		q.logger.Warningf("failed to upload part %d of audit log %s (%w)", partNumber, name, err)
		return 0, "", fmt.Errorf("failed to upload part %d of audit log file %s (%w)", partNumber, name, err)
	}
	etag = *response.ETag
	q.logger.Debugf("completed upload of part %d of audit log %s", partNumber, name)
	return contentLength, etag, nil
}

func (q *uploadQueue) processSingleUpload(s3Connection *s3.S3, name string, handle *os.File, metadata queueEntryMetadata) (int64, error) {
	q.logger.Debugf("processing single upload for audit log %s...", name)
	stat, err := handle.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to upload audit log %s (%w)", name, err)
	}
	contentLength := stat.Size()
	_, err = s3Connection.PutObject(&s3.PutObjectInput{
		ACL:           q.acl,
		Body:          handle,
		Bucket:        aws.String(q.bucket),
		ContentLength: aws.Int64(contentLength),
		ContentType:   aws.String("application/octet-stream"),
		Key:           aws.String(name),
		Metadata:      metadata.ToMap(q.metadataUsername, q.metadataIP),
	})
	if err != nil {
		q.logger.Debugf("single upload failed for audit log %s (%w)", name, err)
	} else {
		q.logger.Debugf("single upload complete for audit log %s", name)
	}
	return contentLength, err
}

func (q *uploadQueue) finalizeUpload(s3Connection *s3.S3, name string, uploadID string, completedParts []*s3.CompletedPart) error {
	q.logger.Debugf("finalizing multipart upload for audit log %s...", name)
	_, err := s3Connection.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket: aws.String(q.bucket),
		Key:    aws.String(name),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
		UploadId: aws.String(uploadID),

		ExpectedBucketOwner: nil,
		RequestPayer:        nil,
	})
	if err != nil {
		q.logger.Warningf("finalizing multipart upload failed for audit log %s (%w)", name, err)
	} else {
		q.logger.Debugf("finalizing multipart upload complete for audit log %s", name)
	}
	return err
}

func (q *uploadQueue) processShouldAbort(s3Connection *s3.S3, name string, failures int, uploadID *string) bool {
	abort := func() {
		if uploadID != nil {
			if err := q.abortSpecificMultipartUpload(name, s3Connection, &name, uploadID); err != nil {
				q.logger.Warningf("failed to abort multipart upload for %s (%v)", err)
			}
		}
		q.queue.Delete(name)
	}
	if failures > 20 {
		q.logger.Warningf("failed to upload audit log %s for 20 times in a row, giving up", failures)
		abort()
		return true
	}
	if failures > 3 {
		select {
		case <-q.ctx.Done():
			q.logger.Warningf("failed to upload audit log %s 3 times and shutdown is requested, giving up", failures)
			abort()
			return true
		default:
		}
	}
	return false
}

func (q *uploadQueue) uploadLoop(s3Connection *s3.S3, name string, entry *queueEntry) {
	defer q.wg.Done()
	var uploadID *string = nil
	uploadedBytes := int64(0)
	var errorHappened bool
	var completedParts []*s3.CompletedPart
	failures := 0
	for {
		if q.processShouldAbort(s3Connection, name, failures, uploadID) {
			break
		}

		entry.waitPartAvailable()
		q.workerSem <- 42
		errorHappened = false

		stat, err := entry.readHandle.Stat()
		if err != nil {
			q.logger.Warningf("failed to stat audit queue file %s before upload (%w)", name, err)
			errorHappened = true
		}

		if !errorHappened {
			var finished bool
			errorHappened, finished, uploadedBytes, completedParts, uploadID = q.processUpload(
				entry,
				uploadedBytes,
				s3Connection,
				name,
				stat.Size()-uploadedBytes,
				uploadID,
				stat,
				completedParts,
			)
			if finished {
				<-q.workerSem
				break
			}
		}

		<-q.workerSem
		if errorHappened || entry.finished {
			// If an error happened, retry immediately.
			// Also go back if the entry is finished to finish uploading the parts.
			entry.markPartAvailable()
		}
		if errorHappened {
			failures++
			time.Sleep(10 * time.Second)
		} else {
			failures = 0
		}
	}
}

func (q *uploadQueue) upload(name string) error {
	rawEntry, ok := q.queue.Load(name)
	if !ok {
		return fmt.Errorf("no such queue entry: %s", name)
	}
	s3Connection := s3.New(q.awsSession)
	entry := rawEntry.(*queueEntry)
	q.wg.Add(1)
	go q.uploadLoop(s3Connection, name, entry)
	return nil
}

func (q *uploadQueue) processUpload(
	entry *queueEntry,
	uploadedBytes int64,
	s3Connection *s3.S3,
	name string,
	remainingBytes int64,
	uploadID *string,
	stat os.FileInfo,
	completedParts []*s3.CompletedPart,
) (bool, bool, int64, []*s3.CompletedPart, *string) {
	if entry.finished && uploadedBytes == 0 {
		// If the entry is finished and nothing has been uploaded yet, upload it as a single file.
		partBytes, err := q.processSingleUpload(s3Connection, name, entry.readHandle, entry.metadata)
		if err != nil {
			q.logger.Warningf("failed to upload audit log %s (%w)", name, err)
			return true, false, uploadedBytes, completedParts, uploadID
		}
		uploadedBytes = uploadedBytes + partBytes
	} else if (entry.finished && remainingBytes > 0) || remainingBytes >= int64(q.partSize) {
		// If the entry is finished and there are bytes remaining, upload. Otherwise, we only upload if
		// more than the part size is available.
		if uploadID == nil {
			var err error
			uploadID, err = q.initializeMultiPartUpload(s3Connection, name, entry.metadata)
			if err != nil {
				return true, false, uploadedBytes, completedParts, uploadID
			}
		}
		if uploadID != nil {
			uploadedBytes, completedParts = q.doMultipartUpload(entry, uploadedBytes, s3Connection, name, stat, uploadID, completedParts)
		}
	} else if entry.finished && remainingBytes == 0 {
		//If the entry is finished and no data is left to be uploaded, finalize the upload.
		if uploadID != nil {
			err := q.finalizeUpload(s3Connection, name, *uploadID, completedParts)
			if err != nil {
				return true, false, uploadedBytes, completedParts, uploadID
			}
		}
		if err := entry.remove(); err != nil {
			q.logger.Warningf("failed to remove queue entry (%w)", err)
		}
		q.queue.Delete(name)
		return false, true, uploadedBytes, completedParts, uploadID
	}
	return false, false, uploadedBytes, completedParts, uploadID
}

func (q *uploadQueue) doMultipartUpload(entry *queueEntry, uploadedBytes int64, s3Connection *s3.S3, name string, stat os.FileInfo, uploadID *string, completedParts []*s3.CompletedPart) (int64, []*s3.CompletedPart) {
	partNumber := uploadedBytes / int64(q.partSize)
	startingByte := partNumber * int64(q.partSize)
	endingByte := (partNumber + 1) * int64(q.partSize)
	if stat.Size() < endingByte {
		endingByte = stat.Size()
	}

	if entry.finished && stat.Size()-endingByte < int64(q.partSize) {
		endingByte = stat.Size()
	}

	partBytes, etag, err := q.processMultiPartUploadPart(s3Connection, name, *uploadID, partNumber, entry.readHandle, startingByte, endingByte)
	if err == nil {
		uploadedBytes = uploadedBytes + partBytes
		completedParts = append(completedParts, &s3.CompletedPart{
			ETag:       aws.String(etag),
			PartNumber: aws.Int64(partNumber),
		})
	}
	return uploadedBytes, completedParts
}

func (q *uploadQueue) abortMultiPartUpload(name string) error {
	s3Connection := s3.New(q.awsSession)
	multiPartUpload, err := s3Connection.ListMultipartUploads(&s3.ListMultipartUploadsInput{
		Bucket: aws.String(q.bucket),
		Prefix: aws.String(name),

		Delimiter:           nil,
		EncodingType:        nil,
		ExpectedBucketOwner: nil,
		KeyMarker:           nil,
		MaxUploads:          nil,
		UploadIdMarker:      nil,
	})
	if err != nil {
		return fmt.Errorf("failed to list existing multipart upload for audit log %s (%w)", name, err)
	}
	for _, upload := range multiPartUpload.Uploads {
		if *upload.Key == name {
			q.logger.Debugf("aborting previous multipart upload ID %s for audit log %s...", *(upload.UploadId), name)
			if err := q.abortSpecificMultipartUpload(name, s3Connection, upload.Key, upload.UploadId); err != nil {
				return err
			}
		}
	}
	return nil
}

func (q *uploadQueue) abortSpecificMultipartUpload(name string, s3Connection *s3.S3, key *string, uploadID *string) error {
	if _, err := s3Connection.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
		Bucket:   aws.String(q.bucket),
		Key:      key,
		UploadId: uploadID,
	}); err != nil {
		return fmt.Errorf("failed to abort  %s (%w)", name, err)
	}
	return nil
}
