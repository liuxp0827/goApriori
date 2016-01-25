package Apriori

import (
	"testing"
)

func TestApriori(t *testing.T) {
	apriori := NewApriori()
	var err error
	apriori.record, err = apriori.getRecord("simple.txt")
	if err != nil {
		t.Error(err)
		return
	}



	cItemSet := apriori.findFirstCandidate()          // 获取第一次的备选集

	lItemSet := apriori.getSupportedItemSet(cItemSet) //获取备选集cItemset满足支持的集合

	for !apriori.endTag {
		ckItemset := apriori.getNextCandidate(lItemSet)

		lkItemset := apriori.getSupportedItemSet(ckItemset)
		apriori.getConfidencedItemset(lkItemset, lItemSet, apriori.dKCountMap, apriori.dCountMap) // 获取备选集cItemset满足置信度的集合

		if len(apriori.confItemSet) != 0 { // 满足置信度的集合不为空
			apriori.result(apriori.confItemSet) // 打印满足置信度的集合
		}
		apriori.confItemSet = apriori.confItemSet[:0]
		cItemSet = ckItemset
		lItemSet = lkItemset
		apriori.dCountMap = make(map[int]int)
		for k, v := range apriori.dKCountMap {
			apriori.dCountMap[k] = v
		}
	}
}
