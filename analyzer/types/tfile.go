package types

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

var tFileMethodsDescriptions map[string]string = map[string]string{
	"exists":       "file.exists() bool : Checks that the file exists",
	"isDir":        "file.isDir() bool : Checks that the file is a directory",
	"parentExists": "file.parentExists() bool : Checks that the parent directory exists",
	"size":         "file.size() int : Gets the size of the file in bytes",
	"perms":        "file.perms() string : Gets the permissions of the file as a string (e.g. \"0644\")",
	"user":         "file.user() string : Gets the user that owns the file",
	"group":        "file.group() string : Gets the group that owns the file",
	"toString":     "file.toString() string : Converts the value to a string",
}

type tFile struct {
	path string
}

func fileFactory(value interface{}) (IType, error) {
	if path, ok := value.(string); ok {
		return &tFile{path: path}, nil
	}

	return nil, fmt.Errorf("value is not a file path (string)")
}

func (t tFile) TypeName() string {
	return "file"
}

func (t tFile) Value() interface{} {
	return t.path
}

func (t tFile) GetMethod(method string) Method {
	tFileMethods := map[string]Method{
		"exists": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("file.exists expects 0 arguments")
			}

			// Check if the file exists
			exists, err := t.exists()
			if err != nil {
				return nil, err
			}

			if !exists {
				return &tBool{value: false}, fmt.Errorf("file does not exist")
			}

			return &tBool{value: true}, nil
		},
		"isDir": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("file.isDir expects 0 arguments")
			}

			// Check if the file exists
			exists, err := t.exists()
			if err != nil {
				return nil, err
			} else if !exists {
				return &tBool{value: false}, fmt.Errorf("file does not exist")
			}

			// Check if the file is a directory
			isDir, err := t.isDir()
			if err != nil {
				return nil, err
			}

			if !isDir {
				return &tBool{value: false}, fmt.Errorf("file is not a directory")
			}

			return &tBool{value: true}, nil
		},
		"parentExists": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("file.parentExists expects 0 arguments")
			}

			// Check if the file exists
			exists, err := t.exists()
			if err != nil {
				return nil, err
			} else if !exists {
				return &tBool{value: false}, fmt.Errorf("file does not exist")
			}

			// Check if the parent directory exists
			parentExists, err := t.parentExists()
			if err != nil {
				return nil, err
			}

			if !parentExists {
				return &tBool{value: false}, fmt.Errorf("parent directory does not exist")
			}

			return &tBool{value: true}, nil
		},
		"size": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("file.size expects 0 arguments")
			}

			// Check if the file exists
			exists, err := t.exists()
			if err != nil {
				return nil, err
			} else if !exists {
				return nil, fmt.Errorf("file does not exist")
			}

			// Get the size of the file
			size, err := t.size()
			if err != nil {
				return nil, err
			}

			return &tInt{value: int(size)}, nil
		},
		"perms": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("file.perms expects 0 arguments")
			}

			// Check if the file exists
			exists, err := t.exists()
			if err != nil {
				return nil, err
			} else if !exists {
				return nil, fmt.Errorf("file does not exist")
			}

			// Get the permissions of the file
			perms, err := t.perms()
			if err != nil {
				return nil, err
			}

			return &tString{value: fmt.Sprintf("%04o", perms&os.ModePerm)}, nil
		},
		"user": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("file.user expects 0 arguments")
			}

			// Check if the file exists
			exists, err := t.exists()
			if err != nil {
				return nil, err
			} else if !exists {
				return nil, fmt.Errorf("file does not exist")
			}

			// Get the user that owns the file
			user, err := t.user()
			if err != nil {
				return nil, err
			}

			return &tString{value: user}, nil
		},
		"group": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("file.group expects 0 arguments")
			}

			// Check if the file exists
			exists, err := t.exists()
			if err != nil {
				return nil, err
			} else if !exists {
				return nil, fmt.Errorf("file does not exist")
			}

			// Get the group that owns the file
			group, err := t.group()
			if err != nil {
				return nil, err
			}

			return &tString{value: group}, nil
		},
		"toString": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("file.toString expects 0 arguments")
			}

			return &tString{value: t.path}, nil
		},
	}

	if _, ok := tFileMethods[method]; !ok {
		return func(args []IType) (IType, error) {
			return nil, fmt.Errorf("file does not have method %s", method)
		}
	}

	return tFileMethods[method]
}

func (t tFile) exists() (bool, error) {
	if _, err := os.Stat(t.path); err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (t tFile) parentExists() (bool, error) {
	parentPath := filepath.Dir(t.path)
	if _, err := os.Stat(parentPath); err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (t tFile) isDir() (bool, error) {
	fileInfo, err := os.Stat(t.path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func (t tFile) size() (int64, error) {
	fileInfo, err := os.Stat(t.path)
	if err != nil {
		return -1, err
	}
	return fileInfo.Size(), nil
}

func (t tFile) perms() (os.FileMode, error) {
	fileInfo, err := os.Stat(t.path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Mode().Perm(), nil
}

func (t tFile) user() (string, error) {
	fileInfo, err := os.Stat(t.path)
	if err != nil {
		return "", err
	}
	uid := fileInfo.Sys().(*syscall.Stat_t).Uid
	u, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

func (t tFile) group() (string, error) {
	fileInfo, err := os.Stat(t.path)
	if err != nil {
		return "", err
	}
	gid := fileInfo.Sys().(*syscall.Stat_t).Gid
	g, err := user.LookupGroupId(strconv.Itoa(int(gid)))
	if err != nil {
		return "", err
	}
	return g.Name, nil
}
