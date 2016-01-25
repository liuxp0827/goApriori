package Apriori

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	MIN_SUPPORT float64 = 0.2 // 最小支持度
	MIN_CONF    float64 = 0.8 // 最小置信度
)

type Apriori struct {
	endTag      bool
	dCountMap   map[int]int // k-1 频繁集的记数表
	dKCountMap  map[int]int // k 频繁集的记数表
	record      [][]string  // 数据记录表
	lable       int         // 用于输出时的一个标记, 记录当前在打印第几级关联集
	confCount   []float64   // 置信度记录表
	confItemSet [][]string  // 满足支持度的集合
}

func NewApriori() *Apriori {
	return &Apriori{
		endTag:      false,
		dCountMap:   make(map[int]int),
		dKCountMap:  make(map[int]int),
		record:      make([][]string, 0),
		lable:       1,
		confCount:   make([]float64, 0),
		confItemSet: make([][]string, 0),
	}
}

func (this *Apriori) result(confItemset2 [][]string) {
	log.Println("********* 频繁模式挖掘结果 ***********")
	for i := 0; i < len(confItemset2); i++ {
		var ret string
		j := 0
		for j = 0; j < len(confItemset2[i])-3; j++ {
			ret +=confItemset2[i][j]+" "
		}

		ret+="---> "
		ret+=confItemset2[i][j]+" "
		j++
		ret+="相对支持度: "+confItemset2[i][j]
		j++
		ret+=", 自信度: "+confItemset2[i][j]
		j++
		log.Println(ret)
	}

}

// 根据lkItemset，lItemset，dkCountMap2，dCountMap2求出满足自信度的集合
func (this *Apriori) getConfidencedItemset(lkItemSet, lItemSet [][]string, dkCountMap2, dCountMap2 map[int]int) {
	for i := 0; i < len(lkItemSet); i++ {
		this.getConfItem(lkItemSet[i], lItemSet, dkCountMap2[i], dCountMap2)
	}
}

// 检验集合 list 是否满足最低自信度要求
// 若满足则在全局变量 confItemset 添加 list
func (this *Apriori) getConfItem(list []string, lItemSet [][]string, count int, dCountMap2 map[int]int) {
	for i := 0; i < len(list); i++ {
		testList := make([]string, 0)
		for j := 0; j < len(list); j++ {
			if i != j {
				testList = append(testList, list[j])
			}
		}
		index := this.findConf(testList, lItemSet) //查找testList中的内容在lItemset的位置
		conf := float64(count) * 1.0 / float64(dCountMap2[index])
		if conf > MIN_CONF { //满足自信度要求
			testList = append(testList, list[i])
			relativeSupport := float64(count) * 1.0 / float64(len(this.record)-1)
			testList = append(testList, strconv.FormatFloat(relativeSupport, 'f', -1, 64))
			testList = append(testList, strconv.FormatFloat(conf, 'f', -1, 64))
			this.confItemSet = append(this.confItemSet, testList)
		}
	}

}

// 查找testList中的内容在lItemset的位置
func (this *Apriori) findConf(testList []string, lItemSet [][]string) int {
	for i := 0; i < len(lItemSet); i++ {
		notHaveTag := false
		for j := 0; j < len(testList); j++ {
			if !this.haveThisItem(testList[j], lItemSet[i]) {
				notHaveTag = true
				break
			}
		}

		if !notHaveTag {
			return i
		}
	}
	return -1
}

//检验 list 中是否包含 str
func (this *Apriori) haveThisItem(str string, list []string) bool {
	for i := 0; i < len(list); i++ {
		if str == list[i] {
			return true
		}
	}
	return false
}

// 获取数据库记录
func (this *Apriori) getRecord(filename string) ([][]string, error) {
	record := make([][]string, 0)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil || err == io.EOF {
			break
		}
		log.Println(string(line))
		lineStr := strings.Split(string(line), "  ")
		lineList := make([]string, 0)
		for i := 0; i < len(lineStr); i++ {
			if strings.HasSuffix(lineStr[i], "T") ||
				strings.HasSuffix(lineStr[i], "YES") {
				lineList = append(lineList, record[0][i])
			} else if strings.HasSuffix(lineStr[i], "F") ||
				strings.HasSuffix(lineStr[i], "NO") {

			} else {
				lineList = append(lineList, lineStr[i])
			}
		}
		record = append(record, lineList)
	}

	file.Close()
	return record, nil
}

