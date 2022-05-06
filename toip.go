package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dlclark/regexp2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func ReadLine(filename string) ([]string, error) {

	var result []string
	// pip begins

	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 {

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		// pipe ends
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer func() {
			f.Close()
			// fmt.Println("文件关闭成功")
		}()
		reader := bufio.NewReader(f)

		for {
			// 这里文本最后一行读不到，需要处理,已处理
			// 或者最后一行置空
			// line,err :=buf.ReadString('\n')
			line, _, err := reader.ReadLine()

			if err != nil {
				if err == io.EOF {
					// 读取结束，报EOF

					// fmt.Println("读取结束")
					break
				}
				return nil, err
			}
			linestr := string(line)
			result = append(result, linestr)
		}
	}
	var temp []string
	for _, value := range result {

		// 处理两头空白字符
		value = strings.TrimSpace(value)
		// 抛弃空行
		expr := `^$`
		reg, _ := regexp2.Compile(expr, 0)
		if isMatch, _ := reg.MatchString(value); !isMatch {
			temp = append(temp, value)

		}
	}
	result = temp
	return result, nil
}

func UrlToIps(urls []string, keep_port bool) []string {
	// 处理URL为IP:PORT列表
	// 正则表达式稍有问题
	// URL后面如果没有/则无法匹配
	// http://123.123.123.123 无法匹配
	expr := `(?<=://).+?(?=/)`
	reg, _ := regexp2.Compile(expr, 0)

	// 对于URL后方无/的，主动添加/
	// http://123.123.123.123
	// 变为
	// http://123.123.123.123/
	// 同时针对单行为IP/域名的情况，主动添加http://xxxx/
	//
	expr_http_finder := `^(http://|https://)`
	reg2, _ := regexp2.Compile(expr_http_finder, 0)

	var temp []string
	for _, value := range urls {
		// 任何value，后方均➕/
		value = value + "/"
		// 查找开头是否为http://或https://，没有则加上

		if isMatch, _ := reg2.MatchString(value); !isMatch {
			value = "http://" + value
		}

		match, _ := reg.FindStringMatch(value)
		if !keep_port {
			ipPort := match.String()
			ipPort = strings.Split(ipPort, ":")[0]
			// 处理 IP:PORT列表为IP列表
			temp = append(temp, ipPort)
		} else {
			temp = append(temp, match.String())
		}

	}
	return temp

}

func IpWriteToFile(ips []string, save_to_file string) bool {
	var tempString string
	// var tempStringCsv string

	// expr := `.txt`
	if !strings.HasSuffix(save_to_file, ".txt") {
		save_to_file += ".txt"
	}
	for _, value := range ips {

		tempString += value + "\n"
		// tempStringCsv += value.Url + "," + value.Ip + "," + value.Location.CountryCode + "," + value.Location.RegionName + "\n"
		// 写入txt
		err := ioutil.WriteFile(save_to_file, []byte(tempString), 0666)
		if err != nil {
			fmt.Println(err)
		}
		// // 写入CSV文件
		// err1 := ioutil.WriteFile("results.csv", []byte(tempStringCsv), 0666)
		// if err1 != nil {
		// 	fmt.Println(err1)
		// }
	}
	return true
}

func DomainToIp(ips []string) ([]string, int, []string) {
	// 正则表达式排除ip，对域名进行处理
	expr := `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
	reg, err := regexp2.Compile(expr, 0)
	if err != nil {
		fmt.Println(err)
	}
	var tempIps []string
	var cdn_count int = 0
	var cdn_domains []string
	var temp_port string
	var temp_domain_port []string
	var contain_port bool = false
	for _, value := range ips {

		if isMatch, _ := reg.MatchString(value); !isMatch {

			// keep_port
			if strings.Contains(value, ":") {
				s := `true`
				contain_port, _ = strconv.ParseBool(s)
				temp_domain_port = strings.Split(value, ":")
				temp_port = temp_domain_port[1]
				value = temp_domain_port[0]
			}

			// Domain to ip
			ns, err := net.LookupHost(value)
			if err != nil {
				fmt.Println(err)
			}
			if contain_port {
				tempIps = append(tempIps, ns[0]+":"+temp_port)
			} else {
				tempIps = append(tempIps, ns[0])
			}

			if len(ns) > 1 {
				cdn_count += 1
				cdn_domains = append(cdn_domains, value)

			}
		} else {
			tempIps = append(tempIps, value)
		}
	}
	return tempIps, cdn_count, cdn_domains
	// return tempIps
}

func main() {
	// 处理URL成为IP:PORT
	// 或IP形式
	url_file_path := flag.String("l", "url.txt", "url文件路径")

	save_to_file := flag.String("o", "", "保存到文件")
	domain_to_ip := flag.Bool("n", false, "对域名进行DNS解析")
	silence := flag.Bool("q", false, "直接输出IP/域名，其他内容不输出")
	keep_port := flag.Bool("p", false, "保留PORT")

	// save_region_code := flag.String("s", "", "要保存的IP归属省份")

	flag.Parse()

	// *url_file_path = strings.ToLower(*url_file_path)
	// *save_to_file = strings.ToUpper(*save_to_file)
	// *save_region_code = strings.ToUpper(*save_region_code)

	urls, err := ReadLine(*url_file_path)
	if err != nil {
		fmt.Println(err)
	}
	misc := UrlToIps(urls, *keep_port)

	var ips []string = misc
	args := os.Args[1:]

	var (
		is_domian_to_ip     bool
		are_we_need_silence bool
	)

	for _, value := range args {
		if value == "-n" {
			s := "true"
			is_domian_to_ip, _ = strconv.ParseBool(s)
			break
		}

	}
	for _, value := range args {
		if value == "-q" {
			s := "true"
			are_we_need_silence, _ = strconv.ParseBool(s)
			break
		}

	}

	var cdn_count int
	var cdn_domains []string

	if is_domian_to_ip {

		ips, cdn_count, cdn_domains = DomainToIp(misc)
	} else {
		ips = misc
	}

	if *domain_to_ip {
		//强行使用flag
	}

	if *silence {
		//强行使用flag
	}

	if !are_we_need_silence {

		// fmt.Println("\n")

		fmt.Println("url总数为: ", len(ips), " 个")

	}
	for _, value := range ips {
		fmt.Println(value)
	}

	if !are_we_need_silence {
		fmt.Println("\n")
		if is_domian_to_ip {
			if cdn_count != 0 {
				fmt.Println("以下域名使用了CDN，注意识别")
				for _, value := range cdn_domains {
					fmt.Println(value)
				}
			}

		}

		fmt.Println("\n\n")
		if *save_to_file != "" {
			IpWriteToFile(ips, *save_to_file)
			fmt.Println("ip写入完毕。。。")
		}
		fmt.Println("\n")

	} else {

		if *save_to_file != "" {
			IpWriteToFile(ips, *save_to_file)

		}
	}

}
