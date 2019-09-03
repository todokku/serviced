package stats

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func Str2UInt32(number string) (num uint32, err error) {
	num64, err := strconv.ParseUint(number, 10, 32)
	if err == nil {
		num = uint32(num64)
	}
	return
}

func UserIDLookup(userName string) (uid uint32, err error) {
	userInfo, err := user.Lookup(userName)
	if err != nil {
		return
	}
	uid, err = Str2UInt32(userInfo.Uid)
	return
}

func GetLimits(username string) (soft uint32, hard uint32, err error) {
	do := func(cmd string) (limit uint32, err error) {
		c := exec.Command("su", "-", username, "-c", "/bin/bash", "-c", cmd)
		out, err := c.Output()
		if err != nil {
			return
		}
		limitStr := strings.TrimSpace(string(out))
		if limitStr == "unlimited" {
			//limit = 1<<32 - 1
			limit = 0
		} else {
			limit, err = Str2UInt32(limitStr)
		}
		return
	}

	soft, err = do("ulimit -u -S")
	hard, err1 := do("ulimit -u -H")
	if err == nil && err1 != nil {
		err = err1
	}
	return
}

// GetNumberOfThreads take PID of process and returns the number of its threads.
// This code is equivalent to `awk '{ if ($1 == "Threads:") print($2)}' /proc/$pid/status`
func GetNumberOfThreads(pid uint32) (num uint32, err error) {
	f, err := os.Open("/proc/" + strconv.FormatUint(uint64(pid), 10) + "/status")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		txt := s.Text()
		if strings.HasPrefix(txt, "Threads:") {
			words := strings.Fields(txt)
			if len(words) == 2 {
				num, err = Str2UInt32(words[1])
				return
			}
		}
	}
	return 0, errors.New("no data, no error, maybe something wrong with the algorithm")
}

// TotalThreadsNumber returns the number of threads of all processes started by the user $uid
func TotalThreadsNumberByUser(uid uint32) (threads uint32, err error) {
	d, err := os.Open("/proc")
	if err != nil {
		return
	}
	fsObjects, err := d.Readdir(-1)
	err1 := d.Close()
	if err != nil || err1 != err1 {
		return
	}

	for _, obj := range fsObjects {
		pid, err := Str2UInt32(obj.Name())
		stat, ok := obj.Sys().(*syscall.Stat_t)
		if err == nil && obj.IsDir() && ok && stat.Uid == uid {
			thr, err := GetNumberOfThreads(pid)
			if err == nil {
				threads += thr
			}
		}
	}
	return
}
