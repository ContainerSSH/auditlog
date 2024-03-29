# Message / error codes

| Code | Explanation |
|------|-------------|
| `AUDIT_S3_CANNOT_CLOSE_METADATA_FILE_HANDLE` | ContainerSSH could not close the metadata file in the local folder. This typically happens when the local folder is on an NFS share. (This is NOT supported.) |
| `AUDIT_S3_CLOSE_FAILED` | ContainerSSH failed to close an audit log file in the local directory. This usually happens when the local directory is on an NFS share. (This is NOT supported.) |
| `AUDIT_S3_FAILED_CREATING_METADATA_FILE` | ContainerSSH failed to create the metadata file for the S3 upload in the local temporary directory. Check if the local directory specified is writable and has enough disk space. |
| `AUDIT_S3_FAILED_METADATA_JSON_ENCODING` | ContainerSSH failed to encode the metadata file. This is a bug, please report it. |
| `AUDIT_S3_FAILED_READING_METADATA_FILE` | ContainerSSH failed to read the metadata file for the S3 upload in the local temporary directory. Check if the local directory specified is readable and the files have not been corrupted. |
| `AUDIT_S3_FAILED_STAT_QUEUE_ENTRY` | ContainerSSH failed to stat the queue file. This usually happens when the local directory is being manually manipulated. |
| `AUDIT_S3_FAILED_WRITING_METADATA_FILE` | ContainerSSH failed to write the local metadata file. Please check if your disk has enough disk space. |
| `AUDIT_S3_MULTIPART_ABORTING` | ContainerSSH is aborting a multipart upload. Check the log message for details. |
| `AUDIT_S3_MULTIPART_FAILED_ABORT` | ContainerSSH failed aborting a multipart upload from a previously crashed ContainerSSH run. |
| `AUDIT_S3_MULTIPART_FAILED_LIST` | ContainerSSH failed to list multipart uploads on the object storage bucket. This is needed to abort uploads from a previously crashed ContainerSSH run. |
| `AUDIT_S3_MULTIPART_PART_UPLOADING` | ContainerSSH is uploading a part of an audit log to the S3-compatible object storage. |
| `AUDIT_S3_MULTIPART_PART_UPLOAD_COMPLETE` | ContainerSSH completed the upload of an audit log part to the S3-compatible object storage. |
| `AUDIT_S3_MULTIPART_PART_UPLOAD_FAILED` | ContainerSSH failed to upload a part to the S3-compatible object storage. Check the message for details. |
| `AUDIT_S3_MULTIPART_UPLOAD` | ContainerSSH is starting a new S3 multipart upload. |
| `AUDIT_S3_MULTIPART_UPLOAD_FINALIZATION_FAILED` | ContainerSSH has uploaded all audit log parts, but could not finalize the multipart upload. |
| `AUDIT_S3_MULTIPART_UPLOAD_FINALIZED` | ContainerSSH has uploaded all audit log parts and has successfully finalized the upload. |
| `AUDIT_S3_MULTIPART_UPLOAD_FINALIZING` | ContainerSSH has uploaded all audit log parts and is now finalizing the multipart upload. |
| `AUDIT_S3_MULTIPART_UPLOAD_INITIALIZATION_FAILED` | ContainerSSH failed to initialize a new multipart upload to the S3-compatible object storage. Check if the S3 configuration is correct and the provided S3 access key and secrets have permissions to start a multipart upload. |
| `AUDIT_S3_NO_SUCH_QUEUE_ENTRY` | ContainerSSH was trying to upload an audit log from the metadata file, but the audit log does not exist. |
| `AUDIT_S3_RECOVERING` | ContainerSSH found a previously aborted multipart upload locally and is now attempting to recover the upload. |
| `AUDIT_S3_REMOVE_FAILED` | ContainerSSH failed to remove an uploaded audit log from the local directory. This usually happens on Windows when a different process has the audit log open. (This is not a supported setup.) |
| `AUDIT_S3_SINGLE_UPLOAD` | ContainerSSH is uploading the full audit log in a single upload to the S3-compatible object storage. This happens when the audit log size is below the minimum size for a multi-part upload. |
| `AUDIT_S3_SINGLE_UPLOAD_COMPLETE` | ContainerSSH successfully uploaded the audit log as a single upload. |
| `AUDIT_S3_SINGLE_UPLOAD_FAILED` | ContainerSSH failed to upload the audit log as a single upload. |
| `AUDIT_STORAGE_CLOSE_FAILED` | ContainerSSH failed to close the audit log storage handler. |

