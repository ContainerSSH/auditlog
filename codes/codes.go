package codes

// ContainerSSH is starting a new S3 multipart upload.
const MMultipartUpload = "AUDIT_S3_MULTIPART_UPLOAD"

// ContainerSSH failed to initialize a new multipart upload to the S3-compatible object storage. Check if the S3
// configuration is correct and the provided S3 access key and secrets have permissions to start a multipart upload.
const EMultipartUploadInitializationFailed = "AUDIT_S3_MULTIPART_UPLOAD_INITIALIZATION_FAILED"

// ContainerSSH is uploading a part of an audit log to the S3-compatible object storage.
const MMultipartPartUploading = "AUDIT_S3_MULTIPART_PART_UPLOADING"

// ContainerSSH failed to upload a part to the S3-compatible object storage. Check the message for details.
const EMultipartPartUploadFailed = "AUDIT_S3_MULTIPART_PART_UPLOAD_FAILED"

// ContainerSSH completed the upload of an audit log part to the S3-compatible object storage.
const MMultipartPartUploadComplete = "AUDIT_S3_MULTIPART_PART_UPLOAD_COMPLETE"

// ContainerSSH has uploaded all audit log parts and is now finalizing the multipart upload.
const MMultipartUploadFinalizing = "AUDIT_S3_MULTIPART_UPLOAD_FINALIZING"

// ContainerSSH has uploaded all audit log parts, but could not finalize the multipart upload.
const EMultipartUploadFinalizationFailed = "AUDIT_S3_MULTIPART_UPLOAD_FINALIZATION_FAILED"

// ContainerSSH has uploaded all audit log parts and has successfully finalized the upload.
const EMultipartUploadFinalized = "AUDIT_S3_MULTIPART_UPLOAD_FINALIZED"

// ContainerSSH failed aborting a multipart upload from a previously crashed ContainerSSH run.
const EMultipartFailedAbort = "AUDIT_S3_MULTIPART_FAILED_ABORT"

// ContainerSSH failed to list multipart uploads on the object storage bucket. This is needed to abort uploads from a
// previously crashed ContainerSSH run.
const EMultipartFailedList = "AUDIT_S3_MULTIPART_FAILED_LIST"

// ContainerSSH is uploading the full audit log in a single upload to the S3-compatible object storage. This happens
// when the audit log size is below the minimum size for a multi-part upload.
const MSingleUpload = "AUDIT_S3_SINGLE_UPLOAD"

// ContainerSSH failed to upload the audit log as a single upload.
const MSingleUploadFailed = "AUDIT_S3_SINGLE_UPLOAD_FAILED"

// ContainerSSH successfully uploaded the audit log as a single upload.
const MSingleUploadComplete = "AUDIT_S3_SINGLE_UPLOAD_COMPLETE"

// ContainerSSH failed to create the metadata file for the S3 upload in the local temporary directory. Check if the
// local directory specified is writable and has enough disk space.
const EFailedCreatingMetadataFile = "AUDIT_S3_FAILED_CREATING_METADATA_FILE"

// ContainerSSH failed to read the metadata file for the S3 upload in the local temporary directory. Check if the
// local directory specified is readable and the files have not been corrupted.
const EFailedReadingMetadataFile = "AUDIT_S3_FAILED_READING_METADATA_FILE"

// ContainerSSH could not close the metadata file in the local folder. This typically happens when the local folder is
// on an NFS share. (This is NOT supported.)
const ECannotCloseMetadataFileHandle = "AUDIT_S3_CANNOT_CLOSE_METADATA_FILE_HANDLE"

// ContainerSSH failed to encode the metadata file. This is a bug, please report it.
const EFailedMetadataJSONEncoding = "AUDIT_S3_FAILED_METADATA_JSON_ENCODING"

// ContainerSSH failed to write the local metadata file. Please check if your disk has enough disk space.
const EFailedWritingMetadataFile = "AUDIT_S3_FAILED_WRITING_METADATA_FILE"

// ContainerSSH failed to stat the queue file. This usually happens when the local directory is being manually
// manipulated.
const EFailedQueueStat = "AUDIT_S3_FAILED_STAT_QUEUE_ENTRY"

// ContainerSSH was trying to upload an audit log from the metadata file, but the audit log does not exist.
const ENoSuchQueueEntry = "AUDIT_S3_NO_SUCH_QUEUE_ENTRY"

// ContainerSSH failed to remove an uploaded audit log from the local directory. This usually happens on Windows when a
// different process has the audit log open. (This is not a supported setup.)
const ERemoveAuditLogFailed = "AUDIT_S3_REMOVE_FAILED"

// ContainerSSH failed to close an audit log file in the local directory. This usually happens when the local directory
// is on an NFS share. (This is NOT supported.)
const ECloseAuditLogFileFailed = "AUDIT_S3_CLOSE_FAILED"

// ContainerSSH is aborting a multipart upload. Check the log message for details.
const EMultipartAborting = "AUDIT_S3_MULTIPART_ABORTING"

// ContainerSSH found a previously aborted multipart upload locally and is now attempting to recover the upload.
const MRecovering = "AUDIT_S3_RECOVERING"

// ContainerSSH failed to close the audit log storage handler.
const EAuditLogStorageCloseFailed = "AUDIT_STORAGE_CLOSE_FAILED"
