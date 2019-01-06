package hpcache

type FileTime struct {
	CreateTime int64
	AccessTime int64
	ModifyTime int64
}

func GetFileTime(file string) (*FileTime, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		return &FileTime{
			CreateTime: time.Unix(stat.Atim.Sec, stat.Atim.Nsec).UnixNano(),
			AccessTime: time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec).UnixNano(),
			ModifyTime: time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec).UnixNano(),
		}, nil
	}
	return nil, fmt.Errorf("not support file info in current platform")
}
