package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode"
)

type paserStat struct {
	CurrSection    string
	InSection      bool
	IsEOF          bool
	TmpData        map[string]string
	TmpSectionData map[string]map[string]string
}

func (c *Config) Reset(opts Options) error {
	isAuto := c.opts.AutoReload

	c.opts.AutoReload = opts.AutoReload
	if opts.CheckDuration > 0 {
		c.opts.CheckDuration = opts.CheckDuration
	}

	if !isAuto && c.opts.AutoReload {
		go c.autoReload()
	}
	return nil
}

func (c *Config) Get(key string) (string, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *Config) SectionGet(section, key string) (string, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	val, ok := c.sectionData[section][key]
	return val, ok
}

func (c *Config) SectionsGet(section string) map[string]string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	tmp := make(map[string]string)
	if val, ok := c.sectionData[section]; ok {
		for k, v := range val {
			tmp[k] = v
		}
	}
	return tmp
}

func (c *Config) new() error {
	f, err := os.Open(c.opts.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = c.parse(f)
	if err != nil {
		return err
	}

	fInfo, err := f.Stat()
	if err != nil {
		return err
	}

	c.Name = fInfo.Name()
	c.modTime = fInfo.ModTime().Unix()

	conf_map_lock.Lock()
	defer conf_map_lock.Unlock()

	if _, ok := conf_map[c.Name]; ok {
		return fmt.Errorf("Repeat filename '%s' on Config.new()", c.Name)
	}
	conf_map[c.Name] = c

	return nil
}

func (c *Config) reload() error {
	f, err := os.Open(c.opts.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = c.parse(f)
	if err != nil {
		return err
	}
	fInfo, err := f.Stat()
	if err != nil {
		return err
	}
	c.modTime = fInfo.ModTime().Unix()

	return nil
}

func (c *Config) parse(f io.Reader) error {
	reader := bufio.NewReader(f)

	ps := &paserStat{
		InSection:      false,
		IsEOF:          false,
		TmpData:        make(map[string]string),
		TmpSectionData: make(map[string]map[string]string),
	}

	//循环读文件每行
	for !ps.IsEOF {
		line, err := reader.ReadBytes('\n')

		if err != nil {
			if err == io.EOF {
				ps.IsEOF = true
			} else {
				return err
			}
		}

		line = bytes.TrimLeftFunc(line, unicode.IsSpace)

		//注释跳过
		if len(line) == 0 || line[0] == '#' || line[0] == ';' {
			continue
		}

		//Section
		if line[0] == '[' {
			closeIdx := bytes.LastIndexByte(line, ']')
			if closeIdx == -1 {
				return fmt.Errorf("unclosed section: %s", line)
			}

			sectionName := string(line[1:closeIdx])
			ps.CurrSection = sectionName
			ps.InSection = true
			ps.TmpSectionData[sectionName] = make(map[string]string)
			continue
		}

		//解析key-val
		lineStr := strings.TrimSpace(string(line))

		idx := strings.Index(lineStr, "=")
		if idx == -1 {
			return fmt.Errorf("invalid key-value: %s", lineStr)
		}

		key := strings.TrimSpace(lineStr[:idx])
		val := strings.TrimSpace(lineStr[idx+1:])
		if len(key) == 0 {
			continue
		}

		if ps.InSection {
			ps.TmpSectionData[ps.CurrSection][key] = val
		} else {
			ps.TmpData[key] = val
		}
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.data = ps.TmpData
	c.sectionData = ps.TmpSectionData
	return nil

}

func (c *Config) autoReload() {
	if !c.opts.AutoReload {
		return
	}
	defer func() {
		go c.autoReload()
	}()
	time.Sleep(c.opts.CheckDuration)

	c.autoLock.Wait()
	c.autoLock.Add(1)
	defer c.autoLock.Done()

	fInfo, err := os.Stat(c.opts.FilePath)
	if err != nil {
		return
	}

	if fInfo.ModTime().Unix() > c.modTime {
		err := c.reload()
		fmt.Println("reload done", err)
	}

}