// 求出cItemset中满足最低支持度集合
func (this *Apriori) getSupportedItemSet(cItemSet [][]string) [][]string {
	end := true
	supportedItemset := make([][]string, 0)
	k := 0
	for i := 0; i < len(cItemSet); i++ {
		count := this.countFrequent(cItemSet[i])                   //统计记录数
		if count >= int(MIN_SUPPORT*float64(len(this.record)-1)) { // count值大于支持度与记录数的乘积，即满足支持度要求
			if len(cItemSet[0]) == 1 {
				this.dCountMap[k] = count
				k++
			} else {
				this.dKCountMap[k] = count
				k++
			}
			supportedItemset = append(supportedItemset, cItemSet[i])
			end = false
		}
	}
	this.endTag = end
	return supportedItemset
}

// 统计数据库记录record中出现list中的集合的个数
func (this *Apriori) countFrequent(list []string) int {
	count := 0
	for i := 1; i < len(this.record); i++ {
		notHavaThisList := false
		for k := 0; k < len(list); k++ {
			thisRecordHave := false
			for j := 1; j < len(this.record[i]); j++ {
				if list[k] == this.record[i][j] {
					thisRecordHave = true
				}
			}
			if !thisRecordHave { // 扫描一遍记录表的一行,
				// 发现 list[i] 不在记录表的第j行中，即list不可能在 j 行中
				notHavaThisList = true
				break
			}
		}
		if !notHavaThisList {
			count++
		}
	}
	return count
}

// 根据cItemset求出下一级的备选集合组,
// 求出的备选集合组中的每个集合的元素的个数比cItemset中的集合的元素大1
func (this *Apriori) getNextCandidate(cItemSet [][]string) [][]string {
	nextItemSet := make([][]string, 0)
	for i := 0; i < len(cItemSet); i++ {
		tempList := make([]string, 0)
		for k := 0; k < len(cItemSet[i]); k++ {
			tempList = append(tempList, cItemSet[i][k])
		}

		for h := i + 1; h < len(cItemSet); h++ {
			for j := 0; j < len(cItemSet[h]); j++ {

				tempList = append(tempList, cItemSet[h][j])
				if this.isSubSetInc(tempList, cItemSet) { // tempList的子集全部在cItemset中
					copyValueHelpList := make([]string, 0)
					for p := 0; p < len(tempList); p++ {
						copyValueHelpList = append(copyValueHelpList, tempList[p])
					}
					if this.isHave(copyValueHelpList, nextItemSet) {
						nextItemSet = append(nextItemSet, copyValueHelpList)
					}
				}
				tempList = tempList[:len(tempList)-1]
			}
		}
	}
	return nextItemSet
}

// 检验 nextItemset 中是否包含 copyValueHelpList
func (this *Apriori) isHave(copyValueHelpList []string, nextItemset [][]string) bool {
	for i := 0; i < len(nextItemset); i++ {
		length := len(copyValueHelpList)
		if length == len(nextItemset[i]) {
			ii := 0
			for ; ii < length; ii++ {
				if copyValueHelpList[ii] != nextItemset[i][ii] {
					break
				}
			}

			if ii == length {
				return false
			}
		}
	}
	return true
}

// 检验 tempList 是不是 cItemset 的子集
func (this *Apriori) isSubSetInc(tempList []string, cItemSet [][]string) bool {
	haveTag := false
	for i := 0; i < len(tempList); i++ { // k 集合 tempList 的子集是否都在 k-1 级频繁级中
		testList := make([]string, 0)
		for j := 0; j < len(tempList); j++ {
			if i != j {
				testList = append(testList, tempList[j])
			}
		}

		length := len(testList)
		for k := 0; k < len(cItemSet); k++ {
			/***** 子集存在于k-1频繁集中 *****/
			if length == len(cItemSet[k]) {
				ii := 0
				for ; ii < length; ii++ {
					if testList[ii] != cItemSet[k][ii] {
						break
					}
				}

				if ii == length {
					haveTag = true
					// 子集存在于k-1频繁集中
					break
				}
			}
			/***** 子集存在于k-1频繁集中 *****/
		}

		if !haveTag {
			return false
		}
	}

	return haveTag
}

// 根据数据库记录求出第一级备选集
func (this *Apriori) findFirstCandidate() [][]string {
	tableList := make([][]string, 0)
	lineList := make([]string, 0)
	size := 0
	for i := 1; i < len(this.record); i++ {
		for j := 1; j < len(this.record[i]); j++ {
			if len(lineList) == 0 {
				lineList = append(lineList, this.record[i][j])
			} else {
				haveThisItem := false
				size = len(lineList)
				for k := 0; k < size; k++ {
					if lineList[k] == this.record[i][j] {
						haveThisItem = true
						break
					}
				}

				if !haveThisItem {
					lineList = append(lineList, this.record[i][j])
				}
			}
		}
	}

	for i := 0; i < len(lineList); i++ {
		helpList := make([]string, 0)
		helpList = append(helpList, lineList[i])
		tableList = append(tableList, helpList)
	}

	return tableList
}
