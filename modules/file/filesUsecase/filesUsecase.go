package filesusecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/file"
)

type IFileUsecase interface {
	UploadImageGCP(req []*file.FileReq) ([]*file.FileRes, error)
	UploadImageLocal(req []*file.FileReq) ([]*file.FileRes, error)
	DeleteImageGCP(req []*file.DeleteFileReq) error
	DeleteImageStorage(req []*file.DeleteFileReq) error
}

type FileUsecase struct {
	cfg config.IConfig
}

type filePublic struct {
	bucket      string
	destination string
	file        *file.FileRes
}

func NewFileUsecase(cfg config.IConfig) IFileUsecase {
	return &FileUsecase{
		cfg: cfg,
	}
}

func (f *filePublic) makePublic(ctx context.Context, client *storage.Client) error {
	acl := client.Bucket(f.bucket).Object(f.destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %w", err)
	}
	fmt.Printf("Blob %v is now publicly accessible.\n", f.destination)
	return nil
}

// Upload to GCP
func (u *FileUsecase) uploadWorker(ctx context.Context, client *storage.Client, jobs <-chan *file.FileReq, results chan<- *file.FileRes, errCh chan<- error) {

	for job := range jobs {

		container, err := job.File.Open()
		if err != nil {
			errCh <- err
			return
		}

		b, err := ioutil.ReadAll(container)
		if err != nil {
			errCh <- err
			return
		}

		buf := bytes.NewBuffer(b)

		// Upload an object with storage.Writer.
		wc := client.Bucket(u.cfg.App().GcpBucket()).Object(job.Destination).NewWriter(ctx)

		if _, err = io.Copy(wc, buf); err != nil {
			errCh <- fmt.Errorf("io.Copy: %w", err)
			return
		}
		// Data can continue to be added to the file until the writer is closed.
		if err := wc.Close(); err != nil {
			errCh <- fmt.Errorf("Writer.Close: %w", err)
			return
		}
		fmt.Printf("%v uploaded to %v.\n", job.FileName, job.Extension)

		newFile := &filePublic{
			file: &file.FileRes{
				FileName: job.FileName,
				Url:      fmt.Sprintf("https://storage.googleapis.com/%s%s", u.cfg.App().GcpBucket(), job.Destination),
			},
			bucket:      u.cfg.App().GcpBucket(),
			destination: job.Destination,
		}
		if err := newFile.makePublic(ctx, client); err != nil {
			errCh <- err
			return
		}

		errCh <- nil
		results <- newFile.file
	}

}

// Upload to GCP
func (u *FileUsecase) UploadImageGCP(req []*file.FileReq) ([]*file.FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	jobsCh := make(chan *file.FileReq, len(req))
	resultCh := make(chan *file.FileRes, len(req))
	errCh := make(chan error, len(req))

	result := make([]*file.FileRes, 0)

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorker := 5
	for i := 0; i < numWorker; i++ {
		go u.uploadWorker(ctx, client, jobsCh, resultCh, errCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errCh
		if err != nil {
			return nil, err
		}

		res := <-resultCh
		result = append(result, res)
	}

	return result, nil
}

func deleteFile(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to delete the file is aborted
	// if the object's generation number does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", object, err)
	}
	fmt.Fprintf(w, "Blob %v deleted.\n", object)
	return nil
}

func (u *FileUsecase) deleteFileWorker(ctx context.Context, client *storage.Client, jobs <-chan *file.DeleteFileReq, errCh chan<- error) {
	for job := range jobs {
		o := client.Bucket(u.cfg.App().GcpBucket()).Object(job.Destination)
		// Optional: set a generation-match precondition to avoid potential race
		// conditions and data corruptions. The request to delete the file is aborted
		// if the object's generation number does not match your precondition.
		attrs, err := o.Attrs(ctx)
		if err != nil {
			errCh <- fmt.Errorf("object.Attrs: %w", err)
			return
		}
		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		if err := o.Delete(ctx); err != nil {
			errCh <- fmt.Errorf("Object(%q).Delete: %w", job.Destination, err)
			return
		}
		fmt.Printf("Blob %v deleted.\n", job.Destination)
		errCh <- nil
	}
}

func (u *FileUsecase) DeleteImageGCP(req []*file.DeleteFileReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	jobsCh := make(chan *file.DeleteFileReq, len(req))
	errCh := make(chan error, len(req))
	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorker := 5
	for i := 0; i < numWorker; i++ {
		go u.deleteFileWorker(ctx, client, jobsCh, errCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errCh
		return err
	}

	return nil
}

// File Systems
// Upload to local File System
func (u *FileUsecase) uploadToLocalWorker(ctx context.Context, jobs <-chan *file.FileReq, results chan<- *file.FileRes, errs chan<- error) {
	for job := range jobs {
		cotainer, err := job.File.Open()
		if err != nil {
			errs <- err
			return
		}
		b, err := ioutil.ReadAll(cotainer)
		if err != nil {
			errs <- err
			return
		}

		// Upload an object to storage
		dest := fmt.Sprintf("./assets/images/%s", job.Destination)
		if err := os.WriteFile(dest, b, 0777); err != nil {
			if err := os.MkdirAll("./assets/images/"+strings.Replace(job.Destination, job.FileName, "", 1), 0777); err != nil {
				errs <- fmt.Errorf("mkdir \"./assets/images/%s\" failed: %v", err, job.Destination)
				return
			}
			if err := os.WriteFile(dest, b, 0777); err != nil {
				errs <- fmt.Errorf("write file failed: %v", err)
				return
			}
		}

		newFile := &filePublic{
			file: &file.FileRes{
				FileName: job.FileName,
				Url:      fmt.Sprintf("http://%s:%d/%s", u.cfg.App().Host(), u.cfg.App().Port(), job.Destination),
			},
			destination: job.Destination,
		}

		errs <- nil
		results <- newFile.file
	}
}

// Upload to local File System
func (u *FileUsecase) UploadImageLocal(req []*file.FileReq) ([]*file.FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	jobsCh := make(chan *file.FileReq, len(req))
	resultCh := make(chan *file.FileRes, len(req))
	errCh := make(chan error, len(req))

	res := make([]*file.FileRes, 0)
	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorker := 5
	for i := 0; i < numWorker; i++ {
		go u.uploadToLocalWorker(ctx, jobsCh, resultCh, errCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errCh
		if err != nil {
			return nil, err
		}

		result := <-resultCh
		res = append(res, result)

	}

	return res, nil
}

func (u *FileUsecase) deleteFromStorageFileWorkers(ctx context.Context, jobs <-chan *file.DeleteFileReq, errs chan<- error) {
	for job := range jobs {
		if err := os.Remove("./assets/images/" + job.Destination); err != nil {
			errs <- fmt.Errorf("remove file: %s failed: %v", job.Destination, err)
			return
		}
		errs <- nil
	}
}

func (u *FileUsecase) DeleteImageStorage(req []*file.DeleteFileReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	jobsCh := make(chan *file.DeleteFileReq, len(req))
	errsCh := make(chan error, len(req))

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.deleteFromStorageFileWorkers(ctx, jobsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		return err
	}
	return nil
}
