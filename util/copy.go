package util

import "os"
import "io"
import "io/ioutil"
import "log"
import "path"

// Copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}
	}

	return
}

// Recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return &CustomError{"Source is not a directory"}
	}

	// create only if necessary

	_, err = os.Open(dest)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dest, fi.Mode())
		if err != nil {
			return err
		}
	}

	entries, err := ioutil.ReadDir(source)

	for _, entry := range entries {
		if path.Ext(entry.Name()) != ".go" {
			sfp := source + "/" + entry.Name()
			dfp := dest + "/" + entry.Name()
			if entry.IsDir() {
				if entry.Name() != ".recime" && entry.Name() != "node_modules"{
					err = CopyDir(sfp, dfp)
					if err != nil {
						log.Println(err)
					}
				}
			} else {
				// perform copy
				err = CopyFile(sfp, dfp)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	return
}

// A struct for returning custom error messages
type CustomError struct {
	What string
}

// Returns the error message defined in What as a string
func (e *CustomError) Error() string {
	return e.What
}
