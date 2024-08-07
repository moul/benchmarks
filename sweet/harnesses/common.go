// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harnesses

import (
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/benchmarks/sweet/common/fileutil"
	"golang.org/x/benchmarks/sweet/common/log"
)

func gitShallowClone(dir, url, ref string) error {
	// Git 2.46+ has a global --no-advice flag, but that's extremely recent as of this writing.
	cmd := exec.Command("git", "-c", "advice.detachedHead=false", "clone", "--depth", "1", "-b", ref, url, dir)
	log.TraceCommand(cmd, false)
	cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	return err
}

func gitRecursiveCloneToCommit(dir, url, branch, hash string) error {
	cloneCmd := exec.Command("git", "clone", "--recursive", "--shallow-submodules", "-b", branch, url, dir)
	log.TraceCommand(cloneCmd, false)
	cloneCmd.Stderr = os.Stderr
	if _, err := cloneCmd.Output(); err != nil {
		return err
	}
	checkoutCmd := exec.Command("git", "-C", dir, "-c", "advice.detachedHead=false", "checkout", hash)
	log.TraceCommand(checkoutCmd, false)
	checkoutCmd.Stderr = os.Stderr
	_, err := checkoutCmd.Output()
	return err
}

func gitCloneToCommit(dir, url, branch, hash string) error {
	cloneCmd := exec.Command("git", "clone", "-b", branch, url, dir)
	log.TraceCommand(cloneCmd, false)
	cloneCmd.Stderr = os.Stderr
	if _, err := cloneCmd.Output(); err != nil {
		return err
	}
	checkoutCmd := exec.Command("git", "-C", dir, "-c", "advice.detachedHead=false", "checkout", hash)
	log.TraceCommand(checkoutCmd, false)
	checkoutCmd.Stderr = os.Stderr
	_, err := checkoutCmd.Output()
	return err
}

func copyFile(dst, src string) error {
	log.CommandPrintf("cp %s %s", src, dst)
	return fileutil.CopyFile(dst, src, nil, nil)
}

func makeWriteable(dir string) error {
	log.CommandPrintf("chmod -R a+w %s", dir)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode()&0222 == 0222 {
			return nil
		}
		return os.Chmod(path, info.Mode()|0222)
	})
}

func symlink(dst, src string) error {
	log.CommandPrintf("ln -s %s %s", src, dst)
	return os.Symlink(src, dst)
}
