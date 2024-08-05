package scanner

import (
	"context"
	"log"
	"os"
	//"fmt"
	"strings"
	"path/filepath"
	"path"

	"github.com/karrick/godirwalk"
	"golang.org/x/sync/errgroup"
)

func skip(needle string, haystack []string) bool {
	for _, f := range haystack {
		if f == needle {
			return true
		}
	}
	return false
}

// walkone descends a single directory tree looking for git repos
func walkone(ctx context.Context, dir string, config *Config, results chan string) error {
	err := godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted:            true,
		ScratchBuffer:       make([]byte, godirwalk.MinimumScratchBufferSize),
		FollowSymbolicLinks: config.FollowSymlinks,
		ErrorCallback: func(path string, err error) godirwalk.ErrorAction {
			patherr, ok := err.(*os.PathError)
			if ok {
				switch patherr.Unwrap().Error() {
				case "no such file or directory":
					// might be symlink pointing to non-existent file
					return godirwalk.SkipNode

				case "too many levels of symbolic links":
					// skip invalid symlinks
					return godirwalk.SkipNode
				}
			}
			log.Printf("ERROR: %s: %v", path, err)
			return godirwalk.Halt
		},
		Callback: func(path string, ent *godirwalk.Dirent) error {

			// early exit?

			select {
			case <-ctx.Done():
				return filepath.SkipDir
			default:
			}

			// process all the SkipThis rules first

			if skip(path, config.ScanDirs.Exclude) {
				return godirwalk.SkipThis
			}
			if ent.IsSymlink() && !config.FollowSymlinks {
				return godirwalk.SkipThis
			}

			// then process non-matching rules which still descend

			if ent.Name() != ".git" {
				return nil
			}
			isDir, _ := ent.IsDirOrSymlinkToDir()
			if !isDir {
				return nil
			}

			results <- filepath.Dir(path)
			return godirwalk.SkipThis // don't descend further
		},
	})
	return err
}

// Walk finds all git repositories in the directories specified in config
func Walk(ctx context.Context, config *Config, results chan string, ignore_dir_errors bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	completeIncludeList := config.ScanDirs.Include

	var errors errgroup.Group

	for i := range config.ScanDirs.Include {
		j := i // copy loop variable
		globPath := config.ScanDirs.Include[j]

		if(string(globPath[len(globPath)-1:]) == "*"){
			log.Printf("GLOBPATH: %s", globPath)
			parent := filepath.Dir(globPath)
			//log.Printf("PARENT: %s", parent)
			//log.Printf("BASEGLOB: %s", path.Base(globPath))
			//log.Printf("BASEGLOB2: %s", path.Base(globPath[0:len(globPath)-1]))
			baseGlob := path.Base(globPath[0:len(globPath)-1])

			entries, err := os.ReadDir(parent)
			if err != nil {
				log.Fatal(err)
			}

			for _, e := range entries {
				if strings.HasPrefix(e.Name(), baseGlob) {
					completeIncludeList = append(completeIncludeList, parent + "/" + e.Name())
					//fmt.Println(e.Name())
				}
			}

		}
	}

	//fmt.Printf("%v", completeIncludeList)
	for i := range completeIncludeList {
		j := i // copy loop variable
		globPath := completeIncludeList[j]

		errors.Go(func() error {
			err := walkone(ctx, globPath, config, results)
			if err == filepath.SkipDir {
				cancel()
			} else if err != nil {
				if ignore_dir_errors {
					log.Printf("ERROR: %s: %v", j, err)
					return nil
				} else {
					return err
				}
			}
			return nil
		})
	}

	err := errors.Wait()
	close(results)
	return err
}
